package contracts

import (
	"github.com/google/uuid"
	"time"
)

type Vote struct {
	ID            uuid.UUID         `json:"id"`
	Name          string            `json:"name"`
	CreatedAt     time.Time         `json:"created_at"`
	EndDate       time.Time         `json:"end_date"`
	IsPublic      bool              `json:"is_public"`
	Options       map[string]string `json:"options"`
	DisplayResult map[string]int8   `json:"display_result"`
}

type CreateVote struct {
	Name     string            `json:"name"`
	EndDate  time.Time         `json:"end_date"`
	IsPublic bool              `json:"is_public"`
	Options  map[string]string `json:"options"`
}
