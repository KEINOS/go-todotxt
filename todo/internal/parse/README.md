# Package Overview

`parse` gives read-only access to todo.txt task strings.

It reads the raw text, splits it into `segment.Segments`, memoizes only the
expensive lookups (inline comments, comment/done flags, header skip indexes,
and key/value maps), and provides simple getter methods. Use the [`task` package](../task/) when you need to change
the task text.

## Constructor

```go
FromTaskString(todoTxt string, opts ...Option) (*Parsed, error)
```

- Splits the todoTxt into segments using `strings.Fields`.
- Applies functional options such as `WithAllowedCtrlChars` if provided.
- Returns a `Parsed` that exposes read-only getters.

## Common Getters

- `String()` gives back the original task text.
- `IsDone()`, `IsCommentLine()`, `HasInlineComment()` expose status flags.
- `Priority()`, `DateCompleted()`, `DateCreated()` return header values.
- `Contexts()`, `Projects()` return slices of tags.
- `KeyValues()` returns a slice of key-value pairs, preserving order and duplicates.
- `Description()` strips markers and inline comments while keeping tag order.
- `Comment()` returns the comment text (line or inline).
- `Components()` bundles the above into a single struct.

Only the heavier getters cache their results lazily (`HasInlineComment`,
`InlineCommentPos`, `IsCommentLine`, `IsDone`, `KeyValues`, header skip index).
Lightweight helpers such as `Priority`, `DateCompleted`, `Contexts`, or
`Description` intentionally recompute on every call to keep the cache surface
small. Calling `Parse()` again clears every cached field so later reads stay
consistent with the updated raw text.

## Data Flow Cheat Sheet

```text
Raw todo.txt formatted task string
        │
        ▼
Parse(input) ── strings.Fields ──► segment.Segments
        │            │
        │            ├── set input to originalText as is
        │            │
        │            └── selected lookups cached lazily (status flags, key-values)
        ▼
Parsed
```

- `Parsed.Segments` is the canonical token list.
- Cache fields such as `inlineComment`, `indexSkip`, and `keyValueCache` stay
        empty until a getter needs them.
- `Parse()` calls `resetCaches()`, so every read after parsing starts from clean
        data and caches again as needed.

## Notes

- This package never mutates the task text.
  - Use `SetText(newText)` only when an external caller updates the raw string,
        then call `Parse()` to re-sync caches.
- Control characters are rejected by default.
  - Allow them explicitly with `WithAllowedCtrlChars` only when you trust the
        input source.
- When you need to modify tasks, switch to the `task` package; it reuses
        `Parsed` under the hood but preserves spacing during edits.
