package domain

type CallResult string

const (
	INIT CallResult = "INIT"
	QUEUEING CallResult = "QUEUEING"
	SUCCESS  CallResult = "SUCCESS"
	FAIL   CallResult = "FAIL"
	NOT_ANSWER CallResult = "NOT_ANSWER"
	CANT_CONNECT CallResult = "CANT_CONNECT"
)

type CallLog struct {
	ID          uint
	PhoneNumber string
	Metadata    map[string]interface{}
	CallResult  CallResult
	CreatedAt   int64
	UpdatedAt   int64
	CallTime    int64
	ResultTime  int64
	PickupTime  *int64
	HangupTime  *int64
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