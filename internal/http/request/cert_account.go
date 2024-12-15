package request

type CertAccountCreate struct {
	CA          string `form:"ca" json:"ca" validate:"required|in:googlecn,google,letsencrypt,buypass,zerossl,sslcom"`
	Email       string `form:"email" json:"email" validate:"required"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type" validate:"required|in:P256,P384,2048,3072,4096"`
}

type CertAccountUpdate struct {
	ID          uint   `form:"id" json:"id" validate:"required|exists:cert_accounts,id"`
	CA          string `form:"ca" json:"ca" validate:"required|in:googlecn,google,letsencrypt,buypass,zerossl,sslcom"`
	Email       string `form:"email" json:"email" validate:"required"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type" validate:"required|in:P256,P384,2048,3072,4096"`
}
