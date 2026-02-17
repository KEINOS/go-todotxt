# Package Overview

`segment` holds the `Segment` type which are the tokens split by spaces of a task (todo.txt formatted task) string. Each `Segment` has methods to identify/predicate its type.

## Key Predicates

- `IsComment()` → segment starts with `#` (line comment or inline comment)
- `IsMarkCompletion()` → exactly `x`
- `IsMarkPriority()` → pattern `(A)` with an uppercase letter
- `IsDate()` → `YYYY-MM-DD`
- `IsTagProject()` → `+project`
- `IsTagContext()` → `@context`
- `IsKeyValue()` → `key:value` (splits at first `:`, value may contain more)
- `IsPlainText()` → none of the above types

Each predicate works on an individual segment produced by `parse`. Use them to
categorise tokens without re-implementing the matching logic.
