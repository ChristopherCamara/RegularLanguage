package nfa

func EpsilonBasis() *NFA {
	newNFA := New()
	startState := newNFA.addState(true, false)
	endState := newNFA.addState(false, true)
	newNFA.addEpsilonTransition(startState, endState)
	return newNFA
}

func SymbolBasis(symbol string) *NFA {
	newNFA := New()
	startState := newNFA.addState(true, false)
	endState := newNFA.addState(false, true)
	newNFA.addTransition(startState, symbol, endState)
	return newNFA
}
