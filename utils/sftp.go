package utils

import (
	"fmt"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type queue struct {
	uploads []string
	deletes []string
}

type SFTP struct {
	queue *queue
	conn  *sftp.Client
}

// Create connection to SFTP server
func (s *SFTP) Connect(cfg *Config) error {
	if cfg.AuthType == "password" {
		fmt.Println("connected with password")
	} else if cfg.AuthType == "key" {
		// Try to connect to server using private key
		if err := connectWithKey(cfg.Host, cfg.User, cfg.PrivateKeyPath, s); err != nil {
			return err
		}
	}

	return nil
}

func connectWithPassword(host, user, password string) (*sftp.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use proper key checking in production!
	}

	// Connect to SSH
	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}

	// Create SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func connectWithKey(host, user, keyPath string, s *SFTP) error {
	// Read private key file
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return err
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}

	s.conn = client
	return nil
}
