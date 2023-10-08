// Package regex named provides named submatches for Go's regexp package.
//
// The package extends the regexp package with the following methods:
//
// 	FindNamed
// 	FindIndexNamed
// 	FindStringNamed
// 	FindStringIndexNamed
// 	FindAllNamed
// 	FindAllIndexNamed
// 	FindAllStringNamed
// 	FindAllStringIndexNamed
//
// These methods work like the corresponding methods in the regexp, replacing
// the slices returned by the corresponding methods for maps indexed by the
// names of the groups.
//
// RegexNamed are created with the Compile and MustCompile functions, which
// work like the corresponding functions in the regexp package.
//
// For example:
//
// 	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
//  m0, m := re.FindStringNamed("foo 42")
//
// m0 will be "foo 42" and m will be a map[string]string with the following
// values:
//
// 	m["name"] == "foo"
// 	m["age"] == "42"
//
// If a group is not matched, the corresponding value in the map will be an
// empty string.

package regex_named

import (
	"errors"
	"regexp"
	"unicode"
	"unicode/utf8"
)

type RegexNamed struct {
	namedMap map[string]int
	*regexp.Regexp
}

const errorString = string(unicode.ReplacementChar)

func parseBytes(in []byte, until rune) ([]string, error) {
	reNamedMatch := regexp.MustCompile(`^\?P\<(.*?)\>`)
	reNoCaptureMatch := regexp.MustCompile(`^\?:`)
	if len(in) == 0 {
		return nil, nil
	}
	nextrun, runlen := utf8.DecodeRune(in)
	in = in[runlen:]
	if nextrun == until {
		return parseBytes(in, 0)
	}
	var prependReturn []string
	switch nextrun {
	case '\\':
		if len(in) == 0 {
			return nil, errors.New("malformed - trailing \\")
		}
		// Skip scaped rune
		if nextrun, runlen := utf8.DecodeRune(in); nextrun != utf8.RuneError {
			in = in[runlen:]
		} else {
			return nil, errors.New("incorrect rune")
		}
		return parseBytes(in, 0)
	case '(':
		untilRecursive := ')'
		if m := reNamedMatch.FindSubmatchIndex(in); m != nil {
			// Named pattern ?P<name>, return name
			prependReturn = []string{string(in[m[2]:m[3]])}
			// skip ?p<name>
			in = in[m[3]+1:]
		} else if m := reNoCaptureMatch.FindSubmatchIndex(in); m != nil {
			// no capture, skip ?: part
			in = in[m[1]+1:]
		} else {
			prependReturn = []string{errorString}
		}
		if rec, err := parseBytes(in, untilRecursive); err == nil {
			return append(prependReturn, rec...), nil
		} else {
			return nil, err
		}
	default:
		return parseBytes(in, 0)
	}
}

func parse(in string) ([]string, error) {
	return parseBytes([]byte(in), 0)
}

func buildMap(namedMatches []string) (map[string]int, error) {
	r := make(map[string]int)
	for i, name := range namedMatches {
		if name != errorString {
			if _, ok := r[name]; ok {
				return nil, errors.New("duplicate name")
			}
			r[name] = i + 1
		}
	}
	return r, nil
}

func Compile(re string) (RegexNamed, error) {
	compiledRe, err := regexp.Compile(re)
	if err != nil {
		return RegexNamed{nil, nil}, err
	}
	if parsed, err := parse(re); err != nil {
		return RegexNamed{nil, nil}, err
	} else {
		if map_, err := buildMap(parsed); err != nil {
			return RegexNamed{nil, nil}, err
		} else {
			return RegexNamed{map_, compiledRe}, nil
		}
	}
}

func MustCompile(re string) RegexNamed {
	r, err := Compile(re)
	if err != nil {
		panic(err)
	}
	return r
}

func composeMap[T any](match []T, pos int) T {
	return match[pos]
}

func composeIndex[T any](match []T, pos int) []T {
	return match[pos*2 : pos*2+2]
}

func mapRe[T, S any](re *RegexNamed, match []T, f func([]T, int) S) (S, map[string]S) {
	if match == nil {
		return *new(S), nil
	}
	rv := make(map[string]S)
	for k, v := range re.namedMap {
		rv[k] = f(match, v)
	}
	return f(match, 0), rv
}

func mapReAll[T, S any](re *RegexNamed, match [][]T, f func([]T, int) S) ([]S, []map[string]S) {
	rv := make([]map[string]S, 0)
	rv0 := make([]S, 0)
	for _, m := range match {
		base, mm := mapRe(re, m, f)
		rv0 = append(rv0, base)
		rv = append(rv, mm)
	}
	return rv0, rv
}

// FindNamed returns a map of named submatches matched by re in b.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindSubmatch for a description of the return value.
func (re *RegexNamed) FindNamed(s []byte) ([]byte, map[string][]byte) {
	return mapRe(re, re.FindSubmatch(s), composeMap)
}

// FindIndexNamed returns a map of named index pairs identifying the
// matched subexpressions matched by re in b.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindSubmatchIndex for a description of the return value.
func (re *RegexNamed) FindIndexNamed(s []byte) ([]int, map[string][]int) {
	return mapRe(re, re.FindSubmatchIndex(s), composeIndex)
}

// FindStringNamed returns a map of named submatches matched by re in s.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindStringSubmatch for a description of the return value.
func (re *RegexNamed) FindStringNamed(s string) (string, map[string]string) {
	return mapRe(re, re.FindStringSubmatch(s), composeMap)
}

// FindStringIndexNamed returns a map of named index pairs identifying the
// matched subexpressions matched by re in s.
// The match itself is returned as the first element of the result.
// If there are no matches, nil is returned.
// See (*Regexp).FindStringSubmatchIndex for a description of the return value.
func (re *RegexNamed) FindStringIndexNamed(s string) ([]int, map[string][]int) {
	return mapRe(re, re.FindStringSubmatchIndex(s), composeIndex)
}

// FindAllNamed is the 'All' version of FindNamed; it returns a slice of all
// successive maps of named submatches matched by re in b.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllSubmatch for a description of the return value.
func (re *RegexNamed) FindAllNamed(b []byte, n int) ([][]byte, []map[string][]byte) {
	return mapReAll(re, re.FindAllSubmatch(b, n), composeMap)
}

// FindAllIndexNamed is the 'All' version of FindIndexNamed; it returns a slice
// of all successive maps of named index pairs identifying the successive
// matches of re in b.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllSubmatchIndex for a description of the return value.
func (re *RegexNamed) FindAllIndexNamed(b []byte, n int) ([][]int, []map[string][]int) {
	return mapReAll(re, re.FindAllSubmatchIndex(b, n), composeIndex)
}

// FindAllStringNamed is the 'All' version of FindStringNamed; it returns a
// slice of all successive maps of named submatches matched by re in s.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllStringSubmatch for a description of the return value.
func (re *RegexNamed) FindAllStringNamed(s string, n int) ([]string, []map[string]string) {
	return mapReAll(re, re.FindAllStringSubmatch(s, n), composeMap)
}

// FindAllStringIndexNamed is the 'All' version of FindStringIndexNamed; it
// returns a slice of all successive maps of named index pairs identifying the
// successive matches of re in s.
// The match itself is returned as the first element of the result.
// A return value of nil indicates no match.
// See (*Regexp).FindAllStringSubmatchIndex for a description of the return value.
func (re *RegexNamed) FindAllStringIndexNamed(s string, n int) ([][]int, []map[string][]int) {
	return mapReAll(re, re.FindAllStringSubmatchIndex(s, n), composeIndex)
}
