package regex_named

import (
	"testing"
)

func TestFindStringNamed(t *testing.T) {
	re := MustCompile(`(?P<name>\w+) (?P<age>\d+)`)
	m := re.FindStringNamed("foo 42", "name")
	if m != "foo" {
		t.Errorf("Expected 'foo', got '%s'", m)
	}
	m = re.FindStringNamed("foo 42", "age")
	if m != "42" {
		t.Errorf("Expected '42', got '%s'", m)
	}
	m = re.FindStringNamed("foo 42", "foo")
	if m != "" {
		t.Errorf("Expected '', got '%s'", m)
	}
}
