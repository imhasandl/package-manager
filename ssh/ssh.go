package ssh

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func createSSHClient() (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			ssh.Password("password"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", "localhost:22", config)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к SSH: %v", err)
	}

	return client, nil
}

func UploadFile(archivePath, remotePath string) error {
	sshClient, err := createSSHClient()
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("не удалось создать SFTP клиент: %v", err)
	}
	defer sftpClient.Close()

	localFile, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	remoteFile, err := sftpClient.Create(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		return err
	}

	fmt.Printf("ZIP файл загружен: %s -> %s\n", archivePath, remotePath)
	return nil
}

func DownloadFile(remotePath, localPath string) error {
	sshClient, err := createSSHClient()
	if err != nil {
		return err
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("не удалось создать SFTP клиент: %v", err)
	}
	defer sftpClient.Close()

	err = os.MkdirAll(filepath.Dir(localPath), 0755)
	if err != nil {
		return err
	}

	remoteFile, err := sftpClient.Open(remotePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, remoteFile)
	if err != nil {
		return err
	}

	fmt.Printf("ZIP файл скачан: %s -> %s\n", remotePath, localPath)
	return nil
}
