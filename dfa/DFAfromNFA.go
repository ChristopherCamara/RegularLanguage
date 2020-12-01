package dfa

import (
	"github.com/ChristopherCamara/RegularLangauge/internal/intArray"
	"github.com/ChristopherCamara/RegularLangauge/nfa"
)

func FromNFA(NFA *nfa.NFA) *DFA {
	epsilonClosures := NFA.GetEpsilonClosures()
	//collapsedIndex := 0
	collapsedStates := make(map[int][]int, 0)
	collapsedTransitions := make(map[int]map[string]int, 0)
	collapsedStates[0] = epsilonClosures[NFA.StartStates[0]]
	queue := []int{0}
	visited := []int{0}
	currentState := collapsedStates[0]
	currentCollapsedIndex := queue[0]
	for currentState != nil {
		for _, symbol := range NFA.Alphabet {
			transitionStates := make([]int, 0)
			for _, state := range currentState {
				if NFA.Transitions[state][symbol] != nil {
					for _, transitionState := range NFA.Transitions[state][symbol] {
						for _, closureState := range epsilonClosures[transitionState] {
							if intArray.IndexOf(closureState, transitionStates) == -1 {
								transitionStates = append(transitionStates, closureState)
							}
						}
						if intArray.IndexOf(transitionState, transitionStates) == -1 {
							transitionStates = append(transitionStates, transitionState)
						}
					}
				}
			}
			if len(transitionStates) != 0 {
				collapsed := len(collapsedStates)
				for index, currentCollapsed := range collapsedStates {
					if intArray.Equals(currentCollapsed, transitionStates) {
						collapsed = index
					}
				}
				//collapsed is either a new entry for the map collapsedStates or an existing entry in collapsedStates
				collapsedStates[collapsed] = transitionStates
				if collapsedTransitions[currentCollapsedIndex] == nil {
					collapsedTransitions[currentCollapsedIndex] = make(map[string]int, 0)
				}
				collapsedTransitions[currentCollapsedIndex][symbol] = collapsed
				if intArray.IndexOf(collapsed, visited) == -1 {
					queue = append(queue, collapsed)
					visited = append(visited, collapsed)
				}
			}
		}
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = collapsedStates[queue[0]]
			currentCollapsedIndex = queue[0]
		} else {
			currentState = nil
		}
	}
	dfa := New()
	dfa.Alphabet = NFA.Alphabet
	for i := 0; i < len(collapsedStates); i++ {
		newState := dfa.addState(false, false)
		for _, state := range collapsedStates[i] {
			if intArray.IndexOf(state, NFA.StartStates) != -1 {
				if intArray.IndexOf(newState, dfa.StartStates) == -1 {
					dfa.StartStates = append(dfa.StartStates, newState)
				}
			}
			if intArray.IndexOf(state, NFA.AcceptStates) != -1 {
				if intArray.IndexOf(newState, dfa.AcceptStates) == -1 {
					dfa.AcceptStates = append(dfa.AcceptStates, newState)
				}
			}
		}
		for symbol, transition := range collapsedTransitions[i] {
			dfa.AddTransition(newState, symbol, transition)
		}
	}
	return dfa
}
