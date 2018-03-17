package main

import (
	"fmt"
	"os/exec"
	"time"
	"bytes"
	"strings"
)

func main() {
	cmd := exec.Command("tail", "-f", "config.json")
	go func() {
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Start()
		outStr := strings.TrimSpace(out.String())
		fmt.Println(outStr)
	}()
	time.Sleep(time.Second*1)
	err := cmd.Process.Kill()
	if err != nil {
		fmt.Println(err)
	}
}