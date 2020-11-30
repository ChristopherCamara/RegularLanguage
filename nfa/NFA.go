package nfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLangauge/internal/intArray"
)

type NFA struct {
	stateCount         int
	Alphabet           []string
	states             []int
	startStates        []int
	acceptStates       []int
	transitions        map[int]map[string][]int
	epsilonTransitions map[int][]int
}

func New() *NFA {
	newNFA := new(NFA)
	newNFA.stateCount = 0
	newNFA.Alphabet = make([]string, 0)
	newNFA.states = make([]int, 0)
	newNFA.startStates = make([]int, 0)
	newNFA.acceptStates = make([]int, 0)
	newNFA.transitions = make(map[int]map[string][]int, 0)
	newNFA.epsilonTransitions = make(map[int][]int, 0)
	return newNFA
}

func (nfa *NFA) removeState(removeState int) {
	intArray.Remove(removeState, &nfa.states)
	intArray.Remove(removeState, &nfa.startStates)
	intArray.Remove(removeState, &nfa.acceptStates)
	delete(nfa.transitions, removeState)
	delete(nfa.epsilonTransitions, removeState)
	for i := 0; i < len(nfa.states); i++ {
		currentState := nfa.states[i]
		for symbol, transitions := range nfa.transitions[currentState] {
			for index, transition := range transitions {
				if transition > removeState {
					nfa.transitions[currentState][symbol][index] = transition - 1
				}
			}
		}
		for index, transition := range nfa.epsilonTransitions[currentState] {
			if transition > removeState {
				nfa.epsilonTransitions[currentState][index] = transition - 1
			}
		}
		if currentState > removeState {
			nfa.states[i]--
			nfa.transitions[currentState-1] = nfa.transitions[currentState]
			nfa.epsilonTransitions[currentState-1] = nfa.epsilonTransitions[currentState]
		}
	}
	for i := 0; i < len(nfa.startStates); i++ {
		if nfa.startStates[i] > removeState {
			nfa.startStates[i]--
		}
	}
	for i := 0; i < len(nfa.acceptStates); i++ {
		if nfa.acceptStates[i] > removeState {
			nfa.acceptStates[i]--
		}
	}
}

func (nfa *NFA) addState(isStart, isAccept bool) int {
	index := nfa.stateCount
	nfa.states = append(nfa.states, index)
	nfa.transitions[index] = make(map[string][]int, 0)
	nfa.epsilonTransitions[index] = make([]int, 0)
	if isStart {
		nfa.startStates = append(nfa.startStates, index)
	}
	if isAccept {
		nfa.acceptStates = append(nfa.acceptStates, index)
	}
	nfa.stateCount++
	return index
}

func (nfa *NFA) addEpsilonTransition(sourceState, targetState int) {
	nfa.epsilonTransitions[sourceState] = append(nfa.epsilonTransitions[sourceState], targetState)
}

func (nfa *NFA) addTransition(sourceState int, symbol string, targetState int) {
	nfa.transitions[sourceState][symbol] = append(nfa.transitions[sourceState][symbol], targetState)
}

func (nfa *NFA) merge(other *NFA) map[int]int {
	//map otherStates to new states in nfa
	newStates := make(map[int]int, 0)
	for _, otherState := range other.states {
		newStates[otherState] = nfa.addState(false, false)
	}
	//update epsilon transitions to new states
	for otherState, transitionStates := range other.epsilonTransitions {
		for _, transitionState := range transitionStates {
			nfa.addEpsilonTransition(newStates[otherState], newStates[transitionState])
		}
	}
	//update transitions to new states
	for otherState, transitions := range other.transitions {
		for symbol, transitionStates := range transitions {
			for _, transitionState := range transitionStates {
				nfa.addTransition(newStates[otherState], symbol, newStates[transitionState])
			}
		}
	}
	return newStates
}

func (nfa *NFA) Concat(other *NFA) {
	newStates := nfa.merge(other)
	//move all transitions from other startStates to our acceptStates
	for _, otherStart := range other.startStates {
		for _, otherTransition := range other.epsilonTransitions[otherStart] {
			for _, acceptState := range nfa.acceptStates {
				nfa.addEpsilonTransition(acceptState, newStates[otherTransition])
			}
		}
		for symbol, otherTransitions := range other.transitions[otherStart] {
			for _, otherTransition := range otherTransitions {
				for _, acceptState := range nfa.acceptStates {
					nfa.addTransition(acceptState, symbol, newStates[otherTransition])
				}
			}
		}
		//change our accept states to the other accept states
		nfa.acceptStates = make([]int, 0)
		for _, otherAccept := range other.acceptStates {
			nfa.acceptStates = append(nfa.acceptStates, newStates[otherAccept])
		}
		//remove the otherStartStates now that its ported over
		for _, otherStart := range other.startStates {
			nfa.removeState(newStates[otherStart])
		}
	}
}

func (nfa *NFA) Union(other *NFA) {
	newNFA := New()
	newStart := newNFA.addState(true, false)
	newStates := newNFA.merge(nfa)
	newOtherStates := newNFA.merge(other)
	newAccept := newNFA.addState(false, true)
	for _, startState := range nfa.startStates {
		newNFA.addEpsilonTransition(newStart, newStates[startState])
	}
	for _, otherStart := range other.startStates {
		newNFA.addEpsilonTransition(newStart, newOtherStates[otherStart])
	}
	for _, acceptState := range nfa.acceptStates {
		newNFA.addEpsilonTransition(newStates[acceptState], newAccept)
	}
	for _, otherAccept := range other.acceptStates {
		newNFA.addEpsilonTransition(newOtherStates[otherAccept], newAccept)
	}
	*nfa = *newNFA
}

func (nfa *NFA) Closure() {
	newNFA := New()
	newStart := newNFA.addState(true, false)
	newStates := newNFA.merge(nfa)
	newAccept := newNFA.addState(false, true)
	newNFA.addEpsilonTransition(newStart, newAccept)
	for _, startState := range nfa.startStates {
		newNFA.addEpsilonTransition(newStart, newStates[startState])
	}
	for _, acceptState := range nfa.acceptStates {
		newNFA.addEpsilonTransition(newStates[acceptState], newAccept)
		for _, startState := range nfa.startStates {
			newNFA.addEpsilonTransition(newStates[acceptState], newStates[startState])
		}
	}
	*nfa = *newNFA
}

func (nfa *NFA) Print() {
	fmt.Println("~~~NFA~~~")
	fmt.Print("start states: ")
	intArray.Print(nfa.startStates)
	for _, state := range nfa.states {
		fmt.Printf("state %d:\n", state)
		if nfa.transitions[state] != nil {
			for symbol, states := range nfa.transitions[state] {
				fmt.Printf("\t%s -> ", symbol)
				intArray.Print(states)
			}
		}
		if len(nfa.epsilonTransitions[state]) != 0 {
			fmt.Print("\t(empty) -> ")
			intArray.Print(nfa.epsilonTransitions[state])
		}
	}
	fmt.Print("accept states: ")
	intArray.Print(nfa.acceptStates)
}
