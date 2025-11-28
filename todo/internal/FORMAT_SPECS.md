<!-- markdownlint-disable MD033 -->
# todo.txt Format Specification

- Official spec: [https://github.com/todotxt/todo.txt](https://github.com/todotxt/todo.txt?tab=readme-ov-file)
  - Supported version: As of commit [e79d866](https://github.com/todotxt/todo.txt/tree/e79d866af927013dc4fcfbe27ed51254a2cac394) (2025-11-26)

## Table of Contents

- [Sample](#sample-task-string)
- [Overview](#overview)
- [Incomplete Tasks](#incomplete-tasks-3-format-rules)
- [Complete Tasks](#complete-tasks-2-format-rules)
- [Additional Metadata](#additional-file-format-definitions)
- [Comments](#comments)
- [Components](#components)

## Sample Task String

```todotxt
# Project tasks for home renovation
x (A) 2016-05-20 2016-04-30 measure space for +chapelShelving @chapel due:2016-05-30
(A) Thank Mom for the meatballs @phone # Call before 8pm
(B) Schedule Goodwill pickup +GarageSale @phone
Post signs around the neighborhood +GarageSale
@GroceryStore pies
```

## Overview

- Plain text, UTF-8 encoded.
- Line breaks: `\n` (LF) or `\r\n` (CRLF).
- One task per line.
- Leading/trailing whitespaces ignored.
- Completion: `x` (lowercase) followed by space.
- Priority: `(A-Z)` in parentheses. Uppercase letters only.
- Dates: `YYYY-MM-DD`.
- Projects: `+` + non-whitespace.
- Contexts: `@` + non-whitespace.
- Metadata: `key:value`.
- **Extended Rules of this package:**
  - Lines starting with `#` are treated as comments (see [Comments](#comments)).
  - Any text after " #" are treated as inline comment.
  - Unless explicitly allowed, control characters are not permitted.
  - Each task line should not be longer than 64 * 1024 Bytes (64KB. Default size of bufio.MaxScanTokenSize).

## Incomplete Tasks: 3 Format Rules

### Rule 1: Priority first if present

Priority `(A-Z)` appears first, followed by space.

### Rule 2: Creation date after priority

Optional `YYYY-MM-DD` directly after priority/space, or first if no priority. Locale is not considered.

### Rule 3: Projects/Contexts anywhere after

- `@context` or `+project` anywhere after priority/date.
- Multiple allowed.

## Complete Tasks: 2 Format Rules

### Rule 1: Start with `x`

- Lowercase `x` followed by space marks completion.

### Rule 2: Completion date after `x`

- `YYYY-MM-DD` directly after `x`.
- If creation date is present, completion date must come first and creation date second.
- Completed tasks with only creation date are not allowed.

## Additional File Format Definitions

Use `key:value` for extra metadata, e.g., `due:2010-01-02`.

## Comments

**Note:** Comments are an extension to the original todo.txt format specification and are not part of the [official standard](https://github.com/todotxt/todo.txt?tab=readme-ov-file).

This package treats lines or segments starting with `#` as comments by default:

- **Comment line**: A line beginning with `#` is treated as a comment and ignored during task processing.
  - Example: `# This is a comment line`
- **Inline comment**: Text following `#` within a task line is treated as a comment.
  - Example: `(A) Buy milk @store # Don't forget the organic one`

Comment lines and inline comments are preserved when reading/writing todo.txt files, but are not parsed as task components.

## Components

| Component | Description | Optional | Example |
| :--: | :-- | :--: | :-- |
| `x` | Completion marker | Yes | `x` |
| `(A)` | Priority | Yes | `(A)` |
| `2016-05-20` | Completion date | Yes | `2016-05-20` |
| `2016-04-30` | Creation date | Yes | `2016-04-30` |
| Description | Task text with tags | No | `measure space +project @context due:2016-05-30` |

### Tags in Description

| Tag | Description | Example |
| :--: | :-- | :-- |
| `+project` | Project/Category | `+chapelShelving` |
| `@context` | Context (Like hash-tags) | `@chapel` |
| `key:value` | Custom metadata | `due:2016-05-30`<br>`location:office` |
