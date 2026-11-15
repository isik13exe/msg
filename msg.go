package msg

import (
	"errors"
	"io"
	"os"
	"reflect"
	"runtime"
	"strconv"

	"golang.org/x/term"
)

// IsTerm is what decides if the ansi characters gets written to not
var IsTerm bool = true

// Spliter is a helpful struct for saving output to a file or secondary Writer
type Spliter struct {
	Value string
	Wr    io.Writer
	Extra io.Writer
}

// Color is a helpful struct for managing colors for ansi output
type Color struct {
	R uint8
	G uint8
	B uint8
}

// String create ansi character from the Color c
func (c Color) String() string {
	return Rgb(c.R, c.G, c.B)
}

var (
	Red   = Color{255, 10, 10}
	Blue  = Color{10, 10, 255}
	Green = Color{10, 255, 10}
)

// Write makes the Splitter write into its main output and secondary output
func (e *Spliter) Write(p []byte) (n int, err error) {
	e.Value = string(p)
	e.Wr.Write(p)
	if e.Extra != nil {
		e.Extra.Write(p)
	}
	return len(p), nil
}

// ErrorSplitter provides a tool to split the stream to a file and stderr
func ErrorSplitter(filename string) *Spliter {
	if filename != "" {
		f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		return &Spliter{
			Wr:    os.Stderr,
			Extra: f,
		}
	}
	return &Spliter{
		Wr: os.Stderr,
	}
}

// FileSplitter provides a tool to split the stream to a file and stdout
func FileSplitter(filename string) *Spliter {
	if filename != "" {
		f, _ := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		return &Spliter{
			Wr:    os.Stdout,
			Extra: f,
		}
	}
	return &Spliter{
		Wr: os.Stdout,
	}
}

// Check checks if the output is a terminal or pipe/file. It should be used before using color-based functions
func Check() {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		IsTerm = false
	}
}

// Rgb create ansi character for the specified r(ed) g(reen) b(lue) values
func Rgb(r, g, b uint8) string {
	sh, _ := Sprint("\x1b[38;2;%;%;%m", r, g, b)
	return sh
}

// Reset create ansi character used to reset color of output
func Reset() string {
	return "\x1b[0m"
}

// Reds appends Red ansi chacter and reset to the strign s
func Reds(s string) string {
	return Red.String() + s + Reset()
}

// Blues appends Blue ansi chacter and reset to the strign s
func Blues(s string) string {
	return Blue.String() + s + Reset()
}

// Greens appends Green ansi chacter and reset to the strign s
func Greens(s string) string {
	return Green.String() + s + Reset()
}

// ERROR is similair to Print but appends a red "[ERROR] " at the beggining
func ERROR(a ...any) {
	if IsTerm {
		Print(append([]any{Reds("[ERROR]")}, a...)...)
	} else {
		Print(append([]any{"[ERROR]"}, a...)...)
	}
}

// ERRORf is similair to Printf but appends a red "[ERROR] " at the beggining
func ERRORf(format string, data ...any) {
	if IsTerm {
		Printf(Reds("[ERROR] ")+format, data...)
	} else {
		Printf("[ERROR] "+format, data...)
	}
}

// ERRORln is similair to Println but appends a red "[ERROR] " at the beggining
func ERRORln(a ...any) {
	if IsTerm {
		Println(append([]any{Reds("[ERROR]")}, a...)...)
	} else {
		Println(append([]any{"[ERROR]"}, a...)...)
	}
}

// IsERROR prints if the error isn't nil then returns the error to be handled. just too lazy to rewrite ERRORln after each error
func IsERROR(err error) error {
	if err != nil {
		ERRORln(err)
	}
	return err
}

// SUCCESS is similair to Print but appends a green "[SUCCESS] " at the beggining
func SUCCESS(a ...any) {
	if IsTerm {
		Print(append([]any{Greens("[SUCCESS]")}, a...)...)
	} else {
		Print(append([]any{"[SUCCESS]"}, a...)...)
	}
}

// SUCCESSf is similair to Printf but appends a green "[SUCCESS] " at the beggining
func SUCCESSf(format string, data ...any) {
	if IsTerm {
		Printf(Greens("[SUCCESS]  ")+format, data...)
	} else {
		Printf("[SUCCESS]  "+format, data...)
	}
}

// SUCCESSln is similair to Println but appends a green "[SUCCESS] " at the beggining
func SUCCESSln(a ...any) {
	if IsTerm {
		Println(append([]any{Greens("[SUCCESS]")}, a...)...)
	} else {
		Println(append([]any{"[SUCCESS]"}, a...)...)
	}
}

