package d7024e

const ALPHA = 3

type Kademlia struct {
	rt        *RoutingTable
	dataStore map[string][]byte
}

func NewKademliaNode(IPAddress string) Kademlia {
	k := new(Kademlia)
	k.rt = NewRoutingTable(NewContact(NewKademliaID(IPAddress), IPAddress))

	return k
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {

	return kademlia.rt.FindClosestContacts(target.ID, ALPHA)

}

func (kademlia *Kademlia) LookupData(hash string) []byte {
	return kademlia.dataStore[hash]
	// TODO(maybe): check if data exists
}

func (kademlia *Kademlia) Store(data []byte) {
	kademlia.dataStore[NewKademliaID(string(data)).String()] = data
}
