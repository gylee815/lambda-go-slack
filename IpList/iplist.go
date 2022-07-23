package Iplist

import (
	"fmt"
	"encoding/json"
)

func GetSite(ip string) string {
	ip_list := `
		{
			"3.36.166.134": "aws1",
			"13.125.90.18": "aws2",
			"13.125.96.236": "aws3"
		}
	`

	var tmp map[string]interface{}
	json.Unmarshal([]byte(ip_list), &tmp)
	
	tmpValue, ok := tmp[ip].(string)
	if !ok {
		fmt.Printf("dateValue is not a string\n")
		return "fail"
	}

	return tmpValue
}