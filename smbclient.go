package smbget

import (
	"fmt"
	"io"
	iofs "io/fs"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type SMBClient struct {
	share_name  string
	server_addr string

	dial *smb2.Dialer

	share *smb2.Share
}

func NewSMBClient(username string, password string, domain string, server_addr string, share_name string) SMBClient {

	dial := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     username,
			Password: password,
			Domain:   domain,
		},
	}

	s := SMBClient{
		share_name:  share_name,
		server_addr: server_addr,
		dial:        dial,
	}

	logger.Info("Created new SMBClient", slog.String("object", fmt.Sprintf("%+v", s)))

	return s
}

func (c *SMBClient) Connect() (err error) {

	if c.share != nil {
		return nil
	}

	conn, err := net.Dial("tcp", c.server_addr+":445")

	if err != nil {
		return err
	}

	logger.Info("Connected to SMB Server")

	session, err := c.dial.Dial(conn)

	if err != nil {
		return err
	}

	logger.Info("Connected to SMB Session")

	c.share, err = session.Mount(c.share_name)

	if err != nil {
		return err
	}

	logger.Info("Connected to SMB Share")

	return nil
}

func (c *SMBClient) Disconnect() {

	c.share.Umount()
}

func (c *SMBClient) GetFile(remotePath string, localPath string, progBarPad int) error {

	err := c.connect()

	if err != nil {
		return err
	}

	remote_file, err := c.share.Open(remotePath)

	if err != nil {
		return err
	}

	defer remote_file.Close()

	stats, err := c.share.Stat(remotePath)

	if err != nil {
		return err
	}

	size := stats.Size()
	padding := strings.Repeat(" ", progBarPad)

	bar := GetProgressBar(size, filepath.Base(remotePath), padding)

	local_file, err := os.Create(localPath)

	if err != nil {
		return err
	}

	defer local_file.Close()

	absLocalPath, _ := filepath.Abs(localPath)

	logger.Info("Trying to transfer file", slog.String("remote", remotePath), slog.String("local", absLocalPath))

	_, err = io.Copy(io.MultiWriter(local_file, bar), remote_file)

	if err != nil {
		return err
	}

	logger.Info("Successfully transferred file", slog.Float64("time", bar.State().SecondsSince), slog.Float64("size", bar.State().CurrentBytes))

	return bar.Close()
}

func (c *SMBClient) ReadFile(filePath string) (string, error) {

	err := c.connect()

	if err != nil {
		return "", err
	}

	file, err := c.share.Open(filePath)

	if err != nil {
		return "", err
	}

	defer file.Close()

	content, err := io.ReadAll(file)

	if err != nil {
		return "", err
	}

	return string(content), nil

}

func (c *SMBClient) GetDirectory(remotePath string, localPath string, excludeExt []string) (map[string][]string, error) {

	err := c.connect()

	if err != nil {
		return nil, err
	}

	stats, err := c.share.Stat(remotePath)

	if err != nil {
		return nil, fmt.Errorf("%s path does not exist", remotePath)
	}

	if !stats.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", remotePath)
	}

	matches, err := iofs.Glob(c.share.DirFS(remotePath), "*")

	if err != nil {
		return nil, err
	}

	filtered_matches := make([]string, 0)

	status := make(map[string][]string)

	longest := 0

	for _, ext := range excludeExt {
		for _, match := range matches {
			if !strings.HasSuffix(match, ext) {

				filtered_matches = append(filtered_matches, filepath.Join(remotePath, match))

				if len(match) > longest {
					longest = len(match)
				}
			}
		}
	}

	logger.Info("Found valid files for transfer", slog.Int("count", len(filtered_matches)))

	for _, match := range filtered_matches {

		err = c.GetFile(match, filepath.Join(localPath, filepath.Base(match)), longest-len(filepath.Base(match)))

		if err != nil {
			status["failed"] = append(status["failed"], filepath.Base(match))
			panic(err)
		} else {
			status["success"] = append(status["success"], filepath.Base(match))
		}
	}

	return status, nil
}
