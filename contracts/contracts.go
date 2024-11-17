package contracts

import (
	"github.com/google/uuid"
	"time"
)

type Vote struct {
	ID        uuid.UUID         `json:"id"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"created_at"`
	EndDate   time.Time         `json:"end_date"`
	IsPublic  bool              `json:"is_public"`
	Options   map[string]string `json:"options"`
}

type CreateVote struct {
	Name     string            `json:"name"`
	EndDate  time.Time         `json:"end_date"`
	IsPublic bool              `json:"is_public"`
	Options  map[string]string `json:"options"`
}

type VoteStats struct {
	Name    string         `json:"name"`
	Options []VoteStatItem `json:"options"`
}

type VoteStatItem struct {
	OptionId   string `json:"option_id"`
	OptionName string `json:"option_name"`
	Votes      int    `json:"votes"`
}
