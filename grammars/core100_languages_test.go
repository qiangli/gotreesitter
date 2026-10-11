package grammars

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

// core100ExpectedLen is the approved inventory size: the original 100 minus
// caddy, disassembly and ebnf, which the fork dropped as GPL-3.0 (see FORK.md).
const core100ExpectedLen = 97

var forbiddenGrammars = []string{"caddy", "disassembly", "ebnf", "jq", "nim"}

func TestCore100LanguageNames(t *testing.T) {
	names := Core100LanguageNames()
	if len(names) != core100ExpectedLen {
		t.Fatalf("expected Core100 length %d, got %d", core100ExpectedLen, len(names))
	}

	seen := make(map[string]struct{}, len(names))
	for _, name := range names {
		if _, dup := seen[name]; dup {
			t.Fatalf("duplicate Core100 language %q", name)
		}
		seen[name] = struct{}{}
	}

	entries := AllLanguages()
	entryByName := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		entryByName[entry.Name] = struct{}{}
	}

	for _, name := range names {
		if _, ok := entryByName[name]; !ok {
			t.Fatalf("Core100 language %q is not registered", name)
		}
	}

	for _, name := range forbiddenGrammars {
		if _, ok := seen[name]; ok {
			t.Errorf("non-permissive grammar %q must not be in Core100", name)
		}
		if _, ok := entryByName[name]; ok {
			t.Errorf("non-permissive grammar %q must not be registered", name)
		}
	}
}

func TestCore100AllowListMatchesLanguageNames(t *testing.T) {
	f, err := os.Open("update_tier1_core100.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	allow := map[string]struct{}{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		allow[line] = struct{}{}
	}
	if err := sc.Err(); err != nil {
		t.Fatal(err)
	}

	names := Core100LanguageNames()
	if len(allow) != len(names) {
		t.Fatalf("update_tier1_core100.txt has %d entries, Core100 has %d", len(allow), len(names))
	}
	for _, name := range names {
		if _, ok := allow[name]; !ok {
			t.Errorf("Core100 language %q missing from update_tier1_core100.txt", name)
		}
	}
}
