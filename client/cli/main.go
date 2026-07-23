package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/RAN-GAN/rendis/client/golang"
)

func main() {
	host := flag.String("h", "127.0.0.1", "Server hostname")
	port := flag.String("p", "8080", "Server port")
	key := flag.String("a", "", "Server key (password)")
	urlFlag := flag.String("u", "", "Full server URL (e.g. wss://app.onrender.com)")
	flag.Parse()

	var url string
	if *urlFlag != "" {
		url = *urlFlag
	} else {
		url = fmt.Sprintf("ws://%s:%s", *host, *port)
	}

	client, err := rendis.New(url, *key)
	if err != nil {
		log.Fatalf("Could not connect to rendis server at %s: %v", url, err)
	}
	defer client.Close()

	promptStr := fmt.Sprintf("%s:%s> ", *host, *port)
	if *urlFlag != "" {
		promptStr = fmt.Sprintf("%s> ", *urlFlag)
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:          promptStr,
		HistoryFile:     "/tmp/readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		log.Fatalf("Error initializing readline: %v", err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // EOF or Ctrl+C
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.ToLower(line) == "exit" || strings.ToLower(line) == "quit" {
			break
		}

		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		resp, err := client.Do(parts...)
		if err != nil {
			fmt.Printf("(error) %v\n", err)
			continue
		}

		printResponse(resp)
	}
}

func printResponse(resp any) {
	switch v := resp.(type) {
	case nil:
		fmt.Println("(nil)")
	case string:
		fmt.Printf("%q\n", v)
	case int64:
		fmt.Printf("(integer) %d\n", v)
	case error:
		fmt.Printf("(error) %v\n", v)
	default:
		fmt.Printf("%v\n", v)
	}
}
