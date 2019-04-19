package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Tv struct {
	Process *exec.Cmd
}

func (t *Tv) isOpen() bool {
	if t.Process == nil {
		return false
	}
	if t.Process.ProcessState == nil {
		// ProcessState contains information about an exited process,
		// available after a call to Wait or Run.
		return true
	}
	return !t.Process.ProcessState.Exited()
}

func (t *Tv) open() error {
	t.Process = exec.Command("/usr/bin/kodi")
	t.Process.Stderr = os.Stderr
	t.Process.Stdout = os.Stdout
	err := t.Process.Start()
	go func() {
		err := t.Process.Wait()
		if err != nil {
			log.Printf("tv daemon exit with error %v", err)
		} else {
			log.Printf("tv daemon exit")
		}
	}()
	return err
}

type tvRpc struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int64       `json:"id"`
}

func (t *Tv) close() error {
	var rpc tvRpc
	rpc.Jsonrpc = "2.0"
	rpc.Method = "Application.Quit"
	rpc.Params = map[string]string{}
	rpc.ID = time.Now().Unix()
	data, err := json.Marshal(rpc)
	if err != nil {
		return err
	}
	_, err = httppost("http://127.0.0.1:8080/jsonrpc", data)
	return err
}

func (t *Tv) pause() error {
	var rpc tvRpc
	rpc.Jsonrpc = "2.0"
	rpc.Method = "Input.ExecuteAction"
	rpc.Params = map[string]string{
		"action": "pause",
	}
	rpc.ID = time.Now().Unix()
	data, err := json.Marshal(rpc)
	if err != nil {
		return err
	}
	_, err = httppost("http://127.0.0.1:8080/jsonrpc", data)
	return err
}

type Screen struct {
}

func (s *Screen) isOpen() bool {
	output, err := exec.Command("/usr/bin/tvservice", "-s").Output()
	if err != nil {
		log.Println(err)
		return false
	}
	outstr := string(output)
	return strings.Contains(outstr, "0x12000a")
}

func (s *Screen) open() error {
	tvservice := exec.Command("/usr/bin/tvservice", "-p")
	fbset16 := exec.Command("/bin/fbset", "-depth", "16")
	fbset32 := exec.Command("/bin/fbset", "-depth", "32")
	xrefresh := exec.Command("/usr/bin/xrefresh", "-display", ":0")
	err := tvservice.Start()
	go func() {
		time.Sleep(3)
		fbset16.Start()
		time.Sleep(1)
		fbset32.Start()
		time.Sleep(1)
		xrefresh.Start()
		fbset16.Wait()
		fbset32.Wait()
		xrefresh.Wait()
	}()
	tvservice.Wait()
	return err
}

func (s *Screen) close() error {
	tvservice := exec.Command("/usr/bin/tvservice", "-o")
	err := tvservice.Start()
	tvservice.Wait()
	return err
}
