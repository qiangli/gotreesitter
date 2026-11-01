---
id: cfc13b0c3e44
kind: bug
title: S557 restore imported TSX grammar smoke parity
seq: 1
status: assigned
priority: p1
created: 2026-11-01T17:04:58.293621Z
assignee: s557-test-readiness
sprint: 557
sprint_id: f5ccf563-f5d1-5b94-8be4-9388e51ef86d
sprint_title: 'Bashy v1.0.0 preparation: split private/public development and retire qiangli/dhnt'
---

Pinned TSX grammar imports and generates successfully but existing JSX-return smoke has zero error-free/parity samples against required floor one. Preserve grammar pin and floor. Add failure diagnostics and focused regression, run bounded hosted TSX gate, locate first lexer/scanner/LR divergence and fix it without global algorithm changes unsupported by evidence.

Diagnostic checkpoints use `[skip ci]` only to suppress the unrelated automatic
full race matrix while the explicit manual TSX container job locates the first
divergence. The full requested public module validation remains required after
the repair; the focused diagnostic is not aggregate acceptance.

Hosted diagnostic 41329763032 reproduced the unchanged pinned smoke: generated
JSX text consumes the semicolon after a closing element, then becomes invalid
after reductions. The reference consumes the semicolon normally. The bounded
repair opts the stateless TSX scanner into existing DFA relex recovery only for
one unshifted live branch with an invalid external token and a valid positive
width replacement. Other scanner capabilities and LR splitting stay unchanged.
Targeted synthetic safety guards and pinned smoke will run hosted before broader
TSX corpus validation and aggregate module acceptance.

Repair checkpoint 66339c2 passed hosted run 41330365940: synthetic safety
regressions, unchanged imported smoke (1/1 no-error, 1/1 parity), and deep
reference parity for the smoke. Independent review added explicit opt-out and
following-source preservation checks. The next focused gate runs all existing
TSX snippet and direct-C-named deep regression cases against the embedded oracle.

Expanded hosted run https://github.com/qiangli/gotreesitter/actions/runs/41330511422
passed at f1f80fc: all 18 TSX snippet cases, seven deep regression cases, unchanged
import smoke (1/1 no-error and parity), and all synthetic scanner/cursor/GLR guard
checks. No tests skipped. Full public module validation and manager acceptance
remain required; this record does not declare aggregate release acceptance.
