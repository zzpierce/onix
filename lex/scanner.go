package lex

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
	StateLit
	StateStr
	StateChar
	StateOpr
	StateNumber
	StateBrace
	StateDot
)

func isSpace(c rune) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func isDecimal(c rune) bool {
	return c >= '0' && c <= '9'
}

func isOpr(c rune) bool {
	return c == '+' || c == '-' || c == '*' || c == '/' || c == '%' || c == '&' ||
		c == '|' || c == '~' || c == '^' || c == ':' || c == '=' || c == '>' || c == '<'
}

func isBrace(c rune) bool {
	return c == '(' || c == ')' || c == '[' || c == ']' || c == '{' || c == '}'
}

func isLit(c rune) bool {
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

func isDot(c rune) bool {
	return c == '.'
}

func (s *Scanner) actionSpace() error {
	e := s.next()
	if isSpace(e) {
		return nil
	}
	s.cur = append(s.cur, e)
	if isDecimal(e) {
		s.nextstate(StateNumber)
		return nil
	}
	if isOpr(e) {
		s.nextstate(StateOpr)
		return nil
	}
	if isBrace(e) {
		s.nextstate(StateBrace)
		return nil
	}
	if e == '"' {
		s.nextstate(StateStr)
		return nil
	}
	if e == '\'' {
		s.nextstate(StateChar)
		return nil
	}
	if isLit(e) {
		s.nextstate(StateLit)
		return nil
	}
	return fmt.Errorf("ActionSpace meet unknown rune: %v", string(e))
}

func (s *Scanner) actionBrace() error {
	s.out()
	s.nextstate(StateSpace)
	return nil
}

func (s *Scanner) actionLit() error {
	e := s.next()
	if isSpace(e) {
		s.out()
		s.nextstate(StateSpace)
		return nil
	}
	if isOpr(e) {
		s.outthen(e)
		s.nextstate(StateOpr)
		return nil
	}
	if isBrace(e) {
		s.outthen(e)
		s.nextstate(StateBrace)
		return nil
	}
	if isDot(e) {
		s.outthen(e)
		s.nextstate(StateDot)
		return nil
	}
	if isLit(e) || isDecimal(e) {
		s.cur = append(s.cur, e)
		return nil
	}
	return fmt.Errorf("ActionLit meet unknown rune: %v", e)
}

func (s *Scanner) actionOpr() error {
	e := s.next()
	if isSpace(e) {
		s.out()
		s.nextstate(StateSpace)
		return nil
	}
	if isBrace(e) {
		s.outthen(e)
		s.nextstate(StateBrace)
		return nil
	}
	if isLit(e) {
		s.outthen(e)
		s.nextstate(StateLit)
		return nil
	}
	if isOpr(e) {
		s.cur = append(s.cur, e)
		return nil
	}
	return fmt.Errorf("ActionOpr meet unknown rune: %v", string(e))
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
	if isSpace(e) {
		s.out()
		s.nextstate(StateSpace)
		return nil
	}
	if isDecimal(e) {
		s.cur = append(s.cur, e)
		return nil
	}
	if isBrace(e) {
		s.outthen(e)
		s.nextstate(StateBrace)
		return nil
	}
	if isOpr(e) {
		s.outthen(e)
		s.nextstate(StateOpr)
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
	return fmt.Errorf("ActionDigit unknown rune: %v", e)
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
	case StateDot:
		err = s.actionDot()
	default:
		return fmt.Errorf("step state error :%v", s.state)
	}
	return err
}

func InitScan(path string) (*Scanner, error) {
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
