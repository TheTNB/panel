package data

import (
	"gorm.io/gorm"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/internal/http/request"
)

type certDNSRepo struct {
	db *gorm.DB
}

func NewCertDNSRepo(db *gorm.DB) biz.CertDNSRepo {
	return &certDNSRepo{
		db: db,
	}
}

func (r certDNSRepo) List(page, limit uint) ([]*biz.CertDNS, int64, error) {
	var certDNS []*biz.CertDNS
	var total int64
	err := r.db.Model(&biz.CertDNS{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&certDNS).Error
	return certDNS, total, err
}

func (r certDNSRepo) Get(id uint) (*biz.CertDNS, error) {
	certDNS := new(biz.CertDNS)
	err := r.db.Model(&biz.CertDNS{}).Where("id = ?", id).First(certDNS).Error
	return certDNS, err
}

func (r certDNSRepo) Create(req *request.CertDNSCreate) (*biz.CertDNS, error) {
	certDNS := &biz.CertDNS{
		Name: req.Name,
		Type: req.Type,
		Data: req.Data,
	}

	if err := r.db.Create(certDNS).Error; err != nil {
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

	return r.db.Save(cert).Error
}

func (r certDNSRepo) Delete(id uint) error {
	return r.db.Model(&biz.CertDNS{}).Where("id = ?", id).Delete(&biz.CertDNS{}).Error
}
