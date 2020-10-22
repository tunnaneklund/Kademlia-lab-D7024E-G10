package main

import (
	"fmt"
	"log"
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
	}

	// startar app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
