<!-- markdownlint-disable MD033 -->
# todo.txt Format Specification

- Official spec: [https://github.com/todotxt/todo.txt](https://github.com/todotxt/todo.txt?tab=readme-ov-file)
  - Supported version: As of commit [e79d866](https://github.com/todotxt/todo.txt/tree/e79d866af927013dc4fcfbe27ed51254a2cac394) (2025-11-26)

## Table of Contents

- [Sample Task String](#sample-task-string)
- [Overview](#overview)
- [Incomplete Tasks: 3 Format Rules](#incomplete-tasks-3-format-rules)
- [Complete Tasks: 2 Format Rules](#complete-tasks-2-format-rules)
- [Additional File Format Definitions](#additional-file-format-definitions)
- [Components](#components)
- [Extended Rules](#extended-rules)
  - [Comments](#comments)
  - [Task Length Limit](#task-length-limit)

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
- Projects: whitespace + `+` + non-whitespace.
- Contexts: whitespace + `@` + non-whitespace.
- Metadata: whitespace + `key:value`.
- **Extended Rules of this package:**
  - Lines starting with `#` are treated as comments (see [Comments](#comments)).
  - Any text after "#" with any Unicode space character are treated as inline comment.
  - Unless explicitly allowed, control characters are not permitted.
  - Each task line should not be longer than 8192 Bytes (8KB).

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

## Components

Components are the elements that make up a todo.txt task line.

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

## Extended Rules

Extended rules that are not part of the [official standard](https://github.com/todotxt/todo.txt?tab=readme-ov-file) but are defined in this package.

### Comments

This package treats lines or segments starting with `#` as comments by default:

- **Comment line**: A line beginning with `#` is treated as a comment and ignored during task processing.
  - Example: `# This is a comment line`
- **Inline comment**: Text following `#` within a task line is treated as a comment.
  - Example: `(A) Buy milk @store # Don't forget the organic one`
- **Not a comment**: If `#` appears in the middle of a word or tag and has no leading whitespace, it is not treated as a comment.
  - Example: `fix Issue#123 in the codebase`

### Task Length Limit

This package restricts the maximum length of each task line to 8KB (8,192 Bytes).

| Encoding Type | Max Characters | Max Length (Bytes) |
| :--- | :--- | :--- |
| ASCII | around 8,000 characters | 8,000 Bytes |
| Multi-byte UTF-8 | around 2,700 characters | 8,100 Bytes |
| Emoji mixed | around 2,000 - 4,000 characters | 8,000 Bytes |

Since todo.txt is primarily designed for plain text tasks, being portable and lightweight is essential.

Based on the below common platform limits and i18n considerations, we've chosen 8KB is a reasonable limit.

| System / Platform | Typical Line Length Limit |
| :--- | :--- |
| Recommended Email subject line | 78 characters |
| RFC 5322 (2.1.1, IMF Line Length Limit) | 998 characters |
| Common Japanese paragraph length | 200-400 characters |
| Twitter/X | 280 characters |
| Git commit message | 72 characters/line |
| GitHub issue title | 256 characters |
| GitHub issue body | ~65,000 characters |
| Jira task title | 255 characters |
| Todoist task | 500 characters |
| Google Keep note | ~20,000 characters |
| Markdown paragraph | ~1,000 characters |
