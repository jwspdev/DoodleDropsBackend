package requests

type PaginateRequest struct {
	PageNumber int `json:"page_number"`
	Limit      int `json:"limit"`
}
