package percona

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/app"
	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/internal/data"
	"github.com/TheTNB/panel/internal/service"
	"github.com/TheTNB/panel/pkg/db"
	"github.com/TheTNB/panel/pkg/io"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/str"
	"github.com/TheTNB/panel/pkg/systemctl"
	"github.com/TheTNB/panel/pkg/types"
)

type Service struct {
	settingRepo biz.SettingRepo
}

func NewService() *Service {
	return &Service{
		settingRepo: data.NewSettingRepo(),
	}
}

// GetConfig 获取配置
func (s *Service) GetConfig(w http.ResponseWriter, r *http.Request) {
	config, err := io.Read(app.Root + "/server/mysql/conf/my.cnf")
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 Percona 配置失败")
		return
	}

	service.Success(w, config)
}

// UpdateConfig 保存配置
func (s *Service) UpdateConfig(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[UpdateConfig](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	if err := io.Write(app.Root+"/server/mysql/conf/my.cnf", req.Config, 0644); err != nil {
		service.Error(w, http.StatusInternalServerError, "写入 Percona 配置失败")
		return
	}

	if err := systemctl.Reload("mysqld"); err != nil {
		service.Error(w, http.StatusInternalServerError, "重载 Percona 失败")
		return
	}

	service.Success(w, nil)
}

// Load 获取负载
func (s *Service) Load(w http.ResponseWriter, r *http.Request) {
	rootPassword, err := s.settingRepo.Get(biz.SettingKeyPerconaRootPassword)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 Percona root密码失败")
		return

	}
	if len(rootPassword) == 0 {
		service.Error(w, http.StatusUnprocessableEntity, "Percona root密码为空")
		return
	}

	status, _ := systemctl.Status("mysqld")
	if !status {
		service.Success(w, []types.NV{})
		return
	}

	raw, err := shell.Execf(`mysqladmin -u root -p "%s" extended-status`, rootPassword)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 Percona 负载失败")
		return
	}

	var load []map[string]string
	expressions := []struct {
		regex string
		name  string
	}{
		{`Uptime\s+\|\s+(\d+)\s+\|`, "运行时间"},
		{`Queries\s+\|\s+(\d+)\s+\|`, "总查询次数"},
		{`Connections\s+\|\s+(\d+)\s+\|`, "总连接次数"},
		{`Com_commit\s+\|\s+(\d+)\s+\|`, "每秒事务"},
		{`Com_rollback\s+\|\s+(\d+)\s+\|`, "每秒回滚"},
		{`Bytes_sent\s+\|\s+(\d+)\s+\|`, "发送"},
		{`Bytes_received\s+\|\s+(\d+)\s+\|`, "接收"},
		{`Threads_connected\s+\|\s+(\d+)\s+\|`, "活动连接数"},
		{`Max_used_connections\s+\|\s+(\d+)\s+\|`, "峰值连接数"},
		{`Key_read_requests\s+\|\s+(\d+)\s+\|`, "索引命中率"},
		{`Innodb_buffer_pool_reads\s+\|\s+(\d+)\s+\|`, "Innodb索引命中率"},
		{`Created_tmp_disk_tables\s+\|\s+(\d+)\s+\|`, "创建临时表到磁盘"},
		{`Open_tables\s+\|\s+(\d+)\s+\|`, "已打开的表"},
		{`Select_full_join\s+\|\s+(\d+)\s+\|`, "没有使用索引的量"},
		{`Select_full_range_join\s+\|\s+(\d+)\s+\|`, "没有索引的JOIN量"},
		{`Select_range_check\s+\|\s+(\d+)\s+\|`, "没有索引的子查询量"},
		{`Sort_merge_passes\s+\|\s+(\d+)\s+\|`, "排序后的合并次数"},
		{`Table_locks_waited\s+\|\s+(\d+)\s+\|`, "锁表次数"},
	}

	for _, expression := range expressions {
		re := regexp.MustCompile(expression.regex)
		matches := re.FindStringSubmatch(raw)
		if len(matches) > 1 {
			d := map[string]string{"name": expression.name, "value": matches[1]}
			if expression.name == "发送" || expression.name == "接收" {
				d["value"] = str.FormatBytes(cast.ToFloat64(matches[1]))
			}

			load = append(load, d)
		}
	}

	// 索引命中率
	readRequests := cast.ToFloat64(load[9]["value"])
	reads := cast.ToFloat64(load[10]["value"])
	load[9]["value"] = fmt.Sprintf("%.2f%%", readRequests/(reads+readRequests)*100)
	// Innodb 索引命中率
	bufferPoolReads := cast.ToFloat64(load[11]["value"])
	bufferPoolReadRequests := cast.ToFloat64(load[12]["value"])
	load[10]["value"] = fmt.Sprintf("%.2f%%", bufferPoolReadRequests/(bufferPoolReads+bufferPoolReadRequests)*100)

	service.Success(w, load)
}

