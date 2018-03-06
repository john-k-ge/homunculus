package uaa

import (
	"log"
	"os"
	"testing"
)

const (
	goodUaaHost = "uaa.system.aws-usw02-pr.ice.predix.io"
	badUaaHost  = "localhost"
	badUid      = "unga"
	badPass     = "bunga"
)

var (
	goodUid  = os.Getenv("CFUID")
	goodPass = os.Getenv("CFPASS")
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
