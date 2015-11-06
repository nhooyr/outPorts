package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type portRange struct {
	min uint16
	max uint16
}

const (
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	NORMAL = "\033[0m"
	BOLD   = "\033[1m"
)

var (
	successMsg, failureMsg     string
	printSuccess, printFailure = true, true
	wg                         = sync.WaitGroup{}
	d                          = net.Dialer{Timeout: time.Second * 3}
	out                        = make(chan string)
	exit                       = make(chan struct{})
)

func main() {
	log.SetPrefix("outPorts: ")
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Stderr.WriteString(`Examples
  outPorts
        check from port 1 to 65535
  outPorts 20-30 40-50
        check from port 20 to 30 and then 40-50
  outPorts 20-10 40-10
        check from port 20 to 10 and then 40 to 10
  outPorts 25
        check port 25
  outPorts 20-25f
        check from port 20-25 and only display failure (carries onto next port(s))
  outPorts 20-25s
        check from port 20-25 and only display success (carries onto next port(s))
  outPorts 20-25b
        check from port 20-25 and display both (needed to reset)
`)
	}
	color := flag.Bool("c", false, "add color/bold for success/failure")
	flag.Parse()
	if *color == true {
		successMsg = GREEN + BOLD + "%s" + NORMAL
		failureMsg = RED + BOLD + "%s" + NORMAL
	} else {
		successMsg = "%s"
		failureMsg = "%s"
	}
	var (
		min, max uint16
		w        int
	)
	go printLoop()
	for _, arg := range os.Args[1:] {
		if arg == "-c" {
			continue
		}
		switch arg[len(arg)-1] {
		case 's':
			printFailure = false
			printSuccess = true
			arg = arg[0 : len(arg)-1]
		case 'f':
			printSuccess = false
			printFailure = true
			arg = arg[0 : len(arg)-1]
		case 'b':
			printSuccess = true
			printFailure = true
			arg = arg[0 : len(arg)-1]
		}
		if arg == "all" {
			min = 1
			max = 65355
		} else if i := strings.Index(arg, "-"); i != -1 {
			if tmp, err := strconv.ParseUint(arg[:i], 10, 16); err != nil {
				log.Fatal(err)
			} else {
				min = uint16(tmp)
			}
			if tmp, err := strconv.ParseUint(arg[i+1:], 10, 16); err != nil {
				log.Fatal(err)
			} else {
				max = uint16(tmp)
				if min > max {
					tmp := max
					max = min
					min = tmp
				}
			}
		} else {
			if tmp, err := strconv.ParseUint(arg, 10, 16); err != nil {
				log.Fatal(err)
			} else {
				min = uint16(tmp)
				max = min
			}
		}
		if w = int((max - min + 1)); w > 1024 {
			w = 1024
		}
		wg.Add(w)
		in := make([]chan uint16, w)
		for i := 0; i < w; i++ {
			in[i] = make(chan uint16, w)
			go worker(in[i])
		}
		for i := 0; min <= max; min++ {
			in[i] <- min
			if min == 65535 {
				break
			}
			i++
			if i == w {
				i = 0
			}
		}
		for i := 0; i < w; i++ {
			close(in[i])
		}
		wg.Wait()
	}
	close(out)
	<-exit
}

// keeps output from writing over each other; actually happens when its outputting so fast
func printLoop() {
	for m := range out {
		fmt.Print(m)
	}
	exit <- struct{}{}
}

func worker(in <-chan uint16) {
	defer wg.Done()
	for port := range in {
		addr := fmt.Sprintf("portquiz.net:%d", port)
		c, err := d.Dial("tcp", addr)
		if err != nil {
			if printFailure {
				out <- fmt.Sprintf(failureMsg+" on port %d\n", "failure", port)
			}
		} else {
			if printSuccess {
				c.Close()
				out <- fmt.Sprintf(successMsg+" on port %d\n", "success", port)
			}
		}
	}
}