// INFO is similair to Print but appends a blue "[INFO] " at the beggining
func INFO(a ...any) {
	if IsTerm {
		Print(append([]any{Blues("[INFO]")}, a...)...)
	} else {
		Print(append([]any{"[INFO]"}, a...)...)
	}
}

// INFOf is similair to Printf but appends a blue "[INFO] " at the beggining
func INFOf(format string, data ...any) {
	if IsTerm {
		Printf(Blues("[INFO]  ")+format, data...)
	} else {
		Printf("[INFO]  "+format, data...)
	}
}

// INFOln is similair to Println but appends a blue "[INFO] " at the beggining
func INFOln(a ...any) {
	if IsTerm {
		Println(append([]any{Blues("[INFO]")}, a...)...)
	} else {
		Println(append([]any{"[INFO]"}, a...)...)
	}
}

// Prints the data in the specified literal color from r(ed) g(reen) b(lue)
func Colorlitln(r uint8, g uint8, b uint8, data ...any) {
	holder, _ := Sprintln(data...)
	if IsTerm {
		Printf(Rgb(r, g, b)+"%"+Reset(), holder)
	} else {
		Printf("%", holder)
	}
}

// Prints the data in the specified Color
func Colorln(c Color, data ...any) error {
	holder, ferr := Sprintln(data...)
	if ferr != nil {
		return ferr
	}
	if IsTerm {
		if ferr := Printf(c.String()+"%\n", holder); ferr != nil {
			return ferr
		}
	} else {
		if ferr := Printf("%\n", holder); ferr != nil {
			return ferr
		}
	}
	return nil
}

type Pc int

const (
	NONE Pc = iota
	DEREF
	PTR
)

// Prints the data into the writer out
func Fprint(out io.Writer, a ...any) error {
	return Fprintf(out, gen_format(a...), a...)
}

// Prints the data into the writer out in a formatted string
func Fprintf(out io.Writer, format string, a ...any) error {
	data, err := Sprint(format, a...)
	if err != nil {
		return err
	}
	n, serr := out.Write([]byte(data))
	if serr != nil || n != len(data) {
		content, err := Sprint("Couldn't print because: %!=% or %", n, len(data), serr)
		return errors.New(content + err.Error())
	}
	return nil
}

// Prints the data into the writer out with an appended newline
func Fprintln(out io.Writer, a ...any) error {
	return Fprintf(out, gen_format(a...)+"\n", a...)
}

/*
Scanning is not supported yet

	func Fscan(in io.Reader, a ...any) error {
		return Fscanf(in, gen_format(a...), a...)
	}

	func Fscanln(in io.Reader, a ...any) error {
		return Fscanf(in, gen_format(a...), a...)
	}

	func Fscanf(in io.Reader, a ...any) error {
		out := make([]byte, 2048)
		err := read(in, &out)
		if err != nil {
			return err
		}
		return nil
	}

	func Mscanf(in io.Reader, output *string) error {
		out := make([]byte, 2048)
		err := read(in, &out)
		if err != nil {
			return err
		}
		*output = string(out)
		return nil
	}

	func read(in io.Reader, data *[]byte) error {
		n, serr := in.Read(*data)
		if serr != nil {
			content, err := Sprint("Couldn't scan because: %\n read just % bytes.", serr, n)
			return errors.New(content + err.Error())
		}
		(*data)[n-1] = 0
		return nil
	}
*/

// Sprint formats the provided data a according to the format string
func Sprint(format string, a ...any) (string, error) {
	return interpol(format, a...)
}

// Sprint formats the provided data a according to the format string appending a newline at the end
func Sprintln(a ...any) (string, error) {
	sh, ferr := interpol(gen_format(a...), a...)
	return sh + "\n", ferr
}

// Println prints the data into the standered output with an appended newline
func Println(a ...any) error {
	return Fprintln(os.Stdout, a...)
}

// Printf prints the data into the standered output in a formatted way
func Printf(format string, a ...any) error {
	err := Fprintf(os.Stdout, format, a...)
	if err != nil {
		panic(err)
	}
	return nil
}

// Print prints the data into the standered output in a formatted way
func Print(a ...any) error {
	return Fprint(os.Stdout, a...)
}

// countf counts the number of % in a string, it returns the number of % relevent to formatting and an array of indices of %% which are not relevent to formatting but should be proccessed into one %
func countf(f []rune) (int, []int, error) {
	l := len(f)
	if l == 0 {
		return 0, []int{}, nil
	}
	count := 0
	ind := []int{}
	for i := 0; i < l; i++ {
		if f[i] == '%' {
			if i < l-1 && f[i+1] == '%' {
				ind = append(ind, i)
				i++
				continue
			} else {
				count++
			}
		}
		if f[i] == '{' && i > 0 && f[i-1] != '_' {
			return 0, nil, Errorf("found relevent '{' in '%' processing.")
		}
	}
	return count, ind, nil
}

