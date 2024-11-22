package search

import (
	"auth-service/pkg/util"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticsearchClient() (*elasticsearch.Client, error) {
	env := util.GetConfig(".")

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{env.EsHost},
		Username:  env.EsUsername,
		Password:  env.EsPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao criar cliente Elasticsearch: %v", err)
	}

	res, err := es.Cluster.Health(es.Cluster.Health.WithTimeout(10 * time.Second))
	if err != nil {
		return nil, fmt.Errorf("erro ao verificar saúde do Elasticsearch: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("erro de resposta do Elasticsearch: %s", res.String())
	}

	return es, nil
}
