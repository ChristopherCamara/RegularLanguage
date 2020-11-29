package dfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLangauge/internal/intArray"
)

type DFA struct {
	RootState *State
	Alphabet  []string
}

func (dfa *DFA) Print() {
	fmt.Println("~~~DFA~~~")
	fmt.Printf("start at state %d\n", dfa.RootState.Index)
	visited := []int{dfa.RootState.Index}
	queue := []*State{dfa.RootState}
	currentState := queue[0]
	for currentState != nil {
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

func ToMinDFA(dfa *DFA, alphabet []string) *DFA {
	sinkState := CreateState(false)
	sinkState.Index = -1
	for _, symbol := range alphabet {
		sinkState.Transition[symbol] = sinkState
	}
	statePartitions := make([][]int, 0)
	statePartitions = append(statePartitions, make([]int, 0))
	statePartitions[0] = append(statePartitions[0], sinkState.Index)
	statePartitions = append(statePartitions, make([]int, 0))
	queue := []*State{dfa.RootState}
	visited := []int{dfa.RootState.Index, sinkState.Index}
	currentState := dfa.RootState
	for currentState != nil {
		for _, symbol := range alphabet {
			if currentState.Transition[symbol] == nil {
				currentState.Transition[symbol] = sinkState
			}
		}
		if !currentState.IsEnd {
			statePartitions[0] = append(statePartitions[0], currentState.Index)
		} else {
			statePartitions[1] = append(statePartitions[1], currentState.Index)
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
	if len(statePartitions[0]) == 0 {
		statePartitions = statePartitions[1:]
	} else if len(statePartitions[1]) == 0 {
		statePartitions = statePartitions[:1]
	}
	numPartitions := 0
	for len(statePartitions) != numPartitions {
		splitFlag := false
		numPartitions = len(statePartitions)
		for currentPartitionIndex := 0; currentPartitionIndex < numPartitions; currentPartitionIndex++ {
			for i := 0; i < len(statePartitions[currentPartitionIndex])-1; i++ {
				for j := i + 1; j < len(statePartitions[currentPartitionIndex]); j++ {
					firstState := findStateByIndex(dfa.RootState, statePartitions[currentPartitionIndex][i])
					secondState := findStateByIndex(dfa.RootState, statePartitions[currentPartitionIndex][j])
					for k := 0; k < numPartitions; k++ {
						if k == currentPartitionIndex {
							continue
						}
						if distinguishable(firstState, secondState, statePartitions[k]) {
							if j == len(statePartitions[currentPartitionIndex])-1 {
								statePartitions[currentPartitionIndex] = statePartitions[currentPartitionIndex][:j]
							} else {
								statePartitions[currentPartitionIndex] = append(statePartitions[currentPartitionIndex][:j], statePartitions[currentPartitionIndex][j+1:]...)
							}
							if !splitFlag {
								statePartitions = append(statePartitions, make([]int, 0))
								splitFlag = true
							}
							statePartitions[numPartitions] = append(statePartitions[numPartitions], secondState.Index)
							j = len(statePartitions[currentPartitionIndex])
							i = -1
							break
						}
					}
				}
			}
		}
	}
	if len(statePartitions[0]) == 1 {
		statePartitions = statePartitions[1:]
	} else {
		sinkIndex := intArray.IndexOf(sinkState.Index, statePartitions[0])
		if sinkIndex == len(statePartitions[0])-1 {
			statePartitions[0] = statePartitions[0][:sinkIndex]
		} else {
			statePartitions[0] = append(statePartitions[0][:sinkIndex], statePartitions[0][sinkIndex+1:]...)
		}
	}
	minDFA := new(DFA)
	minStates := make([]*State, len(statePartitions))
	for index, partition := range statePartitions {
		newState := CreateState(false)
		newState.Index = partition[0]
		minStates[index] = newState
	}
	for i := 0; i < len(statePartitions); i++ {
		for _, stateIndex := range statePartitions[i] {
			currentState := findStateByIndex(dfa.RootState, stateIndex)
			if currentState.IsEnd {
				minStates[i].IsEnd = true
			}
			for symbol, targetState := range currentState.Transition {
				if targetState.Index == sinkState.Index {
					continue
				}
				for j := 0; j < len(statePartitions); j++ {
					if intArray.IndexOf(targetState.Index, statePartitions[j]) != -1 {
						minStates[i].Transition[symbol] = minStates[j]
						break
					}
				}
			}
		}
	}
	minDFA.RootState = minStates[0]
	return minDFA
}