// reduce process the %% in string f according to the indeces ind that should be provided by countf to reduce then into one %
func reduce(f []rune, ind []int) string {
	holder := []rune{}
	old := 0
	for _, here := range ind {
		holder = append(holder, f[old:here]...)
		old = here + 1
	}
	holder = append(holder, f[old:]...)
	return string(holder)
}

// raw_reduce assumes that there is no operative % and { in the string, just escaped ones with %%
func raw_reduce(f string) string {
	h := []rune(f)
	res := []rune{}
	l := len(h)
	for i := 0; i < l; i++ {
		if h[i] == '%' {
			res = append(res, '%')
			i++
			continue
		}
		if h[i] == '_' && i < l-1 && h[i+1] == '{' {
			res = append(res, '{')
			i++
			continue
		}
		res = append(res, h[i])
	}
	return string(res)
}

const (
	SIMPLE = iota
	POS
	PER
)

func interpol(f string, a ...any) (string, error) {
	mode := SIMPLE
	b := []byte(f)
	l := len(b)
	for i := 0; i < l; i++ {
		if b[i] == '%' {
			if i < l-1 && b[i+1] != '%' {
				mode = PER
				break
			} else {
				i++
				continue
			}
		}
		if b[i] == '{' && i > 0 && b[i-1] != '_' {
			mode = POS
			break
		}
	}

	switch mode {
	case SIMPLE:
		if len(a) > 0 {
			return f, Errorf("No formating but % data provided", len(a))
		}
		return raw_reduce(f), nil
	case PER:
		return per_print(f, a...)
	case POS:
		return pos_print(f, a...)
	default:
		panic("unreachable")
	}
}

// per_print iterates the string f to proccess the data array a and bake it all into one string for the case of using %
func per_print(f string, a ...any) (string, error) {
	holder := []rune(f)
	c, ind, err := countf(holder)
	if err != nil {
		return f, err
	}
	if c == 0 && len(a) == 0 {
		return reduce(holder, ind), nil
	}
	if len(f) == 0 {
		return "", errors.New("empty format string")
	}
	data := ""
	msg := []rune{}
	n := 0
	i := 0
	l := len(holder)
	deref := false
	for ; i < l; i++ {
		if holder[i] == '%' {
			if i < l-1 {
				if holder[i+1] == '%' {
					msg = append(msg, holder[i])
					i++
					continue
				}
			}
			if n >= len(a) {
				data += string(msg)
				data += string(f[i])
				hold, _ := Sprint("(<=ex-error %(%%)>%(args))", n+1, len(a))
				data += hold + string(holder[i+1:])
				err := "Unmatched number of %% and data: %%=% data=% too few data"
				sh, ferr := Sprint(err, n+1, len(a))
				if ferr != nil {
					sh = sh + " " + ferr.Error()
				}
				return data, errors.New(sh)
			}
			data += string(msg)
			msg = []rune{}
			hold, err := format(a[n], deref)
			if err != nil {
				data += string(msg)
				data += string(f[i])
				holder, _ := Sprint("(<=%-error)", err.Error()[:5])
				data += holder + f[i+1:]
				return data, err
			}
			data += hold
			if deref {
				deref = false
			}
			n++
			continue
		}
		if holder[i] == '*' && i < l-1 && holder[i+1] == '%' {
			deref = true
			continue
		}
		msg = append(msg, holder[i])
	}
	if len(msg) > 0 {
		data += string(msg)
	}
	if n != len(a) {
		sh, _ := Sprint("Unmatched number of %% and data: %<% too much data", n, len(a))
		holder, _ := Sprint("(<=be-[ERROR]  %(%%)<%(args))", n, len(a))
		return data + holder, errors.New(sh)
	}
	return data, nil
}

// format is a layer between having a data interface a and a reflect.Value v for format_v
func format(a any, deref bool) (string, error) {
	return format_v(reflect.ValueOf(a), deref)
}

