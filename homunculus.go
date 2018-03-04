package homunculus

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"errors"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/john-k-ge/homunculus/cf"
	"github.com/john-k-ge/homunculus/cs"
)

type Homunculus struct {
	//creds      string
	ceilings   map[string]int64
	conditions cs.ConditionSet
	cf         cf.CfClient
}

type HConfig struct {
	CFUser    string
	CFPass    string
	RedisHost string
	RedisPort string
	RedisPass string
}

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
	pc, err := services.WithNameUsingPattern("cache")
	if err != nil {
		log.Printf("Error ")
	}

	var conditions cs.ConditionSet

	switch {
	case len(config.RedisHost) == 0:
		fallthrough
	case len(config.RedisPort) == 0:
		fallthrough
	case len(config.RedisPass) == 0:
		log.Printf("No Redis info passed.  Using local")
		conditions = cs.NewLocalConditionSet()
	default:
		conditions, err = cs.NewRemoteCondtionSet(config.RedisHost, config.RedisPort, config.RedisPass, idx)
	}

	if err != nil {
		log.Printf("Failed to created ConditionSet: %v", err)
		return nil, err
	}

	h := Homunculus{
		ceilings:   make(map[string]int64),
		conditions: conditions,
		// Need to generate a cfConfig -> cfClient
	}

	return &h, nil
}

func (h *Homunculus) AddCondition(cond string, max int64) error {
	h.ceilings[cond] = max
	err := h.conditions.SaveCondition(cond, 0)
	if err != nil {
		log.Printf("Failed to conditions condition %v: %v", cond, err)
		return err
	}
	return nil
}

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

func (h *Homunculus) die(cond string) error {
	log.Printf("`die` called due to failed condition '%v'", fmt.Sprintf("%v.%v", cond, h.idx))
	// make cf kill app instance call

	return nil
}
