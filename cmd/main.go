package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"unstable-channel/client"
	"unstable-channel/proxy"
	"unstable-channel/server"
)

func main() {
	mode := flag.String("mode", "client", "Run mode: 'server', 'client', or 'proxy'")

	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	if *mode == "client" {
		client, err := client.NewClient(8000)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating client: %v\n", err)
		}
		fmt.Printf("local: %s, remote: %s\n", client.Laddr.String(), client.Raddr.String())

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

			// TODO: send input to server
			client.Send(line)
		}

		err = client.Kill()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down client: %v", err)
			os.Exit(1)
		}
	}

	if *mode == "server" {
		server, err := server.NewServer(8080)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating server: %v\n", err)
		}

		// TODO: print received packets from server.Start()
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

	// TODO: implement if *mode == "proxy"
	if *mode == "proxy" {
		proxy, err := proxy.NewProxy(8000)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error creating proxy: %v\n", err)
		}

		// TODO: print received packets from server.Start()
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
