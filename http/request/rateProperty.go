package request

type RateInput struct {
	Score      uint   `json:"score" binding:"required,max=5"`
	Comment    string `json:"comment" binding:"required"`
	PropertyID uint   `json:"property_id" binding:"required"`
	UserID     uint   `json:"user_id" binding:"required"`
}
