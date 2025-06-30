package service

import (
	"call-api/model"
	"call-api/repository"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type CallLogService struct {
	Repo   *repository.CallLogRepository
	Queue  *amqp.Channel
}

func NewCallLogService(repo *repository.CallLogRepository, conn *amqp.Connection) *CallLogService {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	// Declare queue nếu chưa có
	_, err = ch.QueueDeclare("call_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &CallLogService{Repo: repo, Queue: ch}
}

func (s *CallLogService) CreateAndEnqueue(log *model.CallLog) error {
	if err := s.Repo.Create(log); err != nil {
		return err
	}

	// Push to RabbitMQ
	body, _ := json.Marshal(log)
	return s.Queue.Publish(
		"",            // exchange
		"call_queue",  // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
