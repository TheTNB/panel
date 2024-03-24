package internal

import (
	requests "panel/app/http/requests/cert"
	"panel/app/models"
	"panel/pkg/acme"
)

type Cert interface {
	UserStore(request requests.UserStore) error
	UserUpdate(request requests.UserUpdate) error
	UserShow(ID uint) (models.CertUser, error)
	UserDestroy(ID uint) error
	DNSStore(request requests.DNSStore) error
	DNSUpdate(request requests.DNSUpdate) error
	DNSShow(ID uint) (models.CertDNS, error)
	DNSDestroy(ID uint) error
	CertStore(request requests.CertStore) error
	CertUpdate(request requests.CertUpdate) error
	CertShow(ID uint) (models.Cert, error)
	CertDestroy(ID uint) error
	ObtainAuto(ID uint) (acme.Certificate, error)
	ObtainManual(ID uint) (acme.Certificate, error)
	ManualDNS(ID uint) ([]acme.DNSRecord, error)
	Renew(ID uint) (acme.Certificate, error)
	Deploy(ID, WebsiteID uint) error
}
