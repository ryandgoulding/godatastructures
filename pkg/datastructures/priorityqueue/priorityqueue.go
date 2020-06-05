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

/*
Package priorityqueue implements a generic PriorityQueue.  This implementation was adapted from:
https://golang.org/pkg/container/heap/.
*/

package priorityqueue

import (
	"bytes"
	"container/heap"
	"encoding/json"
)

// An Item is something we manage in a Priority queue.
type Item struct {
	Value    interface{} // The Value of the item; arbitrary.
	Priority float64     // The Priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index    int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, Priority so we use greater than here.
	return pq[i].Priority > pq[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Marshal a PriorityQueue in priorityqueue order.  Warning, this method is not terribly efficient, as iterating over a
// heap-based PriorityQueue is destructive.  Thus, O(n) auxillary space is required to store the item references and
// O(n) time complexity is needed to re-construct the priority queue post destruction.  There are likely more efficient
// implementations, but in this case n is expected to remain sufficiently small, so this implementation is "good
// enough".
func (pq *PriorityQueue) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("[")
	pqLen := pq.Len()
	var pqCopy []*Item
	for i := 0; i < pqLen; i++ {
		item := heap.Pop(pq).(*Item)
		json, err := json.Marshal(*item)
		if err != nil {
			return nil, err
		}
        buffer.WriteString(string(json))
		if i < pqLen - 1 {
			buffer.WriteByte(',')
		}
		pqCopy = append(pqCopy, item)
	}
	buffer.WriteString("]")
	for _, item := range pqCopy {
		pq.Push(item)
	}
	return buffer.Bytes(), nil
}
