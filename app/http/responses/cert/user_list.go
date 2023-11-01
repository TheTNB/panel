package responses

import "panel/app/models"

type UserList struct {
	Total int64             `json:"total"`
	Items []models.CertUser `json:"items"`
}
