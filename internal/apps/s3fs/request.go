package s3fs

type Create struct {
	Ak     string `form:"ak" json:"ak" validate:"required"`
	Sk     string `form:"sk" json:"sk" validate:"required"`
	Bucket string `form:"bucket" json:"bucket" validate:"required"`
	URL    string `form:"url" json:"url" validate:"required"`
	Path   string `form:"path" json:"path" validate:"required"`
}

type Delete struct {
	ID int64 `form:"id" json:"id" validate:"required"`
}
