package data

import (
	"github.com/samber/do/v2"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
)

type certDNSRepo struct{}

func NewCertDNSRepo() biz.CertDNSRepo {
	return do.MustInvoke[biz.CertDNSRepo](injector)
}

func (r certDNSRepo) List(page, limit uint) ([]*biz.CertDNS, int64, error) {
	var certDNS []*biz.CertDNS
	var total int64
	err := app.Orm.Model(&biz.CertDNS{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&certDNS).Error
	return certDNS, total, err
}

func (r certDNSRepo) Get(id uint) (*biz.CertDNS, error) {
	certDNS := new(biz.CertDNS)
	err := app.Orm.Model(&biz.CertDNS{}).Where("id = ?", id).First(certDNS).Error
	return certDNS, err
}

func (r certDNSRepo) Create(req *request.CertDNSCreate) (*biz.CertDNS, error) {
	certDNS := &biz.CertDNS{
		Name: req.Name,
		Type: req.Type,
		Data: req.Data,
	}

	if err := app.Orm.Create(certDNS).Error; err != nil {
		return nil, err
	}

	return certDNS, nil
}

func (r certDNSRepo) Update(req *request.CertDNSUpdate) error {
	cert, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	cert.Name = req.Name
	cert.Type = req.Type
	cert.Data = req.Data

	return app.Orm.Save(cert).Error
}

func (r certDNSRepo) Delete(id uint) error {
	return app.Orm.Model(&biz.CertDNS{}).Where("id = ?", id).Delete(&biz.CertDNS{}).Error
}
