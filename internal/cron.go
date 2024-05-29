package internal

import "github.com/TheTNB/panel/app/models"

type Cron interface {
	AddToSystem(cron models.Cron) error
	DeleteFromSystem(cron models.Cron) error
}
