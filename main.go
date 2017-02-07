package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kpawlik/exportms/gratka"
	//	"github.com/kpawlik/exportms/otodom"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kpawlik/exportms/utils"
)

const (
	// DateFormat main date format
	DateFormat = "2006_01_02_15"
	//GratkaExportType default export type
	GratkaExportType = "full"
	// CreadentalsFile name of file with credentials
	CreadentalsFile = "enc"
)

var (
	wg          sync.WaitGroup
	workDir     string
	noOfWorkers int
	sendOnly    bool
	exportTypes string
	tmpPrefix   string
	testMode    bool
	// DbUser database access user name
	DbUser string
	//DbPass database access passworld
	DbPass string
	// DbHost database host name
	DbHost string
	// DbName name of database
	DbName string
	// Domain database Domain
	Domain string
	// OfflineCode offline code
	OfflineCode string
	// FTPHost ftp host name
	FTPHost string
	// FTPLogin ftp login
	FTPLogin string
	// FTPPass ftp password
	FTPPass     string
	credentials utils.Credentials
	settings    *config
	//
	exports map[string]exportFunc
)

type config struct {
	workDir, tmpPrefix string
	sendOnly           bool
	exports            string
	noOfWorkers        int
	testMode           bool
}

type exportFunc func(string, *config) error

func init() {
	var (
		err error
	)
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&workDir, "o", "", "")
	flag.IntVar(&noOfWorkers, "n", 1, "")
	flag.BoolVar(&sendOnly, "s", false, "")
	flag.BoolVar(&testMode, "tm", false, "")
	flag.StringVar(&exportTypes, "e", "", "")
	flag.Parse()
	workDir, err = filepath.Abs(strings.ToLower(workDir))
	utils.LogErr(err, "Dest folder path")
	if len(workDir) == 0 {
		fmt.Println("Set out forlder path")
		os.Exit(0)
	}

	if runtime.GOOS == "windows" {
		tmpPrefix = `c:\tmp`
	} else {
		tmpPrefix = `/tmp/`
	}
	f, err := os.Create(`ms-expot.log`)
	if err != nil {
		panic(err)
	}
	wr := io.MultiWriter(f, os.Stdout)
	log.SetOutput(wr)
	settings = &config{workDir: workDir,
		noOfWorkers: noOfWorkers,
		tmpPrefix:   tmpPrefix,
		sendOnly:    sendOnly,
		exports:     exportTypes,
		testMode:    testMode}
	log.Println("Settings:")
	log.Printf("No of concurrent workers: %d\n", noOfWorkers)
	log.Printf("Work dir: %s\n", workDir)

	exports = map[string]exportFunc{
		"gratka": exportGratka,
		"otodom": exportGratka}

}

// validPaths validates path to folders, if necessery then
// it will remover and recreate folders
func validPaths(conf *config) (err error) {
	workDir := conf.workDir
	tmpPrefix := conf.tmpPrefix
	removeFirst := !conf.sendOnly
	_, err = os.Stat(workDir)
	dirExists := err == nil
	if removeFirst && dirExists && !strings.HasPrefix(workDir, tmpPrefix) {
		return utils.Errorf("Dest folder '%s' already exists is not in tmp folder so need to be removed manually\n", workDir)
	}
	if removeFirst && dirExists {
		if err = os.RemoveAll(workDir); err != nil {
			return
		}
	}
	if err = os.MkdirAll(workDir, os.ModePerm); err != nil {
		return
	}
	return os.MkdirAll(filepath.Join(workDir, "images"), os.ModePerm)
}

// setCredentials sets login/pass/ftp etc data for db acceess and exprot to ftp
func setCredentials(exportType string) {
	var (
		err error
	)
	if credentials == nil {
		if credentials, err = utils.GetCredentials(CreadentalsFile); err != nil {
			panic(err)
		}
	}
	if exportType == "db" {
		dbData := credentials[exportType]
		DbUser = dbData["user"]
		DbPass = dbData["pass"]
		DbHost = dbData["host"]
		DbName = dbData["name"]
		return
	}
	data := credentials[exportType]
	Domain = data["domain"]
	OfflineCode = data["offlineCode"]
	FTPHost = data["ftpHost"]
	FTPLogin = data["ftpLogin"]
	FTPPass = data["ftpPass"]
}

func exportGratka(name string, conf *config) (err error) {
	var (
		gratkaDir string
		zipPath   string
		ftp       *utils.FTP
	)
	setCredentials(name)
	zipPath = gratka.ZipPath(conf.workDir, Domain)
	if !conf.sendOnly {
		startTime := time.Now()
		if gratkaDir, err = gratka.Convert(workDir, OfflineCode, Domain); err != nil {
			return
		}
		log.Printf("Convertion time: %s\n", time.Now().Sub(startTime))
		if zipPath, err = gratka.Zip(workDir, gratkaDir, Domain); err != nil {
			return
		}
	}

	log.Printf("--- Send to  '%s'---\n", name)
	ftp = utils.NewFTP(FTPHost, FTPLogin, FTPPass, conf.testMode)
	return gratka.SendToFtp(Domain, zipPath, ftp)

}

func main() {
	var (
		err      error
		offersNo int
		types    []string
	)
	if len(settings.exports) == 0 {
		types = []string{"gratka", "otodom"}
	} else {
		types = strings.Split(settings.exports, ",")
		for _, typ := range types {
			if _, ok := exports[typ]; !ok {
				log.Fatalf("Unsupported export name %s", typ)
			}
		}
	}
	setCredentials("db")
	startTime := time.Now()
	// validate paths, remove / create dirctory
	if err = validPaths(settings); err != nil {
		utils.LogErrf(err, "Settings paths")
	}

	if !settings.sendOnly {
		// dump to XML
		log.Printf("Start getting data from DB\n")
		offersNo, err = dumpAsXML(settings)
		utils.LogErr(err, "Dump XML")
		log.Printf("Downloaded %d offers from DB (in %v)\n", offersNo, time.Now().Sub(startTime))
	}

	for _, name := range types {
		if err = exports[name](name, settings); err != nil {
			utils.LogErrf(err, "Error durring process %s", name)
		}
	}
	log.Printf("All done in %v\n", time.Now().Sub(startTime))
}
