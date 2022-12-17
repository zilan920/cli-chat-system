package main

import (
	"bufio"
	"fmt"
	"github.com/zilan920/cli-chat-system/pkg/chat"
	"os"
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
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		chatSvc := chat.NewService()
		for {
			fmt.Print(">>-")
			text, err := reader.ReadString('\n')
			fmt.Print("  <<-")
			if err != nil {
				fmt.Println("  something went error", err)
				continue
			}
			funcName, args := getCmd(strings.TrimSuffix(text, "\n"))
			if funcName != "" {
				err, output := chatSvc.CallCmd(funcName, args)
				if err != nil {
					fmt.Println("  error: ", err.Error())
					continue
				}
				fmt.Println(output)
			}
		}
		wg.Done()
	}()
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
