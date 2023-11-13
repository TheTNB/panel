package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"

	"panel/pkg/tools"
)

type AuthMethod int8

const (
	PASSWORD AuthMethod = iota + 1
	PUBLICKEY
)

type SSHClientConfig struct {
	AuthMethod AuthMethod
	HostAddr   string
	User       string
	Password   string
	KeyPath    string
	Timeout    time.Duration
}

func SSHClientConfigPassword(hostAddr, user, Password string) *SSHClientConfig {
	return &SSHClientConfig{
		Timeout:    time.Second * 5,
		AuthMethod: PASSWORD,
		HostAddr:   hostAddr,
		User:       user,
		Password:   Password,
	}
}

func SSHClientConfigPulicKey(hostAddr, user, keyPath string) *SSHClientConfig {
	return &SSHClientConfig{
		Timeout:    time.Second * 5,
		AuthMethod: PUBLICKEY,
		HostAddr:   hostAddr,
		User:       user,
		KeyPath:    keyPath,
	}
}

func NewSSHClient(conf *SSHClientConfig) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         conf.Timeout,
		User:            conf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	switch conf.AuthMethod {
	case PASSWORD:
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	case PUBLICKEY:
		signer, err := getKey(conf.KeyPath)
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	c, err := ssh.Dial("tcp", conf.HostAddr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func getKey(keyPath string) (ssh.Signer, error) {
	key, err := tools.Read(keyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey([]byte(key))
}
