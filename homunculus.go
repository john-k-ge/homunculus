package homunculus

import (
	"log"
	"os"

	"github.com/john-k-ge/homunculus/cf"
	"github.com/john-k-ge/homunculus/cs"
)

// Homunculus : Defines the Homunculus watcher
type Homunculus struct {
	ceilings   map[string]int64
	conditions cs.ConditionSet
	cf         *cf.CfClient
}

// NewHomunculus : Creates a new Homunculus
func NewHomunculus(config *cf.CfEnv) (*Homunculus, error) {
	var conds cs.ConditionSet
	conds, err := cs.NewRemoteCondtionSet(config)

	if err != nil {
		log.Printf("Failed to create remote condition set: %v", err)
		log.Print("Defaulting to local.")
		conds = cs.NewLocalConditionSet()
	}

	cfClient, err := cf.NewCfClient(cf.CfConfig{
		Api:     config.APIHost,
		Uaa:     config.UAAHost,
		AppGuid: config.Guid,
		Uid:     config.CFUser,
		Pass:    config.CFPass,
		Index:   config.Inx,
	})

	if err != nil {
		log.Printf("Failed to create CFClient: %v", err)
		return nil, err
	}

	h := Homunculus{
		ceilings:   make(map[string]int64),
		conditions: conds,
		cf:         cfClient,
	}

	return &h, nil
}

func (h *Homunculus) AddBulkConditions(conds map[string]int64) error {
	log.Printf("Processing %v bulk conditions", len(conds))
	for c, m := range conds {
		err := h.AddCondition(c, m)
		if err != nil {
			log.Printf("Failed to process bulk: c/m `%v/%v`: %v", c, m, err)
			return err
		}
	}
	return nil
}

// AddCondition : Registers a new condition and corresponding max val
func (h *Homunculus) AddCondition(cond string, max int64) error {
	log.Printf("Processing %v, %v", cond, max)
	h.ceilings[cond] = max
	err := h.conditions.SaveCondition(cond, 0)
	if err != nil {
		log.Printf("Failed to add condition %v: %v", cond, err)
		return err
	}
	return nil
}

// Increment : Increment a given condition and die if necessary
func (h *Homunculus) Increment(cond string) (int64, error) {
	log.Printf("Incrementing %v", cond)

	current, err := h.conditions.IncrementCondition(cond)
	if err != nil {
		log.Printf("Failed to pre-check `%v` before increment: %v", cond, err)
		return 0, err
	}
	if current >= h.ceilings[cond] {
		log.Printf("%v: '%v' hit max '%v'.  Calling die()", cond, current, h.ceilings[cond])
		h.die(cond)
		return current, nil
	}
	log.Printf("'%v' is now %v", cond, current)
	return current, nil
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
