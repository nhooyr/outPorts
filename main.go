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
	var (
		min, max uint16
	)
	log.SetPrefix("outPorts: ")
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "No port range specified, using range 1-65535, \"outPorts -h\" for more info")
		min = 1
		max = 65535
	} else if strings.ContainsRune(os.Args[1], 'h')  {
		fmt.Fprintln(os.Stderr, "Usage of outPorts:")
		fmt.Fprintln(os.Stderr, "  outPorts [-u] min[-max]")
		fmt.Fprintln(os.Stderr, "      -u    use udp instead of tcp")
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
	d := net.Dialer{Timeout: time.Second * 1}
	wg := sync.WaitGroup{}
	check := func(port uint16) {
		defer wg.Done()
		addr := fmt.Sprintf("portquiz.net:%d", port)
		c, err := d.Dial("tcp", addr)
		if err != nil {
			fmt.Println("failure on port", port)
		} else {
			c.Close()
			fmt.Println("success on port", port)
		}
	}
	for ; min <= max; min++ {
		wg.Add(1)
		go check(min)
		time.Sleep(time.Millisecond * 3)
	}
	wg.Wait()
}
