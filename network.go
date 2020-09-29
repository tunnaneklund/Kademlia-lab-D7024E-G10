package d7024e

import (
	"encoding/json"
	"net"
)

// network
type Network struct {

	kademlia Kademlia

}

func NewNetwork(ip string) Network {
	n := new(Network)
	n.kademlia = NewKademliaNode(ip)
}


func Listen(ip string, port int) {



	// TODO
	// Ping: respond
	// Find contact: lookup contact, otherwise: ask closer contacts about contact, and respond with it.
	// Find data message: local lookup, otherwise send to closer contacts
	// StoreMessage: store locally.

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	var message interface{}

	dec.Decode(&message)

	switch message.Type {
	case "ping" : {
		enc.Encode(message{"pingresponse", "ok"})
	} case "findcontact" : {


	}

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
	// This tells a node to store the data. The node has to be found first with FindContactMessage
}
