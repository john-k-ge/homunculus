package cf

import (
	"log"
	"os"
	"testing"
)

const (
	goodApi   = "api.system.aws-usw02-pr.ice.predix.io"
	badApi    = "localhost"
	goodGuid  = "60b31aa7-6ee8-42a7-97b6-bfb514b42f04"
	badGuid   = "herp-derp-unga-bunga"
	goodIndex = 0
	badIndex  = "9001"

	uaaUrl = "uaa.system.aws-usw02-pr.ice.predix.io"
)

var (
	uid  = os.Getenv("CFUID")
	pass = os.Getenv("CFPASS")
)

func Test_GoodNewCfClient(t *testing.T) {
	good := CfConfig{
		Uaa:  uaaUrl,
		Uid:  uid,
		Pass: pass,

		Api:     goodApi,
		AppGuid: goodGuid,
		Index:   goodIndex,
	}

	t.Run("Happy path test", func(t *testing.T) {
		_, err := NewCfClient(good)
		if err != nil {
			log.Printf("Client failed: %v", err)
			t.Fail()
		}
		log.Print("Success")
	})

}

func Test_GoodStopApp(t *testing.T) {
	good := CfConfig{
		Uaa:  uaaUrl,
		Uid:  uid,
		Pass: pass,

		Api:     goodApi,
		AppGuid: goodGuid,
		Index:   goodIndex,
	}

	t.Run("Happy path stop app", func(t *testing.T) {
		c, _ := NewCfClient(good)
		err := c.StopCFApp()
		if err != nil {
			log.Printf("Failed to stop app: %v", err)
			t.Fail()
		}
	})
}
