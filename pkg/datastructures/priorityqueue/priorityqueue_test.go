// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package priorityqueue_test

import (
	"container/heap"
	"encoding/json"
	"github.com/ryandgoulding/godatastructures/pkg/datastructures/priorityqueue"
	"os"
	"testing"
)

var testCases []*testCase


type testCase struct {
	description string
	pq               *priorityqueue.PriorityQueue
	raw              map[string]float64
	expectedPopOrder *[]string
}

func (t *testCase) getExpectedLength() int {
	return len(t.raw)
}

func generatePriorityQueue(raw map[string]float64) *priorityqueue.PriorityQueue {
	pq := priorityqueue.PriorityQueue{}
	for value, priority := range raw {
		item := &priorityqueue.Item{
			Value: value,
			Priority: priority,
		}
		heap.Push(&pq, item)
	}
	return &pq
}

func createTestCase(description string, raw map[string]float64, expectedPopOrder *[]string) *testCase {
	pq := generatePriorityQueue(raw)
	testCase := &testCase{
		description: description,
		pq: pq,
		raw: raw,
		expectedPopOrder: expectedPopOrder,
	}
	return testCase
}

func setupEmptyPriorityQueueTestCase() {
	description := "Empty PriorityQueue"
	raw := map[string]float64{}
	expectedPopOrder := &[]string{}
	testCase := createTestCase(description, raw, expectedPopOrder)
	testCases = append(testCases, testCase)
}

func setupSmallPriorityQueueTestCase() {
	description := "Small PriorityQueue"
	raw := map[string]float64{"apple": 10.0, "banana": 5.0, "carrot": 11.0, "danish": 0.0}
	expectedPopOrder := &[]string{"carrot", "apple", "banana", "danish"}
	testCase := createTestCase(description, raw, expectedPopOrder)
	testCases = append(testCases, testCase)
}

func setupLargePriorityQueueTestCase() {
	description := "Large PriorityQueue"
	raw := map[string]float64{}
	expectedPopOrder := []string{}
	for i := 0; i < 1000; i++ {
		str := string(i)
		expectedPopOrder = append([]string{str}, expectedPopOrder...)
		val := float64(i)
		raw[str] = val
	}
	testCase := createTestCase(description, raw, &expectedPopOrder)
	testCases = append(testCases, testCase)
}

func setup() {
	setupEmptyPriorityQueueTestCase()
	setupSmallPriorityQueueTestCase()
	setupLargePriorityQueueTestCase()
}

func TestPriorityQueue_Len(t *testing.T) {
	for _, testCase := range testCases {
		expectedLength := testCase.getExpectedLength()
		actualLength := testCase.pq.Len()
		if expectedLength != actualLength {
			t.Fatalf("%s: Expected: %d Actual: %d", testCase.description, expectedLength, actualLength)
		}
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	for _, testCase := range testCases {
		pq := testCase.pq
		for i := 0; i < pq.Len(); i++ {
			actual := heap.Pop(pq).(*priorityqueue.Item).Value
			expected := (*testCase.expectedPopOrder)[i]
			if expected != actual {
				t.Fatalf("%s: Expected: %s Actual %s", testCase.description, expected, actual)
			}
		}
	}
}

func TestPriorityQueue_MarshalJSON(t *testing.T) {
	// Simple re-entrance test.  MarshalJSON is destructive in this instance, since popping from a PriorityQueue is
	// destructive.  Ensure that subsequent calls return the same result.
	for _, testCase := range testCases {
		firstJsonBytes, err := json.MarshalIndent(testCase.pq, "", "  ")
		if err != nil {
			t.Fatalf("Unexpected error marshaling JSON: %s", err)
		}
		firstJsonString := string(firstJsonBytes)
		secondJsonBytes, err := json.MarshalIndent(testCase.pq, "", "  ")
		if err != nil {
			t.Fatalf("Unexpected error marshaling JSON: %s", err)
		}
		secondJsonString := string(secondJsonBytes)
		if firstJsonString != secondJsonString {
			t.Fatalf("Expected a match for: %s %s", firstJsonString, secondJsonString)
		}
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}
