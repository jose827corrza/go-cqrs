package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"jose827corrza/go-cqrs/models"

	elastic "github.com/elastic/go-elasticsearch/v7"
)

type ElasticSearchRepository struct {
	client *elastic.Client
}

func NewElasticSearchRepository(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return nil, err
	}
	return &ElasticSearchRepository{client: client}, nil
}

func (e *ElasticSearchRepository) Close() {
	//
}

func (e *ElasticSearchRepository) IndexFeed(ctx context.Context, feed models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := e.client.Index(
		"feeds",
		bytes.NewReader(body),
		e.client.Index.WithDocumentID(feed.ID),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithRefresh("wait_for"),
	)
	return err
}

func (e *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) (results []models.Feed, err error) {
	var buf bytes.Buffer

	//El por que se define como esta abajo es el como se define una struct de datos que go no sabe como vendriasn ni el tipo
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzziness":        3,     // busquedas mas rapidas : gou ->go
				"cutoff_frequency": 0.001, //Cuantas veces debe repetirse para que se devuelvan los docs
			},
		},
	}
	if err = json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}
	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("feeds"),
		e.client.Search.WithBody(&buf),
		e.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			results = nil
		}
	}()

	if res.IsError() {
		return nil, errors.New(res.String())
	}
	var eRes map[string]interface{}
	//eRes representacion en JSON e lo decodificado de los resuslts de busqueda
	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	//Desde aca se hace la serializacion para pasar del tipo de elasticSearch a []models.Feed
	var feeds []models.Feed

	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		feed := models.Feed{}
		source := hit.(map[string]interface{})["_source"]
		marshal, err := json.Marshal(source)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(marshal, &feed); err == nil {
			feeds = append(feeds, feed)
		}
	}
	return feeds, nil
}
