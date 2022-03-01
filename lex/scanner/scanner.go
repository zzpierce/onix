package scanner

import (
	"fmt"
	"os"
)

// scanner split the original file content word by word
// according to lexical meaning
type Scanner struct {
	state  int
	buf    []rune
	cur    []rune
	poo    string
	offset int
	dotd   int // dots in digits
}

const (
	StateError = iota
	StateSpace
	StateNumber
	StateChar
	StateStr
	StateLit
	StateOpr
	StateBrace
	StateDotLike
)

func IsSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func IsDecimal(c rune) bool {
	return c >= '0' && c <= '9'
}

func IsOpr(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%' || c == '&' ||
		c == '|' || c == '~' || c == '^' || c == ':' || c == '=' || c == '>' ||
		c == '<' || c == '!'
}

func IsBrace(c rune) bool {
	return c == '(' || c == ')' || c == '[' || c == ']' || c == '{' || c == '}'
}

func IsLit(c rune) bool {
	if c >= 'a' && c <= 'z' {
		return true
	}
	if c >= 'A' && c <= 'Z' {
		return true
	}
	if c == '_' {
		return true
	}
	return false
}

func isDotLike(c rune) bool {
	return c == '.' || c == ','
}

const (
	AfterSpace = 1 << iota
	AfterNumber
	AfterChar
	AfterStr
	AfterLit
	AfterOpr
	AfterBrace
	AfterDotLike
)

func (s *Scanner) afterAction(next rune, after int) bool {
	if (after&AfterSpace == AfterSpace) && IsSpace(next) {
		s.out()
		s.nextstate(StateSpace)
		return true
	}
	if (after&AfterNumber == AfterNumber) && IsDecimal(next) {
		s.outthen(next)
		s.nextstate(StateNumber)
		return true
	}
	if (after&AfterChar == AfterChar) && next == '\'' {
		s.outthen(next)
		s.nextstate(StateChar)
		return true
	}
	if (after&AfterStr == AfterStr) && next == '"' {
		s.outthen(next)
		s.nextstate(StateStr)
		return true
	}
	if (after&AfterLit == AfterLit) && IsLit(next) {
		s.outthen(next)
		s.nextstate(StateLit)
		return true
	}
	if (after&AfterOpr == AfterOpr) && IsOpr(next) {
		s.outthen(next)
		s.nextstate(StateOpr)
		return true
	}
	if (after&AfterBrace == AfterBrace) && IsBrace(next) {
		s.outthen(next)
		s.nextstate(StateBrace)
		return true
	}
	if (after&AfterDotLike == AfterDotLike) && isDotLike(next) {
		s.outthen(next)
		s.nextstate(StateDotLike)
		return true
	}
	return false
}

func (s *Scanner) actionSpace() error {
	e := s.next()
	if IsSpace(e) {
		return nil
	}
	if s.afterAction(e, AfterNumber|AfterChar|AfterStr|AfterLit|AfterOpr|AfterBrace|AfterDotLike) {
		return nil
	}
	return fmt.Errorf("ActionSpace meet unknown rune: %v", string(e))
}

func (s *Scanner) actionStr() error {
	e := s.next()
	if e == '"' {
		s.outwith(e)
		s.nextstate(StateSpace)
		return nil
	}
	s.cur = append(s.cur, e)
	return nil
}

func (s *Scanner) actionDigit() error {
	e := s.next()
	if IsDecimal(e) {
		s.cur = append(s.cur, e)
		return nil
	}
	if e == '.' {
		if s.dotd == 0 {
			s.cur = append(s.cur, e)
			return nil
		} else {
			return fmt.Errorf("ActionDigit multiple dots in one digit")
		}
	}
	if s.afterAction(e, AfterSpace|AfterOpr|AfterBrace) {
		return nil
	}
	return fmt.Errorf("ActionDigit unknown rune: %v", e)
}

func (s *Scanner) actionLit() error {
	e := s.next()
	if IsLit(e) || IsDecimal(e) {
		s.cur = append(s.cur, e)
		return nil
	}
	if s.afterAction(e, AfterSpace|AfterOpr|AfterBrace|AfterDotLike) {
		return nil
	}
	return fmt.Errorf("ActionLit meet unknown rune: %v", e)
}

func (s *Scanner) actionOpr() error {
	e := s.next()
	if IsOpr(e) {
		s.cur = append(s.cur, e)
		return nil
	}
	if s.afterAction(e, AfterSpace|AfterChar|AfterStr|AfterLit|AfterBrace) {
		return nil
	}
	return fmt.Errorf("ActionOpr meet unknown rune: %v", string(e))
}

func (s *Scanner) actionBrace() error {
	s.out()
	s.nextstate(StateSpace)
	return nil
}
func (s *Scanner) actionDot() error {
	s.out()
	s.nextstate(StateSpace)
	return nil
}

func (s *Scanner) step() error {
	var err error
	switch s.state {
	case StateSpace:
		err = s.actionSpace()
	case StateLit:
		err = s.actionLit()
	case StateBrace:
		err = s.actionBrace()
	case StateOpr:
		err = s.actionOpr()
	case StateStr:
		err = s.actionStr()
	case StateNumber:
		err = s.actionDigit()
	case StateDotLike:
		err = s.actionDot()
	default:
		return fmt.Errorf("step state error :%v", s.state)
	}
	return err
}

func NewScanner(path string) (*Scanner, error) {
	by, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := &Scanner{
		state: StateSpace,
		buf:   []rune(string(by)),
	}
	return s, nil
}

func (s *Scanner) Literal() (string, error) {
	defer func() {
		s.poo = ""
	}()
	for s.poo == "" {
		if s.ScanFin() {
			if s.poo == "" && len(s.cur) > 0 {
				s.out()
			}
			break
		}
		err := s.step()
		if err != nil {
			return "", err
		}
	}
	return s.poo, nil
}

func (s *Scanner) ScanFin() bool {
	return s.offset == len(s.buf)
}

func (s *Scanner) next() rune {
	e := s.buf[s.offset]
	s.offset++
	return e
}

func (s *Scanner) out() {
	s.poo = string(s.cur)
	s.cur = []rune{}
}

func (s *Scanner) outthen(c rune) {
	s.poo = string(s.cur)
	s.cur = []rune{c}
}

func (s *Scanner) outwith(c rune) {
	s.cur = append(s.cur, c)
	s.poo = string(s.cur)
	s.cur = []rune{}
}

func (s *Scanner) nextstate(state int) {
	s.state = state
}
