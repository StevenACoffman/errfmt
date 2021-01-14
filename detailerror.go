package errfmt

import (
	"fmt"
	"io"
	"strconv"
)

type DetailError struct {
	Msg, Detail string
	Err         error
}

func (e *DetailError) Unwrap() error { return e.Err }

func (e *DetailError) Error() string {
	if e.Err == nil {
		return e.Msg
	}
	return e.Msg + ": " + e.Err.Error()
}

func (e *DetailError) Format(s fmt.State, c rune) {
	if s.Flag('#') && c == 'v' {
		type nomethod DetailError
		fmt.Fprintf(s, "%#v", (*nomethod)(e))
		return
	}
	if !s.Flag('+') || c != 'v' {
		fmt.Fprintf(s, spec(s, c), e.Error())
		return
	}
	fmt.Fprintln(s, e.Msg)
	if e.Detail != "" {
		io.WriteString(s, "\t")
		fmt.Fprintln(s, e.Detail)
	}
	if e.Err != nil {
		if fErr, ok := e.Err.(fmt.Formatter); ok {
			fErr.Format(s, c)
		} else {
			fmt.Fprintf(s, spec(s, c), e.Err)
			io.WriteString(s, "\n")
		}
	}
}

func spec(s fmt.State, c rune) string {
	buf := []byte{'%'}
	for _, f := range []int{'+', '-', '#', ' ', '0'} {
		if s.Flag(f) {
			buf = append(buf, byte(f))
		}
	}
	if w, ok := s.Width(); ok {
		buf = strconv.AppendInt(buf, int64(w), 10)
	}
	if p, ok := s.Precision(); ok {
		buf = append(buf, '.')
		buf = strconv.AppendInt(buf, int64(p), 10)
	}
	buf = append(buf, byte(c))
	return string(buf)
}
