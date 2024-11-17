package entity

import (
	"github.com/guregu/null"
)

type Alert struct {
	ID           string
	ChannelID    string
	Name         string
	AudioID      null.String
	AudioVolume  int
	CommandsIDs  []CommandID
	RewardsIDs   []string
	GreetingsIDs []string
	KeywordsIDs  []string
}
