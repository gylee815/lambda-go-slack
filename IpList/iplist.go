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
	// ip_list := `
	// 	{
	// 		"3.36.166.134": "aws1",
	// 		"13.125.90.18": "aws2",
	// 		"13.125.96.236": "aws3",
	// 		"27.122.242.65": "lge-magok"
	// 	}
	// `

	var tmp map[string]interface{}
	json.Unmarshal([]byte(ip_list), &tmp)

	tmpValue, ok := tmp[ip].(string)
	if !ok {
		fmt.Printf("dateValue is not a string\n")
		return "fail"
	}

	return tmpValue
}
