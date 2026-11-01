package gotreesitter

import "testing"

type statelessRelexTestScanner struct {
	parserTestUnsafeExternalScanner
}

func (statelessRelexTestScanner) SupportsStatelessRelex() bool { return true }

type disabledRelexTestScanner struct {
	parserTestUnsafeExternalScanner
}

func (disabledRelexTestScanner) SupportsStatelessRelex() bool { return false }

func TestStatelessExternalRelex(t *testing.T) {
	for _, tc := range []struct {
		name          string
		scanner       ExternalScanner
		state         StateID
		symbol        Symbol
		originalValid bool
		want          bool
	}{
		{"stale external", statelessRelexTestScanner{}, 2, 4, false, true},
		{"scanner without capability", parserTestUnsafeExternalScanner{}, 2, 4, false, false},
		{"scanner declines capability", disabledRelexTestScanner{}, 2, 4, false, false},
		{"nonexternal lookahead", statelessRelexTestScanner{}, 2, 1, false, false},
		{"external still valid", statelessRelexTestScanner{}, 2, 4, true, false},
		{"replacement has no action", statelessRelexTestScanner{}, 0, 4, false, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lang := buildArithmeticLanguage()
			lang.ExternalScanner = tc.scanner
			lang.ExternalSymbols = []Symbol{4}
			lang.ExternalLexStates = [][]bool{{false}, {true}}
			lang.LexModes[tc.state].ExternalLexState = 1
			for i := range lang.ParseTable {
				lang.ParseTable[i] = append(lang.ParseTable[i], 0)
			}
			if tc.originalValid {
				lang.ParseTable[tc.state][4] = 4
			}
			parser := NewParser(lang)
			lexer := NewLexer(lang.LexStates, []byte("+\n12"))
			lexer.pos, lexer.row, lexer.col = 2, 1, 0
			source := &dfaTokenSource{lexer: lexer, language: lang, state: tc.state, lookupActionIndex: parser.lookupActionIndex}
			original := Token{Symbol: tc.symbol, StartByte: 0, EndByte: 2, EndPoint: Point{Row: 1}}
			tok, ok := parser.tryRelexCurrentStateDFA(original, tc.state, source, []glrStack{{entries: []stackEntry{{state: tc.state}}}})
			if ok != tc.want {
				t.Fatalf("relex success = %v, want %v", ok, tc.want)
			}
			if tc.want {
				if tok.Symbol != 2 || tok.EndByte != 1 || lexer.pos != 1 || lexer.row != 0 || lexer.col != 1 {
					t.Fatalf("replacement/cursor = %+v / %d:%d:%d", tok, lexer.pos, lexer.row, lexer.col)
				}
			} else if lexer.pos != 2 || lexer.row != 1 || lexer.col != 0 {
				t.Fatalf("failed speculation changed cursor: %d:%d:%d", lexer.pos, lexer.row, lexer.col)
			}
			if source.state != tc.state {
				t.Fatal("speculation changed parser state")
			}
			if tc.want {
				source.SetParserState(0)
				next := source.Next()
				if next.Symbol != 1 || next.StartByte != 2 || next.EndByte != 4 || next.StartPoint.Row != 1 {
					t.Fatalf("short replacement lost following source: %+v", next)
				}
			}
		})
	}
}

func TestStatelessExternalRelexAdapterCapability(t *testing.T) {
	for _, tc := range []struct {
		scanner ExternalScanner
		want    bool
	}{
		{statelessRelexTestScanner{}, true}, {disabledRelexTestScanner{}, false}, {parserTestUnsafeExternalScanner{}, false}, {nil, false},
	} {
		adapter := &externalScannerOrderAdapter{inner: tc.scanner}
		if got := adapter.SupportsStatelessRelex(); got != tc.want {
			t.Fatalf("adapter capability = %v, want %v", got, tc.want)
		}
	}
}

func TestStatelessExternalRelexRejectsSharedLookahead(t *testing.T) {
	for _, stacks := range [][]glrStack{
		{{entries: []stackEntry{{state: 2}}}, {entries: []stackEntry{{state: 2}}, shifted: true}},
		{{entries: []stackEntry{{state: 2}}}, {entries: []stackEntry{{state: 0}}}},
		{{entries: []stackEntry{{state: 2}}}, {entries: []stackEntry{{state: 2}}}},
		{{entries: []stackEntry{{state: 2}}, byteOffset: 2}},
		{{entries: []stackEntry{{state: 2}}, shifted: true}},
	} {
		lang := buildArithmeticLanguage()
		lang.ExternalScanner = statelessRelexTestScanner{}
		lang.ExternalSymbols = []Symbol{4}
		lang.LexModes[2].ExternalLexState = 1
		parser := NewParser(lang)
		lexer := NewLexer(lang.LexStates, []byte("+\n"))
		lexer.pos, lexer.row = 2, 1
		source := &dfaTokenSource{lexer: lexer, language: lang, state: 2, lookupActionIndex: parser.lookupActionIndex}
		if _, ok := parser.tryRelexCurrentStateDFA(Token{Symbol: 4, EndByte: 2}, 2, source, stacks); ok {
			t.Fatal("relexed a shared lookahead")
		}
		if lexer.pos != 2 || lexer.row != 1 || lexer.col != 0 || source.state != 2 {
			t.Fatal("guard changed cursor/state")
		}
	}
}
