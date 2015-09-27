package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	log.SetPrefix("outPorts: ")
	log.SetFlags(0)
	var min, max uint16
	if len(os.Args) < 2 {
		min = 1
		max = 65535
	} else if strings.ContainsRune(os.Args[1], 'h') {
		fmt.Fprintln(os.Stderr, `Usage of outPorts:
  outPorts
        check from port 1 to 65535
  outPorts 20-30
        check from port 20 to 30
  outPorts 20-10
        check from port 20 to (20+10)
  outPorts 25
        check port 25`)
		return
	} else if i := strings.Index(os.Args[1], "-"); i != -1 {
		if tmp, err := strconv.ParseUint(os.Args[1][:i], 10, 16); err != nil {
			log.Println(err)
			return
		} else {
			min = uint16(tmp)
		}
		if tmp, err := strconv.ParseUint(os.Args[1][i+1:], 10, 16); err != nil {
			log.Println(err)
			return
		} else {
			max = uint16(tmp)
			if min > max {
				max += min
			}
		}
	} else {
		if tmp, err := strconv.ParseUint(os.Args[1], 10, 16); err != nil {
			log.Println(err)
			return
		} else {
			min = uint16(tmp)
			max = min
		}
	}
	d := net.Dialer{Timeout: time.Second * 3}
	wg := sync.WaitGroup{}
	out := make(chan string)
	check := func(port uint16) {
		defer wg.Done()
		addr := fmt.Sprintf("portquiz.net:%d", port)
		c, err := d.Dial("tcp", addr)
		if err != nil {
			out <- fmt.Sprintf("\033[31m\033[1m%s\033[0m on port %d\n", "failure", port)
		} else {
			c.Close()
			out <- fmt.Sprintf("\033[32m\033[1m%s\033[0m on port %d\n", "success", port)
		}
	}
	exit := make(chan struct{})
	go printLoop(out, exit)
	for ; min <= max; min++ {
		wg.Add(1)
		go check(min)
		if min == 65535 {
			break
		}
		time.Sleep(time.Millisecond * 3)
	}
	wg.Wait()
	<-exit
}

// keeps output from writing over each other; actually happens when its outputting so fast
func printLoop(out <-chan string, exit chan<- struct{}) {
	for {
		fmt.Print(<-out)
		select {
		case exit <- struct{}{}:
		default:
			// not exiting
		}
	}
}
