package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type queue struct {
	Uploads []string
	Deletes []string
}

type SFTP struct {
	Queue     *queue
	QueueMu   sync.Mutex
	Conn      *sftp.Client
	SshClient *ssh.Client
	RemoteDir string // Store the remote directory for uploads
}

func NewSftp() *SFTP {
	return &SFTP{
		Queue: &queue{
			Uploads: make([]string, 0),
			Deletes: make([]string, 0),
		},
	}
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

	// Store remote directory for uploads
	s.RemoteDir = cfg.RemoteDir

	return nil
}

// ensureConnected checks if the SFTP client is connected
func (s *SFTP) ensureConnected() error {
	if s.Conn == nil {
		return fmt.Errorf("not connected to SFTP server")
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

	s.SshClient = sshConn
	s.Conn = sftpClient
	return nil
}

// Upload uploads a single file to the remote server
// localPath: relative path to the local file (e.g., "main.go" or "utils/config.go")
// Returns error if upload fails
func (s *SFTP) Upload(localPath string) error {
	// Ensure we're connected
	if err := s.ensureConnected(); err != nil {
		return err
	}

	// Open local file for reading
	localFile, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file %s: %w", localPath, err)
	}
	defer localFile.Close()

	// Construct remote path (remoteDir + localPath)
	remotePath := filepath.Join(s.RemoteDir, localPath)

	// Ensure remote directory exists
	remoteDir := filepath.Dir(remotePath)
	if err := s.Conn.MkdirAll(remoteDir); err != nil {
		return fmt.Errorf("failed to create remote directory %s: %w", remoteDir, err)
	}

	// Create remote file
	remoteFile, err := s.Conn.Create(remotePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file %s: %w", remotePath, err)
	}
	defer remoteFile.Close()

	// Copy file contents
	bytesWritten, err := io.Copy(remoteFile, localFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents to %s: %w", remotePath, err)
	}

	fmt.Printf("✓ Uploaded %s (%d bytes) -> %s\n", localPath, bytesWritten, remotePath)
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

	s.SshClient = sshConn
	s.Conn = sftpClient
	return nil
}

// Close properly closes both SFTP and SSH connections
func (s *SFTP) Close() error {
	var errs []error

	if s.Conn != nil {
		if err := s.Conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing SFTP client: %w", err))
		}
		s.Conn = nil
	}

	if s.SshClient != nil {
		if err := s.SshClient.Close(); err != nil {
			errs = append(errs, fmt.Errorf("error closing SSH client: %w", err))
		}
		s.SshClient = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing connections: %v", errs)
	}
	return nil
}
