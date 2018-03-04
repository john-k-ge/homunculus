package cf

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/john-k-ge/homunculus/uaa"
)

const (
	appPath = "/v2/apps/%v/instances/%v" // e.g. /v2/apps/f68b80bc-ab1d-46ee-8b42-94337c96e143/instances/0
)

type CfClient struct {
	apiUrl    string
	appGuid   string
	index     int
	uaaClient *uaa.UaaClient
}

type CfConfig struct {
	Api, Uaa, AppGuid, Uid, Pass, Index string
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

	indxInt, err := strconv.Atoi(config.Index)
	if err != nil {
		log.Printf("Failed to parse index value `%v`: %v", config.Index, err)
		return nil, err
	}

	cf := CfClient{
		apiUrl:  config.Api,
		appGuid: config.AppGuid,
		index:   indxInt,
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

// REST call stuff for CF

/*

DO NOT RESTAGE!!  INSTEAD RESTART AN APP INSTANCE!!

`DELETE /v2/apps/f68b80bc-ab1d-46ee-8b42-94337c96e143/instances/0`

func (app *CFApp) Restage() error {
	client := uaaIntegration.GetPlatformUaaClient()

	appGuid, err := cfCalls.PostCfObject(client, fmt.Sprintf(restageAppPrfx, app.Guid), "")

	if err != nil || len(appGuid) <= 0 {
		log.Printf("Could not restage app %v:  %v", app.Name, err)
		return err
	}

	app.Running = true

	log.Printf("Restaged app with Guid: %v", appGuid)
	return nil
}

*/
