# go-htmx-pccs

An application for assisting with Phoenix Command Combat System wargaming/roleplaying games,
A work in progress replacement for my first PCCS app started in late 2018: https://github.com/mp40/PCCS-Helper

## How to run

```bash
$ go build
$ go run .
```

## How to test

```bash
$ go test ./...
```

## Stack

Go stdlib + HTMX + SQLite + vanilla CSS. No package.json.

## Design Choices

- Minimal dependencies
- Minimal build requirements and steps
- Long term survivability over latest 'best practice'
- Most of a page content is static, reactive frameworks/libraries are unneeded dead weight
- Portable domain logic
- Do not sacrifice simplicity for scalability, it is a hobby app for use with a few mates
