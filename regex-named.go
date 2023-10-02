package regex_named

import (
	"errors"
	"regexp"
	"unicode"
	"unicode/utf8"
)

type RegexNamed struct {
	named_map map[string]int
	*regexp.Regexp
}

const error_string = string(unicode.ReplacementChar)

func parseBytes(in []byte, until rune) ([]string, error) {
	re_named_match := regexp.MustCompile(`^\?P\<(.*?)\>`)
	re_no_capture_match := regexp.MustCompile(`^\?:`)
	if len(in) == 0 {
		return nil, nil
	}
	nextrun, runlen := utf8.DecodeRune(in)
	in = in[runlen:]
	if nextrun == until {
		return parseBytes(in, 0)
	}
	var prepend_return []string
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
		until_recursive := ')'
		if m := re_named_match.FindSubmatchIndex(in); m != nil {
			// Named pattern ?P<name>, return name
			prepend_return = []string{string(in[m[2]:m[3]])}
			// skip ?p<name>
			in = in[m[3]+1:]
		} else if m := re_no_capture_match.FindSubmatchIndex(in); m != nil {
			// no capture, skip ?: part
			in = in[m[1]+1:]
		} else {
			prepend_return = []string{error_string}
		}
		if rec, err := parseBytes(in, until_recursive); err == nil {
			return append(prepend_return, rec...), nil
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
		if name != error_string {
			r[name] = i + 1
		}
	}
	return r
}

func Compile(re string) (RegexNamed, error) {
	compiled_re, err := regexp.Compile(re)
	if err != nil {
		return RegexNamed{nil, nil}, err
	}
	if parsed, err := parse(re); err == nil {
		return RegexNamed{buildMap(parsed), compiled_re}, nil
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

func (re *RegexNamed) FindStringNamed(s, name string) string {
	idx, ok := re.named_map[name]
	if !ok {
		return ""
	}
	match := re.FindStringSubmatch(s)
	return match[idx]
}