// format_v stringefies the data stored in reflec.Value v to be baked into the resulted string of interpol
func format_v(v reflect.Value, deref bool) (string, error) {
	//println(v.Kind().String())
	var perr error = nil
	switch v.Kind() {
	case reflect.Pointer:
		if !deref {
			vptr := v.Pointer()
			if vptr != 0 {
				return "0x" + strconv.FormatUint(uint64(vptr), 16), nil
			} else {
				return "<nil>", nil
			}
		} else {
			if v.Elem().Kind() == reflect.Struct {
				str := v.MethodByName("String")
				if str.IsValid() {
					vstr := str.Call(nil)
					if len(vstr) == 0 {
						return "", errors.New("Method String doesn't return shit")
					}
					if vstr[0].Kind() != reflect.String {
						return "", errors.New("Method String returns shit that isn't string")
					}
					return vstr[0].String(), nil
				}
				err := v.MethodByName("Error")
				if err.IsValid() {
					verr := err.Call(nil)
					if len(verr) == 0 {
						return "", errors.New("Method Error doesn't return shit")
					}
					if verr[0].Kind() != reflect.String {
						return "", errors.New("Method Error returns shit that isn't string")
					}
					return verr[0].String(), nil
				}
				sh, ferr := format_v(v.Elem(), false)
				return "&(" + sh + ")", ferr
			} else {
				sh, ferr := format_v(v.Elem(), false)
				return sh, ferr
			}
		}
	case reflect.Interface:
		return format_v(v.Elem(), deref)
	case reflect.String:
		if deref {
			perr = errors.New("Deref tried to deref a string")
		}
		return v.String(), perr
	case reflect.Float32:
		if deref {
			perr = errors.New("Deref tried to deref a Float32")
		}
		return strconv.FormatFloat(v.Float(), 'g', -1, 32), perr
	case reflect.Float64:
		if deref {
			perr = errors.New("Deref tried to deref a Float64")
		}
		return strconv.FormatFloat(v.Float(), 'g', -1, 64), perr
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if deref {
			perr = errors.New("Deref tried to deref an Int")
		}
		return strconv.FormatInt(v.Int(), 10), perr
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if deref {
			perr = errors.New("Deref tried to deref an Uint")
		}
		return strconv.FormatUint(v.Uint(), 10), perr
	case reflect.Func:
		if deref {
			return strip_name(runtime.FuncForPC(v.Pointer()).Name()), perr
		}

		return "0x" + strconv.FormatUint(uint64(v.Pointer()), 16), perr
	case reflect.Struct:
		if deref {
			perr = errors.New("Deref tried to deref a Struct")
		}
		str := v.MethodByName("String")
		if str.IsValid() {
			vstr := str.Call(nil)
			if len(vstr) == 0 {
				return "", errors.New("Method String doesn't return shit")
			}
			if vstr[0].Kind() != reflect.String {
				return "", errors.New("Method String returns shit that isn't string")
			}
			return vstr[0].String(), nil
		}
		err := v.MethodByName("Error")
		if err.IsValid() {
			verr := err.Call(nil)
			if len(verr) == 0 {
				return "", errors.New("Method Error doesn't return shit")
			}
			if verr[0].Kind() != reflect.String {
				return "", errors.New("Method Error returns shit that isn't string")
			}
			return verr[0].String(), nil
		}
		s := "{"
		tv := v.Type()
		n := tv.NumField()
		for i := 0; i < n; i++ {
			if i == 0 {
				sh, ferr := format_v(v.Field(i), false)
				if ferr != nil {
					return s + sh, ferr
				}
				s += sh
				continue
			}
			sh, ferr := format_v(v.Field(i), false)
			if ferr != nil {
				return s + sh, ferr
			}
			s += ", " + sh
		}
		return s + "}", perr
	case reflect.Array, reflect.Slice:
		if deref {
			perr = errors.New("Deref tried to deref an Array")
		}
		s := "["
		for i := 0; i < v.Len(); i++ {
			if i == 0 {
				sh, ferr := format_v(v.Index(i), false)
				if ferr != nil {
					return s + sh, ferr
				}
				s += sh
				continue
			}
			sh, ferr := format_v(v.Index(i), deref)
			if ferr != nil {
				return s + sh, ferr
			}
			s += ", " + sh
		}
		s += "]"
		return s, perr
	case reflect.Map:
		if deref {
			perr = errors.New("Deref tried to deref a Map")
		}
		s := "map["
		if v.Len() == 0 {
			return s + "]", nil
		}
		for i, kv := range v.MapKeys() {
			if i == 0 {
				sh, serr := format_v(kv, false)
				if serr != nil {
					return s, serr
				}
				vsh, vserr := format_v(v.MapIndex(kv), false)
				if vserr != nil {
					return s, vserr
				}
				s += sh + ":" + vsh
				continue
			}
			sh, serr := format_v(kv, false)
			if serr != nil {
				return s, serr
			}
			vsh, vserr := format_v(v.MapIndex(kv), false)
			if vserr != nil {
				return s, vserr
			}
			s += ", " + sh + ":" + vsh
		}
		return s + "]", perr
	case reflect.Bool:
		if deref {
			perr = errors.New("Deref tried to deref a Map")
		}
		if v.Bool() {
			return "true", perr
		}
		return "false", perr
	case reflect.Invalid:
		return "<nil>", perr
	}
	return "", nil
}

