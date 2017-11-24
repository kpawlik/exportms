package utils

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
)

// ColID type for index columns
type ColID int

// Row type for database result
type Row []interface{}

// Errorf create new error
func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// LogErr log eror as panic
func LogErr(err error, msg string) {
	if err != nil {
		log.Panicf("!! Error !!\n%s, %v\n", msg, err)
	}
}

// LogFatalf log fortmated fatal error
func LogFatalf(err error, format string, args ...interface{}) {
	if err != nil {
		log.Fatalf("%s, %v", fmt.Sprintf(format, args...), err)
	}
}

// LogErrf log fortmated fatal error
func LogErrf(err error, format string, args ...interface{}) {
	if err != nil {
		log.Printf("%s, %v", fmt.Sprintf(format, args...), err)
	}
}

//LogErrfInfo log info about error
func LogErrfInfo(err error, format string, args ...interface{}) {
	if err != nil {
		log.Printf("%s, %v", fmt.Sprintf(format, args...), err)
	}
}

func printRow(row Row) string {
	s := ""
	for _, c := range row {
		switch c.(type) {
		case *int:
			s = fmt.Sprintf("%s, %d", s, *c.(*int))
		case *string:
			s = fmt.Sprintf("%s, %s", s, *c.(*string))
		case *sql.NullInt64:
			dbInt := *c.(*sql.NullInt64)
			if dbInt.Valid {
				s = fmt.Sprintf("%s, %d", s, dbInt.Int64)
			} else {
				s = fmt.Sprintf("%s, NULL", s)
			}
		case *sql.NullString:
			dbS := *c.(*sql.NullString)
			if dbS.Valid {
				s = fmt.Sprintf("%s, %s", s, dbS.String)
			} else {
				s = fmt.Sprintf("%s, NULL", s)
			}
		default:
			fmt.Println(reflect.TypeOf(c))
		}
	}
	return s
}

// SaveImageFromBytes store image bytes into file
func SaveImageFromBytes(buff []byte, path string) (err error) {
	out := make([]byte, len(buff), len(buff))
	if _, err = hex.Decode(out, buff); err != nil {
		return
	}
	if err = ioutil.WriteFile(path, out, os.ModePerm); err != nil {
		return
	}
	return
}

//GetString return column value as a string
func GetString(row Row, index ColID) (str string) {
	val := row[index]
	sval := *val.(*sql.NullString)
	return sval.String
}

//GetInt64 return column value as a int64
func GetInt64(row Row, index ColID) (num int64, ok bool) {
	val := row[index]
	sval := *val.(*sql.NullInt64)
	ok, num = sval.Valid, sval.Int64
	return
}

// DecodeStr decodes string
func DecodeStr(lan string) string {
	dlan := "5CFCFF55ADCAC5D5145F8EAEF27402555D6DE7D5C4341B0DA9FA54F5C9953CE1AC936507A88EC807EB69F103A74D57B91E48495DA3F83662AD33C7243A57180BFF786D6C1068C66F4F27DA8656431FCF7C6AA9EEDB66FB16C06A7CEB5F9299B1FE60ABE235BF59C0E30D999089090B18C72387372309C77A4D4128E69FD4176258F093B866142FAB2B6A7693E754ADF04AD878E91A51E5E1DE6978924BAC13CBC4B8813A110BFECA8F8CD42324CA7AA75A0D31B94890A2C31E97FD4DEC7F24A1E022BEBFC25672995BC22DEEF2F4C8876364623C0C115EB3BD702C641AD3FE9431839BC022A71B76DD741FA7F343E6ECF390B9EEBA26975F70805E63434A533B172375CFE3A46EAD6A265BFCA603204BB444C72C9B36E779D32EDBE5F298CE10E658CC9CB9D7C08C707E93887184D93A5825FC4506C804B36BD23562D36FFF70E6853DD942A95E6C7E2D71C09F158F6791BCB0826A809FA297B20B71997559A970377B37FC6D96C841E887027F22E188FD722E42591163BFA762DFBFF52F2D05F10E394A4B6DCA366852439B66328049"
	out := ""

	for i := 0; i < len(lan)/2; i++ {
		start := i * 2
		end := (i * 2) + 2
		ds := dlan[start:end]
		cs := lan[start:end]
		ids := hexdec(ds) - 49
		ics := hexdec(cs)
		out = fmt.Sprintf("%s%c", out, byte(ics-ids))
	}
	return out
}

func hexdec(s string) uint64 {
	d := uint64(0)
	for i := 0; i < len(s); i++ {
		x := uint64(s[i])
		if x >= 'a' {
			x -= 'a' - 'A'
		}
		d1 := x - '0'
		if d1 > 9 {
			d1 = 10 + d1 - ('A' - '0')
		}
		if 0 > d1 || d1 > 15 {
			panic("hexdec")
		}
		d = (16 * d) + d1
	}
	return d
}

// GetLastIDFileName return  last load id
func GetLastIDFileName(folder string) string {
	var (
		content []byte
		err     error
		idPath  string
	)
	idPath = filepath.Join(folder, "last.id")
	if _, err = os.Stat(idPath); err != nil {
		return ""
	}
	if content, err = ioutil.ReadFile(idPath); err != nil {
		return ""
	}
	return filepath.Join(folder, string(content))

}
