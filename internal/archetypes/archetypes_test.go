package archetypes

import (
	"os"
	"testing"
)

func TestGetScope_Devops(t *testing.T) {
	scope, err := GetScope("devops")
	if err != nil {
		t.Fatalf("GetScope(devops) returned unexpected error: %v", err)
	}

	expected := []string{"infra", "ci-cd", "pipeline", "workflow", "ops"}
	if len(scope) != len(expected) {
		t.Fatalf("GetScope(devops) = %v (len %d); want %v (len %d)", scope, len(scope), expected, len(expected))
	}
	for i, tag := range expected {
		if scope[i] != tag {
			t.Errorf("GetScope(devops)[%d] = %q; want %q", i, scope[i], tag)
		}
	}
}

func TestGetScope_Architect(t *testing.T) {
	scope, err := GetScope("architect")
	if err != nil {
		t.Fatalf("GetScope(architect) returned unexpected error: %v", err)
	}

	expected := []string{"application-code", "go", "typescript", "svelte", "architecture", "design", "fullstack"}
	if len(scope) != len(expected) {
		t.Fatalf("GetScope(architect) = %v (len %d); want %v (len %d)", scope, len(scope), expected, len(expected))
	}
	for i, tag := range expected {
		if scope[i] != tag {
			t.Errorf("GetScope(architect)[%d] = %q; want %q", i, scope[i], tag)
		}
	}
}

func TestGetScope_CastrojoDocs(t *testing.T) {
	scope, err := GetScope("castrojo-docs")
	if err != nil {
		t.Fatalf("GetScope(castrojo-docs) returned unexpected error: %v", err)
	}

	expected := []string{"docs", "documentation", "writing"}
	if len(scope) != len(expected) {
		t.Fatalf("GetScope(castrojo-docs) = %v (len %d); want %v (len %d)", scope, len(scope), expected, len(expected))
	}
	for i, tag := range expected {
		if scope[i] != tag {
			t.Errorf("GetScope(castrojo-docs)[%d] = %q; want %q", i, scope[i], tag)
		}
	}
}

func TestGetScope_NonexistentSlug(t *testing.T) {
	// AC4: unknown slug returns empty slice, not error
	scope, err := GetScope("nonexistent-archetype-xyz")
	if err != nil {
		t.Fatalf("GetScope(nonexistent) returned unexpected error: %v", err)
	}
	if len(scope) != 0 {
		t.Errorf("GetScope(nonexistent) = %v; want empty slice", scope)
	}
}

func TestGetScope_NoScopeSection(t *testing.T) {
	// AC4: an archetype with no ## Scope section returns empty slice, not error.
	// "other" is a minimal archetype that has no Scope section.
	scope, err := GetScope("other")
	if err != nil {
		t.Fatalf("GetScope(other) returned unexpected error: %v", err)
	}
	if len(scope) != 0 {
		t.Errorf("GetScope(other) = %v; want empty slice (no scope section)", scope)
	}
}

func TestGetScope_CEO(t *testing.T) {
	scope, err := GetScope("ceo")
	if err != nil {
		t.Fatalf("GetScope(ceo) returned unexpected error: %v", err)
	}

	expected := []string{"delegation", "review", "planning"}
	if len(scope) != len(expected) {
		t.Fatalf("GetScope(ceo) = %v (len %d); want %v (len %d)", scope, len(scope), expected, len(expected))
	}
	for i, tag := range expected {
		if scope[i] != tag {
			t.Errorf("GetScope(ceo)[%d] = %q; want %q", i, scope[i], tag)
		}
	}
}

// TestGetScope_StopsAtNextHeading verifies that GetScope() returns ONLY the tags
// defined under "## Scope" and stops when it encounters the next "##" heading.
//
// Fixture: internal/archetypes/testdata/scope-boundary.md
//
//	## Scope
//	foo, bar
//
//	## Other Section
//	baz, qux
//
// Expected: ["foo", "bar"] — NOT ["foo", "bar", "baz", "qux"]
func TestGetScope_StopsAtNextHeading(t *testing.T) {
	// Point the overrides dir at the testdata directory so GetScope can read
	// the dedicated scope-boundary.md fixture without embedding it.
	orig := GetOverridesDir()
	SetOverridesDir("testdata")
	t.Cleanup(func() {
		SetOverridesDir(orig)
	})

	// Verify the fixture file exists so a missing file gives a clear failure.
	if _, err := os.Stat("testdata/scope-boundary.md"); err != nil {
		t.Fatalf("fixture file testdata/scope-boundary.md not found: %v", err)
	}

	scope, err := GetScope("scope-boundary")
	if err != nil {
		t.Fatalf("GetScope(scope-boundary) unexpected error: %v", err)
	}

	// AC3: must return ONLY the tags before the next ## heading.
	expected := []string{"foo", "bar"}
	if len(scope) != len(expected) {
		t.Fatalf("GetScope(scope-boundary) = %v (len %d); want %v (len %d) — parser may have leaked tags from the '## Other Section' boundary",
			scope, len(scope), expected, len(expected))
	}
	for i, tag := range expected {
		if scope[i] != tag {
			t.Errorf("GetScope(scope-boundary)[%d] = %q; want %q", i, scope[i], tag)
		}
	}

	// Explicit negative assertion: "baz" and "qux" appear ONLY after "## Other Section"
	// and must NOT be present in the returned tags.
	forbidden := []string{"baz", "qux"}
	for _, tag := range scope {
		for _, f := range forbidden {
			if tag == f {
				t.Errorf("GetScope(scope-boundary) returned %q which is from the post-boundary '## Other Section'; parser did not stop at ##", tag)
			}
		}
	}
}
