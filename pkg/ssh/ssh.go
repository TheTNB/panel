package ssh

import (
	"time"

	"golang.org/x/crypto/ssh"
)

type AuthMethod int8

const (
	PASSWORD AuthMethod = iota + 1
	PUBLICKEY
)

type ClientConfig struct {
	AuthMethod AuthMethod
	HostAddr   string
	User       string
	Password   string
	Key        string
	Timeout    time.Duration
}

func ClientConfigPassword(hostAddr, user, Password string) *ClientConfig {
	return &ClientConfig{
		Timeout:    time.Second * 5,
		AuthMethod: PASSWORD,
		HostAddr:   hostAddr,
		User:       user,
		Password:   Password,
	}
}

func ClientConfigPublicKey(hostAddr, user, key string) *ClientConfig {
	return &ClientConfig{
		Timeout:    time.Second * 5,
		AuthMethod: PUBLICKEY,
		HostAddr:   hostAddr,
		User:       user,
		Key:        key,
	}
}

func NewSSHClient(conf *ClientConfig) (*ssh.Client, error) {
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
	c, err := ssh.Dial("tcp", conf.HostAddr, config) // TODO support ipv6
	if err != nil {
		return nil, err
	}

	return c, nil
}

func parseKey(key string) (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(key))
}
