package kafka

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"

	"orders/cmd/generator"
	r "orders/internal/repository"

	"github.com/segmentio/kafka-go"
)

const (
	topic   string = "orders"
	address string = "localhost:9092"
)

func CreateReader() *kafka.Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{address},
		Topic:     topic,
		Partition: 0,
	})
	return r
}

func CreateTopic() {
	conn, err := kafka.Dial("tcp", address)
	if err != nil {
		log.Fatalln("Error creating kafka connection:", err)
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
}

func StartConsuming(r *kafka.Reader, repo *r.Repository) {
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		log.Printf("New message at topic/partition/offset %v/%v/%v: %s = %s\n",
			m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))

		ctx := context.Background()

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
	}
}
