# msg

A lightweight Go library for colorful, terminal-aware message output with a custom Jai-inspired format system.

## Features

- **ANSI color support** — RGB colors with automatic terminal detection
- **Custom formatting** — simple `%`-based interpolation for any Go type inspired from [Jai](https://inductive.no/jai/)
- **Positional formatting** — `{n}` placeholders to reuse or reorder values
- **Pre-styled log helpers** — `ERROR`, `SUCCESS`, `INFO` with `Print`/`Printf`/`Println` variants
- **Output splitting** — tee output to files or secondary writers
- **Terminal detection** — automatically strips ANSI codes when piping to files
- **Reflection-based serialization** — works with structs, maps, slices, pointers, interfaces, and more

## Installation

```bash
go get github.com/isik13exe/msg
```

Requires `golang.org/x/term`:

```bash
go get golang.org/x/term
```

## Quick Start

```go
package main

import (
    msg "github.com/isik13exe/msg"
)

func main() {
    // Check terminal before using colors
    msg.Check()

    msg.INFO("Starting up...")
    msg.SUCCESSln("Configuration loaded")
    msg.ERRORln("Something went wrong")

    msg.INFOf("Listening on port %", 8080)
    msg.SUCCESSf("Connected to %", "database")

    // Colored output
    msg.Colorln(msg.Red, "This is red")
    msg.Colorlitln(0, 128, 255, "This is custom blue")
}
```

## Output Helpers

Each log level has three variants, mirroring `fmt`:

| Variant | Description |
|---------|-------------|
| `ERROR(...)`  / `ERRORf` / `ERRORln` | Red `[ERROR]` prefix |
| `SUCCESS(...)` / `SUCCESSf` / `SUCCESSln` | Green `[SUCCESS]` prefix |
| `INFO(...)` / `INFOf` / `INFOln` | Blue `[INFO]` prefix |


```go
msg.IsERROR(err) // prints error if non-nil, returns it
```

## Formatting

The library implements its own lightweight formatter. `%` is the sole format verb; it works with **any** type via reflection.

```go
msg.Print("Hello, %", "world")            // Hello, world
msg.Print("% + % = %", 1, 2, 3)           // 1 + 2 = 3
msg.Printf("%", []int{1, 2, 3})           // [1, 2, 3]
```

### Supported Types

- Strings, integers, unsigned integers, floats, bools
- Pointers (`0x...` or dereferenced with `*%`)
- Structs (uses `String()` or `Error()` methods if implemented)
- Slices, arrays, maps
- Functions (name when dereferenced, address otherwise)
- Interfaces
- `nil` renders as `<nil>` to be consistent with Go's `fmt` package
- Arrays, Slices, Maps, Structs use a comma-separated list for better readability

### Positional Formatting

Use `{n}` to refer to a value by its index (starting at `0`), so a value can be reused or placed out of order. Spaces inside the braces are ignored (`{ 0 }` equals `{0}`).

```go
msg.Print("x = {0}, x + x = {1}", 2, 4)   // x = 2, x + x = 4
msg.Print("Say {1}, {0}!", "world", "Hello")  // Say Hello, world!
```

Rules:

- `%` and `{n}` can't be mixed in the same format string.
- The highest index used must be the last value provided, so every value has to be accounted for (`{0}` with two values is an error).
- Prefix with `*` to dereference a pointer, same as `*%`:

```go
name := "Go"
msg.Print("Deref: *{0}, ptr: {0}", &name)    // Deref: Go, ptr: 0xc000...
```

### Literal Percent

Use `%%` to emit a literal `%`.

```go
msg.Print("100%% done")   // 100% done
```

### Dereference Verb

Prefix a verb with `*` to dereference pointers:

```go
name := "Go"
msg.Print("*%", &name)    // Go
msg.Print("%",  &name)    // 0xc000...
```

### Custom Stringers

Any type implementing `String() string` (or `error`) will use that output:

```go
type User struct{ Name string }

func (u User) String() string { return "User(" + u.Name + ")" }

msg.Print("%", User{"alice"})   // User(alice)
```

## Colors

```go
msg.Red      // Color{255, 10, 10}
msg.Green    // Color{10, 255, 10}
msg.Blue     // Color{10, 10, 255}

msg.Rgb(r, g, b)              // ANSI escape string
msg.Reds("text")              // Red text with reset
msg.Greens("text")            // Green text with reset
msg.Blues("text")             // Blue text with reset
msg.Reset()                   // Reset sequence
msg.Colorln(msg.Red, "text")  // Print colored line
```

## Splitting Output

Write to multiple destinations at once — handy for logging to both console and file.

```go
// stdout + file
sp := msg.FileSplitter("output.log")
sp.Write([]byte("hello\n"))

// stderr
es := msg.ErrorSplitter("error.log")
es.Write([]byte("oops\n"))
```

The `Spliter` also stores the last written value in its `Value` field.

## Custom Writers

```go
msg.Fprint(os.Stdout, "hello ", "world")
msg.Fprintf(w, "value: %", 42)
msg.Fprintln(w, "line")
msg.Sprint("formatted %", val)   // returns string
msg.Sprintln(val)                // returns string + "\n"
msg.Errorf("failed: %", err)     // returns error
```

## Terminal Detection

Call `msg.Check()` once at startup:

```go
func main() {
    msg.Check()
    msg.IsTerm = false // or force-disable manually
}
```

`IsTerm` defaults to `true`; `Check()` sets it to `false` when stdout is not a TTY (e.g., piped or redirected), ensuring clean output in scripts.

## Limitations

- Supports only linux for color output
- Scanning (`Fscan`, `Fscanf`, `Fscanln`) is not yet implemented

## License

MIT — see [LICENSE](LICENSE) for details.
