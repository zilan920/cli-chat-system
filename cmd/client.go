package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:2022")
	defer func() {
		conn.Close()
	}()
	if err != nil {
		fmt.Println("open connection error", err)
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		readerCmd := bufio.NewReader(os.Stdin)
		for {
			text, err := readerCmd.ReadString('\n')
			if err != nil {
				fmt.Println("read error", err)
			}
			_, _ = conn.Write([]byte(text))
		}
		wg.Done()
	}()
	go func() {
		reader := bufio.NewReader(conn)
		for {
			output, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("read error", err)
			}
			fmt.Println(output)
		}
		wg.Done()
	}()
	wg.Wait()
}
