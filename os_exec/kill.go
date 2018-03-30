package main

import (
	"fmt"
	"os/exec"
	"time"
	"bytes"
	"strings"
)

func f1()  {
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

func f2()  {
	go func() {
		time.Sleep(time.Second * 1)
		fmt.Println("bbb")
	}()
	fmt.Println("aaa")
	return
}

func f3()  {
	cmd := exec.Command("sleep", "4")
	if err := cmd.Start(); err != nil {
		fmt.Println(err)
	}

	// Wait for the process to finish or kill it after a timeout:
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(3 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("failed to kill process: ", err)
		}
		fmt.Println("hhhh")
		cmd.Wait()
		fmt.Println("process killed as timeout reached")
	case err := <-done:
		if err != nil {
			fmt.Println("process finished with error = %v", err)
			return
		}
		fmt.Println("process finished successfully")
	}
}

func main() {
	f3()
	time.Sleep(time.Second * 700)
}