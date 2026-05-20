package search

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

const PostIndex = "posts"
const maxResultWindow = 10000

type Service struct {
	es  *elasticsearch.Client
	log zerolog.Logger
}

func NewService(es *elasticsearch.Client, log zerolog.Logger) *Service {
	s := &Service{
		es:  es,
		log: log.With().Str("module", "search").Logger(),
	}
	if es != nil {
		s.ensureIndex()
	}
	return s
}

type IndexedPost struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Author    string   `json:"author"`
	PostType  string   `json:"post_type"`
	Tags      []string `json:"tags"`
	CreatedAt int64    `json:"created_at"`
}

func (s *Service) ensureIndex() {
	mapping := `{
		"settings": {
			"number_of_shards": 1,
			"number_of_replicas": 0,
			"analysis": {
				"analyzer": {
					"devix_analyzer": {
						"type": "custom",
						"tokenizer": "standard",
						"filter": ["lowercase", "stop", "snowball"]
					}
				}
			}
		},
		"mappings": {
			"properties": {
				"id": { "type": "keyword" },
				"title": { "type": "text", "analyzer": "devix_analyzer", "boost": 2.0 },
				"content": { "type": "text", "analyzer": "devix_analyzer" },
				"author": { "type": "keyword" },
				"post_type": { "type": "keyword" },
				"tags": { "type": "keyword" },
				"created_at": { "type": "date" }
			}
		}
	}`

	res, err := s.es.Indices.Exists([]string{PostIndex})
	if err == nil && res.StatusCode == 404 {
		res, err = s.es.Indices.Create(
			PostIndex,
			s.es.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			s.log.Error().Err(err).Msg("failed to create index")
		} else {
			s.log.Info().Msg("created elasticsearch index: " + PostIndex)
		}
	}
}

func (s *Service) IndexPost(ctx context.Context, post IndexedPost) error {
	if s.es == nil {
		return nil
	}
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      PostIndex,
		DocumentID: post.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document: %s", res.String())
	}

	return nil
}

func (s *Service) SearchPosts(ctx context.Context, query string, limit, offset int) ([]uuid.UUID, error) {
	if s.es == nil {
		return nil, nil
	}
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	if offset+limit > maxResultWindow {
		return nil, fmt.Errorf("search offset exceeds elasticsearch result window")
	}

	var buf bytes.Buffer
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":     query,
				"fields":    []string{"title^3", "content", "tags^2"},
				"fuzziness": "AUTO",
			},
		},
		"size": limit,
		"from": offset,
	}

	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := s.es.Search(
		s.es.Search.WithContext(ctx),
		s.es.Search.WithIndex(PostIndex),
		s.es.Search.WithBody(&buf),
		s.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	var ids []uuid.UUID
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		idStr := source["id"].(string)
		id, _ := uuid.Parse(idStr)
		ids = append(ids, id)
	}

	return ids, nil
}

func (s *Service) DeletePost(ctx context.Context, postID string) error {
	if s.es == nil {
		return nil
	}
	req := esapi.DeleteRequest{
		Index:      PostIndex,
		DocumentID: postID,
		Refresh:    "true",
	}

	res, err := req.Do(ctx, s.es)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