// ErrorLog 获取错误日志
func (s *Service) ErrorLog(w http.ResponseWriter, r *http.Request) {
	log, err := shell.Execf("tail -n 100 %s/server/mysql/mysql-error.log", app.Root)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, log)
		return
	}

	service.Success(w, log)
}

// ClearErrorLog 清空错误日志
func (s *Service) ClearErrorLog(w http.ResponseWriter, r *http.Request) {
	if _, err := shell.Execf("echo '' > %s/server/mysql/mysql-error.log", app.Root); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

// SlowLog 获取慢查询日志
func (s *Service) SlowLog(w http.ResponseWriter, r *http.Request) {
	log, err := shell.Execf("tail -n 100 %s/server/mysql/mysql-slow.log", app.Root)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, log)
		return
	}

	service.Success(w, log)
}

// ClearSlowLog 清空慢查询日志
func (s *Service) ClearSlowLog(w http.ResponseWriter, r *http.Request) {
	if _, err := shell.Execf("echo '' > %s/server/mysql/mysql-slow.log", app.Root); err != nil {
		service.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	service.Success(w, nil)
}

// GetRootPassword 获取root密码
func (s *Service) GetRootPassword(w http.ResponseWriter, r *http.Request) {
	rootPassword, err := s.settingRepo.Get(biz.SettingKeyPerconaRootPassword)
	if err != nil {
		service.Error(w, http.StatusInternalServerError, "获取 Percona root密码失败")
		return
	}
	if len(rootPassword) == 0 {
		service.Error(w, http.StatusUnprocessableEntity, "Percona root密码为空")
		return
	}

	service.Success(w, rootPassword)
}

// SetRootPassword 设置root密码
func (s *Service) SetRootPassword(w http.ResponseWriter, r *http.Request) {
	req, err := service.Bind[SetRootPassword](r)
	if err != nil {
		service.Error(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	oldRootPassword, _ := s.settingRepo.Get(biz.SettingKeyPerconaRootPassword)
	mysql, err := db.NewMySQL("root", oldRootPassword, s.getSock(), "unix")
	if err != nil {
		// 尝试安全模式直接改密
		if err = db.MySQLResetRootPassword(req.Password); err != nil {
			service.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err = mysql.UserPassword("root", req.Password); err != nil {
			service.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	if err = s.settingRepo.Set(biz.SettingKeyPerconaRootPassword, req.Password); err != nil {
		service.Error(w, http.StatusInternalServerError, fmt.Sprintf("设置保存失败: %v", err))
		return
	}

	service.Success(w, nil)
}

func (s *Service) getSock() string {
	if io.Exists("/tmp/mysql.sock") {
		return "/tmp/mysql.sock"
	}
	if io.Exists(app.Root + "/server/mysql/config/my.cnf") {
		config, _ := io.Read(app.Root + "/server/mysql/config/my.cnf")
		re := regexp.MustCompile(`socket\s*=\s*(['"]?)([^'"]+)`)
		matches := re.FindStringSubmatch(config)
		if len(matches) > 2 {
			return matches[2]
		}
	}
	if io.Exists("/etc/my.cnf") {
		config, _ := io.Read("/etc/my.cnf")
		re := regexp.MustCompile(`socket\s*=\s*(['"]?)([^'"]+)`)
		matches := re.FindStringSubmatch(config)
		if len(matches) > 2 {
			return matches[2]
		}
	}

	return "/tmp/mysql.sock"
}
