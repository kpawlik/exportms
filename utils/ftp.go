package utils

import (
	"os"

	"github.com/dutchcoders/goftp"
	"github.com/mitchellh/ioprogress"
)

// FTP sype and methods
type FTP struct {
	host     string
	user     string
	pass     string
	testMode bool
}

// NewFTP new ftp object
func NewFTP(host, user, pass string, testMode bool) *FTP {
	return &FTP{host: host, user: user, pass: pass, testMode: testMode}
}

// SendFile send file to ftp
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
	info, _ := os.Stat(localFile)
	defer file.Close()
	reader := &ioprogress.Reader{Reader: file,
		Size:     info.Size(),
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, ioprogress.DrawTextFormatBytes)}
	if f.testMode {
		return
	}
	if err = ftp.Mkd(remoteDir); err != nil {
		return
	}
	if err = ftp.Cwd(remoteDir); err != nil {
		return
	}
	ftp.Stor(remoteFile, reader)

	return
}
