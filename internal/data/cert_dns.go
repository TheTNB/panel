package data

import (
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/internal/panel"
)

type certDNSRepo struct{}

func NewCertDNSRepo() biz.CertDNSRepo {
	return &certDNSRepo{}
}

func (r certDNSRepo) List(page, limit uint) ([]*biz.CertDNS, int64, error) {
	var certDNS []*biz.CertDNS
	var total int64
	err := panel.Orm.Model(&biz.CertDNS{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&certDNS).Error
	return certDNS, total, err
}

func (r certDNSRepo) Get(id uint) (*biz.CertDNS, error) {
	certDNS := new(biz.CertDNS)
	err := panel.Orm.Model(&biz.CertDNS{}).Where("id = ?", id).First(certDNS).Error
	return certDNS, err
}

func (r certDNSRepo) Create(req *request.CertDNSCreate) (*biz.CertDNS, error) {
	certDNS := &biz.CertDNS{
		Name: req.Name,
		Type: req.Type,
		Data: req.Data,
	}

	if err := panel.Orm.Create(certDNS).Error; err != nil {
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

	return panel.Orm.Save(cert).Error
}

func (r certDNSRepo) Delete(id uint) error {
	return panel.Orm.Model(&biz.CertDNS{}).Where("id = ?", id).Delete(&biz.CertDNS{}).Error
}
