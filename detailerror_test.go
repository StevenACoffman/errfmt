package errfmt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestDetailErrorFormat(t *testing.T) {
	nested := &DetailError{
		Msg:    `reading "file"`,
		Detail: `cmd/prog/reader.go:122`,
		Err: &DetailError{
			Msg:    "parsing line 23",
			Detail: "iff x > 3 {\n\tcmd/prog/parser.go:85",
			Err: &DetailError{
				Msg:    "syntax error",
				Detail: "cmd/prog/parser.go:214",
			},
		},
	}
	for _, test := range []struct {
		Err    error
		format string
		want   string
	}{
		{
			Err:    errors.New("x"),
			format: "%v",
			want:   `x`,
		},
		{
			Err:    errors.New("x"),
			format: "%+v",
			want:   `x`,
		},
		{
			Err:    &DetailError{Msg: "m"},
			format: "%v",
			want:   `m`,
		},
		{
			Err:    &DetailError{Msg: "m"},
			format: "%+v",
			want:   "m\n",
		},
		{
			Err:    &DetailError{Msg: "m", Detail: "d\ne"},
			format: "%v",
			want:   "m",
		},
		{
			Err:    &DetailError{Msg: "m", Detail: "d\ne"},
			format: "%s",
			want:   "m",
		},
		{
			Err:    &DetailError{Msg: "m", Detail: "d\ne"},
			format: "%+v",
			want:   "m\n\td\ne\n",
		},
		{
			Err:    &DetailError{Msg: "m", Detail: "d", Err: io.ErrUnexpectedEOF},
			format: "%+v",
			want:   "m\n\td\nunexpected EOF\n",
		},
		{
			Err:    &os.PathError{Op: "op", Path: "path", Err: os.ErrNotExist},
			format: "%v",
			want:   "op path: file does not exist",
		},
		{
			Err: &DetailError{
				Msg:    "m",
				Detail: "d",
				Err: &os.PathError{
					Op:   "op",
					Path: "path",
					Err:  os.ErrNotExist,
				},
			},
			format: "%+v",
			want:   "m\n\td\nop path: file does not exist\n",
		},
		{
			Err:    nested,
			format: "%v",
			want:   `reading "file": parsing line 23: syntax error`,
		},
		{
			Err:    nested,
			format: "%+v",
			want: `reading "file"
	cmd/prog/reader.go:122
parsing line 23
	iff x > 3 {
	cmd/prog/parser.go:85
syntax error
	cmd/prog/parser.go:214
`,
		},
		{
			Err:    &DetailError{Msg: "m"},
			format: "%5s",
			want:   "    m",
		},
		{
			Err:    &DetailError{Msg: "m"},
			format: "%X",
			want:   "6D",
		},
		{
			Err:    io.ErrUnexpectedEOF,
			format: "%+15v",
			want:   " unexpected EOF",
		},
		{
			Err:    &DetailError{Msg: "m", Err: io.ErrUnexpectedEOF},
			format: "%+15v",
			want:   "m\n unexpected EOF\n",
		},
	} {
		got := fmt.Sprintf(test.format, test.Err)
		if got != test.want {
			t.Errorf("%q on %#v:\ngot  %q\nwant %q", test.format, test.Err, got, test.want)
		}
	}
}

type specfmt struct{}

func (specfmt) Format(s fmt.State, c rune) {
	io.WriteString(s, spec(s, c))
}

func TestSpec(t *testing.T) {
	for _, format := range []string{
		"%s", "%v", "%.2d", "%5.2X", "%8g", "%+v", "%+-#o",
	} {
		got := fmt.Sprintf(format, specfmt{})
		if got != format {
			t.Errorf("got %q, want %q", got, format)
		}
	}
}
