package nfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLangauge/dfa"
	"github.com/ChristopherCamara/RegularLangauge/internal/intArray"
)

type NFA struct {
	RootState *State
	EndState  *State
	Alphabet  []string
}

func New(rootState, endState *State) *NFA {
	nfa := new(NFA)
	nfa.RootState = rootState
	nfa.EndState = endState
	return nfa
}

func (nfa *NFA) AssignStateIndices() {
	index := 1
	nfa.RootState.Index = 0
	queue := []*State{nfa.RootState}
	currentState := queue[0]
	for currentState != nil {
		for _, nextState := range currentState.EpsilonTransition {
			if nextState.Index == -1 {
				nextState.Index = index
				index++
				queue = append(queue, nextState)
			}
		}
		for _, nextState := range currentState.Transition {
			if nextState.Index == -1 {
				nextState.Index = index
				index++
				queue = append(queue, nextState)
			}
		}
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = queue[0]
		} else {
			currentState = nil
		}
	}
}

func (nfa *NFA) Print() {
	fmt.Println("~~~nfa~~~")
	fmt.Printf("start at state %d\n", nfa.RootState.Index)
	visited := []int{nfa.RootState.Index}
	queue := []*State{nfa.RootState}
	currentState := queue[0]
	for currentState != nil {
		for _, nextState := range currentState.EpsilonTransition {
			if intArray.IndexOf(nextState.Index, visited) == -1 {
				queue = append(queue, nextState)
				visited = append(visited, nextState.Index)
			}
		}
		for _, nextState := range currentState.Transition {
			if intArray.IndexOf(nextState.Index, visited) == -1 {
				queue = append(queue, nextState)
				visited = append(visited, nextState.Index)
			}
		}
		currentState.print()
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = queue[0]
		} else {
			currentState = nil
		}
	}
}

func ToDFA(nfa *NFA) *dfa.DFA {
	newDFA := new(dfa.DFA)
	newDFA.RootState = dfa.CreateState(false)
	rootClosure, rootVisited := make([]int, 0), make([]int, 0)
	epsilonClosure(nfa.RootState, &rootClosure, &rootVisited)
	nfa.RootState.Closure = append(nfa.RootState.Closure, rootClosure...)
	//dfa.RootState.copyOf(nfa.RootState)
	newDFA.RootState = Merge(newDFA.RootState, nfa.RootState)
	currentState := newDFA.RootState
	visited := []int{newDFA.RootState.Index}
	queue := []*dfa.State{newDFA.RootState}
	for currentState != nil {
		for _, closureIndex := range currentState.Closure {
			closureState := findStateByIndex(nfa.RootState, closureIndex)
			if closureState.IsEnd {
				currentState.IsEnd = true
			}
			for _, targetState := range closureState.Transition {
				//copyState := createDFAState(false)
				//copyState.copyOf(targetState)
				//currentState.Transition[symbol] = copyState
				currentState = Merge(currentState, targetState)
			}
		}
		for _, targetState := range currentState.Transition {
			if intArray.IndexOf(targetState.Index, visited) == -1 {
				visited = append(visited, targetState.Index)
				queue = append(queue, targetState)
			}
		}
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = queue[0]
			closure, closureVisited := make([]int, 0), make([]int, 0)
			epsilonClosure(findStateByIndex(nfa.RootState, currentState.Index), &closure, &closureVisited)
			currentState.Closure = append(currentState.Closure, closure...)
		} else {
			currentState = nil
		}
	}
	return newDFA
}
