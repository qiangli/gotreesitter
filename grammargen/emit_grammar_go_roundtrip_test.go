package grammargen

import (
	"os"
	"testing"
)

// The tracked *_grammar.go files are EmitGrammarGo output; re-emitting them
// must be byte-identical so the generator stays the reproducible source.
func TestEmitGrammarGoReproducesTrackedGrammars(t *testing.T) {
	cases := []struct {
		file, fn string
		build    func() *Grammar
	}{
		{"go_grammar.go", "GoGrammar", GoGrammar},
		{"fortran_grammar.go", "FortranGrammar", FortranGrammar},
		{"javascript_grammar.go", "JavaScriptGrammar", JavaScriptGrammar},
	}
	for _, c := range cases {
		t.Run(c.fn, func(t *testing.T) {
			got, err := EmitGrammarGo(c.build(), "grammargen", c.fn)
			if err != nil {
				t.Fatalf("EmitGrammarGo: %v", err)
			}
			want, err := os.ReadFile(c.file)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != string(want) {
				t.Fatalf("re-emitted %s differs from tracked file (got %d bytes, want %d)", c.file, len(got), len(want))
			}
		})
	}
}
