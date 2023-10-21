// Package regexp named provides named submatches for Go's regexp package.
//
// The package extends the regexp package with the following methods:
//
//	FindNamed
//	FindIndexNamed
//	FindStringNamed
//	FindStringIndexNamed
//	FindAllNamed
//	FindAllIndexNamed
//	FindAllStringNamed
//	FindAllStringIndexNamed
//
// These methods work like the corresponding methods in the regexp, replacing
// the slices returned by the corresponding methods for maps indexed by the
// names of the groups.
//
// RegexpNamed are created with the Compile and MustCompile functions, which
// work like the corresponding functions in the regexp package.
//
// For example:
//
//		re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
//	 m0, m := re.FindStringNamed("foo 42")
//
// m0 will be "foo 42" and m will be a map[string]string with the following
// values:
//
//	m["name"] == "foo"
//	m["age"] == "42"
//
// If a group is not matched, the corresponding value in the map will be an
// empty string.
package regexp_named

import (
	"errors"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type RegexpNamed struct {
	namedMap map[string]int
	*regexp.Regexp
}

const UnnamedCapture = string(unicode.ReplacementChar)

var reNamedMatch = regexp.MustCompile(`^\?P\<(.*?)\>`)
var reNoCaptureMatch = regexp.MustCompile(`^\?:`)

// parseBytes parses a regular expression, returning a slice
// of strings containing the names of the groups in the regular
// expression:
//   - If the i'th group is named, its name is returned in the
//     i position of the slice.
//   - If the i'th group is not named, the i position of the
//     slice is set to UnnamedCapture.
//   - Non capturing groups are ignored.
//
// If the length of the input is 0, nil is returned.
// If the regular expression is malformed (invalid rune found or
// the regexp ends in a backslash), an error is returned.
//
// Examples:
//
//	parseBytes([]byte(`(?P<name>\w+) (?P<age>\d+)`), 0)
//
// will return
//
//	[]string{"name", "age"}, nil
//
// while
//
//	parseBytes([]byte(`(?P<name>\w+) (?:\w+) (\d+)`), 0)
//
//	will return
//
//	[]string{name, UnnamedCapture}, nil
func parseBytes(input []byte) ([]string, error) {
	if len(input) == 0 {
		return nil, nil
	}
	nextrun, runlen := utf8.DecodeRune(input)
	// Advance input to next runee
	input = input[runlen:]
	switch nextrun {
	case '\\':
		// Scape character, skip next rune
		if len(input) == 0 {
			return nil, errors.New("error parsing named regexp: trailing backslash at end of expression")
		}
		if nextrun, runlen := utf8.DecodeRune(input); nextrun != utf8.RuneError {
			// effectively skip next rune
			input = input[runlen:]
		} else {
			return nil, errors.New("error parsing named regexp: incorrect rune after backslash")
		}
		return parseBytes(input)
	case '(':
		var groupName []string
		if m := reNamedMatch.FindSubmatchIndex(input); m != nil {
			// Named pattern ?P<name>
			// return name of the group found
			groupName = []string{string(input[m[2]:m[3]])}
			// skip "?p<name>", "(" already skipped
			input = input[m[3]+1:]
		} else if m := reNoCaptureMatch.FindSubmatchIndex(input); m != nil {
			// no capturing group
			// nothing to return
			// skip "?:", "(" already skipped
			input = input[m[1]+1:]
		} else {
			// capture with no name
			// return unnamedCapture
			groupName = []string{UnnamedCapture}
			// "(" already skipped
		}
		// Parse the rest
		if recursiveResult, err := parseBytes(input); err == nil {
			// Prepend the named match found to the rest of the named groups names
			// that are parsed recursively
			return append(groupName, recursiveResult...), nil
		} else {
			return nil, err
		}
	default:
		return parseBytes(input)
	}
}

func buildMap(namedMatches []string) (map[string]int, error) {
	r := make(map[string]int)
	for i, name := range namedMatches {
		if name != UnnamedCapture {
			if _, ok := r[name]; ok {
				return nil, errors.New("error parsing named regexp: duplicate named group")
			}
			r[name] = i + 1
		}
	}
	return r, nil
}

// Compile is the 'Compile' version of regexp.Compile; it returns a RegexpNamed
// object that can be used to match against text.
//
// The regular expression syntax is the same as that of the regexp package,
// but allows matching against named submatches, using the methods
// FindNamed, FindIndexNamed, FindStringNamed, FindStringIndexNamed,
// FindAllNamed, FindAllIndexNamed, FindAllStringNamed and
// FindAllStringIndexNamed; see the documentation of those methods in
// this package for details.
//
// The methods of regexp package can be used with the RegexpNamed type.
//
// If the expression is malformed, or if a named group is duplicated, an
// error is returned.
//
// See regexp.Compile for more information.
func Compile(re string) (RegexpNamed, error) {
	compiledRe, err := regexp.Compile(re)
	if err != nil {
		return RegexpNamed{nil, nil}, err
	}
	if parsed, err := parseBytes([]byte(re)); err != nil {
		return RegexpNamed{nil, nil}, err
	} else {
		if map_, err := buildMap(parsed); err != nil {
			return RegexpNamed{nil, nil}, err
		} else {
			return RegexpNamed{map_, compiledRe}, nil
		}
	}
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular
// expressions.
func MustCompile(re string) RegexpNamed {
	r, err := Compile(re)
	if err != nil {
		panic(`regexp_named: Compile(` + quote(re) + `): ` + err.Error())
	}
	return r
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}

func getResult[T any](match []T, pos int) T {
	return match[pos]
}

func getResultIndex[T any](match []T, pos int) []T {
	// match[pos*2], match[pos*2+1]
	return match[pos*2 : pos*2+2]
}

// Helper function to map the result of the regexp.findSubmatch_ functions
// to a map of the names of the named groups of the regexp.
// It takes a list of submatches, and returns the result of apply resultFunc
// to all the submatches.
func mapRe[T, S any](re *RegexpNamed, submatches []T, resultFunc func([]T, int) S) (S, map[string]S) {
	if submatches == nil {
		return *new(S), nil
	}
	rv := make(map[string]S)
	for k, v := range re.namedMap {
		rv[k] = resultFunc(submatches, v)
	}
	return resultFunc(submatches, 0), rv
}

// Helper function to map the result of the regexp.findSubmatch_All functions
// to a map of the names of the named groups of the regexp.
// Calls mapRe for each match.
func mapReAll[T, S any](re *RegexpNamed, matches [][]T, composeFunc func([]T, int) S) ([]S, []map[string]S) {
	rv := make([]map[string]S, 0)
	rv0 := make([]S, 0)
	for _, submatches := range matches {
		base, namedSubmatches := mapRe(re, submatches, composeFunc)
		rv0 = append(rv0, base)
		rv = append(rv, namedSubmatches)
	}
	return rv0, rv
}

// FindNamed returns a map of named submatches matched by re in b.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindSubmatch for a description of the return value.
func (re *RegexpNamed) FindNamed(s []byte) ([]byte, map[string][]byte) {
	return mapRe(re, re.FindSubmatch(s), getResult)
}

// FindIndexNamed returns a map of named index pairs identifying the
// matched subexpressions matched by re in b.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindSubmatchIndex for a description of the return value.
func (re *RegexpNamed) FindIndexNamed(s []byte) ([]int, map[string][]int) {
	return mapRe(re, re.FindSubmatchIndex(s), getResultIndex)
}

// FindStringNamed returns a map of named submatches matched by re in s.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindStringSubmatch for a description of the return value.
func (re *RegexpNamed) FindStringNamed(s string) (string, map[string]string) {
	return mapRe(re, re.FindStringSubmatch(s), getResult)
}

// FindStringIndexNamed returns a map of named index pairs identifying the
// matched subexpressions matched by re in s.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindStringSubmatchIndex for a description of the return value.
func (re *RegexpNamed) FindStringIndexNamed(s string) ([]int, map[string][]int) {
	return mapRe(re, re.FindStringSubmatchIndex(s), getResultIndex)
}

// FindAllNamed is the 'All' version of FindNamed; it returns a slice of all
// successive maps of named submatches matched by re in b.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllSubmatch for a description of the return value.
func (re *RegexpNamed) FindAllNamed(b []byte, n int) ([][]byte, []map[string][]byte) {
	return mapReAll(re, re.FindAllSubmatch(b, n), getResult)
}

// FindAllIndexNamed is the 'All' version of FindIndexNamed; it returns a slice
// of all successive maps of named index pairs identifying the successive
// matches of re in b.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllSubmatchIndex for a description of the return value.
func (re *RegexpNamed) FindAllIndexNamed(b []byte, n int) ([][]int, []map[string][]int) {
	return mapReAll(re, re.FindAllSubmatchIndex(b, n), getResultIndex)
}

// FindAllStringNamed is the 'All' version of FindStringNamed; it returns a
// slice of all successive maps of named submatches matched by re in s.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllStringSubmatch for a description of the return value.
func (re *RegexpNamed) FindAllStringNamed(s string, n int) ([]string, []map[string]string) {
	return mapReAll(re, re.FindAllStringSubmatch(s, n), getResult)
}

// FindAllStringIndexNamed is the 'All' version of FindStringIndexNamed; it
// returns a slice of all successive maps of named index pairs identifying the
// successive matches of re in s.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllStringSubmatchIndex for a description of the return value.
func (re *RegexpNamed) FindAllStringIndexNamed(s string, n int) ([][]int, []map[string][]int) {
	return mapReAll(re, re.FindAllStringSubmatchIndex(s, n), getResultIndex)
}
