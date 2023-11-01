package responses

import "panel/app/models"

type CertList struct {
	Total int64         `json:"total"`
	Items []models.Cert `json:"items"`
}
