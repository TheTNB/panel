package data

import (
	"context"
	"errors"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/http/request"
	"github.com/TheTNB/panel/internal/panel"
	"github.com/TheTNB/panel/pkg/acme"
	"github.com/TheTNB/panel/pkg/cert"
)

type certAccountRepo struct{}

func NewCertAccountRepo() biz.CertAccountRepo {
	return &certAccountRepo{}
}

func (r certAccountRepo) List(page, limit uint) ([]*biz.CertAccount, int64, error) {
	var accounts []*biz.CertAccount
	var total int64
	err := panel.Orm.Model(&biz.CertAccount{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&accounts).Error
	return accounts, total, err
}

func (r certAccountRepo) Get(id uint) (*biz.CertAccount, error) {
	account := new(biz.CertAccount)
	err := panel.Orm.Model(&biz.CertAccount{}).Where("id = ?", id).First(account).Error
	return account, err
}

func (r certAccountRepo) Create(req *request.CertAccountCreate) (*biz.CertAccount, error) {
	account := new(biz.CertAccount)
	account.CA = req.CA
	account.Email = req.Email
	account.Kid = req.Kid
	account.HmacEncoded = req.HmacEncoded
	account.KeyType = req.KeyType

	var err error
	var client *acme.Client
	switch account.CA {
	case "letsencrypt":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CALetsEncrypt, nil, acme.KeyType(account.KeyType))
	case "buypass":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CABuypass, nil, acme.KeyType(account.KeyType))
	case "zerossl":
		eab, eabErr := r.getZeroSSLEAB(account.Email)
		if eabErr != nil {
			return nil, eabErr
		}
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAZeroSSL, eab, acme.KeyType(account.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CASSLcom, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	case "google":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAGoogle, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	default:
		return nil, errors.New("CA 提供商不支持")
	}

	if err != nil {
		return nil, errors.New("向 CA 注册账号失败，请检查参数是否正确")
	}

	privateKey, err := cert.EncodeKey(client.Account.PrivateKey)
	if err != nil {
		return nil, errors.New("获取私钥失败")
	}
	account.PrivateKey = string(privateKey)

	if err = panel.Orm.Create(account).Error; err != nil {
		return nil, err
	}

	return account, nil
}

func (r certAccountRepo) Update(req *request.CertAccountUpdate) error {
	account, err := r.Get(req.ID)
	if err != nil {
		return err
	}

	account.CA = req.CA
	account.Email = req.Email
	account.Kid = req.Kid
	account.HmacEncoded = req.HmacEncoded
	account.KeyType = req.KeyType

	var client *acme.Client
	switch account.CA {
	case "letsencrypt":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CALetsEncrypt, nil, acme.KeyType(account.KeyType))
	case "buypass":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CABuypass, nil, acme.KeyType(account.KeyType))
	case "zerossl":
		eab, eabErr := r.getZeroSSLEAB(account.Email)
		if eabErr != nil {
			return eabErr
		}
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAZeroSSL, eab, acme.KeyType(account.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CASSLcom, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	case "google":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAGoogle, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	default:
		return errors.New("CA 提供商不支持")
	}

	if err != nil {
		return errors.New("向 CA 注册账号失败，请检查参数是否正确")
	}

	privateKey, err := cert.EncodeKey(client.Account.PrivateKey)
	if err != nil {
		return errors.New("获取私钥失败")
	}
	account.PrivateKey = string(privateKey)

	return panel.Orm.Save(account).Error
}

func (r certAccountRepo) Delete(id uint) error {
	return panel.Orm.Model(&biz.CertAccount{}).Where("id = ?", id).Delete(&biz.CertAccount{}).Error
}

// getZeroSSLEAB 获取 ZeroSSL EAB
func (r certAccountRepo) getZeroSSLEAB(email string) (*acme.EAB, error) {
	type data struct {
		Success    bool   `json:"success"`
		EabKid     string `json:"eab_kid"`
		EabHmacKey string `json:"eab_hmac_key"`
	}
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	resp, err := client.R().SetFormData(map[string]string{
		"email": email,
	}).SetResult(&data{}).Post("https://api.zerossl.com/acme/eab-credentials-email")
	if err != nil || !resp.IsSuccess() {
		return &acme.EAB{}, errors.New("获取ZeroSSL EAB失败")
	}
	eab := resp.Result().(*data)
	if !eab.Success {
		return &acme.EAB{}, errors.New("获取ZeroSSL EAB失败")
	}

	return &acme.EAB{KeyID: eab.EabKid, MACKey: eab.EabHmacKey}, nil
}
