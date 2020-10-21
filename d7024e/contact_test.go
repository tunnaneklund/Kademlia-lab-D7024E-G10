package d7024e


import (
	"testing"
)

func TestDelete(t *testing.T) {
	candidates := ContactCandidates{}
	c1 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("BBBBBBBB00000000000000000000000000000000"), "localhost:8001")
	c3 := NewContact(NewKademliaID("CCCCCCCC00000000000000000000000000000000"), "localhost:8001")
	contacts := []Contact{c1, c2, c3}
	candidates.Append(contacts)

	candidates.Delete(1)

	output1 := candidates.contacts[1]
	expected1 := c3

	output2 := candidates.Len()
	expected2 := 2

	if output1 != expected1 {
		t.Errorf("expected %v from candidates.contacts[1] but was %v", expected1, output1)
	}

	if output2 != expected2 {
		t.Errorf("expected %v from candidates.Len() but was %v", expected2, output2)
	}

}