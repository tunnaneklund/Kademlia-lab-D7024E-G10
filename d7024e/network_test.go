package d7024e

import (
	"testing"
	"fmt"
)

func TestListenBadPort(t *testing.T) {
	
	n := Network{}

	input := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	// expected: out != nil
	output := n.Listen(input)

	if output == nil {
		t.Errorf("expected error on input %v got nil", input)
	}
}

func TestCreatePingMessage(t *testing.T) {
	n := NewNetwork("tmp")
	res := n.createPingMessage()
	
	out1 := res.Type
	out2 := res.Sender.Address

	exp1 := "ping"
	exp2 := "tmp"

	if out1 != exp1 {
		t.Errorf("expected res.Type to be %v but was %v", exp1, out1)
	}

	if out2 != exp2 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp2, out2)
	}
	
}

func TestCreateFindContactMessage(t *testing.T) {
	n := NewNetwork("tmp")
	tar := *NewRandomKademliaID("hello")
	res := n.createFindContactMessage(tar)
	
	out1 := res.Type
	out2 := res.Sender.Address
	out3 := res.Target.String()

	exp1 := "findcontact"
	exp2 := "tmp"
	exp3 := "c45c4531861ea506fe66fe15f53ae6c1739851bc"

	if out1 != exp1 {
		t.Errorf("expected res.Type to be %v but was %v", exp1, out1)
	}

	if out2 != exp2 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp2, out2)
	}
	
	if out3 != exp3 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp3, out3)
	}

}

func TestCreateFindDataMessage(t *testing.T) {
	n := NewNetwork("tmp")
	hash := "heloothisisahash"
	res := n.createFindDataMessage(hash)
	
	out1 := res.Type
	out2 := res.Sender.Address
	out3 := res.Hash

	exp1 := "finddata"
	exp2 := "tmp"
	exp3 := "heloothisisahash"

	if out1 != exp1 {
		t.Errorf("expected res.Type to be %v but was %v", exp1, out1)
	}

	if out2 != exp2 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp2, out2)
	}
	
	if out3 != exp3 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp3, out3)
	}

}


func TestCreateSendStoreMessage(t *testing.T) {
	n := NewNetwork("tmp")
	data := "this is the data"
	res := n.createSendStoreMessage(data)
	
	out1 := res.Type
	out2 := res.Sender.Address
	out3 := res.Data

	exp1 := "storedata"
	exp2 := "tmp"
	exp3 := "this is the data"

	if out1 != exp1 {
		t.Errorf("expected res.Type to be %v but was %v", exp1, out1)
	}

	if out2 != exp2 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp2, out2)
	}
	
	if out3 != exp3 {
		t.Errorf("expected res.Sender.Address to be %v but was %v", exp3, out3)
	}

}


func TestLocalData(t *testing.T) {
	network := NewNetwork("tmp")
	network.storeLocalData("data1")
	network.storeLocalData("data2") 
	network.storeLocalData("data3")

	input := "0574856fed300d91b92e9f61b1574240e2b5e793" // hash for "data2"
	output := network.getLocalData(input)
	expected := "data2"

	if expected != output {
		t.Errorf("Expected %v with input %v, got %v", expected, input, output)
	}

}

func TestFindClosestLocalContacts(t *testing.T) {
	n := NewNetwork("tmp")
	c0 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
	n.rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	n.rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	n.rt.AddContact(c0)
	n.rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))

	input := *NewKademliaID("1111111200000000000000000000000000000000")
	output := n.findClosestLocalContacts(input)


	if output[0].ID != c0.ID {
		t.Errorf("expected ID of nearest contact to be %v but was %v", c0.ID, output[0].ID)
	}
}

func TestNewStatus(t *testing.T) {
	output := newStatus()
	expected := shortlistStatus{false, false, true}
	if output != expected {
		t.Errorf("EXPECTED OUTPUT to be %v but was %v", expected, output)
	}
}

func TestUpdateContacts(t *testing.T) {
	idmap := make(map[KademliaID]shortlistStatus)
	c1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	c2 := NewContact(NewKademliaID("AAAAAAAA00000000000000000000000000000000"), "localhost:8001")
	c3 := NewContact(NewKademliaID("BBBBBBBB00000000000000000000000000000000"), "localhost:8001")

	contactlst := []Contact{c1, c2, c3}

	updateContacts(&contactlst, &idmap)

	expected := shortlistStatus{false, false, true}
	output := idmap[*c1.ID]

	if output != expected {
		t.Errorf("expected output to be %v but was %v", expected, output)
	}
}

func TestGetDataStoreString(t *testing.T) {
	network := NewNetwork("tmp")
	network.storeLocalData("data1")
	network.storeLocalData("data2") 
	network.storeLocalData("data3")

	fmt.Println(network.getDataStoreString())
}

func TestPrintClosestContacts(t *testing.T) {
	n := NewNetwork("tmp")
	c0 := NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")
	n.rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	n.rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	n.rt.AddContact(c0)
	n.rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))

	n.PrintClosestContacts()
} 

func TestCreateFindContactResponse(t *testing.T) {
	n := NewNetwork("tmp")
	res := ResponseMessage{}
	req := RequestMessage{}

	req.Target = *NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	req.Sender = NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002")

	n.createFindContactResponse(&res, req)

	out1 := res.Type
	exp1 := "findcontact"

	out2 := res.Status
	exp2 := "ok"

	out3 := len(res.Contacts)
	exp3 := 0

	if out1 != exp1 {
		t.Errorf("expected res.Type to be %v but was %v", exp1, out1)
	}

	if out2 != exp2 {
		t.Errorf("expected res.Status to be %v but was %v", exp2, out2)
	}
	
	if out3 != exp3 {
		t.Errorf("expected len(res.Contacts) to be %v but was %v", exp3, out3)
	}
	
}

func TestCreateFindDataResponse(t *testing.T) {
	n := NewNetwork("tmp")
	res := ResponseMessage{}
	req := RequestMessage{}

	req.Hash = "asd"	

	n.createFindDataResponse(&res, req)
	if res.Status != "fail" {
		t.Errorf("expected status to be %v but was %v", "fail", res.Status)
	}

	n.dataStore["asd"] = "anton"
	n.createFindDataResponse(&res, req)
	
	if res.Status != "ok" {
		t.Errorf("expected status to be %v but was %v", "fail", res.Status)
	}

	
}
func TestCreateGetCLIResponse(t *testing.T) {
	n := NewNetwork("tmp")
	res := ResponseMessage{}
	req := RequestMessage{}
	data := DataReturn{ "data", "ID", "IP" }

	n.createGetCLIResponse(&res, req, data)

	if res.Status != "ok" {
		t.Errorf("expected res.Status to be ok but was %v", res.Status)
	}
}


