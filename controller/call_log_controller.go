package controller

import (
	"call-api/model"
	"call-api/repository"
	"call-api/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, conn *amqp.Connection) {
	repo := repository.NewCallLogRepository(db)
	svc := service.NewCallLogService(repo, conn)

	r.POST("/calls", func(c *gin.Context) {
		var input model.CallLog
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		now := time.Now().UnixMilli()
		input.CreatedAt = now
		input.UpdatedAt = now
		input.CallResult = "INIT"

		if err := svc.CreateAndEnqueue(&input); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, input)
	})

	r.GET("/calls", func(c *gin.Context) {
		phone := c.Query("phone_number")
		start := c.Query("start_at")
		end := c.Query("end_at")
		metadataField := c.Query("metadata_display_field")

		var logs []model.CallLog
		query := db.Model(&model.CallLog{})

		if phone != "" {
			query = query.Where("phone_number LIKE ?", "%"+phone+"%")
		}
		if start != "" && end != "" {
			startInt, err1 := strconv.ParseInt(start, 10, 64)
			endInt, err2 := strconv.ParseInt(end, 10, 64)
			if err1 != nil || err2 != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start or end time"})
				return
			}
			query = query.Where("created_at BETWEEN ? AND ?", startInt, endInt)
		}

		if err := query.Order("id DESC").Find(&logs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		//metadata _display_field
		result := make([]map[string]interface{}, 0)
		for _, logItem := range logs {
			item := map[string]interface{}{
				"id":           logItem.ID,
				"phone_number": logItem.PhoneNumber,
				"call_result":  logItem.CallResult,
				"created_at":   logItem.CreatedAt,
				"call_time":    logItem.CallTime,
				"result_time":  logItem.ResultTime,
				"pickup_time":  logItem.PickupTime,
				"hangup_time":  logItem.HangupTime,
			}

			var raw map[string]interface{}
			if err := json.Unmarshal(logItem.Metadata, &raw); err != nil {
				raw = map[string]interface{}{}
			}
			if metadataField != "" {
				if val, ok := raw[metadataField]; ok {
					item["metadata"] = map[string]interface{}{metadataField: val}
				} else {
					item["metadata"] = map[string]interface{}{}
				}
			} else {
				item["metadata"] = raw
			}

			result = append(result, item)
		}

		c.JSON(200, result)
	})

	r.GET("/calls/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		var call model.CallLog
		if err := db.First(&call, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Call not found"})
			return
		}

		// Parse metadata
		var meta map[string]interface{}
		_ = json.Unmarshal(call.Metadata, &meta)

		// Trả về đầy đủ
		c.JSON(200, gin.H{
			"id":           call.ID,
			"phone_number": call.PhoneNumber,
			"call_result":  call.CallResult,
			"created_at":   call.CreatedAt,
			"updated_at":   call.UpdatedAt,
			"call_time":    call.CallTime,
			"result_time":  call.ResultTime,
			"pickup_time":  call.PickupTime,
			"hangup_time":  call.HangupTime,
			"metadata":     meta,
		})
	})

	r.PUT("/calls/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		var input model.CallLog
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// Marshal metadata nếu là map
		if m, ok := any(input.Metadata).(map[string]interface{}); ok {
			b, _ := json.Marshal(m)
			input.Metadata = b
		}

		input.UpdatedAt = time.Now().UnixMilli()
		err = db.Model(&model.CallLog{}).Where("id = ?", id).Updates(&input).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Updated"})
	})

	r.DELETE("/calls/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID"})
			return
		}

		err = db.Delete(&model.CallLog{}, id).Error
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"message": "Deleted"})
	})

}
