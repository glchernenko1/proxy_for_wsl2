package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
)

const TRADE = 2

type Ports struct {
	Ports []Port `json:"ports"`
}

type Port struct {
	ListenPort  string `json:"listenport"`
	ConnectPort string `json:"connectport"`
}

func Rider(filePath string) (*Ports, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("File \"Ports.json\" could not be opened \n")
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var ports Ports
	err = json.Unmarshal(byteValue, &ports)
	if err != nil {
		return nil, fmt.Errorf("invalid json")
	}
	return &ports, nil
}

func deleteOldPort(ports chan Port, wg *sync.WaitGroup) {
	for p := range ports {
		cmd := exec.Command("netsh",
			"interface",
			"portproxy",
			"delete",
			"v4tov4",
			fmt.Sprintf("listenport=%s", p.ListenPort),
			"listenaddress=0.0.0.0",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		wg.Done()
	}

}

func ScanIP() string {
	cmd := exec.Command("wsl",
		"hostname",
		"-I",
	)
	bufIP := new(bytes.Buffer)
	cmd.Stdout = bufIP
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	ans := bufIP.String()
	return strings.Split(ans, " ")[0]
}

func addPort(id string, ports chan Port, wg *sync.WaitGroup) {
	for p := range ports {
		cmd := exec.Command("netsh",
			"interface",
			"portproxy",
			"add",
			"v4tov4",
			"listenaddress=0.0.0.0",
			fmt.Sprintf("listenport=%s", p.ListenPort),
			fmt.Sprintf("connectaddress=%s", id),
			fmt.Sprintf("connectport=%s", p.ConnectPort),
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		wg.Done()
	}
}

func (p *Ports) AddPorts() {
	id := ScanIP()
	portChan := make(chan Port, TRADE)
	var wg sync.WaitGroup
	for i := 0; i < cap(portChan); i++ {
		go addPort(id, portChan, &wg)
	}
	for _, port := range p.Ports {
		portChan <- port
		wg.Add(1)
	}
	wg.Wait()
	close(portChan)

}

func (p *Ports) DeleteOldPorts() {
	portChan := make(chan Port, TRADE)
	var wg sync.WaitGroup
	for i := 0; i < cap(portChan); i++ {
		go deleteOldPort(portChan, &wg)
	}
	for _, port := range p.Ports {
		portChan <- port
		wg.Add(1)
	}
	wg.Wait()
	close(portChan)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Path not correctly")
	}
	ports, err := Rider(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	ports.DeleteOldPorts()
	ports.AddPorts()
}
