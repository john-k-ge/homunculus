package cf

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/john-k-ge/homunculus/uaa"
)

const (
	appPath = "/v2/apps/%v/instances/%v" // e.g. /v2/apps/f68b80bc-ab1d-46ee-8b42-94337c96e143/instances/0
	https   = "https://"
)

type CfClient struct {
	apiUrl    string
	appGuid   string
	index     int
	uaaClient *uaa.UaaClient
}

type CfConfig struct {
	Api, Uaa, AppGuid, Uid, Pass string
	Index                        int
}

func validateConfig(c *CfConfig) error {
	if len(c.Api) == 0 || len(c.Uaa) == 0 || len(c.AppGuid) == 0 || len(c.Uid) == 0 || len(c.Pass) == 0 {
		log.Printf("Validation failed")
		log.Printf("Config vals: '%v', '%v', '%v', '%v', '%v'", c.Api, c.Uaa, c.AppGuid, c.Uid, c.Pass)
		return errors.New("invalid config values")
	}
	return nil
}

func NewCfClient(config CfConfig) (*CfClient, error) {

	err := validateConfig(&config)
	if err != nil {
		log.Printf("failed to create CFClient: %v", err)
		return nil, err
	}

	if err != nil {
		log.Printf("Failed to parse index value `%v`: %v", config.Index, err)
		return nil, err
	}

	cf := CfClient{
		apiUrl:  https + config.Api,
		appGuid: config.AppGuid,
		index:   config.Index,
	}

	cf.uaaClient = uaa.NewUaaClient(config.Uaa, config.Uid, config.Pass)

	return &cf, nil
}

func (cf *CfClient) StopCFApp() error {
	client, err := cf.uaaClient.Authenticate()
	if err != nil {
		log.Printf("Failed to create CF client: %v", err)
		return err
	}

	delPath := fmt.Sprintf(appPath, cf.appGuid, cf.index)
	req, err := http.NewRequest(http.MethodDelete, cf.apiUrl+delPath, strings.NewReader(""))
	if err != nil {
		log.Printf("Failed to create http request:  %v", err)
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		return err
	}
	log.Printf("Response code: %v", resp.StatusCode)
	// I should be dead by now
	return nil
}
