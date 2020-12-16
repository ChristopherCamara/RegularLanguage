package RegularLanguage

func EpsilonBasis() *NFA {
	newNFA := NewNFA()
	startState := newNFA.AddState(true, false)
	endState := newNFA.AddState(false, true)
	newNFA.AddEpsilonTransition(startState, endState)
	return newNFA
}

func SymbolBasis(symbol string) *NFA {
	newNFA := NewNFA()
	startState := newNFA.AddState(true, false)
	endState := newNFA.AddState(false, true)
	newNFA.AddTransition(startState, symbol, endState)
	return newNFA
}
