package input

type RateInput struct {
	Score      uint   `json:"score" binding:"required,max=5"`
	Feature    string `json:"feature" binding:"required"`
	PropertyID uint   `json:"property_id" binding:"required"`
}
