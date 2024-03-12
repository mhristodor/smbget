package smbget

import (
	"fmt"
	"io"
	iofs "io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type SMBClient struct{

	share_name string
	server_addr string

	dial *smb2.Dialer

	conn net.Conn
	session *smb2.Session
	share *smb2.Share

}


func NewSMBClient(username string, password string, domain string, server_addr string, share_name string) (SMBClient, error) {

	dial := &smb2.Dialer{
        Initiator: &smb2.NTLMInitiator{
            User:     username,
            Password: password,
            Domain:   domain,
        },
    }

	s := SMBClient{
		share_name: share_name,
		server_addr: server_addr,
		dial: dial,
	}


	err := s.connect()

	if err != nil {
		return s, err
	}

	return s, nil
}

func (c *SMBClient) connect() (err error) {

	if c.share != nil{
		return nil
	}

	if c.session != nil{
		return nil
	}

	if c.conn != nil{
		return nil
	}

	c.conn, err = net.Dial("tcp", c.share_name + ":445")

	if err != nil {
		return err
	}

	c.session, err = c.dial.Dial(c.conn)

	if err != nil {
		return err
	}

	c.share, err = c.session.Mount(c.share_name)

	if err != nil {
		return err
	}

	return nil
}


func (c *SMBClient) Disconnect() error {

	if c.share != nil{
		c.share.Umount()
	}

	if c.session != nil{
		c.session.Logoff()
	}

	if c.conn != nil{
		c.conn.Close()
	}

	return nil
}

func (c *SMBClient) GetFile(remotePath string, localPath string, progBarPad int) error {

	c.connect()
	
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

	io.Copy(io.MultiWriter(local_file, bar), remote_file)
	
	return nil
}

func (c *SMBClient) GetDirectory(remotePath string, localPath string, excludeExt []string) (map[string][]string, error) {
	
	c.connect()
	
	stats, err := c.share.Stat(remotePath)

	if err != nil {
		return nil, err
	}

	if !stats.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", remotePath)
	}


	matches, err := iofs.Glob(c.share.DirFS(remotePath), "*")
	if err != nil {
		panic(err)
	}

	filtered_matches := make([]string, len(matches))

	status := make(map[string][]string)

	longest := 0

	for _, ext := range excludeExt {
		for _, match := range matches {
			if filepath.Ext(match) != ext {
				
				filtered_matches = append(filtered_matches, filepath.Join(remotePath, match))
				
				if len(match) > longest {
					longest = len(match)
				}
			}
		}
	}

	for _, match := range filtered_matches {

		err = c.GetFile(match, filepath.Join(localPath, match), longest - len(filepath.Base(match)))

		if err != nil {
			status["failed"] = append(status["failed"], filepath.Base(match))
		}

		status["success"] = append(status["success"], filepath.Base(match))
	}

	return status, nil
}