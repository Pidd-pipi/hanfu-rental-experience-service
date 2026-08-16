package dto

// CreateHanfuRequest is the payload for adding a hanfu.
type CreateHanfuRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=64"`
	Dynasty     string  `json:"dynasty" binding:"required"`
	Form        string  `json:"form"`
	Color       string  `json:"color"`
	Size        string  `json:"size" binding:"required"`
	PricePerDay float64 `json:"price_per_day" binding:"required,gt=0"`
	Images      string  `json:"images"`
	Stock       int     `json:"stock" binding:"gte=1"`
}

// ListHanfuQuery adds dynasty/size filters to pagination.
type ListHanfuQuery struct {
	PageQuery
	Dynasty string `form:"dynasty"`
	Size    string `form:"size"`
	Form    string `form:"form"`
}
