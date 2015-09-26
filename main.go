package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

func main() {
	log.SetPrefix("outPorts: ")
	var min, max uint16
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "No port range specified, using range 1-65535, \"outPorts -h\" for more info")
		min = 1
		max = 65535
	} else if strings.ContainsRune(os.Args[1], 'h') {
		fmt.Fprintln(os.Stderr, `Usage of outPorts:
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
	// hide cursor
	os.Stdin.Write([]byte{27, 91, 63, 50, 53, 108})
	/// save current termios
	var old syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), ioctlReadTermios, uintptr(unsafe.Pointer(&old)), 0, 0, 0); err != 0 {
		log.Fatalln("not a terminal, got:", err)
	}
	// capture signals
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	cleanup := func() {
		// make cursor visible
		os.Stdin.Write([]byte{27, 91, 51, 52, 104, 27, 91, 63, 50, 53, 104})
		// set tty to normal
		if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), ioctlWriteTermios, uintptr(unsafe.Pointer(&old)), 0, 0, 0); err != 0 {
			log.Fatal(err)
		}
	}
	go func() {
		<-sigs
		// restore text to normal
		os.Stdout.Write([]byte{27, 91, 48, 109})
		cleanup()
		os.Exit(0)
	}()
	// set raw mode
	raw := old
	raw.Lflag &^= syscall.ECHO | syscall.ICANON
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, os.Stdin.Fd(), ioctlWriteTermios, uintptr(unsafe.Pointer(&raw)), 0, 0, 0); err != 0 {
		log.Fatal(err)
	}
	d := net.Dialer{Timeout: time.Second * 3}
	wg := sync.WaitGroup{}
	check := func(port uint16) {
		defer wg.Done()
		addr := fmt.Sprintf("portquiz.net:%d", port)
		c, err := d.Dial("tcp", addr)
		if err != nil {
			fmt.Printf("\033[31m\033[01m%s\033[00m on port %d\n", "failure", port)
		} else {
			c.Close()
			fmt.Printf("\033[32m\033[01m%s\033[00m on port %d\n", "success", port)
		}
	}
	for ; min <= max; min++ {
		wg.Add(1)
		go check(min)
		if min == 65535 {
			break
		}
		time.Sleep(time.Millisecond * 3)
	}
	wg.Wait()
	cleanup()
}
