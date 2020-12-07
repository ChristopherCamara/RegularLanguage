package dfa

import (
	"fmt"
	"github.com/ChristopherCamara/RegularLanguage/internal/intArray"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	"log"
	"strconv"
)

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

func (dfa *DFA) SaveGraphviz(fileName string) {
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
	err = g.RenderFilename(graph, graphviz.SVG, fileName+".svg")
	if err != nil {
		log.Fatal(err)
	}
}

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
