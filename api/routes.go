package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"vote-app/contracts"
	"vote-app/persistance"
)

type VoteController struct {
	db *persistance.RedisCache
}

func AddVoteController(db *persistance.RedisCache) *VoteController {
	return &VoteController{
		db: db,
	}
}

func (h *VoteController) CreateVote(c *gin.Context) {
	var createVote contracts.CreateVote

	if err := c.ShouldBindBodyWithJSON(&createVote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	vote := contracts.Vote{
		ID:            uuid.New(),
		Name:          createVote.Name,
		CreatedAt:     time.Now().UTC(),
		EndDate:       createVote.EndDate,
		IsPublic:      createVote.IsPublic,
		Options:       createVote.Options,
		DisplayResult: make(map[string]int8),
	}

	for key := range createVote.Options {
		vote.DisplayResult[key] = 0
	}

	if err := h.db.CreateVote(&vote); err != nil {
		log.Printf("Error storing vote: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "vote not created"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": vote})
	return
}

func (h *VoteController) GetVotes(c *gin.Context) {

	votes, err := h.db.GetVotes()
	if err != nil {
		log.Printf("Error getting votes: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "no votes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": votes})

	return
}
