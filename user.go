package segmentationService

type User struct {
	id       int    `json:"id" db:"id"`
	Username string `json:"username" binding:"required"`
}
