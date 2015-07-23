package utils

import (
	"github.com/dutchcoders/goftp"
	"os"
)

type FTP struct {
	host     string
	user     string
	pass     string
	testMode bool
}

func NewFTP(host, user, pass string, testMode bool) *FTP {
	return &FTP{host: host, user: user, pass: pass, testMode: testMode}
}

func (f *FTP) SendFile(remoteDir, remoteFile, localFile string) (err error) {
	var (
		file *os.File
		ftp  *goftp.FTP
	)
	if ftp, err = goftp.Connect(f.host); err != nil {
		return
	}
	defer ftp.Close()
	if err = ftp.Login(f.user, f.pass); err != nil {
		return
	}
	if file, err = os.Open(localFile); err != nil {
		return
	}
	defer file.Close()
	if f.testMode {
		return
	}
	if err = ftp.Mkd(remoteDir); err != nil {
		return
	}
	if err = ftp.Cwd(remoteDir); err != nil {
		return
	}
	ftp.Stor(remoteFile, file)

	return
}
