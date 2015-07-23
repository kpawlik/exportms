package utils

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ZipFolder(folder, outName string) (err error) {
	var (
		dirToZip   *os.File
		filesToZip []string
		zFile      io.Writer
		zipFile    *os.File
		zipWriter  *zip.Writer
		content    []byte
	)

	if dirToZip, err = os.Open(folder); err != nil {
		return
	}
	defer dirToZip.Close()
	if filesToZip, err = dirToZip.Readdirnames(-1); err != nil {
		return
	}
	if zipFile, err = os.Create(outName); err != nil {
		return
	}
	defer zipFile.Close()
	zipWriter = zip.NewWriter(zipFile)

	for _, fileName := range filesToZip {
		if content, err = ioutil.ReadFile(filepath.Join(folder, fileName)); err != nil {
			return
		}
		if zFile, err = zipWriter.Create(fileName); err != nil {
			return
		}
		if _, err = zFile.Write(content); err != nil {
			return
		}
	}
	return zipWriter.Close()

}
