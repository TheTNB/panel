package request

type SystemctlService struct {
	Service string `json:"service" validate:"required"`
}
