package handler

import (
	"time"
	"tsuwari/timers/internal/types"
	"tsuwari/twitch"

	model "tsuwari/models"

	"github.com/go-co-op/gocron"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nicklaw5/helix"
	"github.com/satont/tsuwari/nats/bots"
	"github.com/satont/tsuwari/nats/parser"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	twitch *twitch.Twitch
	nats   *nats.Conn
	db     *gorm.DB
	logger *zap.Logger
	store  types.Store
}

func New(
	twitch *twitch.Twitch,
	nats *nats.Conn,
	db *gorm.DB,
	logger *zap.Logger,
	store types.Store,
) *Handler {
	return &Handler{twitch: twitch, nats: nats, db: db, logger: logger, store: store}
}

func (c *Handler) Handle(j gocron.Job) {
	t := c.store[j.Tags()[0]]

	streamData := model.ChannelsStreams{}

	err := c.db.Where(`"userId" = ?`, t.Model.ChannelID).First(&streamData).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if t.Model.MessageInterval > 0 &&
		t.Model.LastTriggerMessageNumber-int32(
			streamData.ParsedMessages,
		)+t.Model.MessageInterval > 0 {
		return
	}

	users, err := c.twitch.Client.GetUsers(&helix.UsersParams{
		IDs: []string{t.Model.ChannelID},
	})

	if err != nil || len(users.Data.Users) == 0 {
		return
	}

	user := users.Data.Users[0]

	rawMessage := t.Model.Responses[t.SendIndex]

	requestBytes, protoError := proto.Marshal(&parser.ParseResponseRequest{
		Sender: &parser.Sender{
			Id:          "",
			Name:        "bot",
			DisplayName: "Bot",
			Badges:      []string{"BROADCASTER"},
		},
		Channel: &parser.Channel{Id: user.ID, Name: user.Login},
		Message: &parser.Message{Text: rawMessage},
	})
	if protoError != nil {
		c.logger.Sugar().Error(err)
		return
	}

	response, natsError := c.nats.Request("parser.parseTextResponse", requestBytes, 5*time.Second)
	if natsError != nil {
		c.logger.Sugar().Error(err)
		return
	}
	responseData := parser.ParseResponseResponse{}

	err = proto.Unmarshal(response.Data, &responseData)

	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	botsRequest := bots.SendMessage{
		ChannelId:   user.ID,
		ChannelName: &user.Login,
		Message:     rawMessage,
	}
	bytes, _ := proto.Marshal(&botsRequest)
	c.nats.Publish("bots.sendMessage", bytes)

	nextIndex := t.SendIndex + 1

	if nextIndex+1 <= len(t.Model.Responses) {
		t.SendIndex = nextIndex
	} else {
		t.SendIndex = 0
	}

	t.Model.LastTriggerMessageNumber = int32(streamData.ParsedMessages)

	err = c.db.
		Model(&model.ChannelsTimers{}).
		Where(`"id" = ?`, t.Model.ID).
		Updates(model.ChannelsTimers{LastTriggerMessageNumber: int32(streamData.ParsedMessages)}).
		Error

	if err != nil {
		c.logger.Sugar().Error(err)
	}
}