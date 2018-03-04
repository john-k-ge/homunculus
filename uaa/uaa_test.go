package uaa

import (
	"log"
	"testing"
)

const (
	uaaHost = "uaa.system.aws-usw02-pr.ice.predix.io"
	uid     = "predix-support_systemuser"
	pass    = ">pTz5yvM97N@y_G"
)

func TestNewUaaClient(t *testing.T) {
	u := NewUaaClient(uaaHost, uid, pass)
	if u == nil {
		log.Printf("Failed to create UAA client with: %v, %v, %v", uaaHost, uid, pass)
		t.Fail()
	}
}

func TestUaaClient_Authenticate(t *testing.T) {
	u := NewUaaClient(uaaHost, uid, pass)
	_, err := u.Authenticate()
	if err != nil {
		log.Printf("Failed to authenticate with %v, %v, %v: %v", uaaHost, uid, pass, err)
		t.Fail()
	}
}
