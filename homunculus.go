package homunculus

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"errors"

	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/john-k-ge/homunculus/cf"
	"github.com/john-k-ge/homunculus/cs"
)

// Homunculus : Defines the Homunculus watcher
type Homunculus struct {
	//creds      string
	ceilings   map[string]int64
	conditions cs.ConditionSet
	cf         *cf.CfClient
}

// HConfig : Defines the Homunculus necessary config properties
type HConfig struct {
	CFUser  string
	CFPass  string
	APIHost string
	UAAHost string
}

// NewHomunculus : Creates a new Homunculus
func NewHomunculus(config *HConfig) (*Homunculus, error) {
	guid := os.Getenv("CF_INSTANCE_GUID")
	if len(guid) == 0 {
		log.Printf("No guid found in CF env")
		return nil, errors.New("no guid found in CF env")
	}
	cfIndex := os.Getenv("CF_INSTANCE_INDEX")
	idx, err := strconv.Atoi(cfIndex)
	if err != nil {
		log.Printf("Could not parse app index: %v", err)
		return nil, err
	}

	appEnv, _ := cfenv.Current()
	services := appEnv.Services
	caches, err := services.WithNameUsingPattern("cache")
	if err != nil {
		log.Printf("Error ")
	}

	var rHost, rPort, rPass string
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

	var conditions cs.ConditionSet
	conditions, err = cs.NewRemoteCondtionSet(rHost, rPort, rPass, idx)

	if err != nil {
		log.Printf("Failed to create remote condition set: %v", err)
		log.Print("Defaulting to local.")
		conditions = cs.NewLocalConditionSet()
	}

	cfClient, err := cf.NewCfClient(cf.CfConfig{
		Api:     config.APIHost,
		Uaa:     config.UAAHost,
		AppGuid: guid,
		Uid:     config.CFUser,
		Pass:    config.CFPass,
	})

	if err != nil {
		log.Printf("Failed to create CFClient: %v", err)
	}

	h := Homunculus{
		ceilings:   make(map[string]int64),
		conditions: conditions,
		cf:         cfClient,
	}

	return &h, nil
}

// AddCondition : Registers a new condition and corresponding max val
func (h *Homunculus) AddCondition(cond string, max int64) error {
	h.ceilings[cond] = max
	err := h.conditions.SaveCondition(cond, 0)
	if err != nil {
		log.Printf("Failed to conditions condition %v: %v", cond, err)
		return err
	}
	return nil
}

// Increment : Increment a given condition and die if necessary
func (h *Homunculus) Increment(cond string) error {
	current, err := h.conditions.IncrementCondition(cond)
	if err != nil {
		log.Printf("Failed to pre-check `%v` before increment: %v", cond, err)
		return err
	}
	if current >= h.ceilings[cond] {
		log.Printf("%v: '%v' hit max '%v'.  Calling die()", cond, current, h.ceilings[cond])
		h.die(cond)
		return nil
	}
	log.Printf("'%v' is now %v", cond, current)
	return nil
}

// die : Try to gracefully shutdown via a CF kill; otherwise just die hard
func (h *Homunculus) die(cond string) error {
	log.Printf("`die` called due to failed condition '%v'", cond)
	err := h.cf.StopCFApp()
	if err != nil {
		log.Printf("Failed to stop app gracefully: %v", err)
		os.Exit(1)
	}
	return nil
}
