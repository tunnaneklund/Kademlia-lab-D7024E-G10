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
			// check if user inputed hash
			if len(commands) > 1 {
				hash := commands[1]
				correct := true
				if !correctHash(hash) {
					fmt.Println("input expects only values 0-9a-fA-F")
					correct = false
				}
				if len(hash) != 40 {
					fmt.Println("expects hash of length 40")
					correct = false
				}
				if correct {
					req := d7024e.RequestMessage{}
					req.Type = "get"
					req.Data = hash
					res, _ := sendTCPRequest(req, "localhost:8080")
					fmt.Printf("status: %v\ndata: %v\nfrom: %v %v\n", res.Status, res.Data, res.SenderID, res.SenderIP)
					break
				}

			} else {
				fmt.Println("Command need hash input")
				break
			}
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

func correctHash(str string) bool {
	lower := strings.ToLower(str)
	for _, c := range lower {
		if !((c >= 97 && c <= 102) || (c >= 48 && c <= 57)) {
			return false
		}
	}
	return true
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
