package d7024e

import (
	"testing"
)

func TestNewKademliaID(t *testing.T) {

	tests := []struct {
		input KademliaID
		expected KademliaID
	}{
		{ *NewKademliaID("0001FF0000000000000000000000000000000000"), KademliaID{0,1,255,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0} },
		{ *NewKademliaID("00000000000000000000000000000000000000001"), KademliaID{0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0} },
	}
	
	for _, test := range tests {
		if test.input != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, test.input)
		}
	}

}


func TestNewRandomKademliaID(t *testing.T) {
	expected := KademliaID{ 177, 166, 22, 74, 100, 64, 176, 32, 223, 254, 252, 124, 228, 219, 191, 51, 232, 96, 78, 5 }
	input := "asd"
	output := *NewRandomKademliaID(input)
	
	if expected != output {
		t.Errorf("Expected %v with input %v, got %v", expected, input, output)
	}
}

func TestLess(t *testing.T) {
	id1 := KademliaID{0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0}
	id2 := KademliaID{0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0}
	id3 := KademliaID{0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0,0}

	if id1.Less(&id2) {
		t.Errorf("expected false on %v < %v got true", id1, id2)
	}

	if !id1.Less(&id3) {
		t.Errorf("expected true on %v < %v got false", id1, id2)
	}

	if id3.Less(&id1) {
		t.Errorf("expected false on %v < %v got true", id3, id1)
	}

}


func TestEquals(t *testing.T) {
	id1 :=  KademliaID{0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0}
	id2 :=  KademliaID{0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0}

	if !id1.Equals(&id2) {
		t.Errorf("expected %v.Equals(%v) to be true but got false", id1, id2)
	}
}

func TestCalcDistance(t *testing.T) {
	id1 := 		KademliaID{0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0,0}
	id2 := 		KademliaID{0,0,0,0,0,0,0,0,0,1,1,0,0,0,0,0,0,0,0,0}
	expected := KademliaID{0,0,0,0,0,0,0,0,0,0,1,0,0,0,0,0,0,0,0,0}

	output := id1.CalcDistance(&id2) 

	if expected != *output {
		t.Errorf("expected %v, got %v with input %v and %v", expected, output, id1, id2)
	}
}



