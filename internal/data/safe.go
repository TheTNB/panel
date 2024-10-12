package data

import (
	"errors"
	"strings"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/os"
	"github.com/TheTNB/panel/pkg/shell"
	"github.com/TheTNB/panel/pkg/systemctl"
)

type safeRepo struct {
	ssh string
}

func NewSafeRepo() biz.SafeRepo {
	var ssh string
	if os.IsRHEL() {
		ssh = "sshd"
	} else {
		ssh = "ssh"
	}

	return &safeRepo{
		ssh: ssh,
	}
}

func (r *safeRepo) GetSSH() (uint, bool, error) {
	out, err := shell.Execf("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	if err != nil {
		return 0, false, err
	}

	running, err := systemctl.Status(r.ssh)
	if err != nil {
		return 0, false, err
	}

	return cast.ToUint(out), running, nil
}

func (r *safeRepo) UpdateSSH(port uint, status bool) error {
	oldPort, err := shell.Execf("cat /etc/ssh/sshd_config | grep 'Port ' | awk '{print $2}'")
	if err != nil {
		return err
	}

	_, _ = shell.Execf("sed -i 's/#Port %s/Port %d/g' /etc/ssh/sshd_config", oldPort, port)
	_, _ = shell.Execf("sed -i 's/Port %s/Port %d/g' /etc/ssh/sshd_config", oldPort, port)

	if !status {
		return systemctl.Stop(r.ssh)
	}

	return systemctl.Restart(r.ssh)
}

func (r *safeRepo) GetPingStatus() (bool, error) {
	out, err := shell.Execf(`firewall-cmd --list-all`)
	if err != nil {
		return true, errors.New(out)
	}

	if !strings.Contains(out, `rule protocol value="icmp" drop`) {
		return true, nil
	}

	return false, nil
}

func (r *safeRepo) UpdatePingStatus(status bool) error {
	var err error
	if status {
		_, err = shell.Execf(`firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'`)
	} else {
		_, err = shell.Execf(`firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'`)
	}
	if err != nil {
		return err
	}

	_, err = shell.Execf(`firewall-cmd --reload`)
	if err != nil {
		return err
	}

	return nil
}
