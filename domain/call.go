package domain

type CallResult string

const (
	INIT         CallResult = "INIT"
	QUEUEING     CallResult = "QUEUEING"
	SUCCESS      CallResult = "SUCCESS"
	FAIL         CallResult = "FAIL"
	NOT_ANSWER   CallResult = "NOT_ANSWER"
	CANT_CONNECT CallResult = "CANT_CONNECT"
)

type CallLog struct {
	ID          uint                   `json:"id"`
	PhoneNumber string                 `json:"phone_number"`
	Metadata    map[string]interface{} `json:"metadata"`
	CallResult  CallResult             `json:"call_result"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
	CallTime    int64                  `json:"call_time"`
	ResultTime  int64                  `json:"result_time"`
	PickupTime  *int64                 `json:"pickup_time"`
	HangupTime  *int64                 `json:"hangup_time"`
}

type CallFilter struct {
	PhoneNumber string
	StartAt     int64
	EndAt       int64
}

type CallRepository interface {
	Create(call *CallLog) error
	Update(id uint, updated *CallLog) error
	Delete(id uint) error
	GetByID(id uint) (*CallLog, error)
	List(filter CallFilter) ([]CallLog, error)
}
