package d7024e


import (
	"testing"
	"fmt"
)

func TestAddContact(t *testing.T) {
	bucket := newBucket()

	a := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	b := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8001")

	bucket.AddContact(b)
	bucket.AddContact(a)
	
	if a != bucket.list.Front().Value {
		t.Error("expected newest contact to be at front of the list but was not")
	}

	bucket.AddContact(b)
	if b != bucket.list.Front().Value {
		t.Error("expected newly updated contact to be at front of the list but was not")
	}
	
}

func TestLen(t *testing.T) {
	input := newBucket()
	expected := 0
	output := input.Len()

	if expected != output {
		t.Errorf("expected %v but got %v with input %v", expected, output, input)
	}
}

func TestString(t *testing.T) {
	bucket := newBucket()

	a := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	b := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8001")

	bucket.AddContact(b)
	bucket.AddContact(a)

	fmt.Println(bucket.String())
}