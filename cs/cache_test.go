package cs

import (
	"fmt"
	"log"
	"testing"

	"github.com/john-k-ge/homunculus/cf"
)

const (
	host      = "0.0.0.0"
	port      = "6379"
	pass      = ""
	index     = 2
	condition = "DB_Errors"
	count     = 0
	bogusCond = "BOGUS"
)

func buildCaches() ([]ConditionSet, error) {
	var caches = []ConditionSet{}

	env := cf.CfEnv{
		RHost:   host,
		RPort:   port,
		RPasswd: pass,
		Inx:     index,
	}

	rc, err := NewRemoteCondtionSet(&env)
	if err != nil {
		return caches, err
	}
	lc := NewLocalConditionSet()
	caches = append(caches, lc)
	caches = append(caches, rc)
	return caches, nil
}

func TestNewRemoteCache(t *testing.T) {
	env := cf.CfEnv{
		RHost:   host,
		RPort:   port,
		RPasswd: pass,
		Inx:     index,
	}
	_, err := NewRemoteCondtionSet(&env)
	if err != nil {
		log.Printf("Failed to create new cache(%v, %v, %v): %v", host, port, pass, err)
		t.Fail()
	}
}

func TestSaveCondition(t *testing.T) {
	conds, err := buildCaches()
	if err != nil {
		log.Printf("Failed to build caches: %v", err)
		t.Fail()
	}

	for i, c := range conds {
		t.Run(fmt.Sprintf("Caching condition: #%v", i), func(t *testing.T) {
			err := c.SaveCondition(condition, count)
			if err != nil {
				log.Printf("Error encountered caching (%v, %v): %v", condition, count, err)
				t.Fail()
			}
		})
		t.Run(fmt.Sprintf("bad condition: #%v", i), func(t *testing.T) {
			err := c.SaveCondition("", 200)
			if err == nil {
				log.Printf("Test should have failed with empty key")
				t.Fail()
			}
			log.Printf("Bad key rejected: %v", err)
		})
	}
}

func TestCheckCondition(t *testing.T) {
	conds, err := buildCaches()
	if err != nil {
		log.Printf("Failed to build caches: %v", err)
		t.Fail()
	}

	for _, c := range conds {
		_ = c.SaveCondition(condition, count)
		t.Run("Fetching condition", func(t *testing.T) {
			val, err := c.CheckCondition(condition)
			if err != nil {
				log.Printf("Failed getting val for %v: %v", condition, err)
				t.Fail()
			}
			if val != count {
				log.Printf("Fetched cond %v does not match %v", val, condition)
				t.Fail()
			}
		})
		t.Run("Non-existent key", func(t *testing.T) {
			_, err := c.CheckCondition(bogusCond)
			if err == nil {
				log.Printf("Test should have failed with bogus key")
				t.Fail()
			}
		})
	}
}

func TestConditionExists(t *testing.T) {
	conds, err := buildCaches()
	if err != nil {
		log.Printf("Failed to build caches: %v", err)
		t.Fail()
	}

	for _, c := range conds {
		_ = c.SaveCondition(condition, count)
		t.Run("Condition exists", func(t *testing.T) {
			exists, err := c.ConditionExists(condition)
			if err != nil {
				log.Printf("Failed to check exists for %v: %v", condition, err)
				t.Fail()
			}
			if !exists {
				log.Printf("%v should exist but doesn't", condition)
				t.Fail()
			}
			exists, err = c.ConditionExists(bogusCond)
			if err != nil {
				log.Printf("Failed to check exists for %v: %v", bogusCond, err)
				t.Fail()
			}
			if exists {
				log.Printf("%v shouldn't exist but does", bogusCond)
				t.Fail()
			}
		})
	}
}

func TestDeleteCondition(t *testing.T) {
	conds, err := buildCaches()
	if err != nil {
		log.Printf("Failed to build caches: %v", err)
		t.Fail()
	}

	for _, c := range conds {
		_ = c.SaveCondition(condition, count)
		t.Run("Deleting condition", func(t *testing.T) {
			before, _ := c.ConditionExists(condition)
			c.DeleteCondition(condition)
			after, _ := c.ConditionExists(condition)
			if !before || after {
				log.Printf("Before: %v, after %v", before, after)
				t.Fail()
			}
		})
	}
}

func TestIncrementCondition(t *testing.T) {
	conds, err := buildCaches()
	if err != nil {
		log.Printf("Failed to build caches: %v", err)
		t.Fail()
	}

	for _, c := range conds {
		_ = c.SaveCondition(condition, count)
		t.Run("Incrementing condition", func(t *testing.T) {
			before, _ := c.CheckCondition(condition)
			c.IncrementCondition(condition)
			after, _ := c.CheckCondition(condition)
			log.Printf("Before: %v, After: %v", before, after)
			if after != (before + 1) {
				t.Fail()
			}
		})
	}
}

func TestClearCache(t *testing.T) {
	conds, err := buildCaches()
	if err != nil {
		log.Printf("Failed to build caches: %v", err)
		t.Fail()
	}

	for _, c := range conds {
		_ = c.SaveCondition(condition, count)
		t.Run("Testing flushall", func(t *testing.T) {
			err := c.ClearConditions()
			if err != nil {
				log.Printf("Failed to clear cache: %v", err)
				t.Fail()
			}
		})
	}
}
