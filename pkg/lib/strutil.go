package lib

import (
	"github.com/mattn/go-runewidth"
	"github.com/muesli/reflow/ansi"
)

func CenterTo(msg string, ln int) string {
	msgLn := ansi.PrintableRuneWidth(msg)
	rem := ln - msgLn
	if rem < 1 {
		return msg
	}

	div := rem / 2
	msg = LPad(msg, msgLn+div)
	msg = RPad(msg, ln)
	return msg
}

func RPad(msg string, ln int) string {
	return RPadWith(msg, ' ', ln)
}

func LPad(msg string, ln int) string {
	return LPadWith(msg, ' ', ln)
}

func LPadWith(msg string, x rune, ln int) string {
	width := ansi.PrintableRuneWidth(msg)
	count := ln - width
	if count > 0 {
		b := make([]rune, count/runewidth.RuneWidth(x))
		for i := range b {
			b[i] = x
		}
		return string(b) + msg
	}
	return msg

}

// RPadWith will append runes up to ln runewidths
func RPadWith(msg string, x rune, ln int) string {
	width := ansi.PrintableRuneWidth(msg)
	count := ln - width
	if count > 0 {
		b := make([]rune, count/runewidth.RuneWidth(x))
		for i := range b {
			b[i] = x
		}
		return msg + string(b)
	}
	return msg

}

// Truncate removes any overflow past a desired length. It's possible for the result
// to be shorter than the desired length.
func Truncate(msg string, ln int) string {
	if ansi.PrintableRuneWidth(msg) <= ln {
		return msg
	}

	r := []rune(msg)
	var i, rw int
	for ; i < len(r); i++ {
		rw += runewidth.RuneWidth(r[i])
		if rw > ln {
			break
		}
	}
	return string(r[0:i])
}

// ExaExactWidth truncate a message to a particular length or add right padding until the length is hit.
func ExactWidth(msg string, ln int) string {
	return RPad(msg, ln)
}
