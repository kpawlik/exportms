package gratka

import (
	"fmt"
	"io"
	"log"
	"github.com/kpawlik/exportms/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	DateTimeFormat      = "2006-01-02 15:04:05"
	DateFormat          = "2006-01-02"
	FtpFolderDateFormat = "20060102_150405"
)

var (
	sectionMap = map[string]string{
		"spr": "1",
		"dow": "4",
		"wyn": "5",
		"kup": "2"}
	subSectionMap = map[string]string{
		"dom":        "2",
		"mieszkanie": "1",
		"lokal":      "4",
		"dzialka":    "3",
		"garaz":      "5",
		"biuro":      "8"}
)

func getSectionId(name, typ string) string {
	var (
		sMap map[string]string
	)
	switch typ {
	case "rubryka":
		sMap = sectionMap
	case "podrubryka":
		sMap = subSectionMap
		name = unifyGroup(name)
	default:
		log.Panicf("Unregognized type: '%s'\n", typ)
	}
	name = strings.ToLower(name)
	for k, v := range sMap {
		if strings.HasPrefix(name, k) {
			return v
		}
	}
	log.Panicf("Unregognized SecionId, type: '%s', value: '%s'\n", name)
	return ""
}

func getFloor(floor string) string {
	var (
		no  int
		err error
	)
	if no, err = strconv.Atoi(floor); err != nil {
		return ""
	}
	no += 1
	switch {
	case no > 22:
		return "22"
	case no < 1:
		return ""
	default:
		return strconv.Itoa(no)
	}
	return ""
}

func getNoOfFloors(floor, offerType string) string {
	no, err := strconv.Atoi(floor)
	switch {
	case err != nil:
		return ""
	case no < 0:
		return ""
	case no > 21:
		return "21"
	case offerType == "dom" && no < 6:
		no += 1
		return strconv.Itoa(no)
	default:
		return floor
	}
	return ""
}

func getNoOfRooms(rooms, offerType string) string {
	no, err := strconv.Atoi(rooms)
	switch {
	case err != nil, no < 1:
		return ""
	case offerType == "mieszkanie":
		no += 1
		if no > 8 {
			return "8"
		}
		return strconv.Itoa(no)
	case offerType == "dom":
		if no > 11 {
			return "11"
		}
	}
	return rooms
}

func unifyGroup(group string) string {

	var (
		prefixes = map[string]string{
			"dom":      "dom",
			"mieszka":  "mieszkanie",
			"buiro":    "biuro",
			"buira":    "biuro",
			"dzia":     "dzialka",
			"gara":     "garaz",
			"lokal":    "lokal",
			"komercyj": "lokal",
			"użytkowe": "lokal",
			"hal":      "lokal",
			"hotel":    "lokal",
			"centru":   "lokal",
			"magaz":    "lokal",
			"zakład":   "lokal",
			"pawilon":  "lokal"}
	)
	group = strings.ToLower(group)
	for k, v := range prefixes {
		if strings.HasPrefix(group, k) {
			return v
		}
	}
	return group
}

func formatFloat(number string) string {
	f, err := strconv.ParseFloat(number, 64)
	switch {
	case err != nil, f == 0.0:
		return ""
	default:
		return fmt.Sprintf("%.2f", f)
	}
}

func formatPrice(number string) string {
	f, err := strconv.ParseFloat(number, 64)
	switch {
	case err != nil, f == 0.0:
		return ""
	default:
		return fmt.Sprintf("%.2f", f/100)
	}
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func SendToFtp(companyDomain, zipPath string, ftp *utils.FTP) (err error) {
	remoteDirName := fmt.Sprintf("%s_%s", companyDomain, time.Now().Format(FtpFolderDateFormat))
	_, fileName := filepath.Split(zipPath)
	if err = ftp.SendFile(remoteDirName, fileName, zipPath); err != nil {
		return
	}
	return
}
