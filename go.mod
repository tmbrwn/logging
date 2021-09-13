module github.com/tmbrwn/logging

go 1.16

require (
	github.com/fatih/color v1.12.0
	github.com/jonboulle/clockwork v0.2.2
)

// this module should stay in v0. v1 is retracted.
retract (
	v1.0.0
	v1.0.1
	v1.1.0
	v1.2.0
	v1.2.1
)
