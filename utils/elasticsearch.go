package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"golang.org/x/net/context"
)

var EsClient *elasticsearch.Client
var EsCtx = context.Background()

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Elasticsearch client: %v", err)
	}

	EsClient = client
	log.Println("Elasticsearch client initialized successfully.")
}

type SMSDocument struct {
	ID          string    `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Message     string    `json:"message"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func IndexSMS(document SMSDocument) error {
	jsonBody, err := json.Marshal(document)
	if err != nil {
		return err
	}

	res, err := EsClient.Index(
		"sms_index",
		bytes.NewReader(jsonBody),
		EsClient.Index.WithDocumentID(document.ID),
		EsClient.Index.WithContext(EsCtx),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return err
	}

	log.Printf("Indexed document ID: %s", document.ID)
	return nil
}
