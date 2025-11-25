# Internal Packages

These packages prepare for future refactoring of the `todo` package, implementing a design shift toward preserving original task string formatting.

The shared `spec` package keeps todo.txt markers consistent across the stack.

## Package Layers

| Package | Responsibility | Note |
| :--: | :-- |
| `spec/` | Shared todo.txt markers used across internal packages | |
| `segment/` | Provides the `Segment` type plus predicate helpers | Segments are tokens split by spaces |
| `parse/` | Tokenizes task strings and exposes read-only getters | |
| `task/` | Modifies task components while preserving original format | |

## Design Philosophy

The todo.txt format emphasizes human readability. Tags can appear anywhere to form natural sentences:

- Old assumption: `(A) buy milk @store +shopping` (tags at end)
- Natural usage: `(A) buy @milk when +shopping` (tags inline)

To preserve user-intended formatting, these packages manipulate the original task string directly rather than reconstructing from parsed components. This prevents position shifts when round-tripping tasks.

- See: [Issue #21](https://github.com/KEINOS/go-todotxt/issues/21)
