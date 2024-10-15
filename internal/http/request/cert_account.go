package request

type CertAccountCreate struct {
	CA          string `form:"ca" json:"ca" validate:"required"`
	Email       string `form:"email" json:"email" validate:"required"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type" validate:"required"`
}

type CertAccountUpdate struct {
	ID          uint   `form:"id" json:"id" validate:"required,exists=cert_accounts id"`
	CA          string `form:"ca" json:"ca" validate:"required"`
	Email       string `form:"email" json:"email" validate:"required"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type" validate:"required"`
}
