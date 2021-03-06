package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
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

type RequestMessage struct {
	Type   string
	Sender Contact
	Target KademliaID // findcontact
	Hash   string     // finddata
	Data   string     // storedata
}

type ResponseMessage struct {
	Type     string
	Sender   Contact
	Status   string    // "ok" | "fail"
	Data     string    // finddata
	Contacts []Contact // findcontact
	SenderID string
	SenderIP string
}

// NewNetwork constructor
func NewNetwork(ip string) Network {
	n := Network{}
	n.cc = make(chan []Contact, ALPHA*2)
	n.rt = NewRoutingTable(NewContact(NewRandomKademliaID(ip), ip))
	n.dataStore = make(map[string]string)
	return n
}

func (network *Network) Listen(port string) error {

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
		return err
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

	req := RequestMessage{}

	res := ResponseMessage{}
	res.Sender = network.rt.me

	dec.Decode(&req)

	switch req.Type {

	// for CLI
	case "exit":
		res.Type = "exit"
		res.Status = "ok"

		enc.Encode(res)

		os.Exit(0)

	// for CLI
	case "put":
		data := req.Data
		id := network.storeData(data)

		res.Type = "put"
		res.Status = "ok"
		res.Data = id

		enc.Encode(res)

	case "ping":

		res.Type = "ping"
		res.Status = "ok"

		enc.Encode(res)
		network.rt.AddContact(req.Sender)

	case "findcontact":
		network.createFindContactResponse(&res, req)	
		enc.Encode(res)

	case "finddata":
		network.createFindDataResponse(&res, req)
		enc.Encode(res)

	case "storedata":

		network.storeLocalData(string(req.Data))

		res.Type = "storedata"
		res.Status = "ok"

		enc.Encode(res)

	// for CLI
	case "get":
		data := network.getData(string(req.Data))
		network.createGetCLIResponse(&res, req, data)
		enc.Encode(res)

	// for CLI
	case "printds":
		res.Type = "printds"
		res.Status = "ok"
		res.Data = network.getDataStoreString()

		enc.Encode(res)

	// for CLI
	case "printrt":
		res.Type = "printrt"
		res.Status = "ok"
		res.Data = network.rt.String()

		enc.Encode(res)
	}
}

func (network Network) createGetCLIResponse(res *ResponseMessage, req RequestMessage, data DataReturn)  {
	res.Type = "get"
	res.Data = data.Data
	res.SenderID = data.From
	res.SenderIP = data.FromIP
	if res.SenderID != "" && res.SenderIP != "" {
		res.Status = "ok"
	} else {
		res.Status = "fail"
	}
}


func (network Network) createFindContactResponse(res *ResponseMessage, req RequestMessage)  {
	target := req.Target
	contacts := network.findClosestLocalContacts(target)
	res.Type = "findcontact"
	res.Status = "ok"
	res.Contacts = contacts

	network.rt.AddContact(req.Sender)
	network.PrintClosestContacts()
}

func (network Network) createFindDataResponse(res *ResponseMessage, req RequestMessage)  {
	data := network.getLocalData(req.Hash)

	res.Type = "finddata"
	if data != "" {
		res.Data = data
		res.Status = "ok"
	} else {
		res.Status = "fail"
	}
}

func (network Network) SendPingMessageIP(address string) string {
	id := NewRandomKademliaID("dummy id")
	c := NewContact(id, address)
	return network.SendPingMessage(&c)

}

func sendTCPRequest(req RequestMessage, contact *Contact) (res ResponseMessage, err error) {
	conn, err := net.Dial("tcp", contact.Address)
	if err != nil {
		return ResponseMessage{}, err
	}

	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	enc.Encode(req)

	dec.Decode(&res)
	return
}

func (network Network) createPingMessage() RequestMessage {
	req := RequestMessage{}
	req.Sender = network.rt.me
	req.Type = "ping"

	return req
}

// SendPingMessage returns the contact that was pinged, can be used to obtain full contact from just ip
func (network Network) SendPingMessage(contact *Contact) string {
	res, err := sendTCPRequest(network.createPingMessage(), contact)
	if err != nil {
		return "fail"
	}

	network.rt.AddContact(res.Sender)

	return res.Status
}

func (network Network) createFindContactMessage(target KademliaID) RequestMessage {
	req := RequestMessage{}
	req.Sender = network.rt.me
	req.Type = "findcontact"
	req.Target = target

	return req
}

