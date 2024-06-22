package types

type S3fsMount struct {
	ID     int64  `json:"id"`
	Path   string `json:"path"`
	Bucket string `json:"bucket"`
	Url    string `json:"url"`
}
