package nfa

func EpsilonBasis() *NFA {
	startState := createState(false)
	endState := createState(true)
	startState.addEpsilonTransition(endState)
	return Create(startState, endState)
}

func SymbolBasis(symbol string) *NFA {
	startState := createState(false)
	endState := createState(true)
	startState.addTransition(endState, symbol)
	return Create(startState, endState)
}

func Concat(firstNFA, secondNFA *NFA) *NFA {
	for _, targetState := range secondNFA.RootState.EpsilonTransition {
		firstNFA.EndState.addEpsilonTransition(targetState)
	}
	for symbol, targetState := range secondNFA.RootState.Transition {
		firstNFA.EndState.Transition[symbol] = targetState
	}
	firstNFA.EndState.IsEnd = false
	return Create(firstNFA.RootState, secondNFA.EndState)
}

func Union(firstNFA, secondNFA *NFA) *NFA {
	newRoot := createState(false)
	newEnd := createState(true)
	newRoot.addEpsilonTransition(firstNFA.RootState)
	newRoot.addEpsilonTransition(secondNFA.RootState)
	firstNFA.EndState.addEpsilonTransition(newEnd)
	firstNFA.EndState.IsEnd = false
	secondNFA.EndState.addEpsilonTransition(newEnd)
	secondNFA.EndState.IsEnd = false
	return Create(newRoot, newEnd)
}

func Closure(nfa *NFA) *NFA {
	newRoot := createState(false)
	newEnd := createState(true)
	newRoot.addEpsilonTransition(newEnd)
	newRoot.addEpsilonTransition(nfa.RootState)
	nfa.EndState.addEpsilonTransition(newEnd)
	nfa.EndState.IsEnd = false
	nfa.EndState.addEpsilonTransition(nfa.RootState)
	return Create(newRoot, newEnd)
}
