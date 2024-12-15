package request

type ID struct {
	ID uint `json:"id" form:"id" query:"id" validate:"required|number"`
}
