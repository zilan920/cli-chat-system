package main

import (
	"bufio"
	"fmt"
	"github.com/zilan920/cli-chat-system/pkg/chat"
	"net"
	"strings"
	"sync"
)

var allowedCmd = map[string]int{
	"login":     1,
	"send":      1,
	"read":      1,
	"reply":     1,
	"forward":   1,
	"broadcast": 1,
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:2022")
	if err != nil {
		fmt.Println("open connection error", err)
		return
	}
	defer func(listener net.Listener) {
		listener.Close()
	}(listener)
	chatSvc := chat.NewService()
	wg := sync.WaitGroup{}
	for {
		wg.Add(1)
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listen error", err)
			return
		}
		go handleConn(conn, chatSvc, &wg)
	}
	wg.Wait()
}

func getCmd(input string) (funcName string, args []string) {
	res := strings.Split(input, " ")
	if len(res) > 0 {
		funcName = res[0]
		args = res[1:]
		if _, ok := allowedCmd[funcName]; ok {
			return funcName, args
		}
	}
	return
}

func handleConn(conn net.Conn, chatSvc chat.Service, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		reader := bufio.NewReader(conn)
		fmt.Println("service started")
		text, err1 := reader.ReadString('\n')
		if err1 != nil {
			fmt.Println("read error", err1)
			continue
		}
		fmt.Println(text)
		funcName, args := getCmd(strings.TrimSuffix(text, "\n"))
		if funcName == "exit" {
			conn.Write([]byte("bye ~\n"))
			wg.Done()
			break
		} else if funcName != "" {
			err, output := chatSvc.CallCmd(funcName, args)
			if err != nil {
				fmt.Println("  error: ", err.Error())
				continue
			}
			conn.Write([]byte(output + "\n"))
			fmt.Println(output)
		}
	}
}
