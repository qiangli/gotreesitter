# qiangli/gotreesitter — permissive-only fork for dhnt/bashy

> **This is a permanent dhnt-internal fork** of
> [odvcencio/gotreesitter](https://github.com/odvcencio/gotreesitter).
> No upstream PR is pursued for it (operator decision 2026-10-06, Sprint 369).
> For library usage, benchmarks, and the query API, read the
> [original README](https://github.com/odvcencio/gotreesitter#readme) —
> everything documented there applies here except the five removed grammars.

## Why this fork exists

gotreesitter embeds pre-compiled tree-sitter grammar parse tables. A table is
data derived from that grammar's `parser.c`, so **each grammar's own license
applies to the bytes linked into any binary that embeds it**. Five of the 206
grammars are non-permissive, and bashy ships permissive-only binaries with no
exceptions:

| grammar     | upstream grammar repository                             | license |
|-------------|---------------------------------------------------------|---------|
| caddy       | https://github.com/opa-oz/tree-sitter-caddy             | GPL-3.0 |
| disassembly | https://github.com/ColinKennedy/tree-sitter-disassembly | GPL-3.0 |
| ebnf        | https://github.com/RubixDev/ebnf                        | GPL-3.0 |
| jq          | https://github.com/nverno/tree-sitter-jq                | GPL-3.0 |
| nim         | https://github.com/alaviss/tree-sitter-nim              | MPL-2.0 |

Upstream offers no exclusion mechanism that removes embedded bytes (subset
tags gate registration only; `GOTREESITTER_GRAMMAR_SET` filters at runtime
over already-linked bytes), so the bytes could only be dropped in a fork.
The other 201 grammars are kept and attributed in
[yoke/THIRD_PARTY_GRAMMARS.md](https://github.com/qiangli/yoke/blob/main/THIRD_PARTY_GRAMMARS.md).

## What changed vs upstream v0.16.0

Base: upstream module `v0.16.0` imported verbatim (first commit in this
repository's history is the pristine snapshot; its `go.sum` line upstream is
`h1:pg29HG6idOpmbh6BkPc6ncFP0C8rH9ISAzFVNdHm2ss=`). On top of it, exactly one
content commit plus docs:

- Deleted the five `grammars/grammar_blobs/*.bin` parse tables (689,429 bytes).
- Deleted their `grammars/*_register.go` / `grammars/*_scanner.go` files,
  subset-registry entries, lock lines, smoke samples, and test references.
- Removed dead code referencing the dropped languages
  (`normalizeNimTopLevelCallEnd`, the `caddy` case in trailing-newline span
  normalization, harness confidence lists).
- **Nothing else changed.** A diff of this repository against upstream
  `v0.16.0` that touches anything outside the five grammars is a defect —
  report it.

Full provenance, license table, and re-sync policy: [FORK.md](FORK.md).

## Consuming it

The module path is kept as upstream's (`github.com/odvcencio/gotreesitter`)
on purpose, so one path `replace` redirects every importer — direct and
transitive — to this blob-free copy:

```go
replace github.com/odvcencio/gotreesitter => ../gotreesitter
```

(yoke and bashy pin the reviewed commit via `.sibling-pins`; standalone
checkouts hydrate it with `scripts/bootstrap-siblings.sh`.)

## Status

Permanent dhnt-internal fork. The `exclude-grammars-tag` branch keeps a
reference implementation of a build-tag exclusion approach, for the record —
it is not proposed upstream.
