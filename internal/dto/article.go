package dto

// CreateArticleRequest is the payload for publishing an article.
type CreateArticleRequest struct {
	Title    string `json:"title" binding:"required,min=2,max=64"`
	Category string `json:"category" binding:"required"`
	Content  string `json:"content" binding:"required,min=10"`
	CoverURL string `json:"cover_url"`
}
