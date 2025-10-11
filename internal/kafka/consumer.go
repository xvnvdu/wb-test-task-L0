package kafka

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"orders/internal/generator"
	"strconv"
	"time"

	repo "orders/internal/repository"

	"github.com/segmentio/kafka-go"
)

const (
	topic   string = "orders"
	address string = "kafka:9092"
)

func CreateReader() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     topic,
		GroupID:   "orders-group",
		Partition: 0,
	})
	return r
}

func CreateTopic() {
	var conn *kafka.Conn
	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		conn, err = kafka.Dial("tcp", address)
		if err == nil {
			break
		}
		log.Printf("Error creating Kafka connection (attempt %d/%d): %v", i+1, maxRetries, err)

		if i == maxRetries-1 {
			log.Fatalln("Failed to connect to Kafka after all attempts.")
		}

		time.Sleep(time.Second * 5)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Fatalln("Error creating kafka controller:", err)
	}

	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalln("Error creating controlerConn:", err)
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatalln("Error creating topic:", err)
	}
	log.Printf("Topic %s created successfuly on %s", topic, address)
}

func StartConsuming(r *kafka.Reader, repo *repo.Repository) {
	ctx := context.Background()
	for {
		m, err := r.FetchMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("New message at topic/partition/offset %v/%v/%v: %s = %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		var orders []*generator.Order
		err = json.Unmarshal(m.Value, &orders)
		if err != nil {
			log.Println("Error unmarshalling orders data:", err)
			continue
		}

		err = repo.SaveToDB(orders, ctx)
		if err != nil {
			log.Printf("Failed to save orders from Kafka message: %v\n", err)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			log.Fatalln("Error committing message:", err)
		}
		log.Printf("Committed message at topic/partition/offset %v/%v/%v\n",
			m.Topic, m.Partition, m.Offset)
	}
}
