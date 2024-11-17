package persistance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"vote-app/contracts"
)

type RedisCache struct {
	Db *redis.Client
}

func (cache *RedisCache) CreateVote(vote *contracts.Vote) error {
	voteJson, _ := json.Marshal(&vote)
	ctx := context.Background()
	key := buildKey(vote.ID.String())

	set := make([]string, 0, len(vote.Options)*2)

	for k := range vote.Options {
		set = append(set, k, "0")
	}

	_, err := cache.Db.Set(ctx, key, voteJson, 0).Result()
	if err != nil {
		log.Fatalf("Failed to set key=%s, error=%s", key, err.Error())
	}

	_, sErr := cache.Db.HSet(ctx, buildStatsKey(key), set).Result()
	if sErr != nil {
		log.Fatalf("Failed to set hashset key=%s, error=%s", key, sErr.Error())
	}
	return nil
}

func (cache *RedisCache) GetVotes() ([]contracts.Vote, error) {
	var cursor uint64
	keys, cursor, err := cache.Db.Scan(context.Background(), cursor, "votes:*", 10).Result()
	if err != nil {
		return nil, err
	}
	pipe := cache.Db.Pipeline()
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

func (cache *RedisCache) GetVote(id string) (*contracts.Vote, error) {
	bytes, err := cache.Db.Get(context.Background(), buildKey(id)).Bytes()
	if err != nil {
		return nil, err
	}

	var vote contracts.Vote

	err = json.Unmarshal(bytes, &vote)
	if err != nil {
		return nil, err
	}

	return &vote, nil
}

func (cache *RedisCache) GetVoteStats(id string) (*contracts.VoteStats, error) {
	vote, err := cache.GetVote(id)
	if err != nil {
		return nil, err
	}

	stats, err := cache.Db.HGetAll(context.Background(), buildStatsKey(buildKey(id))).Result()
	if err != nil {
		return nil, err
	}

	voteStats := &contracts.VoteStats{
		Name:    vote.Name,
		Options: make([]contracts.VoteStatItem, 0, len(stats)),
	}

	for k, v := range stats {
		votes, _ := strconv.Atoi(v)
		voteStats.Options = append(voteStats.Options, contracts.VoteStatItem{
			OptionId:   k,
			OptionName: vote.Options[k],
			Votes:      votes,
		})
	}

	return voteStats, nil
}

func (cache *RedisCache) Vote(id string, optionID string) (*contracts.VoteStats, error) {
	_, err := cache.Db.HIncrBy(context.Background(), buildStatsKey(buildKey(id)), optionID, 1).Result()
	if err != nil {
		return nil, err
	}

	stats, err := cache.GetVoteStats(id)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func buildKey(id string) string {
	return fmt.Sprintf("votes:%s", id)
}
func buildStatsKey(id string) string {
	return fmt.Sprintf("stats:%s", id)
}
