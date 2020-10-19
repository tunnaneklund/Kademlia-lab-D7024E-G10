package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const ALPHA = 3
const K = 10

// network
type Network struct {
	rt        *RoutingTable
	dataStore map[string]string
	cc        chan []Contact
}

// *.Type == "ping" | "findcontact" | "finddata" | "storedata"

type requestMessage struct {
	Type   string
	Sender Contact
	Target KademliaID // findcontact
	Hash   string     // finddata
	Data   string     // storedata
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
	n.cc = make(chan []Contact, ALPHA*2)
	n.rt = NewRoutingTable(NewContact(NewRandomKademliaID(ip), ip))
	return n
}

func (network *Network) Listen(port string) {

	// Port only as argument for local testing
	// Ping: respond
	// Find contact: lookup contact, otherwise: ask closer contacts about contact, and respond with it.
	// Find data message: local lookup, otherwise send to closer contacts
	// StoreMessage: store locally.

	ln, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error1")
		fmt.Println(err)
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
		network.rt.AddContact(req.Sender)

	case "findcontact":

		target := req.Target
		contacts := network.findClosestLocalContacts(target)
		res.Type = "findcontact"
		res.Status = "ok"
		res.Contacts = contacts
		enc.Encode(res)
		network.rt.AddContact(req.Sender)
		network.PrintClosestContacts()

	case "finddata":

		data := network.getLocalData(req.Hash)

		res.Type = "finddata"
		if data != "" {
			res.Data = data
			res.Status = "ok"
		} else {
			res.Status = "fail"
		}

		enc.Encode(res)

	case "storedata":

		network.storeLocalData(string(req.Data))

		res.Type = "storedata"
		res.Status = "ok"

		enc.Encode(res)
	}
}

func (network Network) SendPingMessageIP(address string) string {
	id := NewRandomKademliaID("dummy id")
	c := NewContact(id, address)
	return network.SendPingMessage(&c)

}

// SendPingMessage returns the contact that was pinged, can be used to obtain full contact from just ip
func (network Network) SendPingMessage(contact *Contact) string {
	req := requestMessage{}
	req.Sender = network.rt.me
	req.Type = "ping"

	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		return "fail"
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	enc.Encode(req)

	res := responseMessage{}

	dec.Decode(&res)

	network.rt.AddContact(res.Sender)

	return res.Status
}

// first contact in list posted to cc is "contact"
func (network Network) SendFindContactMessage(contact Contact, target KademliaID) {
	req := requestMessage{}
	req.Sender = network.rt.me
	req.Type = "findcontact"
	req.Target = target

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

	s := make([]Contact, 1)
	s[0] = contact
	network.cc <- append(s, res.Contacts...)

}

func (network *Network) SendFindDataMessage(hash string, contact Contact) string {
	req := requestMessage{}
	req.Sender = network.rt.me
	req.Type = "finddata"
	req.Hash = hash

	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		return ""
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	enc.Encode(req)

	res := responseMessage{}

	dec.Decode(&res)

	return res.Data

}

func (network *Network) SendStoreMessage(data string, contact Contact) {
	req := requestMessage{}
	req.Sender = network.rt.me
	req.Type = "storedata"
	req.Data = data

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

}

func (network *Network) getLocalData(hash string) string {
	return network.dataStore[hash]
}

func (network *Network) storeLocalData(data string) { // needs fix NewKademliaID not working as intended
	network.dataStore[NewRandomKademliaID(data).String()] = data
}

func (network *Network) storeData(data string) {
	id := NewRandomKademliaID(data)
	contacts := network.ContactLookup(*id)
	for _, c := range contacts {
		go network.SendStoreMessage(data, c)
	}
	network.storeLocalData(data)
}

func (network *Network) findClosestLocalContacts(target KademliaID) []Contact {
	return network.rt.FindClosestContacts(&target, K)
}

type shortlistStatus struct {
	queried   bool
	responded bool
	inList    bool
}

func newStatus() shortlistStatus {
	s := shortlistStatus{}
	s.queried = false
	s.responded = false
	s.inList = true
	return s
}

// ContactLookup return a list of the K closest contacts to ID of target
func (network *Network) ContactLookup(target KademliaID) []Contact {
	// used to keep track of contacts
	m := make(map[KademliaID]shortlistStatus)
	// using ContactCandidate to get access to sorting
	shortlist := ContactCandidates{}
	shortlist.Append(network.findClosestLocalContacts(target))
	updateContacts(&shortlist.contacts, &m)
	closestNode := shortlist.contacts[0]
	exitOnNext := false
	for {

		// removes unresponsive contacts and sends rpc find to ALPHA nr of contacts
		counter := ALPHA
		for i, c := range shortlist.contacts {

			if counter < 1 {
				break
			}

			stat := m[*c.ID]
			if stat.queried {
				if !stat.responded {
					m[*c.ID] = shortlistStatus{}
					shortlist.Delete(i)
				}
			} else {
				counter--
				stat.queried = true
				m[*c.ID] = stat
				go network.SendFindContactMessage(c, target)
			}
		}

		// adds new contacts from channel to shortlist, only blocks for 2 sec
		loopVar := true
		for loopVar {
			select {
			case cl := <-network.cc:
				{
					temp := m[*cl[0].ID]
					temp.responded = true
					m[*cl[0].ID] = temp
					network.rt.AddContact(cl[0])
					// removes contacts aleady in shortlist
					for i := len(cl) - 1; i >= 0; i-- {
						if m[*cl[i].ID].inList || (*cl[i].ID).Equals(network.rt.me.ID) {
							cl = append(cl[:i], cl[i+1:]...)
						}
					}

					updateContacts(&cl, &m)
					shortlist.Append(cl)
				}
			case <-time.After(2 * time.Second):
				loopVar = false
			}
		}

		// checks if progress is being made, otherwise return K closest
		shortlist.Sort()

		if shortlist.contacts[0] == closestNode && exitOnNext {
			if shortlist.Len() < K {
				return shortlist.GetContacts(shortlist.Len())
			}
			return shortlist.GetContacts(K)
		} else if shortlist.contacts[0] == closestNode {
			exitOnNext = true
		} else {
			exitOnNext = false
			closestNode = shortlist.contacts[0]
		}
	}

}

// PrintClosestContacts prints contacts
func (network *Network) PrintClosestContacts() {
	fmt.Println(network.findClosestLocalContacts(*network.rt.me.ID))
}

// sets up information on each contact
func updateContacts(cl *[]Contact, m *map[KademliaID]shortlistStatus) {
	for _, c := range *cl {
		(*m)[*c.ID] = newStatus()
	}
}

func (network *Network) JoinNetwork(contactAddress string) {

	// add the initial known node to k-bucket
	network.SendPingMessageIP(contactAddress)

	// perform node lookup on itself
	network.ContactLookup(*network.rt.me.ID)

	network.PrintClosestContacts()
}
