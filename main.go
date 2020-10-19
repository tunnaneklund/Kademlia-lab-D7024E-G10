package main

import (
	"d7024e"
	"fmt"
	"os"
)

func main() {

	// used for testing pings

	// without docker
	//port := os.Args[1]
	//ip := "localhost"

	// with docker
	port := ":8080"
	ip := string(os.Args[1])

	address := ip + port
	fmt.Println(address)

	network := d7024e.NewNetwork(address)

	if len(os.Args) > 2 {
		// without docker
		//otherAddress := localhost + os.Args[2]

		// with docker
		otherAddress := os.Args[2] + port
		network.JoinNetwork(otherAddress)

	}

	network.Listen(port)

}
