package request

type CertAccountCreate struct {
	CA          string `form:"ca" json:"ca"`
	Email       string `form:"email" json:"email"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type"`
}

type CertAccountUpdate struct {
	ID          uint   `form:"id" json:"id"`
	CA          string `form:"ca" json:"ca"`
	Email       string `form:"email" json:"email"`
	Kid         string `form:"kid" json:"kid"`
	HmacEncoded string `form:"hmac_encoded" json:"hmac_encoded"`
	KeyType     string `form:"key_type" json:"key_type"`
}
