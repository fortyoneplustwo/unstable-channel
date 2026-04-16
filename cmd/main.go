package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"unstable-channel/client"
	"unstable-channel/proxy"
	"unstable-channel/server"
)

func main() {
	mode := flag.String("mode", "proxy", "Run mode: 'server', 'client', or 'proxy'")

	flag.Parse()
	args := os.Args[1:]

	reader := bufio.NewReader(os.Stdin)

	if *mode == "client" {
		if len(args) != 4 {
			fmt.Println("Usage: ./cmd -mode client [destination_port] [proxy_port]")
			os.Exit(1)
		}
		destPort, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing arg1: %v", err)
		}
		proxyPort, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing arg2: %v", err)
		}

		client, err := client.New(destPort, proxyPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
		}
		fmt.Printf("Destination port: %s\n", strconv.Itoa(client.Raddr.Port))
		fmt.Printf("Proxy port: %s\n", strconv.Itoa(client.Paddr.Port))

		fmt.Println("Enter text to send and --kill to exit.")

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				break
			}
			if line == "--kill\n" {
				break
			}
			fmt.Printf("You entered: %s", line)

			client.Send(line)
		}

		err = client.Kill()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down client: %v", err)
			os.Exit(1)
		}
	}

	if *mode == "server" {
		if len(args) != 3 {
			fmt.Println("Usage: ./cmd -mode server [local_port]")
			os.Exit(1)
		}
		localPort, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing arg1: %v", err)
		}

		server, err := server.New(localPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating server: %v\n", err)
		}

		go server.Start()

		fmt.Printf("listening on: %s\n", server.Conn.LocalAddr().String())

		fmt.Println("Enter --kill to exit.")

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				break
			}
			if line == "--kill\n" {
				break
			}
			fmt.Printf("You entered: %s", line)
		}

		err = server.Kill()
		if err != nil {
			fmt.Fprintf(os.Stderr, "server cleanup error: %v", err)
			os.Exit(1)
		}
	}

	if *mode == "proxy" {
		if len(args) != 2 && len(args) != 4 {
			fmt.Printf("args len: %d\n", len(args))
			// NOTE: destination port is temporary until we implement packets
			fmt.Println("Usage: ./cmd [local_port] [destination_port]")
			os.Exit(1)
		}

		var arg1 string
		var arg2 string
		switch len(args) {
		case 2:
			arg1 = args[0]
			arg2 = args[1]
		case 4:
			arg1 = args[2]
			arg2 = args[3]
		}

		localPort, err := strconv.Atoi(arg1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing arg1: %v", err)
		}
		destPort, err := strconv.Atoi(arg2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error parsing arg2: %v", err)
		}

		proxy, err := proxy.New(localPort, destPort)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating proxy: %v\n", err)
		}

		go proxy.Start()

		fmt.Printf("listening on: %s\n", proxy.Conn.LocalAddr().String())

		fmt.Println("Enter --kill to exit.")

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
				break
			}
			if line == "--kill\n" {
				break
			}
			fmt.Printf("You entered: %s", line)
		}

		err = proxy.Kill()
		if err != nil {
			fmt.Fprintf(os.Stderr, "proxy cleanup error: %v", err)
			os.Exit(1)
		}
	}

	os.Exit(0)
}
