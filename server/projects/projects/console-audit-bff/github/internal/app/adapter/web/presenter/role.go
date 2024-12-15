package presenter

type (
	Roles map[string]map[string]map[string][]string

	UserRoles struct {
		Roles []string `json:"roles"`
	}

	Role struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

