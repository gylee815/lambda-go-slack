package Iplist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func fileRead(fileName string) []byte {
	absPath, _ := filepath.Abs(fileName)
	dat, err := ioutil.ReadFile(absPath)
	check(err)

	return dat
}

func GetSite(ip string) string {
	ip_list := fileRead("iplist.json")

	var tmp map[string]interface{}
	json.Unmarshal([]byte(ip_list), &tmp)

	siteName, ok := tmp[ip].(string)
	if !ok {
		fmt.Printf("siteName is not a string\n")
		return "fail"
	}

	return siteName
}
