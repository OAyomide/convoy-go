package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"time"

	convoy "github.com/frain-dev/convoy-go"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	URL        = "http://localhost:5005/api/v1"
	projectID  = "01HB8J53CSBC4ZWCJ95TCQ6S43"
	endpointID = "01HCZ0RBYH4MQTTCJKJ9KTS5KB"
	apiKey     = "CO.vMkWVbqa7mFsmeGA.MkU35AfkWF3AcUVvNOqBj94QGZ05jxzjUmH4sgMYcipAji26dnnyNJo5bQkSzUTu"
	kUsername  = "aHVtYW5lLXNsb3RoLTEyMjc5JC1Fof9Z-WROYxBkMGlQZPIpKC1LufOesVCVwGY"
	kPassword  = "ZDkyMzliYzUtMWJkYS00MTkzLWI2NjQtNGM5ZTM0ODQ1YTI0"
)

func main() {
	logger := convoy.NewLogger(os.Stdout, convoy.DebugLevel)
	ctx := context.Background()

	mechanism, err := scram.Mechanism(scram.SHA256, kUsername, kPassword)
	if err != nil {
		log.Fatalln(err)
	}

	sharedTransport := &kafka.Transport{
		SASL: mechanism,
		TLS:  &tls.Config{},
	}

	kClient := &kafka.Client{
		Addr:      kafka.TCP("humane-sloth-12279-us1-kafka.upstash.io:9092"),
		Timeout:   10 * time.Second,
		Transport: sharedTransport,
	}

	ko := &convoy.KafkaOptions{
		Client: kClient,
		Topic:  "demo-topic",
	}

	kc := convoy.New(URL, apiKey, projectID,
		convoy.OptionLogger(logger),
		convoy.OptionKafkaOptions(ko),
	)

	fmt.Println("writing kafka event...")
	err = writeKafkaEvent(ctx, kc)
	if err != nil {
		log.Fatal(err)
	}
}

func writeKafkaEvent(ctx context.Context, c *convoy.Client) error {
	body := &convoy.CreateEventRequest{
		EventType:      "test.customer.event",
		EndpointID:     endpointID,
		IdempotencyKey: "subomi-abcd",
		Data: []byte(`{
						"event_type": "test.event", 
						"data": { 
							"Hello": "World", 
							"Test": "Data" 
						}
					}`),
	}

	return c.Kafka.WriteEvent(ctx, body)
}
