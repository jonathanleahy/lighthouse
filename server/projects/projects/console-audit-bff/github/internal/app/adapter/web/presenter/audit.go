package presenter

import "github.com/pismo/console-audit-bff/internal/app/adapter/web/graphql/entity"

type (
	Audit struct {
		ID           int64         `json:"id"`
		Operation    *Operation    `json:"operation"`
		User         *User         `json:"user"`
		UserAgent    *UserAgent    `json:"user_agent"`
		Localization *Localization `json:"localization"`
		Http         *Http         `json:"http"`
	}

	Operation struct {
		Tenant   string  `json:"tenant"`
		Action   string  `json:"action"`
		Domain   string  `json:"domain"`
		DomainID string  `json:"domain_id"`
		Origin   *string `json:"origin"`
		CID      string  `json:"cid"`
		Date     string  `json:"date"`
	}

	User struct {
		Email string   `json:"email"`
		Roles []string `json:"roles"`
	}

	UserAgent struct {
		Device   *string `json:"device"`
		DeviceIp *string `json:"device_ip"`
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

	SearchAudit struct {
		CurrentPage int      `json:"current_page"`
		Pages       int      `json:"pages"`
		PerPage     int      `json:"per_page"`
		TotalItems  int      `json:"total_items"`
		Items       []*Audit `json:"items"`
	}

	SearchFeature struct {
		CurrentPage int               `json:"current_page"`
		Pages       int               `json:"pages"`
		PerPage     int               `json:"per_page"`
		TotalItems  int               `json:"total_items"`
		Items       []*entity.Feature `json:"items"`
	}

	SearchEndpoint struct {
		CurrentPage int                `json:"current_page"`
		Pages       int                `json:"pages"`
		PerPage     int                `json:"per_page"`
		TotalItems  int                `json:"total_items"`
		Items       []*entity.Endpoint `json:"items"`
	}
)

