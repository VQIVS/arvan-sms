package presenter

type Status string

const (
	Pending Status = "pending"
	Success Status = "success"
	Failed  Status = "failed"
)

type SendSMSReq struct {
	UserID    uint   `json:"user_id"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
}

type SendSMSResp struct {
	ID      uint   `json:"id"`
	Status  Status `json:"status"`
	Message string `json:"message"`
}

type SMSResp struct {
	ID        uint   `json:"id"`
	Recipient string `json:"recipient"`
	Message   string `json:"message"`
	Status    Status `json:"status"`
}
