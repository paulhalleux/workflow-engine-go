package semver

import "testing"

func TestParseValid(t *testing.T) {
	v, err := Parse("v1.2.3-alpha+001")
	if err != nil {
		t.Fatal(err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("unexpected parsed version: %+v", v)
	}
	if v.PreRelease != "alpha" || v.Build != "001" {
		t.Fatalf("unexpected pre/build: %s %s", v.PreRelease, v.Build)
	}
}

func TestBumps(t *testing.T) {
	if s := Ensure(""); s != InitialVersion() {
		t.Fatalf("expected initial for empty, got %s", s)
	}

	s, _ := BumpPatch("1.2.3")
	if s != "1.2.4" {
		t.Fatalf("expected 1.2.4, got %s", s)
	}

	s, _ = BumpMinor("1.2.3")
	if s != "1.3.0" {
		t.Fatalf("expected 1.3.0, got %s", s)
	}

	s, _ = BumpMajor("1.2.3")
	if s != "2.0.0" {
		t.Fatalf("expected 2.0.0, got %s", s)
	}
}
