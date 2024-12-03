package request

type ProcessKill struct {
	PID int32 `json:"pid" validate:"required"`
}
