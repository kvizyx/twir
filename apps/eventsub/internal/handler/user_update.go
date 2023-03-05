package handler

import (
	"encoding/json"
	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/tsuwari/libs/pubsub"
	"go.uber.org/zap"
)

func (c *handler) handleUserUpdate(h *eventsub_bindings.ResponseHeaders, event *eventsub_bindings.EventUserUpdate) {
	bytes, err := json.Marshal(&pubsub.UserUpdateMessage{
		UserID:      event.UserID,
		UserLogin:   event.UserLogin,
		UserName:    event.UserName,
		Description: event.Description,
	})
	if err != nil {
		zap.S().Error(err)
	}

	c.services.PubSub.Publish("user.update", bytes)
}
