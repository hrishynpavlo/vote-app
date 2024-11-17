package persistance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"vote-app/contracts"
)

func CreateVote(vote *contracts.Vote, db *redis.Client) error {
	voteJson, _ := json.Marshal(&vote)
	_, err := db.Set(context.Background(), fmt.Sprintf("votes:%s", vote.ID.String()), voteJson, 0).Result()
	return err
}

func GetVotes(db *redis.Client) ([]contracts.Vote, error) {
	var cursor uint64
	keys, cursor, err := db.Scan(context.Background(), cursor, "votes:*", 10).Result()
	if err != nil {
		return nil, err
	}
	pipe := db.Pipeline()
	for _, key := range keys {
		pipe.Get(context.Background(), key)
	}
	result, err := pipe.Exec(context.Background())
	if err != nil {
		return nil, err
	}
	votes := make([]contracts.Vote, 0, len(result))

	for _, cmd := range result {
		var vote contracts.Vote
		cmdResult, _ := cmd.(*redis.StringCmd).Bytes()
		err := json.Unmarshal(cmdResult, &vote)
		if err != nil {
			return nil, err
		}
		votes = append(votes, vote)
	}

	return votes, nil
}
