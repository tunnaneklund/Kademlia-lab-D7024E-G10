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
				fmt.Println("testprint från CLI")
				return nil
			},
		},
		{
			Name:  "test",
			Usage: "test tcp",
			// Action är vad som händer när vi kör name kommandot
			Action: func(c *cli.Context) error {
				// print för att testa
				req := d7024e.RequestMessage{}
				req.Type = "clitest"
				res, _ := sendTCPRequest(req, "localhost:8080")
				fmt.Println(res.Status)
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
