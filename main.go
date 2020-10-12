package main

import (
	"d7024e"
	"fmt"
	"os"
)

func main() {

	// used for testing pings

	port := ":" + os.Args[1]

	fmt.Println(port)
	network := d7024e.NewNetwork("localhost" + port)

	if len(os.Args) > 2 {
		address := "localhost:" + os.Args[2]
		c := d7024e.NewContact(d7024e.NewRandomKademliaID(), address)
		fmt.Println(network.SendPingMessage(&c))

	}

	network.Listen(port)

}
