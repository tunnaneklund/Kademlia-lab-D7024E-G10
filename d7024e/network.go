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
	dataStore map[string][]byte
	cc        chan []Contact
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
	n.cc = make(chan []Contact, ALPHA*2)
	n.rt = NewRoutingTable(NewContact(NewRandomKademliaID(), ip))
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

func (network Network) SendPingMessage(contact *Contact) string { // TODO: dial
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

	return res.Status
}

// first contact in list posted to cc is "contact"
func (network Network) SendFindContactMessage(contact Contact, target Contact) {
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

func (network *Network) SendFindDataMessage(hash string) {
	// Find closest n contacts, then request data

}

func (network *Network) SendStoreMessage(data []byte) {
	// Find clostest n contacts, then request to store data
}

func (network *Network) getLocalData(hash string) []byte {
	return network.dataStore[hash]
}

func (network *Network) storeLocalData(data []byte) { // needs fix NewKademliaID not working as intended
	network.dataStore[NewKademliaID(string(data)).String()] = data
}

func (network *Network) findClosestLocalContacts(target Contact) []Contact {
	return network.rt.FindClosestContacts(target.ID, K)
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

// ContactLookup return a list of the K closest contacts to an ID
func (network *Network) ContactLookup(target Contact) []Contact {
	// used to keep track of contacts
	m := make(map[*KademliaID]shortlistStatus)
	// using ContactCandidate to get access to sorting
	shortlist := ContactCandidates{}
	shortlist.Append(network.findClosestLocalContacts(target))
	updateContacts(&shortlist.contacts, target, &m)
	closestNode := shortlist.contacts[0]
	exitOnNext := false
	for {

		// removes unresponsive contacts and sends rpc find to ALPHA nr of contacts
		counter := ALPHA
		for i, c := range shortlist.contacts {

			if counter < 1 {
				break
			}

			stat := m[c.ID]
			if stat.queried {
				if !stat.responded {
					m[c.ID] = shortlistStatus{}
					shortlist.Delete(i)
				}
			} else {
				counter--
				stat.queried = true
				m[c.ID] = stat
				go network.SendFindContactMessage(c, target)
			}
		}

		// adds new contacts from channel to shortlist, only blocks for 2 sec
		loopVar := true
		for loopVar {
			select {
			case cl := <-network.cc:
				{
					temp := m[cl[0].ID]
					temp.responded = true
					m[cl[0].ID] = temp

					network.rt.AddContact(cl[0])

					// removes contacts aleady in shortlist
					for i := len(cl) - 1; i >= 0; i-- {
						if m[cl[i].ID].inList {
							cl = append(cl[:i], cl[i+1:]...)
						}
					}

					updateContacts(&cl, target, &m)
					shortlist.Append(cl)
				}
			case <-time.After(2 * time.Second):
				loopVar = false
			}
		}

		// checks if progress is being made, otherwise return K closest
		shortlist.Sort()

		if shortlist.contacts[0] == closestNode && exitOnNext {
			return shortlist.GetContacts(K)
		} else if shortlist.contacts[0] == closestNode {
			exitOnNext = true
		} else {
			exitOnNext = false
		}

	}

}

// sets up information on each contact
func updateContacts(cl *[]Contact, target Contact, m *map[*KademliaID]shortlistStatus) {
	for _, c := range *cl {
		(*m)[c.ID] = newStatus()
	}
}