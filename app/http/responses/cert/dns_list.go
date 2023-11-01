package responses

import "panel/app/models"

type DNSList struct {
	Total int64            `json:"total"`
	Items []models.CertDNS `json:"items"`
}
