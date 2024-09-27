package s3fs

type Create struct {
	Ak     string `form:"ak" json:"ak"`
	Sk     string `form:"sk" json:"sk"`
	Bucket string `form:"bucket" json:"bucket"`
	URL    string `form:"url" json:"url"`
	Path   string `form:"path" json:"path"`
}

type Delete struct {
	ID int64 `form:"id" json:"id"`
}
