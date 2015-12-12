### Go

- upper => Public / lower => Private
- all files named `*_test.go` are ignored by go build, and will only run by go test.
- [build tags](https://golang.org/pkg/go/build/): `// +build foo` =>
- https://github.com/a8m/go-lang-cheat-sheet
- https://golang.org/doc/code.html#Testing


### Atom

- watch syntax in keymap.cson.
- cannot define same section multiple times. restating overrides previous.
- list of all commands: View | Developer | Toggle Developer Tools, console
  type atom.commands.registeredCommands (useful to bind a shortcut to something
  not currently bound).
- can tab switching go on last used order vs. tab order? there's a package to do
  this but it looks sketchy. the last-cursor-position package is good-ish
  enough.
