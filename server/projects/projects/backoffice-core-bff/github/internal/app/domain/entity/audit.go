package entity

const ACTION_CREATE = "create"

type (
	Audit struct {
		Operation    Operation    `json:"operation"`
		User         User         `json:"user"`
		UserAgent    UserAgent    `json:"user_agent"`
		Localization Localization `json:"localization"`
		Http         Http         `json:"http"`
	}

	Operation struct {
		Tenant   string `json:"tenant"`
		Action   string `json:"action"`
		Domain   string `json:"domain"`
		DomainId string `json:"domain_id"`
		CID      string `json:"cid"`
		Date     string `json:"date"`
	}

	User struct {
		Email      string   `json:"email"`
		Permission []string `json:"permissions"`
	}

	UserAgent struct {
		Device   string `json:"device"`
		DeviceIp string `json:"device_ip"`
	}

	Localization struct {
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}

	Http struct {
		Code     int    `json:"code"`
		Request  string `json:"request"`
		Response string `json:"response"`
	}
)

func NewAudit() Audit {
	return Audit{
		Operation:    Operation{},
		User:         User{},
		UserAgent:    UserAgent{},
		Localization: Localization{},
		Http:         Http{},
	}
}

