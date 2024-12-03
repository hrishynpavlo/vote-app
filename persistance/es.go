package persistance

import (
	"bytes"
	"encoding/json"
	elasticsearch8 "github.com/elastic/go-elasticsearch/v8"
	"vote-app/configuration"
	"vote-app/contracts"
)

const voteIndex = "votes"

type ElasticSearchClient struct {
	Client *elasticsearch8.Client
}

func AddElasticSearchDb(cfg *configuration.Configuration) *ElasticSearchClient {
	es, err := elasticsearch8.NewClient(elasticsearch8.Config{Addresses: []string{cfg.ElasticSearchUrl}})
	if err != nil {
		panic(err)
	}
	return &ElasticSearchClient{
		Client: es,
	}
}

func (esService *ElasticSearchClient) CreateVote(vote *contracts.Vote) error {
	data, _ := json.Marshal(vote)
	_, err := esService.Client.Index(voteIndex, bytes.NewReader(data))
	if err != nil {
		return err
	}

	return nil
}
