package regexParser

import "github.com/ChristopherCamara/RegularLangauge/nfa"

func (p *RegexParser) expr() *nfa.NFA {
	term := p.term()
	if p.hasMoreChars() && p.peek() == "|" {
		p.eat("|")
		return nfa.Union(term, p.expr())
	}
	return term
}

func (p *RegexParser) term() *nfa.NFA {
	factor := p.factor()
	if p.hasMoreChars() && p.peek() != ")" && p.peek() != "|" {
		return nfa.Concat(factor, p.term())
	}
	return factor
}

func (p *RegexParser) factor() *nfa.NFA {
	atom := p.atom()
	if p.hasMoreChars() && p.isMetaChar(p.peek()) {
		p.next()
		return nfa.Closure(atom)
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
	return newNFA
}
