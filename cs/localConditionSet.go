package cs

import (
	"errors"
	"log"
	"sync"
)

type LocalCache struct {
	sync.RWMutex
	conditionSet map[string]int64
}

func NewLocalConditionSet() *LocalCache {
	lc := LocalCache{
		conditionSet: make(map[string]int64),
	}
	log.Printf("Creating a localcache: %v", lc)
	return &lc
}

func (lc *LocalCache) SaveCondition(condition string, count int64) error {
	if len(condition) == 0 {
		log.Print("Key must have a value")
		return errors.New("key must have a value")
	}
	lc.RLock()
	lc.conditionSet[condition] = count
	lc.Unlock()
	return nil
}

func (lc *LocalCache) CheckCondition(condition string) (int64, error) {
	v, exists := lc.conditionSet[condition]
	if !exists {
		return 0, errors.New("condition does not exist")
	}
	return v, nil
}

func (lc *LocalCache) ConditionExists(condition string) (bool, error) {
	_, exists := lc.conditionSet[condition]
	return exists, nil
}

func (lc *LocalCache) DeleteCondition(condition string) error {
	lc.RLock()
	delete(lc.conditionSet, condition)
	lc.Unlock()
	return nil
}

func (lc *LocalCache) ClearConditions() error {
	lc.RLock()
	lc.conditionSet = make(map[string]int64)
	lc.Unlock()
	return nil
}

func (lc *LocalCache) IncrementCondition(condition string) (int64, error) {
	_, exists := lc.conditionSet[condition]
	if !exists {
		return 0, errors.New("condition " + condition + " not found")
	}
	lc.RLock()
	lc.conditionSet[condition] = lc.conditionSet[condition] + 1
	lc.Unlock()
	return lc.conditionSet[condition], nil
}

func (lc *LocalCache) Size() (int64, error) {
	return int64(len(lc.conditionSet)), nil
}
