package main

import (
	"bufio"
	"d7024e"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("kademlia>")
		b, _, _ := reader.ReadLine()
		command := string(b)
		commands := strings.Fields(command)

		switch commands[0] {
		case "testprint":
			fmt.Printf("this prints: %v", commands[1])
		case "put":
			req := d7024e.RequestMessage{}
			req.Type = "put"
			req.Data = commands[1]
			res, _ := sendTCPRequest(req, "localhost:8080")
			fmt.Printf("status: %v\nhash: %v\n", res.Status, res.Data)
		case "get":
			hash := commands[1]
			req := d7024e.RequestMessage{}
			req.Type = "get"
			req.Data = hash
			res, _ := sendTCPRequest(req, "localhost:8080")
			fmt.Printf("status: %v\ndata: %v\n", res.Status, res.Data)
		case "exit":
			req := d7024e.RequestMessage{}
			req.Type = "exit"
			res, _ := sendTCPRequest(req, "localhost:8080")
			fmt.Println(res.Status)
			os.Exit(0)
		case "printrt":
			req := d7024e.RequestMessage{}
			req.Type = "printrt"
			res, _ := sendTCPRequest(req, "localhost:8080")
			fmt.Printf("status: %v\ndata: %v\n", res.Status, res.Data)
		case "printds":
			req := d7024e.RequestMessage{}
			req.Type = "printds"
			res, _ := sendTCPRequest(req, "localhost:8080")
			fmt.Printf("status: %v\ndata: %v\n", res.Status, res.Data)
		default:
			fmt.Println("Command not recognized")
		}

	}
}

func sendTCPRequest(req d7024e.RequestMessage, address string) (res d7024e.ResponseMessage, err error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return d7024e.ResponseMessage{}, err
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	enc.Encode(req)

	dec.Decode(&res)
	return
}
