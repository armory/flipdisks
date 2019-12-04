package snake

import (
	"container/heap"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	type action struct {
		command string
		pushArg *autoPlayNode
		popWant *autoPlayNode
	}
	tests := []struct {
		name    string
		pq      priorityQueue
		actions []action
	}{
		{
			name: "priority queue works as expected",
			pq:   []*autoPlayNode{},
			actions: []action{
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 10}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 7}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 20}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 30}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 5}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 10}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 20}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 30}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 6}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 30}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 20}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 8}},
				{command: "push", pushArg: &autoPlayNode{heuristicDistance: 10}},

				{command: "pop", popWant: &autoPlayNode{heuristicDistance: 5}},
				{command: "pop", popWant: &autoPlayNode{heuristicDistance: 6}},
				{command: "pop", popWant: &autoPlayNode{heuristicDistance: 7}},
				{command: "pop", popWant: &autoPlayNode{heuristicDistance: 8}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			heap.Init(&tt.pq)
			for _, a := range tt.actions {
				switch a.command {
				case "push":
					heap.Push(&tt.pq, a.pushArg)
				case "pop":
					got := heap.Pop(&tt.pq).(*autoPlayNode)
					assert.Equal(t, a.popWant.heuristicDistance, got.heuristicDistance)
				}
			}
		})
	}
}
