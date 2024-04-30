package domain

type MessagesSuccess interface {
	GetMessage() string
	GetData() interface{}
}

type successData struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (s *successData) GetMessage() string {
	return s.Message
}

func (s *successData) GetData() interface{} {
	return s.Data
}

func NewStatusCreated(message string, data interface{}) MessagesSuccess {
	return &successData{
		Message: message,
		Data:    data,
	}
}

func NewStatusOk(message string, data interface{}) MessagesSuccess {
	return &successData{
		Message: message,
		Data:    data,
	}
}
