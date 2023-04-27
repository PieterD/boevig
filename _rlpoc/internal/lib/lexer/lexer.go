package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

type StateFunc func(ctx Context, in rune) StateFunc

func zeroStateFunc(_ Context, _ rune) StateFunc {
	return nil
}

type TokenFunc func(token interface{}) error

type gridLexer struct {
	lif       LocationInFile
	token     *strings.Builder
	tokenFunc TokenFunc
	stateFunc StateFunc
	err       error
}

type Context interface {
	// Token returns the currently accumulated token string.
	Token() string
	// Emit will run the TokenFunc provided to Run with the given token.
	Emit(token interface{})
	// Clear clears the currently accumulated token string.
	Clear()
	// Error ends the lexing process, and causes run to return an error with err in its history.
	// The StateFunc returned by Error should be returned immediately.
	Error(err error) StateFunc
}

func Run(inputName string, input io.Reader, initState StateFunc, tokenFunc TokenFunc) error {
	lexer := &gridLexer{
		lif: LocationInFile{
			File: inputName,
			Line: 1,
			Rune: 1,
		},
		token:     &strings.Builder{},
		tokenFunc: tokenFunc,
		stateFunc: initState,
	}
	return lexer.run(input)
}

func (lexer *gridLexer) run(input io.Reader) error {
	reader := bufio.NewReader(input)
	for {
		r, runeSize, err := reader.ReadRune()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("%v: reading rune from input: %w", lexer.lif, err)
		}
		if r == unicode.ReplacementChar && runeSize == 1 {
			return fmt.Errorf("%v: invalid UTF-8 rune", lexer.lif)
		}

		/* Run the state func! */
		newStateFunc := lexer.stateFunc(lexer, r)
		/* Run the state func! */

		if lexer.err != nil {
			return fmt.Errorf("%v: lexer error: %w", lexer.lif, lexer.err)
		}
		if lexer.stateFunc == nil {
			return fmt.Errorf("%v: no state func", lexer.lif)
		}
		lexer.stateFunc = newStateFunc
		if r == '\n' {
			lexer.lif.AdvanceLine()
		} else {
			lexer.lif.AdvanceRune()
		}
	}
	return nil
}

func (lexer *gridLexer) Emit(token interface{}) {
	err := lexer.tokenFunc(token)
	if err != nil {
		lexer.Error(fmt.Errorf("token error while emitting '%s': %w", lexer.token, err))
	}
}

func (lexer *gridLexer) Error(err error) StateFunc {
	if lexer.err == nil {
		lexer.err = err
	}
	return zeroStateFunc
}

func (lexer *gridLexer) Token() string {
	return lexer.token.String()
}

func (lexer *gridLexer) Clear() {
	lexer.token.Reset()
}

var _ Context = &gridLexer{}
