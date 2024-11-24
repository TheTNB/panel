package types

type DatabaseStatus string

const (
	DatabaseStatusValid   DatabaseStatus = "valid"
	DatabaseStatusInvalid DatabaseStatus = "invalid"
)

type Database struct {
	Name     string         `json:"name"`
	ServerID uint           `json:"server_id"`
	Status   DatabaseStatus `json:"status"`
	Remark   string         `json:"remark"`
}
