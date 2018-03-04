package cs

import (
	"errors"
	"log"
)

type LocalCache struct {
	conditionSet map[string]int64
}

func NewLocalConditionSet() *LocalCache {
	lc := LocalCache{
		conditionSet: make(map[string]int64),
	}
	return &lc
}

func (lc *LocalCache) SaveCondition(condition string, count int64) error {
	if len(condition) == 0 {
		log.Print("Key must have a value")
		return errors.New("key must have a value")
	}
	lc.conditionSet[condition] = count
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
	delete(lc.conditionSet, condition)
	return nil
}

func (lc *LocalCache) ClearConditions() error {
	lc.conditionSet = make(map[string]int64)
	return nil
}

func (lc *LocalCache) IncrementCondition(condition string) (int64, error) {
	_, exists := lc.conditionSet[condition]
	if !exists {
		return 0, errors.New("condition " + condition + " not found")
	}
	lc.conditionSet[condition] = lc.conditionSet[condition] + 1
	return lc.conditionSet[condition], nil
}

func (lc *LocalCache) Size() (int64, error) {
	return int64(len(lc.conditionSet)), nil
}
