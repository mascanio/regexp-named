# regex-named.go

`regex-named.go` is a package that extends Go's regexp package, adding `find` methods that returns maps of named capture groups.

## Installation

```bash
go get github.com/mascanio/regex-named.go
```

## Usage

```go
re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
m0, m := re.FindStringNamed("foo 42")

fmt.Println(m0) // foo 42
fmt.Println(m["name"]) // foo
fmt.Println(m["age"]) // 42
```

## Methods
These methods work like the corresponding methods in the regexp, replacing the slices returned by the corresponding methods for maps indexed by the names of the groups.

```go
func (re *RegexNamed) FindNamed(s []byte) ([]byte, map[string][]byte)

func (re *RegexNamed) FindIndexNamed(s []byte) ([]int, map[string][]int)

func (re *RegexNamed) FindStringNamed(s string) (string, map[string]string)

func (re *RegexNamed) FindStringIndexNamed(s string) ([]int, map[string][]int)

func (re *RegexNamed) FindAllNamed(b []byte, n int) ([][]byte, []map[string][]byte)

func (re *RegexNamed) FindAllIndexNamed(b []byte, n int) ([][]int, []map[string][]int)

func (re *RegexNamed) FindAllStringNamed(s string, n int) ([]string, []map[string]string)

func (re *RegexNamed) FindAllStringIndexNamed(s string, n int) ([][]int, []map[string][]int)
```

RegexNamed are created with the Compile and MustCompile functions, which work like the corresponding functions in the regexp package.

```go
func Compile(expr string) (*RegexNamed, error)
func MustCompile(expr string) *RegexNamed
```

## License
regex-named.go is licensed under the MIT license. See the LICENSE file for details.