// first contact in list posted to cc is "contact"
func (network Network) SendFindContactMessage(contact Contact, target KademliaID) {
	req := network.createFindContactMessage(target)

	res, err := sendTCPRequest(req, &contact)
	if err != nil {
		return
	}

	s := make([]Contact, 1)
	s[0] = contact
	network.cc <- append(s, res.Contacts...)
}

func (network *Network) createFindDataMessage(hash string) RequestMessage {
	req := RequestMessage{}
	req.Sender = network.rt.me
	req.Type = "finddata"
	req.Hash = hash

	return req
}

func (network *Network) SendFindDataMessage(hash string, contact Contact, c chan DataReturn) {
	req := network.createFindDataMessage(hash)

	res, err := sendTCPRequest(req, &contact)
	if err != nil || res.Status == "fail" {
		return
	}

	c <- DataReturn{res.Data, contact.ID.String(), contact.Address}
}

func (network *Network) createSendStoreMessage(data string) RequestMessage {
	req := RequestMessage{}
	req.Sender = network.rt.me
	req.Type = "storedata"
	req.Data = data

	return req
}

func (network *Network) SendStoreMessage(data string, contact Contact) {
	req := network.createSendStoreMessage(data)

	_, err := sendTCPRequest(req, &contact)
	if err != nil {
		return
	}
}

func (network *Network) getLocalData(hash string) string {
	return network.dataStore[hash]
}

func (network *Network) storeLocalData(data string) { // needs fix NewKademliaID not working as intended
	network.dataStore[NewRandomKademliaID(data).String()] = data
}

func (network *Network) storeData(data string) string {
	id := NewRandomKademliaID(data)
	contacts := network.ContactLookup(*id)
	for _, c := range contacts {
		go network.SendStoreMessage(data, c)
	}
	if len(contacts) <= 0 {
		network.storeLocalData(data)
	}
	return id.String()
}

type DataReturn struct {
	Data   string
	From   string
	FromIP string
}

// first find closest K nodes to hash
// then probes and ask for data
// returns when it gets data from one of the contacts
// TODO: HANDLE IF NOONE HAS
func (network *Network) getData(hash string) DataReturn {
	id := NewKademliaID(hash)
	contacts := network.ContactLookup(*id)

	res := make(chan DataReturn)
	for _, c := range contacts {
		go network.SendFindDataMessage(hash, c, res)
	}
	select {
	case data := <-res:
		return data
	case <-time.After(2 * time.Second):
		return DataReturn{}
	}

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
	notResponded := 0
	waitTime := 2 * time.Second
	for {

		// removes unresponsive contacts and sends rpc find to ALPHA nr of contacts
		counter := ALPHA
		for i := 0; i < len(shortlist.contacts); i++ {
			c := shortlist.contacts[i]
			if counter < 1 {
				break
			}

			stat := m[*c.ID]
			if stat.queried {
				if !stat.responded {
					m[*c.ID] = shortlistStatus{}
					shortlist.Delete(i)
					i--
				}
			} else {
				counter--
				stat.queried = true
				m[*c.ID] = stat
				notResponded++
				go network.SendFindContactMessage(c, target)
			}
		}

		// adds new contacts from channel to shortlist, only blocks for 2 sec
		loopVar := true
		for loopVar {
			if notResponded > 0 {
				waitTime = 2 * time.Second
			} else {
				waitTime = 0 * time.Millisecond
			}
			select {
			case cl := <-network.cc:
				{
					notResponded--
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
			case <-time.After(waitTime):
				loopVar = false
			}
		}

		// checks if progress is being made, otherwise return K closest
		shortlist.Sort()

		if shortlist.contacts[0] == closestNode {
			if shortlist.Len() < K {
				return shortlist.GetContacts(shortlist.Len())
			}
			return shortlist.GetContacts(K)
		}
		closestNode = shortlist.contacts[0]
	}
}

// ta bort? få högre coverage
// PrintClosestContacts prints contacts
func (network *Network) PrintClosestContacts() {
	fmt.Println(network.findClosestLocalContacts(*network.rt.me.ID))
}

// get string to neatly represent datastore
func (network *Network) getDataStoreString() string {
	b, _ := json.MarshalIndent(network.dataStore, "", "  ")

	return string(b)
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
