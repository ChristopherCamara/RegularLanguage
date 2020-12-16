package RegularLanguage

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLanguage/internal/intArray"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"strconv"
	"strings"
)

type NFA struct {
	nextState          int
	Alphabet           []string
	States             []int
	StartStates        []int
	AcceptStates       []int
	Transitions        map[int]map[string][]int
	EpsilonTransitions map[int][]int
}

type edge struct {
	edge  *cgraph.Edge
	label string
}

func NewNFA() *NFA {
	newNFA := new(NFA)
	newNFA.nextState = 0
	newNFA.Alphabet = make([]string, 0)
	newNFA.States = make([]int, 0)
	newNFA.StartStates = make([]int, 0)
	newNFA.AcceptStates = make([]int, 0)
	newNFA.Transitions = make(map[int]map[string][]int, 0)
	newNFA.EpsilonTransitions = make(map[int][]int, 0)
	return newNFA
}

func (nfa *NFA) SaveGraphviz(fileName string) {
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
	for _, currentState := range nfa.States {
		newNode, err := graph.CreateNode(strconv.Itoa(currentState))
		if err != nil {
			log.Panic(err)
		}
		if intArray.IndexOf(currentState, nfa.StartStates) != -1 {
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
		if intArray.IndexOf(currentState, nfa.AcceptStates) != -1 {
			newNode.SetShape(cgraph.DoubleCircleShape)
		} else {
			newNode.SetShape(cgraph.CircleShape)
		}
		nodeMappings[currentState] = newNode
	}
	for _, currentState := range nfa.States {
		for _, transitionState := range nfa.EpsilonTransitions[currentState] {
			newEdge, err := graph.CreateEdge("", nodeMappings[currentState], nodeMappings[transitionState])
			if err != nil {
				log.Panic(err)
			}
			newEdge.SetLabel("ε")
			if edgeMappings[currentState] == nil {
				edgeMappings[currentState] = make(map[int]*edge)
				edgeMappings[currentState][transitionState] = &edge{edge: newEdge, label: "ε"}
			} else if currentEdge, exists := edgeMappings[currentState][transitionState]; exists {
				currentEdge.label = currentEdge.label + ",ε"
				currentEdge.edge.SetLabel(currentEdge.label)
			} else {
				edgeMappings[currentState][transitionState] = &edge{edge: newEdge, label: "ε"}
			}
		}
		for symbol, transitionStates := range nfa.Transitions[currentState] {
			for _, transitionState := range transitionStates {
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
	}
	err = g.RenderFilename(graph, graphviz.SVG, fileName+".svg")
	if err != nil {
		log.Fatal(err)
	}
}

func (nfa *NFA) RemoveState(removeState int) {
	intArray.Remove(removeState, &nfa.States)
	intArray.Remove(removeState, &nfa.StartStates)
	intArray.Remove(removeState, &nfa.AcceptStates)
	delete(nfa.Transitions, removeState)
	delete(nfa.EpsilonTransitions, removeState)
	for i := 0; i < len(nfa.States); i++ {
		currentState := nfa.States[i]
		for symbol, transitions := range nfa.Transitions[currentState] {
			for index, transition := range transitions {
				if transition > removeState {
					nfa.Transitions[currentState][symbol][index] = transition - 1
				}
			}
		}
		for index, transition := range nfa.EpsilonTransitions[currentState] {
			if transition > removeState {
				nfa.EpsilonTransitions[currentState][index] = transition - 1
			}
		}
		if currentState > removeState {
			nfa.States[i]--
			nfa.Transitions[currentState-1] = nfa.Transitions[currentState]
			nfa.EpsilonTransitions[currentState-1] = nfa.EpsilonTransitions[currentState]
		}
	}
	for i := 0; i < len(nfa.StartStates); i++ {
		if nfa.StartStates[i] > removeState {
			nfa.StartStates[i]--
		}
	}
	for i := 0; i < len(nfa.AcceptStates); i++ {
		if nfa.AcceptStates[i] > removeState {
			nfa.AcceptStates[i]--
		}
	}
	nfa.nextState--
}

func (nfa *NFA) AddState(isStart, isAccept bool) int {
	index := nfa.nextState
	nfa.States = append(nfa.States, index)
	nfa.Transitions[index] = make(map[string][]int, 0)
	nfa.EpsilonTransitions[index] = make([]int, 0)
	if isStart {
		nfa.StartStates = append(nfa.StartStates, index)
	}
	if isAccept {
		nfa.AcceptStates = append(nfa.AcceptStates, index)
	}
	nfa.nextState++
	return index
}

func (nfa *NFA) AddEpsilonTransition(sourceState, targetState int) {
	nfa.EpsilonTransitions[sourceState] = append(nfa.EpsilonTransitions[sourceState], targetState)
}

func (nfa *NFA) AddTransition(sourceState int, symbol string, targetState int) {
	nfa.Transitions[sourceState][symbol] = append(nfa.Transitions[sourceState][symbol], targetState)
}

func (nfa *NFA) merge(other *NFA) map[int]int {
	//map otherStates to new states in nfa
	newStates := make(map[int]int, 0)
	for _, otherState := range other.States {
		newStates[otherState] = nfa.AddState(false, false)
	}
	//update epsilon transitions to new states
	for otherState, transitionStates := range other.EpsilonTransitions {
		for _, transitionState := range transitionStates {
			nfa.AddEpsilonTransition(newStates[otherState], newStates[transitionState])
		}
	}
	//update transitions to new states
	for otherState, transitions := range other.Transitions {
		for symbol, transitionStates := range transitions {
			for _, transitionState := range transitionStates {
				nfa.AddTransition(newStates[otherState], symbol, newStates[transitionState])
			}
		}
	}
	return newStates
}

func (nfa *NFA) Concat(other *NFA) {
	newStates := nfa.merge(other)
	//move all transitions from other startStates to our acceptStates
	for _, otherStart := range other.StartStates {
		for _, otherTransition := range other.EpsilonTransitions[otherStart] {
			for _, acceptState := range nfa.AcceptStates {
				nfa.AddEpsilonTransition(acceptState, newStates[otherTransition])
			}
		}
		for symbol, otherTransitions := range other.Transitions[otherStart] {
			for _, otherTransition := range otherTransitions {
				for _, acceptState := range nfa.AcceptStates {
					nfa.AddTransition(acceptState, symbol, newStates[otherTransition])
				}
			}
		}
		//change our accept states to the other accept states
		nfa.AcceptStates = make([]int, 0)
		for _, otherAccept := range other.AcceptStates {
			nfa.AcceptStates = append(nfa.AcceptStates, newStates[otherAccept])
		}
		//remove the otherStartStates now that its ported over
		for _, otherStart := range other.StartStates {
			nfa.RemoveState(newStates[otherStart])
		}
	}
}

func (nfa *NFA) Union(other *NFA) {
	newNFA := NewNFA()
	newStart := newNFA.AddState(true, false)
	newStates := newNFA.merge(nfa)
	newOtherStates := newNFA.merge(other)
	newAccept := newNFA.AddState(false, true)
	for _, startState := range nfa.StartStates {
		newNFA.AddEpsilonTransition(newStart, newStates[startState])
	}
	for _, otherStart := range other.StartStates {
		newNFA.AddEpsilonTransition(newStart, newOtherStates[otherStart])
	}
	for _, acceptState := range nfa.AcceptStates {
		newNFA.AddEpsilonTransition(newStates[acceptState], newAccept)
	}
	for _, otherAccept := range other.AcceptStates {
		newNFA.AddEpsilonTransition(newOtherStates[otherAccept], newAccept)
	}
	*nfa = *newNFA
}

func (nfa *NFA) Closure() {
	newNFA := NewNFA()
	newStart := newNFA.AddState(true, false)
	newStates := newNFA.merge(nfa)
	newAccept := newNFA.AddState(false, true)
	newNFA.AddEpsilonTransition(newStart, newAccept)
	for _, startState := range nfa.StartStates {
		newNFA.AddEpsilonTransition(newStart, newStates[startState])
	}
	for _, acceptState := range nfa.AcceptStates {
		newNFA.AddEpsilonTransition(newStates[acceptState], newAccept)
		for _, startState := range nfa.StartStates {
			newNFA.AddEpsilonTransition(newStates[acceptState], newStates[startState])
		}
	}
	*nfa = *newNFA
}

func (nfa *NFA) Print() {
	fmt.Println("~~~NFA~~~")
	fmt.Print("start states: ")
	intArray.Print(nfa.StartStates)
	for _, state := range nfa.States {
		fmt.Printf("state %d:\n", state)
		if nfa.Transitions[state] != nil {
			for symbol, states := range nfa.Transitions[state] {
				fmt.Printf("\t%s -> ", symbol)
				intArray.Print(states)
			}
		}
		if len(nfa.EpsilonTransitions[state]) != 0 {
			fmt.Print("\t(empty) -> ")
			intArray.Print(nfa.EpsilonTransitions[state])
		}
	}
	fmt.Print("accept states: ")
	intArray.Print(nfa.AcceptStates)
}

func (nfa *NFA) epsilonClosure(state int, closure *[]int, visited *[]int) {
	*closure = append(*closure, state)
	if nfa.EpsilonTransitions[state] != nil {
		for _, transitionState := range nfa.EpsilonTransitions[state] {
			if intArray.IndexOf(transitionState, *visited) == -1 {
				*visited = append(*visited, transitionState)
				nfa.epsilonClosure(transitionState, closure, visited)
			}
		}
	}
}

func (nfa *NFA) GetEpsilonClosures() map[int][]int {
	epsilonClosures := make(map[int][]int, 0)
	for _, state := range nfa.States {
		closure := make([]int, 0)
		visited := make([]int, 0)
		nfa.epsilonClosure(state, &closure, &visited)
		epsilonClosures[state] = closure

	}
	return epsilonClosures
}