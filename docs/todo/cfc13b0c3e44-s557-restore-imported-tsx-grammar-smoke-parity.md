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
