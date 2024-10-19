package ssh

import (
	"os"
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
	KeyPath    string
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

func ClientConfigPublicKey(hostAddr, user, keyPath string) *ClientConfig {
	return &ClientConfig{
		Timeout:    time.Second * 5,
		AuthMethod: PUBLICKEY,
		HostAddr:   hostAddr,
		User:       user,
		KeyPath:    keyPath,
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
		signer, err := getKey(conf.KeyPath)
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

func getKey(keyPath string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	return ssh.ParsePrivateKey(key)
}
