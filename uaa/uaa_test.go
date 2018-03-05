package uaa

import (
	"log"
	"testing"
)

const (
	goodUaaHost = "uaa.system.aws-usw02-pr.ice.predix.io"
	goodUid     = "predix-support_systemuser"
	goodPass    = ">pTz5yvM97N@y_G"
	badUaaHost  = "localhost"
	badUid      = "unga"
	badPass     = "bunga"
)

func TestUaaClient_AuthenticateGood(t *testing.T) {
	u := NewUaaClient(goodUaaHost, goodUid, goodPass)
	_, err := u.Authenticate()
	if err != nil {
		log.Printf("Failed to authenticate with %v, %v, %v: %v", goodUaaHost, goodUid, goodPass, err)
		t.Fail()
	}
}

func TestUaaClient_AuthenticateBadHost(t *testing.T) {
	u := NewUaaClient(badUaaHost, goodUid, goodPass)
	_, err := u.Authenticate()
	if err == nil {
		log.Print("Bad host should fail")
		t.Fail()
	}
	log.Printf("Err: %v", err)
}

func TestUaaClient_AuthenticateBadUid(t *testing.T) {
	u := NewUaaClient(goodUaaHost, badUid, goodPass)
	_, err := u.Authenticate()
	if err == nil {
		log.Print("Bad uid should fail")
		t.Fail()
	}
	log.Printf("Err: %v", err)
}

func TestUaaClient_AuthenticateBadPass(t *testing.T) {
	u := NewUaaClient(goodUaaHost, goodUid, badPass)
	_, err := u.Authenticate()
	if err == nil {
		log.Print("Bad pass should fail")
		t.Fail()
	}
	log.Printf("Err: %v", err)
}
