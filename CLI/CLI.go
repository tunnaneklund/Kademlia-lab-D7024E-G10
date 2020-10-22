package main

import (
	"d7024e"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Kademlia lab CLI"
	app.Usage = "Ping, put, get, find contact id, print rt & printbuckets"

	// we create our commands
	app.Commands = []cli.Command{
		{
			Name:  "testprint",
			Usage: "den testar en print",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				// print för att testa
				fmt.Printf("this prints: %v", c.Args().Get(0))
				return nil
			},
		},
		{
			Name:  "put",
			Usage: "Takes a single argument, the contents of the file you are uploading, and outputs thehash of the object, if it could be uploaded successfully.",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				data := c.Args().Get(0)
				req := d7024e.RequestMessage{}
				req.Type = "put"
				req.Data = data
				res, _ := sendTCPRequest(req, "localhost:8080")
				fmt.Printf("status: %v\nhash: %v\n", res.Status, res.Data)
				return nil
			},
		},
		{
			Name:  "get",
			Usage: "Takes a hash as its only argument, and outputs the contents of the object and thenode it was retrieved from, if it could be downloaded successfully.",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				hash := c.Args().Get(0)
				req := d7024e.RequestMessage{}
				req.Type = "get"
				req.Data = hash
				res, _ := sendTCPRequest(req, "localhost:8080")
				fmt.Printf("status: %v\ndata: %v\n", res.Status, res.Data)
				return nil
			},
		},
		{
			Name:  "exit",
			Usage: "close node",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				// print för att testa
				req := d7024e.RequestMessage{}
				req.Type = "exit"
				res, _ := sendTCPRequest(req, "localhost:8080")
				fmt.Println(res.Status)
				return nil
			},
		},
		{
			Name:  "printrt",
			Usage: "print routing table",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				// print för att testa
				req := d7024e.RequestMessage{}
				req.Type = "printrt"
				res, _ := sendTCPRequest(req, "localhost:8080")
				fmt.Printf("status: %v\ndata: %v", res.Status, res.Data)
				return nil
			},
		},
		{
			Name:  "printds",
			Usage: "print data store",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				// print för att testa
				req := d7024e.RequestMessage{}
				req.Type = "printds"
				res, _ := sendTCPRequest(req, "localhost:8080")
				fmt.Printf("status: %v\ndata: %v", res.Status, res.Data)
				return nil
			},
		},
	}

	// startar app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
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
