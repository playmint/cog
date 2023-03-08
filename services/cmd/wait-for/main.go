package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type servicesType []string

var timeout int
var services servicesType
var filename string

func (s *servicesType) String() string {
	return fmt.Sprintf("%+v", *s)
}

func (s *servicesType) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

// waitForServices tests and waits on the availability of a TCP host and port
func waitForServices(services []string, timeout time.Duration) error {
	var depChan = make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(services))
	go func() {
		for _, s := range services {
			go func(s string) {
				defer wg.Done()
				for {
					_, err := net.Dial("tcp", s)
					if err == nil {
						return
					}
					time.Sleep(1 * time.Second)
				}
			}(s)
		}
		wg.Wait()
		close(depChan)
	}()

	select {
	case <-depChan: // services are ready
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("services aren't ready in %s", timeout)
	}
}

func waitForFile(filename string, timeout time.Duration) error {
	if filename == "" {
		return nil
	}
	var depChan = make(chan struct{})
	go func() {
		for {
			if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
				time.Sleep(1 * time.Second)
			} else {
				depChan <- struct{}{}
				break
			}
		}
		close(depChan)
	}()

	select {
	case <-depChan: // services are ready
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("file %s did not appear before timeout", filename)
	}
}

func init() {
	flag.IntVar(&timeout, "t", 20, "timeout")
	flag.Var(&services, "it", "<host:port> [host2:port,...] comma seperated list of services")
	flag.StringVar(&filename, "f", "", "<filename> name of a file to wait for existance")
}

func main() {
	flag.Parse()
	if len(services) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if err := waitForServices(services, time.Duration(timeout)*time.Second); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := waitForFile(filename, time.Duration(timeout)*time.Second); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("services are ready!")
	os.Exit(0)
}
