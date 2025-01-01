package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm"

	"github.com/tnb-labs/panel/internal/biz"
	"github.com/tnb-labs/panel/internal/http/request"
	"github.com/tnb-labs/panel/pkg/acme"
	"github.com/tnb-labs/panel/pkg/cert"
)

type certAccountRepo struct {
	db   *gorm.DB
	user biz.UserRepo
}

func NewCertAccountRepo(db *gorm.DB, user biz.UserRepo) biz.CertAccountRepo {
	return &certAccountRepo{
		db:   db,
		user: user,
	}
}

func (r certAccountRepo) List(page, limit uint) ([]*biz.CertAccount, int64, error) {
	var accounts []*biz.CertAccount
	var total int64
	err := r.db.Model(&biz.CertAccount{}).Order("id desc").Count(&total).Offset(int((page - 1) * limit)).Limit(int(limit)).Find(&accounts).Error
	return accounts, total, err
}

func (r certAccountRepo) GetDefault(userID uint) (*biz.CertAccount, error) {
	user, err := r.user.Get(userID)
	if err != nil {
		return nil, err
	}

	account := new(biz.CertAccount)
	if err = r.db.Model(&biz.CertAccount{}).Where("ca = ?", "googlecn").Where("email = ?", user.Email).First(account).Error; err == nil {
		return account, nil
	}

	req := &request.CertAccountCreate{
		CA:      "googlecn",
		Email:   user.Email,
		KeyType: string(acme.KeyEC256),
	}

	return r.Create(req)
}

func (r certAccountRepo) Get(id uint) (*biz.CertAccount, error) {
	account := new(biz.CertAccount)
	err := r.db.Model(&biz.CertAccount{}).Where("id = ?", id).First(account).Error
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
	case "googlecn":
		eab, eabErr := r.getGoogleEAB()
		if eabErr != nil {
			return nil, eabErr
		}
		account.Kid = eab.KeyID
		account.HmacEncoded = eab.MACKey
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAGoogleCN, eab, acme.KeyType(account.KeyType))
	case "google":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAGoogle, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	case "letsencrypt":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CALetsEncrypt, nil, acme.KeyType(account.KeyType))
	case "buypass":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CABuypass, nil, acme.KeyType(account.KeyType))
	case "zerossl":
		eab, eabErr := r.getZeroSSLEAB(account.Email)
		if eabErr != nil {
			return nil, eabErr
		}
		account.Kid = eab.KeyID
		account.HmacEncoded = eab.MACKey
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAZeroSSL, eab, acme.KeyType(account.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CASSLcom, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	default:
		return nil, errors.New("unsupported CA")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to register account: %v", err)
	}

	privateKey, err := cert.EncodeKey(client.Account.PrivateKey)
	if err != nil {
		return nil, errors.New("failed to get private key")
	}
	account.PrivateKey = string(privateKey)

	if err = r.db.Create(account).Error; err != nil {
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
	case "googlecn":
		eab, eabErr := r.getGoogleEAB()
		if eabErr != nil {
			return eabErr
		}
		account.Kid = eab.KeyID
		account.HmacEncoded = eab.MACKey
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAGoogleCN, eab, acme.KeyType(account.KeyType))
	case "google":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAGoogle, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	case "letsencrypt":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CALetsEncrypt, nil, acme.KeyType(account.KeyType))
	case "buypass":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CABuypass, nil, acme.KeyType(account.KeyType))
	case "zerossl":
		eab, eabErr := r.getZeroSSLEAB(account.Email)
		if eabErr != nil {
			return eabErr
		}
		account.Kid = eab.KeyID
		account.HmacEncoded = eab.MACKey
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CAZeroSSL, eab, acme.KeyType(account.KeyType))
	case "sslcom":
		client, err = acme.NewRegisterAccount(context.Background(), account.Email, acme.CASSLcom, &acme.EAB{KeyID: account.Kid, MACKey: account.HmacEncoded}, acme.KeyType(account.KeyType))
	default:
		return errors.New("unsupported CA")
	}

	if err != nil {
		return errors.New("failed to register account")
	}

	privateKey, err := cert.EncodeKey(client.Account.PrivateKey)
	if err != nil {
		return errors.New("failed to get private key")
	}
	account.PrivateKey = string(privateKey)

	return r.db.Save(account).Error
}

func (r certAccountRepo) Delete(id uint) error {
	return r.db.Model(&biz.CertAccount{}).Where("id = ?", id).Delete(&biz.CertAccount{}).Error
}

// getGoogleEAB 获取 Google EAB
func (r certAccountRepo) getGoogleEAB() (*acme.EAB, error) {
	type data struct {
		Message string `json:"message"`
		Data    struct {
			KeyId  string `json:"key_id"`
			MacKey string `json:"mac_key"`
		} `json:"data"`
	}
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	client.SetRetryCount(2)

	resp, err := client.R().SetResult(&data{}).Get("https://gts.rat.dev/eab")
	if err != nil || !resp.IsSuccess() {
		return &acme.EAB{}, fmt.Errorf("failed to get Google EAB: %v", err)
	}
	eab := resp.Result().(*data)
	if eab.Message != "success" {
		return &acme.EAB{}, fmt.Errorf("failed to get Google EAB: %s", eab.Message)
	}

	return &acme.EAB{KeyID: eab.Data.KeyId, MACKey: eab.Data.MacKey}, nil
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
		return &acme.EAB{}, fmt.Errorf("failed to get ZeroSSL EAB: %v", err)
	}
	eab := resp.Result().(*data)
	if !eab.Success {
		return &acme.EAB{}, fmt.Errorf("failed to get ZeroSSL EAB")
	}

	return &acme.EAB{KeyID: eab.EabKid, MACKey: eab.EabHmacKey}, nil
}
