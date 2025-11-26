# Package Overview

`task` lets you edit todo.txt strings without losing their original spacing or
tag order.

It wraps a `Parsed` instance of [`parse` package](../parse/) and keeps the raw task text as the source of
truth until you choose to re-parse it.

## Constructor

```go
New(taskString string, opts ...Option) (*Task, error)
```

- Parses the input string with `parse.FromTaskString`.
- Applies optional modifiers if provided (for example, `WithPriority("A")`).
- Calls `Apply()` to sync the parsed data with any modifications.

If parsing fails (invalid control characters, malformed dates, etc.), `New`
returns the wrapped error.

## Core Methods

- `String()` returns the current raw task text, including any pending edits to be
    applied.
- `InsertAfter(target, segment string)` inserts a token after the first matching
    segment.
- `AppendSegment(segment string)` adds a token at the tail, inserts it before an
    inline comment when present, and prepends it on comment lines so the leading
    `#` stays at the front.
- `RemoveSegment(segment string)` deletes the first exact segment while
    preserving spacing.
- `Apply(opts ...Option)` re-runs `Parse(opts)` only when the task is marked dirty,
    clearing caches and syncing `Parsed.Segments`.

### Convenient Mutators

- `AddContext(context string)` normalizes and adds an `@context`, skipping
    comment lines and duplicate tags.
- `RemoveContext(context string)` removes an `@context` tag from the task.
- `AddProject(project string)` normalizes and adds a `+project`, skipping
    comment lines and duplicate tags.
- `RemoveProject(project string)` removes a `+project` tag from the task.
- `SetPriority(priority string)` validates an `(A)` style mark and inserts or
    replaces it in the correct position.
- `RemovePriority()` removes any existing priority mark.
- `SetTag(key, value string)` sets or updates a key-value pair tag
    (e.g., `due:2024-12-25`). See [Duplicate Keys](#duplicate-keys) for edge-case
    behavior.
- `RemoveTag(key string)` removes a key-value pair tag by key name.
- `Complete()` / `CompleteWithDate(date string)` marks the task as done.
- `Reopen()` marks the task as incomplete.

Most mutators mark the task dirty by updating the raw text. Call `Apply()` after
batching your edits to refresh the parsed view.

## Relationship With `parse` Package

`task.Task` embeds `*parse.Parsed` to reuse the getters while keeping the raw
text (`originalText`) as the single source of truth.

```text
Task.String() ─────────────► Raw text (originalText)
        │                            │ (edit via Task methods)
        ├─ SetText() marks dirty ◄───┘
        │
        └─ Apply()
                 │
                 └─ when dirty ► Parsed.Parse(new text) ► caches reset
```

- `Task` methods edit the raw string (for example, by adding tags) and set the
    `isDirty` flag.
- `Apply()` runs `Parse()` only when the task is dirty, rebuilding
    `Parsed.Segments` and clearing caches so the getters stay in sync.

## Notes and Tips

- Work with `Task` when you must preserve whitespace or inline comment
    positions. The helpers operate on the raw string for that reason.
- Always finish with `Apply()` before reading derived fields like `Contexts()`
    or `DateCreated()`. Otherwise you may read stale caches.
- Options in `task/options.go` mirror many mutators. Mutation options can run
    via `New()` or `Apply()`, but parse-layer options (for example,
    `WithAllowedCtrlChars`) are consumed only during construction. Recreate the
    task if you need to change parse options later.

## When to Use Which Package

- just need to read data → call `parse.FromTaskString`
- need to edit the text but keep spacing → call `task.New`, use the helper
    methods, then `Apply()`

This split keeps reading logic deterministic while letting mutators defer
re-parsing until they actually finalize changes.

## Duplicate Keys

When the same key-value key (e.g., `due:`) appears multiple times in a task:

- `SetTag(key, value)` updates only the **first** occurrence.
- `KeyValues()` returns only the **last** value (since it uses a `map[string]string`).
- `RemoveTag(key)` removes only the **first** occurrence.

This behavior comes from the upstream todo.txt spec which implicitly treats keys
as unique. To store multiple values for a single key, use a comma-separated
format like `tag:value1,value2` or `tag1:value1 tag2:value2`.

**Example:**

```text
Original: "Buy milk due:2024-01-01 due:2024-02-01"
After SetTag("due", "2024-03-01"):
         "Buy milk due:2024-03-01 due:2024-02-01"
KeyValues()["due"] returns: "2024-02-01"  // last occurrence
```

To avoid confusion, ensure your task files do not contain duplicate keys, or
parse and normalize them before editing.
