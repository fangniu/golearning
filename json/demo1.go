package main

import (
	"encoding/json"
	"fmt"
)

type IRResponse struct {
	Code     int    `json:"code"`
	CodeDesc string `json:"codeDesc"`
	Message  string `json:"message"`
	//data     interface{} `json:"data"`
}




func main() {
	//var services []string
	//contents := `["aa", "bb", "cc"]`
	//err := json.Unmarshal([]byte(contents), &services)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(services)
	irr := IRResponse{}
	s1 := `{"code":0,"codeDesc":"Success","data":[{"category":"pornDetection","confidence":0,"label":"porn","subCode":0,"suggestion":"pass"},{"category":"pornDetection","confidence":17,"label":"hot","subCode":0,"suggestion":"pass"}],"message":"No Error"}`
	err := json.Unmarshal([]byte(s1), &irr)
	fmt.Println(s1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(irr.CodeDesc)
	}
}
