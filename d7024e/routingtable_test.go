package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000"))

	c1 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002")
	c2 := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(c1)
	rt.AddContact(c2)

	contacts := rt.FindClosestContacts(NewKademliaID("2111111500000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
		
	}
	
	if contacts[0].ID != c2.ID {
		t.Error("FindClosestContacts on 2111111500000000000000000000000000000000 did not give 2111111400000000000000000000000000000000 as closest")
	}

	if contacts[1].ID != c1.ID {
		t.Error("FindClosestContacts on 2111111500000000000000000000000000000000 did not give 1111111400000000000000000000000000000000 as second closest")
	}
}

func TestGetBucketIndexOnSelf(t *testing.T) {
	me := NewContact(NewKademliaID("0000000000000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(me)

	input := me.ID
	expected := 159
	output := rt.getBucketIndex(input)
	
	if expected != output {
		t.Errorf("Expected %v with input %v, got %v", expected, input, output)
	}
}