// Gen_format exposes gen_format for debugging
func Gen_format(a ...any) string {
	return gen_format(a...)
}

// gen_format generates the formatting string for the array of data a
func gen_format(a ...any) string {
	s := ""
	for i, ar := range a {
		if i == 0 {
			if reflect.ValueOf(ar).Kind() == reflect.Pointer {
				s += "*%"
			} else {
				s += "%"
			}
			continue
		}
		if reflect.ValueOf(ar).Kind() == reflect.Pointer {
			s += " *%"
		} else {
			s += " %"
		}
	}
	return s
}

// Errorf is akin of fmt.Errorf, it creates an error from the formatinf string f and the array of data a
func Errorf(format string, a ...any) error {
	sh, _ := Sprint(format, a...)
	return errors.New(sh)
}

// strip_name is used to get the name of functions/methods/structs without the parant module
func strip_name(s string) string {
	l := len(s)
	i := l - 1
	for ; i >= 0; i-- {
		if s[i] == '.' {
			break
		}
	}
	if i != l-1 {
		return s[i+1:]
	}
	return s
}

type Brace struct {
	Pos   int
	Size  int
	Count int
	Deref bool
}

type Braces struct {
	Ins       []Brace
	Max       int
	Not_empty bool
}

// braces gets all teh instances of used {n}
func braces(f string) (Braces, error) {
	it := Braces{}
	r := []rune(f)
	l := len(r)
	i := 0
	deref := false
	for ; i < l; i++ {
		if r[i] == '%' {
			if i == l-1 || r[i+1] != '%' {
				return Braces{}, Errorf("found %% mixed in at: %", i)
			}
			i++
			continue
		}
		if r[i] == '*' && i < l-1 && r[i+1] == '{' {
			deref = true
			continue
		}
		if r[i] == '{' && i > 0 && r[i-1] != '_' {
			hold := Brace{Pos: i, Deref: deref, Size: 1}
			if deref {
				deref = false
				hold.Pos--
				hold.Size++
			}
			if i == l-1 {
				return Braces{}, Errorf("didn't close brace at: %", i)
			}
			i++
			nhold := []rune{}
			for ; i < l && r[i] == ' '; i++ {
				hold.Size++
			}
			for ; i < l; i++ {
				if r[i] == ' ' || r[i] == '}' {
					break
				}
				hold.Size++
				nhold = append(nhold, r[i])
			}
			for ; i < l && r[i] == ' '; i++ {
				hold.Size++
			}
			if i == l || r[i] != '}' {
				return Braces{}, Errorf("didn't close brace at: %", i)
			}
			hold.Size++
			i++
			n, err := strconv.Atoi(string(nhold))
			if err != nil {
				return Braces{}, Errorf("found weird stuff in braces at: % , found \"%\"", i-2, string(nhold))
			}
			hold.Count = n
			if it.Max < n {
				it.Max = n
			}
			it.Ins = append(it.Ins, hold)
			it.Not_empty = true
		}
	}
	return it, nil
}

// pos_print formats the data into the string using the {n} fromat to use fewer data instances
func pos_print(f string, a ...any) (string, error) {
	if f == "" {
		return "", nil
	}
	in, err := braces(f)
	if err != nil {
		return f, err
	}
	if !in.Not_empty {
		return f, nil
	}
	if in.Max != len(a)-1 {
		return f, Errorf("Mismatch between formatting and values: %(values needed) != %(values provided)", in.Max+1, len(a))
	}
	sa := map[int]string{}
	for i, av := range a {
		str, err := format(av, false)
		if err != nil {
			return f, err
		}
		sa[i] = str
	}
	res := ""
	from := 0
	for i := 0; i < len(in.Ins); i++ {
		str := ""
		if in.Ins[i].Deref {
			str, err = format(a[in.Ins[i].Count], true)
			if err != nil {
				return f, err
			}
		} else {
			str = sa[in.Ins[i].Count]
		}
		res += f[from:in.Ins[i].Pos]
		res += str
		from = in.Ins[i].Pos + in.Ins[i].Size
	}
	if from <= len(f)-1 {
		res += f[from:]
	}
	return res, nil
}
