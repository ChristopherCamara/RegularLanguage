package dfa

import (
	"github.com/ChristopherCamara/RegularLanguage/internal/intArray"
)

func (dfa *DFA) distinguishable(first, second int, otherPartition []int) bool {
	for symbol, targetState := range dfa.Transitions[first] {
		if _, exists := dfa.Transitions[second][symbol]; !exists {
			continue
		}
		if intArray.IndexOf(targetState, otherPartition) != -1 && intArray.IndexOf(dfa.Transitions[second][symbol], otherPartition) == -1 {
			return true
		}
		if intArray.IndexOf(targetState, otherPartition) == -1 && intArray.IndexOf(dfa.Transitions[second][symbol], otherPartition) != -1 {
			return true
		}
	}
	for symbol, targetState := range dfa.Transitions[second] {
		if _, exists := dfa.Transitions[first][symbol]; exists {
			continue
		}
		if intArray.IndexOf(targetState, otherPartition) != -1 && intArray.IndexOf(dfa.Transitions[first][symbol], otherPartition) == -1 {
			return true
		}
		if intArray.IndexOf(targetState, otherPartition) == -1 && intArray.IndexOf(dfa.Transitions[first][symbol], otherPartition) != -1 {
			return true
		}
	}
	return false
}

func (dfa *DFA) ToMinDFA() {
	sinkState := -1
	statePartitions := make([][]int, 0)
	statePartitions = append(statePartitions, make([]int, 0))
	statePartitions = append(statePartitions, make([]int, 0))
	queue := []int{dfa.StartStates[0]}
	visited := []int{dfa.StartStates[0]}
	currentState := queue[0]
	for currentState != -1 {
		for _, symbol := range dfa.Alphabet {
			if _, exists := dfa.Transitions[currentState][symbol]; !exists {
				if sinkState == -1 {
					sinkState = dfa.AddState(false, false)
					for _, symbol := range dfa.Alphabet {
						dfa.AddTransition(sinkState, symbol, sinkState)
					}
				}
				dfa.Transitions[currentState][symbol] = sinkState
			}
		}
		if intArray.IndexOf(currentState, dfa.AcceptStates) == -1 {
			statePartitions[0] = append(statePartitions[0], currentState)
		} else {
			statePartitions[1] = append(statePartitions[1], currentState)
		}
		for _, nextState := range dfa.Transitions[currentState] {
			if intArray.IndexOf(nextState, visited) == -1 {
				queue = append(queue, nextState)
				visited = append(visited, nextState)
			}
		}
		queue = queue[1:]
		if len(queue) != 0 {
			currentState = queue[0]
		} else {
			currentState = -1
		}
	}
	if len(statePartitions[0]) == 0 {
		statePartitions = statePartitions[1:]
	} else if len(statePartitions[1]) == 0 {
		statePartitions = statePartitions[:1]
	}
	numPartitions := 0
	for len(statePartitions) != numPartitions {
		numPartitions = len(statePartitions)
		previousPartitions := make([][]int, numPartitions)
		for i := 0; i < numPartitions; i++ {
			previousPartitions[i] = make([]int, len(statePartitions[i]))
			copy(previousPartitions[i], statePartitions[i])
		}
		splitFlag := false
		for currentPartitionIndex := 0; currentPartitionIndex < numPartitions; currentPartitionIndex++ {
			for i := 0; i < len(statePartitions[currentPartitionIndex])-1; i++ {
				for j := i + 1; j < len(statePartitions[currentPartitionIndex]); j++ {
					firstState := statePartitions[currentPartitionIndex][i]
					secondState := statePartitions[currentPartitionIndex][j]
					for k := 0; k < numPartitions; k++ {
						if k == currentPartitionIndex {
							continue
						}
						if dfa.distinguishable(firstState, secondState, previousPartitions[k]) {
							intArray.Remove(secondState, &statePartitions[currentPartitionIndex])
							if !splitFlag {
								statePartitions = append(statePartitions, make([]int, 0))
								splitFlag = true
							}
							statePartitions[numPartitions] = append(statePartitions[numPartitions], secondState)
							j--
							break
						}
					}
				}
			}
		}
	}
	if sinkState != -1 {
		for i := 0; i < len(statePartitions); i++ {
			if intArray.IndexOf(sinkState, statePartitions[i]) != -1 {
				if len(statePartitions[i]) == 1 {
					if i == len(statePartitions)-1 {
						statePartitions = statePartitions[:i]
					} else {
						statePartitions = append(statePartitions[:i], statePartitions[i+1:]...)
					}
				} else {
					intArray.Remove(sinkState, &statePartitions[i])
				}
				break
			}
		}
	}
	minDFA := New()
	minDFA.Alphabet = dfa.Alphabet
	minStates := make(map[int]int, 0)
	for i := 0; i < len(statePartitions); i++ {
		minStates[i] = minDFA.AddState(false, false)
	}
	for i := 0; i < len(statePartitions); i++ {
		for _, state := range statePartitions[i] {
			if intArray.IndexOf(state, dfa.StartStates) != -1 && intArray.IndexOf(minStates[i], minDFA.StartStates) == -1 {
				minDFA.StartStates = append(minDFA.StartStates, minStates[i])
			}
			if intArray.IndexOf(state, dfa.AcceptStates) != -1 && intArray.IndexOf(minStates[i], minDFA.AcceptStates) == -1 {
				minDFA.AcceptStates = append(minDFA.AcceptStates, minStates[i])
			}
			for symbol, targetState := range dfa.Transitions[state] {
				if targetState == sinkState {
					continue
				}
				for j := 0; j < len(statePartitions); j++ {
					if intArray.IndexOf(targetState, statePartitions[j]) != -1 {
						minDFA.Transitions[minStates[i]][symbol] = minStates[j]
						break
					}
				}
			}
		}
	}
	*dfa = *minDFA
}
