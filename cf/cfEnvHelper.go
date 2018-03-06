package cf

import (
	"fmt"
	"log"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
)

type CfEnv struct {
	AppName string
	Inx     int
	Guid    string
	RHost   string
	RPort   string
	RPasswd string
	CFUser  string
	CFPass  string
	APIHost string
	UAAHost string
}

func GetCFEnvVals() *CfEnv {
	appEnv, _ := cfenv.Current()
	name := appEnv.Name
	idx := appEnv.Index
	guid := appEnv.AppID
	log.Printf("idx: `%v`, guid: `%v`", idx, guid)
	if len(guid) == 0 {
		log.Printf("No guid found in CF env")
	}
	var rHost, rPort, rPass string
	services := appEnv.Services
	caches, err := services.WithLabel("predix-cache")
	if err != nil {
		log.Printf("Error %v", err)

	} else {
		for _, cache := range caches {
			for credKey, credVal := range cache.Credentials {
				switch {
				case strings.EqualFold(credKey, "host"):
					rHost = credVal.(string)
				case strings.EqualFold(credKey, "port"):
					rPort = fmt.Sprint(credVal)
				case strings.EqualFold(credKey, "password"):
					rPass = credVal.(string)
				}
			}
		}
	}

	conf := CfEnv{
		AppName: name,
		Inx:     idx,
		Guid:    guid,
		RHost:   rHost,
		RPort:   rPort,
		RPasswd: rPass,
	}
	return &conf
}
