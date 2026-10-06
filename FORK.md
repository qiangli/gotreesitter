# Fork provenance: qiangli/gotreesitter

This repository is a **minimal pinned fork of
[github.com/odvcencio/gotreesitter](https://github.com/odvcencio/gotreesitter)**,
imported verbatim from the upstream module **v0.16.0** (the first commit in
this repository's history is the pristine module snapshot; its `go.sum` line
upstream is `h1:pg29HG6idOpmbh6BkPc6ncFP0C8rH9ISAzFVNdHm2ss=`). The fork is
hosted as qiangli/gotreesitter but **keeps upstream's module path**
`github.com/odvcencio/gotreesitter` on purpose: a single path `replace`
directive in the consumer's go.mod then redirects EVERY importer to this
blob-free copy — yoke's direct imports and, critically, the transitive
`github.com/qiangli/gfy/pkg/extract` imports (gfy pins odvcencio v0.15.3,
which embeds the same five grammars), so renaming the module path would leave
the GPL/MPL bytes linked through gfy.

## Status: permanent dhnt-internal fork

This fork is maintained for dhnt/bashy use only (operator decision
2026-10-06, Sprint 369). No upstream PR is pursued: an exclusion-tag change
of this kind would not be accepted upstream, so the `replace` directives in
yoke and bashy stay. The `exclude-grammars-tag` branch keeps a reference
implementation of the tag approach. Re-sync policy: rebase onto new upstream
releases only for grammar updates, re-deleting the same five grammars.

## Why the fork exists

gotreesitter embeds pre-compiled tree-sitter grammar parse tables for 206
languages. A grammar table is data derived from that grammar's `parser.c`, so
**each grammar's own license applies to the bytes linked into any binary that
embeds it**. Five of the 206 are non-permissive:

| grammar  | upstream grammar repository                          | license  |
|----------|------------------------------------------------------|----------|
| caddy    | https://github.com/opa-oz/tree-sitter-caddy           | GPL-3.0  |
| disassembly | https://github.com/ColinKennedy/tree-sitter-disassembly | GPL-3.0 |
| ebnf     | https://github.com/RubixDev/ebnf                     | GPL-3.0  |
| jq       | https://github.com/nverno/tree-sitter-jq             | GPL-3.0  |
| nim      | https://github.com/alaviss/tree-sitter-nim           | MPL-2.0  |

bashy's licensing policy (`docs/licensing-supply-chain-policy.md` §1) is
"compiled-in means permissive only". Operator decision 2026-10-02: keep the
other 201 grammars (attributed in `yoke/THIRD_PARTY_GRAMMARS.md`) and drop
these five.

Upstream offers **no per-grammar exclusion that removes the bytes**: the
`grammar_blobs/*.bin` wildcard in `blob_source_embedded.go` embeds every blob,
the explicit embed list in `blob_source_embedded_core.go` still names caddy,
disassembly and ebnf, the `grammar_subset` build tags only gate *registration*
(not embedding), and `GOTREESITTER_GRAMMAR_SET` is a runtime filter over bytes
that are already linked in. The only way to keep the other 201 while shipping
none of these five's bytes is to delete them from a fork — hence this repo.

## What the fork changes (everything else is upstream v0.16.0, untouched)

- deleted for each of the five grammars: `grammar_blobs/<name>.bin`,
  `<name>_register.go`, `<name>_scanner.go` (where present),
  `z_subset_registry_register_<name>.go`, `z_subset_scanner_register_<name>.go`
  (where present), and the `grammars/languages.lock` line.
- removed their entries from the generated/reference sites:
  `registry_builtin_gen.go`, `embedded_grammars_gen.go` (loader funcs),
  `blob_source_embedded_core.go`, `core100_languages.go`, `smoke_samples.go`,
  `linguist_gen.go`, `zzz_scanner_attachments.go`, `languages.lock`, and the
  caddy/nim parser-result compatibility normalizers in
  `parser_result_compat.go` / `parser_result_misc_spans.go`, plus their cases
  in the fork's own tests and `cmd/gen_linguist` / `cmd/harnessgate` tables.
- `grammars/languages.lock` now pins 201 grammars, all permissive.

## Re-sync policy

This fork intentionally tracks upstream v0.16.0 with only the deletions
above (the module path stays odvcencio so the one-line `replace` keeps
working). On an upstream bump: re-import the new version as a pristine commit,
re-apply this deletion list (the five names), and regenerate
`yoke/THIRD_PARTY_GRAMMARS.md` from the fork's `grammars/languages.lock`
(yoke `scripts/grammar-licenses.sh`) — the attribution report and the lock
file must never disagree. If upstream gains an upstreamable per-grammar
exclusion mechanism, prefer retiring this fork for it.
