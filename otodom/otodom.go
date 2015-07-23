package otodom

import (
	"fmt"
	"github.com/kpawlik/exportms/utils"
	"path/filepath"
	"time"
)

const (
	FtpFolderDateFormat = "20060102_150405"
)

func SendToFtp(companyDomain, zipPath string, ftp *utils.FTP) (err error) {
	remoteDirName := fmt.Sprintf("%s_%s", companyDomain, time.Now().Format(FtpFolderDateFormat))
	_, fileName := filepath.Split(zipPath)
	if err = ftp.SendFile(remoteDirName, fileName, zipPath); err != nil {
		return
	}
	return
}
