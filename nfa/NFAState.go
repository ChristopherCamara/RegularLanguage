package nfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLangauge/dfa"
	"github.com/ChristopherCamara/RegularLangauge/internal/intArray"
)

type State struct {
	Index             int
	IsEnd             bool
	Transition        map[string]*State
	EpsilonTransition []*State
	Closure           []int
}

func createState(isEnd bool) *State {
	newState := new(State)
	newState.IsEnd = isEnd
	newState.Index = -1
	newState.Transition = make(map[string]*State)
	newState.EpsilonTransition = make([]*State, 0)
	newState.Closure = make([]int, 0)
	return newState
}

func (s *State) print() {
	fmt.Printf("state %d:\n", s.Index)
	if s.IsEnd {
		fmt.Println("\tIS AN END state")
	}
	for symbol, nextState := range s.Transition {
		fmt.Printf("\t%s -> %d\n", symbol, nextState.Index)
	}
	for _, nextState := range s.EpsilonTransition {
		fmt.Printf("\t(empty) -> %d\n", nextState.Index)
	}
}

func (s *State) addEpsilonTransition(toState *State) {
	s.EpsilonTransition = append(s.EpsilonTransition, toState)
}

func (s *State) addTransition(toState *State, symbol string) {
	s.Transition[symbol] = toState
}

func findStateByIndex(rootState *State, index int) *State {
	queue := []*State{rootState}
	visited := []int{rootState.Index}
	currentState := queue[0]
	for currentState != nil {
		if currentState.Index == index {
			return currentState
		}
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
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = queue[0]
		} else {
			currentState = nil
		}
	}
	return currentState
}

func epsilonClosure(startState *State, closure *[]int, visited *[]int) {
	*closure = append(*closure, startState.Index)
	if len(startState.EpsilonTransition) != 0 {
		for _, nextState := range startState.EpsilonTransition {
			if intArray.IndexOf(nextState.Index, *visited) == -1 {
				*visited = append(*visited, nextState.Index)
				epsilonClosure(nextState, closure, visited)
			}
		}
	}
}

func copyOf(s *dfa.State, target *State) {
	s.Index = target.Index
	s.IsEnd = target.IsEnd
	s.Closure = []int{}
	s.Closure = append(s.Closure, target.Closure...)
	for symbol, targetState := range target.Transition {
		copyState := dfa.CreateState(false)
		copyState.Index = targetState.Index
		copyState.IsEnd = targetState.IsEnd
		copyState.Closure = []int{}
		copyState.Closure = append(copyState.Closure, targetState.Closure...)
		s.Transition[symbol] = copyState
	}
}

func Merge(targetState *dfa.State, fromState *State) *dfa.State {
	targetState.Index = fromState.Index
	targetState.IsEnd = fromState.IsEnd
	targetState.Closure = []int{}
	targetState.Closure = append(targetState.Closure, fromState.Closure...)
	for symbol, transitionState := range fromState.Transition {
		if targetState.Transition[symbol] != nil {
			targetState.Transition[symbol] = Merge(targetState.Transition[symbol], transitionState)
		} else {
			targetState.Transition[symbol] = Merge(dfa.CreateState(false), transitionState)
		}
	}
	return targetState
}
