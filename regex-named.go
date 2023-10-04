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

func buildMap(namedMatches []string) map[string]int {
	r := make(map[string]int)
	for i, name := range namedMatches {
		if name != errorString {
			r[name] = i + 1
		}
	}
	return r
}

func Compile(re string) (RegexNamed, error) {
	compiledRe, err := regexp.Compile(re)
	if err != nil {
		return RegexNamed{nil, nil}, err
	}
	if parsed, err := parse(re); err == nil {
		return RegexNamed{buildMap(parsed), compiledRe}, nil
	} else {
		return RegexNamed{nil, nil}, err
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

func (re *RegexNamed) FindNamed(s []byte) ([]byte, map[string][]byte) {
	return mapRe(re, re.FindSubmatch(s), composeMap)
}

func (re *RegexNamed) FindIndexNamed(s []byte) ([]int, map[string][]int) {
	return mapRe(re, re.FindSubmatchIndex(s), composeIndex)
}

func (re *RegexNamed) FindStringNamed(s string) (string, map[string]string) {
	return mapRe(re, re.FindStringSubmatch(s), composeMap)
}

func (re *RegexNamed) FindStringIndexNamed(s string) ([]int, map[string][]int) {
	return mapRe(re, re.FindStringSubmatchIndex(s), composeIndex)
}

func (re *RegexNamed) FindAllNamed(s []byte, n int) ([][]byte, []map[string][]byte) {
	return mapReAll(re, re.FindAllSubmatch(s, n), composeMap)
}

func (re *RegexNamed) FindAllIndexNamed(s []byte, n int) ([][]int, []map[string][]int) {
	return mapReAll(re, re.FindAllSubmatchIndex(s, n), composeIndex)
}

func (re *RegexNamed) FindAllStringNamed(s string, n int) ([]string, []map[string]string) {
	return mapReAll(re, re.FindAllStringSubmatch(s, n), composeMap)
}

func (re *RegexNamed) FindAllStringIndexNamed(s string, n int) ([][]int, []map[string][]int) {
	return mapReAll(re, re.FindAllStringSubmatchIndex(s, n), composeIndex)
}
