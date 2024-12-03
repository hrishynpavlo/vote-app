package persistance

import "vote-app/contracts"

type Storage struct {
	elastic *ElasticSearchClient
	Redis   *RedisCache
}

func AddStorageDecorator(elastic *ElasticSearchClient, redis *RedisCache) *Storage {
	return &Storage{
		elastic: elastic,
		Redis:   redis,
	}
}

func (s *Storage) CreateVote(vote *contracts.Vote) error {
	err := s.elastic.CreateVote(vote)
	if err != nil {
		return err
	}

	err = s.Redis.CreateVote(vote)
	if err != nil {
		return err
	}

	return nil
}
