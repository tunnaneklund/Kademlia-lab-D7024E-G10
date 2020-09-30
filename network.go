package d7024e

import (
	"encoding/json"
	"net"
)

const ALPHA = 3

// network
type Network struct {
	rt        *RoutingTable
	dataStore map[string][]byte
}

// *.Type == "ping" | "findcontact" | "finddata" | "storedata"

type requestMessage struct {
	Type   string
	Sender Contact
	Target Contact // findcontact
	Hash   string  // finddata
	Data   string  // storedata
}

type responseMessage struct {
	Type     string
	Sender   Contact
	Status   string    // "ok" | "fail"
	Data     string    // finddata
	Contacts []Contact // findcontact
}

// NewNetwork constructor
func NewNetwork(ip string) Network {
	n := Network{}
	n.rt = NewRoutingTable(NewContact(NewKademliaID(ip), ip))
	return n
}

func (network *Network) Listen() {

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
		go network.handleConnection(conn)
	}
}

func (network *Network) handleConnection(conn net.Conn) {

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	req := requestMessage{}

	res := responseMessage{}
	res.Sender = network.rt.me

	dec.Decode(&req)

	switch req.Type {
	case "ping":

		res.Type = "ping"
		res.Status = "ok"

		enc.Encode(res)

	case "findcontact":

		target := req.Target
		contacts := network.findClosestLocalContacts(target)
		res.Type = "findcontact"
		res.Status = "ok"
		res.Contacts = contacts
		enc.Encode(res)

	case "finddata":

		data := network.getLocalData(req.Hash)

		res.Type = "finddata"
		if data != nil {
			res.Data = string(data)
			res.Status = "ok"
		} else {
			res.Status = "fail"
		}

		enc.Encode(res)

	case "storedata":

		network.storeLocalData([]byte(req.Data))

		res.Type = "storedata"
		res.Status = "ok"

		enc.Encode(res)
	}
}

func (network *Network) SendPingMessage(contact *Contact) { // TODO: dial
	req := requestMessage{}
	req.Sender = network.rt.me
	req.Type = "ping"

	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		return
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	enc.Encode(req)

	res := responseMessage{}

	dec.Decode(&res)

	// TODO: What to do on fail/succes?
	// TODO: Detect fail? timeout?
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	req := requestMessage{}
	req.Sender = network.rt.me
	req.Type = "findcontact"
	req.Target = *contact

	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		return
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	enc.Encode(req)

	res := responseMessage{}

	dec.Decode(&res)

	// TODO: What to do with response? return contact list?
}

func (network *Network) SendFindDataMessage(hash string) {
	// Find closest n contacts, then request data

}

func (network *Network) SendStoreMessage(data []byte) {
	// Find clostest n contacts, then request to store data
}

func (network *Network) getLocalData(hash string) []byte {
	return network.dataStore[hash]
}

func (network *Network) storeLocalData(data []byte) {
	network.dataStore[NewKademliaID(string(data)).String()] = data
}

func (network *Network) findClosestLocalContacts(target Contact) []Contact {
	return network.rt.FindClosestContacts(target.ID, ALPHA)
}
