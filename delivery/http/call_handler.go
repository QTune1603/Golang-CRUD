package http

import (
	"Golang-CRUD/domain"
	"Golang-CRUD/usecase"
	"Golang-CRUD/usecase/reader"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CallHandler struct {
	Usecase    *usecase.CallUsecase
	ReaderRepo reader.CallReaderRepository
}

func NewCallHandler(uc *usecase.CallUsecase, rr reader.CallReaderRepository) *CallHandler {
	return &CallHandler{
		Usecase:    uc,
		ReaderRepo: rr,
	}
}

// Create
func (h *CallHandler) Create(c *gin.Context) {
	var input domain.CallLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Created"})
}

// Update
func (h *CallHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input domain.CallLog
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Usecase.Update(uint(id), &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated"})
}

// Delete
func (h *CallHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.Usecase.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

// GetByID
func (h *CallHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	call, err := h.Usecase.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
		return
	}
	c.JSON(http.StatusOK, call)
}

// List
func (h *CallHandler) List(c *gin.Context) {
	var filter domain.CallFilter
	filter.PhoneNumber = c.Query("phone_number")

	start := c.Query("start_at")
	end := c.Query("end_at")
	if start != "" && end != "" {
		startInt, _ := strconv.ParseInt(start, 10, 64)
		endInt, _ := strconv.ParseInt(end, 10, 64)
		filter.StartAt = startInt
		filter.EndAt = endInt
	}

	metaField := c.Query("metadata_display_field")

	calls, err := h.ReaderRepo.GetWithMetadataField(filter, metaField)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response
	result := make([]map[string]interface{}, 0)
	for _, call := range calls {
		item := map[string]interface{}{
			"id":           call.ID,
			"phone_number": call.PhoneNumber,
			"call_result":  call.CallResult,
			"created_at":   call.CreatedAt,
			"updated_at":   call.UpdatedAt,
			"call_time":    call.CallTime,
			"result_time":  call.ResultTime,
			"pickup_time":  call.PickupTime,
			"hangup_time":  call.HangupTime,
		}

		raw := call.Metadata
		if metaField != "" {
			if val, ok := raw[metaField]; ok {
				item["metadata"] = map[string]interface{}{metaField: val}
			} else {
				item["metadata"] = map[string]interface{}{}
			}
		} else {
			item["metadata"] = raw
		}

		result = append(result, item)
	}

	c.JSON(http.StatusOK, result)
}
