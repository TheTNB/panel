package data

import (
	"errors"
	"strings"

	"github.com/spf13/cast"

	"github.com/TheTNB/panel/internal/biz"
	"github.com/TheTNB/panel/pkg/io"
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
	if os.IsRHEL() {
		out, err := shell.Execf(`firewall-cmd --list-all`)
		if err != nil {
			return true, errors.New(out)
		}

		if !strings.Contains(out, `rule protocol value="icmp" drop`) {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		config, err := io.Read("/etc/ufw/before.rules")
		if err != nil {
			return true, err
		}
		if strings.Contains(config, "-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT") {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (r *safeRepo) UpdatePingStatus(status bool) error {
	var out string
	var err error
	if os.IsRHEL() {
		if status {
			out, err = shell.Execf(`firewall-cmd --permanent --remove-rich-rule='rule protocol value=icmp drop'`)
		} else {
			out, err = shell.Execf(`firewall-cmd --permanent --add-rich-rule='rule protocol value=icmp drop'`)
		}
	} else {
		if status {
			out, err = shell.Execf(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/g' /etc/ufw/before.rules`)
		} else {
			out, err = shell.Execf(`sed -i 's/-A ufw-before-input -p icmp --icmp-type echo-request -j ACCEPT/-A ufw-before-input -p icmp --icmp-type echo-request -j DROP/g' /etc/ufw/before.rules`)
		}
	}

	if err != nil {
		return errors.New(out)
	}

	if os.IsRHEL() {
		out, err = shell.Execf(`firewall-cmd --reload`)
	} else {
		out, err = shell.Execf(`ufw reload`)
	}

	if err != nil {
		return errors.New(out)
	}

	return nil
}
