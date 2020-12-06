package dfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLanguage/internal/intArray"
)

type DFA struct {
	nextState    int
	Alphabet     []string
	States       []int
	StartStates  []int
	AcceptStates []int
	Transitions  map[int]map[string]int
}

func New() *DFA {
	newDFA := new(DFA)
	newDFA.nextState = 0
	newDFA.Alphabet = make([]string, 0)
	newDFA.States = make([]int, 0)
	newDFA.StartStates = make([]int, 0)
	newDFA.AcceptStates = make([]int, 0)
	newDFA.Transitions = make(map[int]map[string]int, 0)
	return newDFA
}

func (dfa *DFA) addState(isStart, isAccept bool) int {
	index := dfa.nextState
	dfa.States = append(dfa.States, index)
	dfa.Transitions[index] = make(map[string]int, 0)
	if isStart {
		dfa.StartStates = append(dfa.StartStates, index)
	}
	if isAccept {
		dfa.AcceptStates = append(dfa.AcceptStates, index)
	}
	dfa.nextState++
	return index
}

func (dfa *DFA) AddTransition(sourceState int, symbol string, targetState int) {
	dfa.Transitions[sourceState][symbol] = targetState
}

func (dfa *DFA) Print() {
	fmt.Println("~~~DFA~~~")
	fmt.Print("start states: ")
	intArray.Print(dfa.StartStates)
	for _, state := range dfa.States {
		fmt.Printf("state %d:\n", state)
		if dfa.Transitions[state] != nil {
			for symbol, state := range dfa.Transitions[state] {
				fmt.Printf("\t%s -> %d\n", symbol, state)
			}
		}
	}
	fmt.Print("accept states: ")
	intArray.Print(dfa.AcceptStates)
}
