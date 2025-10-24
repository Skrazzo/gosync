package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type queue struct {
	uploads []string
	deletes []string
}

type SFTP struct {
	queue     *queue
	conn      *sftp.Client
	sshClient *ssh.Client
}

// Create connection to SFTP server
func (s *SFTP) Connect(cfg *Config) error {
	if cfg.AuthType == "password" {
		if err := s.connectWithPassword(cfg.Host, cfg.User, cfg.Password); err != nil {
			return err
		}
	} else if cfg.AuthType == "key" {
		// Try to connect to server using private key
		if err := s.connectWithKey(cfg.Host, cfg.User, cfg.PrivateKeyPath); err != nil {
			return err
		}
	}

	return nil
}

func (s *SFTP) connectWithPassword(host, user, password string) error {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use proper key checking in production!
		Timeout:         30 * time.Second,
	}

	// Connect to SSH
	sshConn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return err
	}

	// Create SFTP client
	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		sshConn.Close()
		return err
	}

	s.sshClient = sshConn
	s.conn = sftpClient
	return nil
}

// ensureConnected checks if the SFTP client is connected
func (s *SFTP) ensureConnected() error {
	if s.conn == nil {
		return fmt.Errorf("not connected to SFTP server")
	}
	return nil
}

// Close properly closes both SFTP and SSH connections
func (s *SFTP) Close() error {
	var errs []error

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing SFTP client: %w", err))
		}
		s.conn = nil
	}

	if s.sshClient != nil {
		if err := s.sshClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing SSH client: %w", err))
		}
		s.sshClient = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}
	return nil
}

func (s *SFTP) connectWithKey(host, user, keyPath string) error {
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
		Timeout:         30 * time.Second,
	}

	sshConn, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return err
	}

	sftpClient, err := sftp.NewClient(sshConn)
	if err != nil {
		sshConn.Close()
		return err
	}

	s.sshClient = sshConn
	s.conn = sftpClient
	return nil
}
