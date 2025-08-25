package entity

type Meta struct {
	TotalData int `json:"total_items"`
	TotalPage int `json:"total_pages"`
	Page      int `json:"page"`
	Limit     int `json:"limit"`
}
