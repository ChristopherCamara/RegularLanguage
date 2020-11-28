package dfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLangauge/internal/intArray"
	"github.com/ChristopherCamara/RegularLangauge/nfa"
)

type State struct {
	Index      int
	IsEnd      bool
	Transition map[string]*State
	Closure    []int
}

func CreateState(isEnd bool) *State {
	newState := new(State)
	newState.IsEnd = isEnd
	newState.Index = -1
	newState.Transition = make(map[string]*State)
	newState.Closure = make([]int, 0)
	return newState
}

func (s *State) print() {
	fmt.Printf("State %d:\n", s.Index)
	if s.IsEnd {
		fmt.Println("\tIS AN END STATE")
	}
	for symbol, nextState := range s.Transition {
		fmt.Printf("\t%s -> %d\n", symbol, nextState.Index)
	}
}

func (s *State) copyOf(target *nfa.State) {
	s.Index = target.Index
	s.IsEnd = target.IsEnd
	s.Closure = []int{}
	s.Closure = append(s.Closure, target.Closure...)
	for symbol, targetState := range target.Transition {
		copyState := CreateState(false)
		copyState.Index = targetState.Index
		copyState.IsEnd = targetState.IsEnd
		copyState.Closure = []int{}
		copyState.Closure = append(copyState.Closure, targetState.Closure...)
		s.Transition[symbol] = copyState
	}
}

func Merge(targetState *State, fromState *nfa.State) *State {
	targetState.Index = fromState.Index
	targetState.IsEnd = fromState.IsEnd
	targetState.Closure = []int{}
	targetState.Closure = append(targetState.Closure, fromState.Closure...)
	for symbol, transitionState := range fromState.Transition {
		if targetState.Transition[symbol] != nil {
			targetState.Transition[symbol] = Merge(targetState.Transition[symbol], transitionState)
		} else {
			targetState.Transition[symbol] = Merge(CreateState(false), transitionState)
		}
	}
	return targetState
}

func findStateByIndex(rootState *State, index int) *State {
	queue := []*State{rootState}
	visited := []int{rootState.Index}
	currentState := queue[0]
	for currentState != nil {
		if currentState.Index == index {
			return currentState
		}
		for _, nextState := range currentState.Transition {
			if intArray.IndexOf(nextState.Index, visited) == -1 {
				queue = append(queue, nextState)
				visited = append(visited, nextState.Index)
			}
		}
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = queue[0]
		} else {
			currentState = nil
		}
	}
	return currentState
}

func distinguishable(first, second *State, otherPartition []int) bool {
	for symbol, targetState := range first.Transition {
		if second.Transition[symbol] == nil {
			continue
		}
		if intArray.IndexOf(targetState.Index, otherPartition) != -1 && intArray.IndexOf(second.Transition[symbol].Index, otherPartition) == -1 {
			return true
		}
		if intArray.IndexOf(targetState.Index, otherPartition) == -1 && intArray.IndexOf(second.Transition[symbol].Index, otherPartition) == 1 {
			return true
		}
	}
	for symbol, targetState := range second.Transition {
		if first.Transition[symbol] == nil {
			continue
		}
		if intArray.IndexOf(targetState.Index, otherPartition) != -1 && intArray.IndexOf(first.Transition[symbol].Index, otherPartition) == -1 {
			return true
		}
		if intArray.IndexOf(targetState.Index, otherPartition) == -1 && intArray.IndexOf(first.Transition[symbol].Index, otherPartition) != -1 {
			return true
		}
	}
	return false
}
