package rule

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type Exists struct {
	DB *gorm.DB
}

func NewExists(db *gorm.DB) *Exists {
	return &Exists{DB: db}
}

// Exists 格式 `exists=categories，id，other_field`
func (r *Exists) Exists(fl validator.FieldLevel) bool {
	requestValue := fl.Field().Interface()
	params := strings.Fields(fl.Param())
	if len(params) < 2 {
		return false
	}

	tableName := params[0]
	fieldNames := params[1:]

	query := r.DB.Table(tableName).Where(fmt.Sprintf("%s = ?", fieldNames[0]), requestValue)
	for _, fieldName := range fieldNames[1:] {
		query = query.Or(fmt.Sprintf("%s = ?", fieldName), requestValue)
	}

	var count int64
	err := query.Count(&count).Error
	if err != nil {
		return false
	}

	return count != 0
}
