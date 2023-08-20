package timers

import (
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func NewGorm(db *gorm.DB) Repository {
	return &gormRepository{
		db,
	}
}

type gormRepository struct {
	db *gorm.DB
}

func (c *gormRepository) convertEntity(entity model.ChannelsTimers) Timer {
	result := Timer{
		ID:        entity.ID,
		Name:      entity.Name,
		ChannelID: entity.ChannelID,
		Interval:  int(entity.TimeInterval),
	}

	for _, r := range entity.Responses {
		result.Responses = append(
			result.Responses,
			TimerResponse{
				ID:         r.ID,
				Text:       r.Text,
				IsAnnounce: r.IsAnnounce,
			},
		)
	}

	return result
}

func (c *gormRepository) GetById(id string) (Timer, error) {
	entity := model.ChannelsTimers{}
	result := Timer{}

	if err := c.db.Where("id = ?", id).Preload("Responses").Find(&entity).Error; err != nil {
		return result, err
	}

	if entity.ID == "" {
		return result, NotFoundError
	}

	return c.convertEntity(entity), nil
}

func (c *gormRepository) GetAll() ([]Timer, error) {
	var timers []model.ChannelsTimers
	if err := c.db.
		Preload("Responses").
		Preload("Channel").
		Where(`enabled = ?`, true).
		Find(&timers).Error; err != nil {
		return nil, err
	}

	result := make([]Timer, 0, len(timers))
	for _, t := range timers {
		if !t.Channel.IsEnabled {
			continue
		}

		result = append(result, c.convertEntity(t))
	}

	return result, nil
}