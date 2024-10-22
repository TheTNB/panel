package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"
)

type AuthMethod string

const (
	PASSWORD  AuthMethod = "password"
	PUBLICKEY AuthMethod = "publickey"
)

type ClientConfig struct {
	AuthMethod AuthMethod    `json:"auth_method"`
	Host       string        `json:"host"`
	User       string        `json:"user"`
	Password   string        `json:"password"`
	Key        string        `json:"key"`
	Timeout    time.Duration `json:"timeout"`
}

func ClientConfigPassword(host, user, Password string) *ClientConfig {
	return &ClientConfig{
		Timeout:    10 * time.Second,
		AuthMethod: PASSWORD,
		Host:       host,
		User:       user,
		Password:   Password,
	}
}

func ClientConfigPublicKey(host, user, key string) *ClientConfig {
	return &ClientConfig{
		Timeout:    10 * time.Second,
		AuthMethod: PUBLICKEY,
		Host:       host,
		User:       user,
		Key:        key,
	}
}

func NewSSHClient(conf ClientConfig) (*ssh.Client, error) {
	if conf.Timeout == 0 {
		conf.Timeout = 10 * time.Second
	}

	config := &ssh.ClientConfig{}
	config.SetDefaults()
	config.Timeout = conf.Timeout
	config.User = conf.User
	config.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	switch conf.AuthMethod {
	case PASSWORD:
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Password)}
	case PUBLICKEY:
		signer, err := parseKey(conf.Key)
		if err != nil {
			return nil, err
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}
	c, err := ssh.Dial("tcp", conf.Host, config) // TODO support ipv6
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseKey(key string) (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(key))
}
