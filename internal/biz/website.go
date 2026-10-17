package biz

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/leonelquinteros/gotext"
	"github.com/spf13/cast"
	"gorm.io/gorm"

	"github.com/acepanel/panel/v3/internal/request"
	"github.com/acepanel/panel/v3/pkg/acme"
	"github.com/acepanel/panel/v3/pkg/types"
)

type WebsiteType string

const (
	WebsiteTypeProxy  WebsiteType = "proxy"
	WebsiteTypePHP    WebsiteType = "php"
	WebsiteTypeStatic WebsiteType = "static"
)

type Website struct {
	ID        uint        `gorm:"primaryKey" json:"id"`
	Name      string      `gorm:"not null;default:'';unique" json:"name"`
	Type      WebsiteType `gorm:"not null;index;default:'static'" json:"type"`
	Status    bool        `gorm:"not null;default:true" json:"status"`
	Path      string      `gorm:"not null;default:''" json:"path"`
	SSL       bool        `gorm:"not null;default:false" json:"ssl"`
	CertID    uint        `gorm:"not null;index;default:0" json:"cert_id"` // 部署的证书 ID
	Remark    string      `gorm:"not null;default:''" json:"remark"`
	ExpireAt  *time.Time  `json:"expire_at"` // 到期时间，nil 表示不限时
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`

	CertExpire string   `gorm:"-:all" json:"cert_expire"` // 仅显示
	PHP        uint     `gorm:"-:all" json:"php"`         // 仅显示
	Domains    []string `gorm:"-:all" json:"domains"`     // 仅显示

	Cert *Cert `gorm:"foreignKey:CertID" json:"cert"`
}

type WebsiteRepo interface {
	GetRewrites() (map[string]string, error)
	UpdateDefaultConfig(ctx context.Context, req *request.WebsiteDefaultConfig) error
	Count() (int64, error)
	Get(id uint) (*types.WebsiteSetting, error)
	GetByName(name string) (*types.WebsiteSetting, error)
	List(typ, keyword string, page, limit uint) ([]*Website, int64, error)
	Create(ctx context.Context, req *request.WebsiteCreate) (*Website, error)
	Update(ctx context.Context, req *request.WebsiteUpdate) (*Website, error)
	SwitchType(ctx context.Context, req *request.WebsiteSwitchType) (*Website, error)
	GetForDelete(id uint) (*Website, error)
	RemoveFiles(ctx context.Context, name string, removePath bool) error
	Delete(website *Website) error
	ReloadWebServer(ctx context.Context) error
	UpdateRemark(id uint, remark string) error
	ResetConfig(ctx context.Context, id uint) error
	Rebuild(ctx context.Context, website *Website) (bool, []string, error)
	UpdateStatus(ctx context.Context, id uint, status bool) error
	UpdateExpireAt(id uint, expireAt *time.Time) error
	UpdateCert(ctx context.Context, req *request.WebsiteUpdateCert) error
}

type WebsiteUsecase struct {
	repo           WebsiteRepo
	log            *slog.Logger
	t              *gotext.Locale
	cert           *CertUsecase
	certAccount    *CertAccountUsecase
	database       *DatabaseUsecase
	databaseUser   *DatabaseUserUsecase
	databaseServer DatabaseServerRepo
	tamper         *TamperUsecase
	stat           *WebsiteStatUsecase
}

func NewWebsiteUsecase(certAccountUsecase *CertAccountUsecase, certUsecase *CertUsecase, databaseUsecase *DatabaseUsecase, databaseUserUsecase *DatabaseUserUsecase, tamperUsecase *TamperUsecase, websiteStatUsecase *WebsiteStatUsecase, t *gotext.Locale, log *slog.Logger, databaseServerRepo DatabaseServerRepo, websiteRepo WebsiteRepo) *WebsiteUsecase {
	return &WebsiteUsecase{
		repo:           websiteRepo,
		log:            log,
		t:              t,
		cert:           certUsecase,
		certAccount:    certAccountUsecase,
		database:       databaseUsecase,
		databaseUser:   databaseUserUsecase,
		databaseServer: databaseServerRepo,
		tamper:         tamperUsecase,
		stat:           websiteStatUsecase,
	}
}

func (uc *WebsiteUsecase) GetRewrites() (map[string]string, error) {
	return uc.repo.GetRewrites()
}

func (uc *WebsiteUsecase) UpdateDefaultConfig(ctx context.Context, req *request.WebsiteDefaultConfig) error {
	return uc.repo.UpdateDefaultConfig(ctx, req)
}

func (uc *WebsiteUsecase) Count() (int64, error) {
	return uc.repo.Count()
}

func (uc *WebsiteUsecase) Get(id uint) (*types.WebsiteSetting, error) {
	return uc.repo.Get(id)
}

func (uc *WebsiteUsecase) GetByName(name string) (*types.WebsiteSetting, error) {
	return uc.repo.GetByName(name)
}

func (uc *WebsiteUsecase) List(typ, keyword string, page, limit uint) ([]*Website, int64, error) {
	return uc.repo.List(typ, keyword, page, limit)
}

func (uc *WebsiteUsecase) Create(ctx context.Context, req *request.WebsiteCreate) (*Website, error) {
	w, err := uc.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	// 记录日志
	uc.log.Info("website created", slog.String("type", OperationTypeWebsite), slog.Uint64("operator_id", operatorID(ctx)), slog.String("name", req.Name), slog.String("website_type", req.Type), slog.String("path", req.Path))

	// 重载 Web 服务器
	if err = uc.repo.ReloadWebServer(ctx); err != nil {
		return nil, err
	}

	// 创建数据库
	name := "local_" + req.DBType
	if req.DB {
		server, err := uc.databaseServer.GetByName(ctx, name)
		if err != nil {
			return nil, errors.New(uc.t.Get("can't find %s database server, please add it first", name))
		}
		if err = uc.database.Create(ctx, &request.DatabaseCreate{
			ServerID:   server.ID,
			Name:       req.DBName,
			CreateUser: true,
			Username:   req.DBUser,
			Password:   req.DBPassword,
			Host:       "localhost",
			Comment:    "website " + req.Name,
		}); err != nil {
			return nil, err
		}
	}

	return w, nil
}

func (uc *WebsiteUsecase) Update(ctx context.Context, req *request.WebsiteUpdate) error {
	website, err := uc.repo.Update(ctx, req)
	if err != nil {
		return err
	}

	// 记录日志
	uc.log.Info("website updated", slog.String("type", OperationTypeWebsite), slog.Uint64("operator_id", operatorID(ctx)), slog.Uint64("id", uint64(req.ID)), slog.String("name", website.Name))

	return uc.repo.ReloadWebServer(ctx)
}

func (uc *WebsiteUsecase) SwitchType(ctx context.Context, req *request.WebsiteSwitchType) error {
	website, err := uc.repo.SwitchType(ctx, req)
	if err != nil {
		return err
	}

	uc.log.Info("website type switched", slog.String("type", OperationTypeWebsite), slog.Uint64("operator_id", operatorID(ctx)), slog.Uint64("id", uint64(req.ID)), slog.String("name", website.Name), slog.String("website_type", req.Type))

	return nil
}

func (uc *WebsiteUsecase) Delete(ctx context.Context, req *request.WebsiteDelete) error {
	website, err := uc.repo.GetForDelete(req.ID)
	if err != nil {
		return err
	}

	// 删站是删文件、删库、删记录、reload 的连续破坏性动作，中途取消只会留下半删的站点
	ctx = context.WithoutCancel(ctx)

	// 清理防篡改规则
	if req.Path && website.Path != "" {
		if rules, listErr := uc.tamper.ListRules(); listErr == nil {
			cleaned := filepath.Clean(website.Path)
			for _, rule := range rules {
				if rule.Path != "" && filepath.Clean(rule.Path) == cleaned {
					_ = uc.tamper.DeleteRule(rule.ID)
				}
			}
		}
	}

	_ = uc.repo.RemoveFiles(ctx, website.Name, req.Path)
	if req.DB {
		if mysql, err := uc.databaseServer.GetByName(ctx, "local_mysql"); err == nil {
			_ = uc.databaseUser.DeleteByNames(ctx, mysql.ID, []string{website.Name})
			_ = uc.database.Delete(ctx, mysql.ID, website.Name)
		}
		if postgres, err := uc.databaseServer.GetByName(ctx, "local_postgresql"); err == nil {
			_ = uc.databaseUser.DeleteByNames(ctx, postgres.ID, []string{website.Name})
			_ = uc.database.Delete(ctx, postgres.ID, website.Name)
		}
	}

	if err = uc.repo.Delete(website); err != nil {
		return err
	}

	// 清理该站点在统计库中的全部数据
	_ = uc.stat.DeleteBySite(website.Name)

	// 记录日志
	uc.log.Info("website deleted", slog.String("type", OperationTypeWebsite), slog.Uint64("operator_id", operatorID(ctx)), slog.Uint64("id", uint64(req.ID)), slog.String("name", website.Name))

	return uc.repo.ReloadWebServer(ctx)
}

func (uc *WebsiteUsecase) UpdateRemark(id uint, remark string) error {
	return uc.repo.UpdateRemark(id, remark)
}

// Rebuild 按当前 Web 服务器重建站点配置
func (uc *WebsiteUsecase) Rebuild(ctx context.Context, website *Website) (bool, []string, error) {
	return uc.repo.Rebuild(ctx, website)
}

func (uc *WebsiteUsecase) ResetConfig(ctx context.Context, id uint) error {
	return uc.repo.ResetConfig(ctx, id)
}

func (uc *WebsiteUsecase) UpdateStatus(ctx context.Context, id uint, status bool) error {
	return uc.repo.UpdateStatus(ctx, id, status)
}

func (uc *WebsiteUsecase) UpdateExpireAt(id uint, expireAt *time.Time) error {
	return uc.repo.UpdateExpireAt(id, expireAt)
}

func (uc *WebsiteUsecase) UpdateCert(ctx context.Context, req *request.WebsiteUpdateCert) error {
	return uc.repo.UpdateCert(ctx, req)
}

func (uc *WebsiteUsecase) ObtainCert(ctx context.Context, id uint, dnsID uint) error {
	website, err := uc.repo.Get(id)
	if err != nil {
		return err
	}

	// 泛域名必须使用 DNS 验证
	hasWildcard := slices.ContainsFunc(website.Domains, func(d string) bool {
		return strings.Contains(d, "*")
	})
	if hasWildcard && dnsID == 0 {
		return errors.New(uc.t.Get("wildcard domains require DNS verification, please select a DNS provider"))
	}

	account, err := uc.certAccount.GetDefault(ctx, cast.ToUint(ctx.Value("user_id")))
	if err != nil {
		return err
	}

	newCert, err := uc.cert.GetByWebsite(website.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newCert, err = uc.cert.Create(ctx, &request.CertCreate{
				Type:        string(acme.KeyEC256),
				Domains:     website.Domains,
				AutoRenewal: true,
				AccountID:   account.ID,
				DNSID:       dnsID,
				WebsiteIDs:  []uint{website.ID},
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	newCert.Domains = website.Domains
	newCert.DNSID = dnsID
	if err = uc.cert.Save(newCert); err != nil {
		return err
	}

	_, err = uc.cert.ObtainAutoWithProgressCallback(ctx, newCert.ID, nil)
	if err != nil {
		return err
	}

	return uc.cert.Deploy(ctx, newCert.ID, []uint{website.ID}, true)
}
