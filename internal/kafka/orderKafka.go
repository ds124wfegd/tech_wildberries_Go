package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
	"github.com/ds124wfegd/tech_wildberries_Go/internal/service"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	service service.OrderService
}

// new kafka consumer
func NewConsumer(brokers []string, topic, groupID string, service service.OrderService) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	return &Consumer{reader: r, service: service}
}

// starting kafka consumer
func (c *Consumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumer")
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping Kafka consumer")
			return
		default:
			m, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("Error reading message: %v", err)
				continue
			}
			var order entity.Order
			if err := json.Unmarshal(m.Value, &order); err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}
			if order.OrderUID == "" {
				log.Printf("Invalid order: missing order_uid")
				continue
			}
			if err := c.service.Ingest(ctx, &order); err != nil {
				log.Printf("Error ingesting order: %v", err)
				continue
			}
			log.Printf("Successfully processed order: %s", order.OrderUID)
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
		return err
	}
	log.Println("Stopping Kafka consumer...")
	return nil
}

/*
// for information :)
type ReaderConfig struct {
    Brokers          []string       // Адреса брокеров Kafka (например, ["localhost:9092"])
    Topic            string         // Название топика, из которого будет читаться
    Partition        int            // Номер партиции (если читается конкретная партиция)
    GroupID          string         // ID consumer-группы (для управления офсетами)
    GroupTopics      []string       // Топики для группового потребления (альтернатива Topic)
    StartOffset      int64          // Начальный офсет (FirstOffset или LastOffset)
    MinBytes         int            // Минимальное количество байт для чтения (дефолт 1)
    MaxBytes         int            // Максимальное количество байт за один запрос (дефолт 1MB)
    MaxWait          time.Duration  // Макс. время ожидания новых данных
    ReadLagInterval  time.Duration  // Интервал проверки отставания (lag)
    CommitInterval   time.Duration  // Интервал коммита офсетов
    HeartbeatInterval time.Duration // Интервал heartbeat для группы
    SessionTimeout   time.Duration  // Таймаут сессии потребителя
    RebalanceTimeout time.Duration  // Таймаут ребалансировки
    RetentionTime    time.Duration  // Время хранения офсетов
    Logger           Logger         // Логгер (может быть nil)
    ErrorLogger      Logger         // Логгер ошибок
    IsolationLevel   IsolationLevel // Уровень изоляции (ReadUncommitted или ReadCommitted)
    QueueCapacity    int            // Размер внутренней очереди сообщений
}
*/
