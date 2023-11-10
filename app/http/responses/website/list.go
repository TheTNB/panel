package responses

import "panel/app/models"

type List struct {
	Total int64            `json:"total"`
	Items []models.Website `json:"items"`
}
