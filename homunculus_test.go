package homunculus

import (
	"log"
	"os"
	"testing"

	"github.com/john-k-ge/homunculus/cf"
)

const (
	dbError = "db-error"
	timeout = "connection-timeout"
)

var (
	cfConf = cf.CfEnv{
		Inx:     0,
		Guid:    "60b31aa7-6ee8-42a7-97b6-bfb514b42f04",
		CFUser:  os.Getenv("CFUID"),
		CFPass:  os.Getenv("CFPASS"),
		APIHost: "api.system.aws-usw02-pr.ice.predix.io",
		UAAHost: "uaa.system.aws-usw02-pr.ice.predix.io",
	}
	sampleConds = map[string]int64{
		dbError: 5,
		timeout: 3,
	}
)

func TestNewHomunculus(t *testing.T) {
	_, err := NewHomunculus(&cfConf)
	if err != nil {
		log.Printf("Failed to create: %v", err)
		t.Fail()
	}
}

func TestHomunculus_AddBulkConditions(t *testing.T) {
	h, _ := NewHomunculus(&cfConf)
	err := h.AddBulkConditions(sampleConds)
	if err != nil {
		log.Printf("Failed to add bulk: %v", err)
		t.Fail()
	}
	size, err := h.conditions.Size()
	if err != nil {
		log.Printf("Failed to get size: %v", err)
		t.Fail()
	}
	if size != int64(len(sampleConds)) {
		log.Printf("not loaded. Sample: %v vs.  Actual: %v", len(sampleConds), size)
		t.Fail()
	}
}

func TestHomunculus_Increment(t *testing.T) {
	h, _ := NewHomunculus(&cfConf)
	_ = h.AddBulkConditions(sampleConds)
	val1, err := h.Increment(dbError)
	if err != nil {
		log.Printf("Failed to increment %v: %v", dbError, err)
		t.Fail()
	}

	val2, _ := h.Increment(dbError)
	if val2 != (val1 + 1) {
		log.Printf("Failed to increment %v: from %v", dbError, val1)
		t.Fail()
	}

}
