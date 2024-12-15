package presenter

type Tracking struct {
	TrackingID    string `json:"tracking_id"`
	CorrelationID string `json:"correlation_id"`
}

type Extensions struct {
	Tenant   string    `json:"tenant"`
	User     *User     `json:"user"`
	Tracking *Tracking `json:"tracking"`
}

type ErrorMessage struct {
	Message    string      `json:"message"`
	Extensions *Extensions `json:"extensions"`
}

