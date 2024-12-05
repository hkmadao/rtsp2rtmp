package result

type Page struct {
	Total int         `json:"total"`
	Page  interface{} `json:"page"`
}
