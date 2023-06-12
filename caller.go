package maleo

import (
	"encoding/json"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode"
)

const sep = string(os.PathSeparator)

type Caller interface {
	// Function returns the function information.
	Function() *runtime.Func
	// Name returns the function name of the caller.
	Name() string
	// ShortName returns only function name of the caller.
	ShortName() string
	// ShortSource returns only the latest three items path in the File Path where the Caller comes from.
	ShortSource() string
	// String Sets this caller as `file_path:line` format.
	String() string
	// Line returns the line number of the caller.
	Line() int
	// File returns the file path of the caller.
	File() string
	// PC returns the program counter of the caller.
	PC() uintptr
	// FormatAsKey Like .String(), but the returned string is usually safe for URL.
	//
	// Default implementation changes runes other than letters, digits, `-` and `.` to `_`.
	FormatAsKey() string
	// Depth returns the depth of the caller that is initially used.
	Depth() int
}

type caller struct {
	pc    uintptr
	file  string
	line  int
	depth int
}

func (c caller) Name() string {
	f := runtime.FuncForPC(c.pc)
	return f.Name()
}

func (c caller) ShortName() string {
	s := strings.Split(c.Name(), "/")
	return s[len(s)-1]
}

func (c caller) Line() int {
	return c.line
}

func (c caller) File() string {
	return c.file
}

func (c caller) PC() uintptr {
	return c.pc
}

func (c caller) MarshalJSON() ([]byte, error) {
	type A struct {
		File string `json:"file"`
		Name string `json:"name"`
	}
	return json.Marshal(A{
		File: c.String(),
		Name: c.ShortName(),
	})
}

func (c caller) Function() *runtime.Func {
	f := runtime.FuncForPC(c.pc)
	return f
}

func (c caller) Depth() int {
	return c.depth
}

// ShortSource returns only the latest three items path in the File Path where the Caller comes from.
func (c caller) ShortSource() string {
	s := strings.Split(c.file, sep)

	for len(s) > 3 {
		s = s[1:]
	}

	return strings.Join(s, sep)
}

// FormatAsKey Like .String(), but runes other than letters, digits, `-` and `.` are set to `_`.
func (c caller) FormatAsKey() string {
	s := &strings.Builder{}
	strLine := strconv.Itoa(c.line)
	s.Grow(len(c.file) + 1 + len(strLine))
	replaceSymbols(s, c.file, '_')
	s.WriteRune('_')
	s.WriteString(strLine)
	return s.String()
}

// String Sets this caller as `file_path:line` format.
func (c caller) String() string {
	s := &strings.Builder{}
	strLine := strconv.Itoa(c.line)
	s.Grow(len(c.file) + 1 + len(strLine))
	s.WriteString(c.file)
	s.WriteRune(':')
	s.WriteString(strLine)
	return s.String()
}

func replaceSymbols(builder *strings.Builder, s string, rep rune) {
	for _, c := range s {
		switch {
		case unicode.In(c, unicode.Digit, unicode.Letter), c == '-', c == '.':
			builder.WriteRune(c)
		default:
			builder.WriteRune(rep)
		}
	}
}

// GetCaller returns the caller information for who calls this function. A value of 1 will return this GetCaller location.
// So you may want the value to be 2 or higher if you wrap this call in another function.
//
// Returns zero value if the caller information cannot be obtained.
func GetCaller(depth int) Caller {
	pc, file, line, _ := runtime.Caller(depth)
	return &caller{
		pc:    pc,
		file:  file,
		line:  line,
		depth: depth,
	}
}
