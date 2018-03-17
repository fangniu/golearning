package main

import (
	"encoding/json"
	"fmt"
)



func main() {
	var services []string
	contents := `["aa", "bb", "cc"]`
	err := json.Unmarshal([]byte(contents), &services)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(services)
}
