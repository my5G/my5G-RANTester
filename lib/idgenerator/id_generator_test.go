package idgenerator_test

import (
	"fmt"
	"math/rand"
	"my5G-RANTester/lib/idgenerator"
	"sync"
	"testing"
	"time"
)

func TestAllocate(t *testing.T) {
	testCases := []struct {
		minValue int64
		maxValue int64
	}{
		{1, 20},
		{11, 50},
		{1, 12345678},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("minValue: %d, maxValue: %d", testCase.minValue, testCase.maxValue), func(t *testing.T) {
			idGenerator := idgenerator.NewGenerator(testCase.minValue, testCase.maxValue)

			for i := testCase.minValue; i <= testCase.maxValue; i++ {
				id, err := idGenerator.Allocate()
				if id != i {
					t.Errorf("expected id: %d, output id: %d", i, id)
					t.FailNow()
				} else if err != nil {
					t.Error(err)
					t.FailNow()
				}
			}

			for i := testCase.minValue; i <= testCase.maxValue; i++ {
				idGenerator.FreeID(i)
			}
		})
	}
}

func TestConcurrency(t *testing.T) {
	var usedMap sync.Map

	idGenerator := idgenerator.NewGenerator(1, 12345678)

	wg := sync.WaitGroup{}
	for routineID := 1; routineID <= 10; routineID++ {
		wg.Add(1)
		go func(routineID int) {
			for i := 0; i < 1000; i++ {
				id, _ := idGenerator.Allocate()
				if value, ok := usedMap.Load(id); ok {
					routineID := value.(int)
					t.Errorf("ID %d has been allocated at routine[%d], concurrent test failed", id, routineID)
				} else {
					usedMap.Store(id, routineID)
				}
			}
			usedMap.Range(func(key, value interface{}) bool {
				id := key.(int64)
				idGenerator.FreeID(id)
				return true
			})
			wg.Done()
		}(routineID)
	}
	wg.Wait()
}

func TestUnique(t *testing.T) {
	testCases := []struct {
		minValue int64
		maxValue int64
	}{
		{1, 10},
		{11, 1567},
	}

	rand.Seed(time.Now().Unix())

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("minValue: %d, maxValue: %d", testCase.minValue, testCase.maxValue), func(t *testing.T) {
			usedMap := make(map[int64]bool)

			valueRange := testCase.maxValue - testCase.minValue + 1
			testRange := int(valueRange * 3)
			idGenerator := idgenerator.NewGenerator(testCase.minValue, testCase.maxValue)

			for i := 0; i < testRange; i++ {
				id, err := idGenerator.Allocate()
				if err != nil {
					t.Error(err)
					t.FailNow()
				}

				if _, ok := usedMap[id]; ok {
					t.Errorf("ID %d has been allocated, test failed", id)
					t.FailNow()
				} else {
					usedMap[id] = true
				}

				// do one free operation from the beginning of second round
				if i >= int(valueRange-1) {
					// retrieve all id to keys
					var keys []int64
					for key := range usedMap {
						keys = append(keys, key)
					}

					keyIdx := rand.Intn(len(keys))
					idToFree := keys[keyIdx]
					idGenerator.FreeID(idToFree)
					delete(usedMap, idToFree)
				}
			}
		})
	}
}

func TestTriggerNoSpaceToAllocateError(t *testing.T) {
	testCases := []struct {
		minValue int64
		maxValue int64
	}{
		{1, 10},
	}

	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("minValue: %d, maxValue: %d", testCase.minValue, testCase.maxValue), func(t *testing.T) {

			valueRange := int(testCase.maxValue - testCase.minValue + 1)
			idGenerator := idgenerator.NewGenerator(testCase.minValue, testCase.maxValue)

			for i := 0; i < valueRange; i++ {
				_, err := idGenerator.Allocate()
				if err != nil {
					t.Error(err)
					t.FailNow()
				}
			}

			// trigger "No available value range to allocate id" error
			_, err := idGenerator.Allocate()
			if err == nil {
				t.Error("expect return error, but error is nil")
				t.FailNow()
			}
		})
	}
}
