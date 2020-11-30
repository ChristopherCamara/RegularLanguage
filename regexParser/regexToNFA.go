package regexParser

import "github.com/ChristopherCamara/RegularLangauge/nfa"

func (p *RegexParser) expr() *nfa.NFA {
	term := p.term()
	if p.hasMoreChars() && p.peek() == "|" {
		p.eat("|")
		term.Union(p.expr())
	}
	return term
}

func (p *RegexParser) term() *nfa.NFA {
	factor := p.factor()
	if p.hasMoreChars() && p.peek() != ")" && p.peek() != "|" {
		factor.Concat(p.term())
	}
	return factor
}

func (p *RegexParser) factor() *nfa.NFA {
	atom := p.atom()
	if p.hasMoreChars() && p.isMetaChar(p.peek()) {
		p.next()
		atom.Closure()
	}
	return atom
}

func (p *RegexParser) atom() *nfa.NFA {
	if p.peek() == "(" {
		p.eat("(")
		expr := p.expr()
		p.eat(")")
		return expr
	}
	return p.char()
}

func (p *RegexParser) char() *nfa.NFA {
	if p.isMetaChar(p.peek()) {
		panic("Unexpected meta char!")
	}
	current := p.next()
	p.Alphabet = append(p.Alphabet, current)
	return nfa.SymbolBasis(current)
}

func (p *RegexParser) ParseToNFA(regex string) *nfa.NFA {
	p.Regex = regex
	p.position = 0
	if p.Regex == "" {
		return nfa.EpsilonBasis()
	}
	newNFA := p.expr()
	newNFA.Alphabet = p.Alphabet
	return newNFA
}
