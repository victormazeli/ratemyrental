package dtos

type PaginationDTO struct {
	TotalDocs  int64       `json:"total_docs"`
	TotalPages int64       `json:"total_pages"`
	Docs       interface{} `json:"docs"`
	Page       int64       `json:"page"`
	Limit      int64       `json:"limit"`
}
