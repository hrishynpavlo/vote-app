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
	db *persistance.Storage
}

func AddVoteController(db *persistance.Storage) *VoteController {
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
		ID:        uuid.New(),
		Name:      createVote.Name,
		CreatedAt: time.Now().UTC(),
		EndDate:   createVote.EndDate,
		IsPublic:  createVote.IsPublic,
		Options:   createVote.Options,
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

	votes, err := h.db.Redis.GetVotes()
	if err != nil {
		log.Printf("Error getting votes: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": "no votes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": votes})

	return
}

func (h *VoteController) GetVoteStats(c *gin.Context) {
	id := c.Param("id")
	voteStats, err := h.db.Redis.GetVoteStats(id)
	if err != nil {
		log.Printf("Error getting vote stats: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "no vote stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": voteStats})
	return
}

func (h *VoteController) Vote(c *gin.Context) {
	id := c.Param("id")
	optionID := c.Param("optionId")
	voteStats, err := h.db.Redis.Vote(id, optionID)
	if err != nil {
		log.Printf("Error during voting: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"message": "vote rejected"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": voteStats})
	return
}
