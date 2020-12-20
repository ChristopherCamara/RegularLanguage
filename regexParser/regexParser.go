package regexparser

//RegexParser struct definition
type RegexParser struct {
	Alphabet []string
	Regex    string
	position int
}

func (p *RegexParser) hasMoreChars() bool {
	if p.position < len([]rune(p.Regex)) {
		return true
	}
	return false
}

func (p *RegexParser) isMetaChar(symbol string) bool {
	switch symbol {
	case "*":
		return true
	default:
		return false
	}
}

func (p *RegexParser) peek() string {
	return string(p.Regex[p.position])
}

func (p *RegexParser) eat(symbol string) {
	if symbol != string(p.Regex[p.position]) {
		panic("Ran into unexpected character!")
	}
	p.position++
}

func (p *RegexParser) next() string {
	current := p.peek()
	p.eat(current)
	return current
}
