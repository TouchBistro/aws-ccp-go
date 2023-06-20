package util

import "testing"

func TestCoalesce(t *testing.T) {

	defaultStr := "default"
	var s *string
	val := Coalesce(s, defaultStr)
	if val != "default" {
		t.Errorf("expected %v, found %v", defaultStr, val)
	}

	str := "str"
	s = &str
	val = Coalesce(s, defaultStr)
	if val != str {
		t.Errorf("expected %v, found %v", str, val)
	}

}
