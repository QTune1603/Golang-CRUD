package consumer

import (
	"Golang-CRUD/domain"
	"Golang-CRUD/internal/infra/repository"
	"encoding/json"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"log"
)

func StartResultUpdater(conn *amqp.Connection, db *gorm.DB) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Không thể mở kênh RabbitMQ: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("call_result_queue", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Không thể khai báo queue: %v", err)
	}

	msgs, err := ch.Consume("call_result_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Không thể lắng nghe queue: %v", err)
	}

	repo := repository.NewCallRepository(db)
	log.Println("Consumer cập nhật kết quả đang lắng nghe call_result_queue...")

	for msg := range msgs {
		var result domain.CallLog
		if err := json.Unmarshal(msg.Body, &result); err != nil {
			log.Printf("Lỗi giải mã message: %v\n", err)
			continue
		}

		// Update một số trường
		updateData := &domain.CallLog{
			CallResult: result.CallResult,
			ResultTime: result.ResultTime,
			PickupTime: result.PickupTime,
			HangupTime: result.HangupTime,
			UpdatedAt:  result.ResultTime,
		}

		err := repo.Update(result.ID, updateData)
		if err != nil {
			log.Printf("Cập nhật DB lỗi ID %d: %v\n", result.ID, err)
		} else {
			log.Printf("Cập nhật thành công kết quả ID %d: %s\n", result.ID, result.CallResult)
		}
	}
}
