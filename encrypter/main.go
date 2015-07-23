package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"github.com/kpawlik/exportms/utils"
	"os"
)

func main() {
	var (
		err         error
		data        map[string]map[string]string
		passwd      string
		encFileName string = "enc"
		bytes       []byte
		decrypted   string
		encrypted   []byte
		jsonStr     string
	)
	if bytes, err = ioutil.ReadFile(`data.json`); err != nil {
		log.Panicln(err)
	}
	jsonStr = string(bytes)
	if err = json.Unmarshal([]byte(jsonStr), &data); err != nil {
		log.Panicf("JSON unmarshal %v\n", err)

	}
	fmt.Print("Podaj hasło: ")
	fmt.Scanln(&passwd)
	if _, err = os.Stat(encFileName); err == nil {
		if bytes, err = ioutil.ReadFile(encFileName); err != nil {
			log.Panicf("Read encrypted file %v\n", err)
		}
		if decrypted, err = utils.Decrypt(bytes, passwd); err != nil {
			log.Panicf("Decrypt data %v\n", err)
		}
		fmt.Printf("Data: \n%s\n", decrypted)
	} else {
		if encrypted, err = utils.Encrypt(jsonStr, passwd); err != nil {
			log.Panicf("Encrypt data %v\n", err)
		}
		if err = ioutil.WriteFile(encFileName, encrypted, os.ModePerm); err != nil {
			log.Panicf("Save encrypted %v\n", err)
		}
	}

}
