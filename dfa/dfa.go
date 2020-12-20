package dfa

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ChristopherCamara/finiteAutomata/internal/intArray"
	"github.com/ChristopherCamara/finiteAutomata/nfa"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

//DFA deterministic finite automata struct definition
type DFA struct {
	nextState    int
	Alphabet     []string
	States       []int
	StartStates  []int
	AcceptStates []int
	Transitions  map[int]map[string]int
}

type edge struct {
	edge  *cgraph.Edge
	label string
}

//New returns ready to use *DFA
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

//SaveGraphviz image file
func (dfa *DFA) SaveGraphviz(fileName string) {
	fileName = strings.ReplaceAll(fileName, "*", "star")
	fileName = strings.ReplaceAll(fileName, "|", " or ")
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	graph.SetLayout("dot")
	graph.SetRankDir(cgraph.LRRank)
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	nodeMappings := make(map[int]*cgraph.Node)
	edgeMappings := make(map[int]map[int]*edge)
	for _, currentState := range dfa.States {
		newNode, err := graph.CreateNode(strconv.Itoa(currentState))
		if err != nil {
			log.Panic(err)
		}
		if intArray.IndexOf(currentState, dfa.StartStates) != -1 {
			nilNode, err := graph.CreateNode("")
			if err != nil {
				log.Panic(err)
			}
			nilNode.SetShape(cgraph.NoneShape)
			_, err = graph.CreateEdge("", nilNode, newNode)
			if err != nil {
				log.Panic(err)
			}
		}
		if intArray.IndexOf(currentState, dfa.AcceptStates) != -1 {
			newNode.SetShape(cgraph.DoubleCircleShape)
		} else {
			newNode.SetShape(cgraph.CircleShape)
		}
		nodeMappings[currentState] = newNode
	}
	for _, currentState := range dfa.States {
		for symbol, transitionState := range dfa.Transitions[currentState] {
			newEdge, err := graph.CreateEdge("", nodeMappings[currentState], nodeMappings[transitionState])
			if err != nil {
				log.Panic(err)
			}
			newEdge.SetLabel(symbol)
			if edgeMappings[currentState] == nil {
				edgeMappings[currentState] = make(map[int]*edge)
				edgeMappings[currentState][transitionState] = &edge{edge: newEdge, label: symbol}
			} else if currentEdge, exists := edgeMappings[currentState][transitionState]; exists {
				currentEdge.label = currentEdge.label + "," + symbol
				currentEdge.edge.SetLabel(currentEdge.label)
			} else {
				edgeMappings[currentState][transitionState] = &edge{edge: newEdge, label: symbol}
			}
		}
	}
	err = g.RenderFilename(graph, graphviz.PNG, fileName+".png")
	if err != nil {
		log.Fatal(err)
	}
}

//AddState to a DFA
func (dfa *DFA) AddState(isStart, isAccept bool) int {
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

//AddTransition from sourceState to targetState with given symbol
func (dfa *DFA) AddTransition(sourceState int, symbol string, targetState int) {
	dfa.Transitions[sourceState][symbol] = targetState
}

//Print out DFA information
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

//Reverse a DFA, which is one of the operations it is closed under
func (dfa *DFA) Reverse() *DFA {
	NFA := nfa.New()
	NFA.Alphabet = dfa.Alphabet
	stateMappings := make(map[int]int)
	for _, state := range dfa.States {
		stateMappings[state] = NFA.AddState(false, false)
		if intArray.IndexOf(state, dfa.StartStates) != -1 {
			NFA.AcceptStates = append(NFA.AcceptStates, stateMappings[state])
		}
		if intArray.IndexOf(state, dfa.AcceptStates) != -1 {
			NFA.StartStates = append(NFA.StartStates, stateMappings[state])
		}
	}
	for _, state := range dfa.States {
		for symbol, transitionState := range dfa.Transitions[state] {
			if _, exists := NFA.Transitions[stateMappings[transitionState]]; !exists {
				NFA.Transitions[stateMappings[transitionState]] = make(map[string][]int, 0)
			}
			NFA.Transitions[stateMappings[transitionState]][symbol] = append(NFA.Transitions[stateMappings[transitionState]][symbol], stateMappings[state])
		}
	}
	newStart := NFA.AddState(false, false)
	for _, startState := range NFA.StartStates {
		NFA.AddEpsilonTransition(newStart, startState)
	}
	NFA.StartStates = []int{newStart}
	NFA.SaveGraphviz("test")
	reverseDFA := FromNFA(NFA)
	reverseDFA.Minimize()
	return reverseDFA
}

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

//Minimize a DFA, transform a DFA to the DFA with minimal states
func (dfa *DFA) Minimize() {
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

//FromNFA create a DFA from a NFA
func FromNFA(NFA *nfa.NFA) *DFA {
	epsilonClosures := NFA.GetEpsilonClosures()
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
		newState := dfa.AddState(false, false)
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
