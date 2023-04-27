package lexer

import (
	"fmt"
)

func MatchString(s string, token interface{}, success StateFunc) StateFunc {
	if s == "" {
		return success
	}
	runes := []rune(s)
	if len(runes) == 0 {
		return success
	}
	pos := 0
	var f StateFunc
	f = func(ctx Context, in rune) StateFunc {
		if pos >= len(runes) {
			return ctx.Error(fmt.Errorf("reached end of runes illegally"))
		}
		if in != runes[pos] {
			return ctx.Error(fmt.Errorf("expected exact string match '%s', found '%s%c'", s, string(runes[:pos]), in))
		}
		pos++
		if pos == len(runes) {
			ctx.Emit(token)
			return success
		}
		return f
	}
	return f
}
