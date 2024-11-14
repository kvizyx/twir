package asynq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/parser-new/internal/task-queue"
	"github.com/satont/twir/libs/twitch"
)

const TaskModUser = "task:mod_user"

type TaskModUserPayload struct {
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
}

func (tp *TaskProcessor) ProcessDistributeModUser(
	ctx context.Context,
	task taskqueue.Task,
) error {
	var payload TaskModUserPayload

	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("unmarshal task payload: %w", err)
	}

	twitchClient, err := twitch.NewUserClientWithContext(
		ctx,
		payload.ChannelID,
		tp.config,
		tp.tokensClient,
	)
	if err != nil {
		return fmt.Errorf("create twitch client: %w", err)
	}

	res, err := twitchClient.AddChannelModerator(
		&helix.AddChannelModeratorParams{
			BroadcasterID: payload.ChannelID,
			UserID:        payload.UserID,
		},
	)
	if err != nil {
		return fmt.Errorf("request to add channel moderator: %w", err)
	}

	if len(res.ErrorMessage) != 0 {
		return fmt.Errorf("add channel moderator: %s", res.ErrorMessage)
	}

	return nil
}

func (td *TaskDistributor) DistributeModUser(
	ctx context.Context,
	payload *taskqueue.TaskModUserPayload,
	opts ...taskqueue.Option,
) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	task := asynq.NewTask(TaskModUser, payloadBytes, td.fromOptions(opts...)...)

	taskInfo, err := td.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("enqueue task: %w", err)
	}

	td.logger.Info(
		"sent task to distribute mod user",
		slog.String("task_id", taskInfo.ID),
	)

	return nil
}
