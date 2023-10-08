package regex_named

import (
	"testing"
)

func TestFindStringNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindStringNamed("foo 42"); m == nil {
		t.Error("Expected match")
	} else {
		if m0 != "foo 42" {
			t.Error("Expected match to be foo 42")
		}
		if m["name"] != "foo" {
			t.Error("Expected name to be foo")
		}
		if m["age"] != "42" {
			t.Error("Expected age to be 42")
		}
	}
}

func TestFindAllStringNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+)? (?P<age>\d+)`)
	if m0, m := re.FindAllStringNamed("foo 42 43", -1); m == nil {
		t.Error("Expected match")
	} else {
		if m0[0] != "foo 42" {
			t.Error("Expected match to be foo 42")
		}
		if m[0]["name"] != "foo" {
			t.Error("Expected name to be foo")
		}
		if m[0]["age"] != "42" {
			t.Error("Expected age to be 42")
		}
		if m0[1] != " 43" {
			t.Error("Expected match to be  43")
		}
		if m[1]["name"] != "" {
			t.Error("Expected name to be bar")
		}
		if m[1]["age"] != "43" {
			t.Error("Expected age to be 43")
		}
	}
}

func TestFindNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindNamed([]byte("foo 42")); m == nil {
		t.Error("Expected match")
	} else {
		if string(m0) != "foo 42" {
			t.Error("Expected match to be foo 42")
		}
		if string(m["name"]) != "foo" {
			t.Error("Expected name to be foo")
		}
		if string(m["age"]) != "42" {
			t.Error("Expected age to be 42")
		}
	}
}

func TestFindAllNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindAllNamed([]byte("foo 42 bar 43"), -1); m == nil {
		t.Error("Expected match")
	} else {
		if string(m0[0]) != "foo 42" {
			t.Error("Expected match to be foo 42")
		}
		if string(m[0]["name"]) != "foo" {
			t.Error("Expected name to be foo")
		}
		if string(m[0]["age"]) != "42" {
			t.Error("Expected age to be 42")
		}
		if string(m0[1]) != "bar 43" {
			t.Error("Expected match to be bar 43")
		}
		if string(m[1]["name"]) != "bar" {
			t.Error("Expected name to be bar")
		}
		if string(m[1]["age"]) != "43" {
			t.Error("Expected age to be 43")
		}
	}
}

func sliceEq[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestFindIndexNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindIndexNamed([]byte("foo 42")); m == nil {
		t.Error("Expected match")
	} else {
		if !sliceEq(m0, []int{0, 6}) {
			t.Error("Expected match to be {0, 6}")
		}
		if !sliceEq(m["name"], []int{0, 3}) {
			t.Error("Expected name to be {0, 3}")
		}
		if !sliceEq(m["age"], []int{4, 6}) {
			t.Error("Expected age to be {4, 6}")
		}
	}
}

func TestFindStringIndexNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindStringIndexNamed("foo 42"); m == nil {
		t.Error("Expected match")
	} else {
		if !sliceEq(m0, []int{0, 6}) {
			t.Error("Expected match to be {0, 6}")
		}
		if !sliceEq(m["name"], []int{0, 3}) {
			t.Error("Expected name to be {0, 3}")
		}
		if !sliceEq(m["age"], []int{4, 6}) {
			t.Error("Expected age to be {4, 6}")
		}
	}
}

func TestFindAllIndexNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindAllIndexNamed([]byte("foo 42 bar 43"), -1); m == nil {
		t.Error("Expected match")
	} else {
		if !sliceEq(m0[0], []int{0, 6}) {
			t.Error("Expected match to be {0, 6}")
		}
		if !sliceEq(m[0]["name"], []int{0, 3}) {
			t.Error("Expected name to be {0, 3}")
		}
		if !sliceEq(m[0]["age"], []int{4, 6}) {
			t.Error("Expected age to be {4, 6}")
		}
		if !sliceEq(m0[1], []int{7, 13}) {
			t.Error("Expected match to be {7, 13}")
		}
		if !sliceEq(m[1]["name"], []int{7, 10}) {
			t.Error("Expected name to be {7, 10}")
		}
		if !sliceEq(m[1]["age"], []int{11, 13}) {
			t.Error("Expected age to be {11, 13}")
		}
	}
}

func TestFindAllStringIndexnamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	if m0, m := re.FindAllStringIndexNamed("foo 42 bar 43", -1); m == nil {
		t.Error("Expected match")
	} else {
		if !sliceEq(m0[0], []int{0, 6}) {
			t.Error("Expected match to be {0, 6}")
		}
		if !sliceEq(m[0]["name"], []int{0, 3}) {
			t.Error("Expected name to be {0, 3}")
		}
		if !sliceEq(m[0]["age"], []int{4, 6}) {
			t.Error("Expected age to be {4, 6}")
		}
		if !sliceEq(m0[1], []int{7, 13}) {
			t.Error("Expected match to be {7, 13}")
		}
		if !sliceEq(m[1]["name"], []int{7, 10}) {
			t.Error("Expected name to be {7, 10}")
		}
		if !sliceEq(m[1]["age"], []int{11, 13}) {
			t.Error("Expected age to be {11, 13}")
		}
	}
}

func TestNoCapture(t *testing.T) {
	re := MustCompile(`(?:\w+) (\d+)`)
	if m0, m := re.FindStringNamed("foo 42"); m == nil {
		t.Error("Expected match")
	} else {
		if m0 != "foo 42" {
			t.Error("Expected match to be foo 42")
		}
		if len(m) != 0 {
			t.Error("Expected no named match")
		}
	}
}

func TestNested(t *testing.T) {
	re := MustCompile(`(?P<a>(?:1(?:2)?)*)(?P<b>3)`)
	if m0, m := re.FindStringNamed("1211121123"); m == nil {
		t.Error("Expected match")
	} else {
		if m0 != "1211121123" {
			t.Error("Expected match to be  1211121123")
		}
		if m["a"] != "121112112" {
			t.Error("Expected a to be 121112112")
		}
		if m["b"] != "3" {
			t.Error("Expected b to be 3")
		}
		if _, ok := m["2"]; ok {
			t.Error("Expected no 2")
		}
	}
}

func TestDuplicatedName(t *testing.T) {
	_, err := Compile(`(?P<name>\w+) (?P<name>\d+)`)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestMalformed(t *testing.T) {
	_, err := Compile(`(?P<name>\w+) (?P<age>\d+`)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestMalformedTrailingBackslash(t *testing.T) {
	_, err := Compile(`(?P<name>\w+)\`)
	if err == nil {
		t.Error("Expected error")
	}
}
