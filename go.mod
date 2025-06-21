module github.com/grackleclub/log

go 1.23.0

retract [v1.0.0, v1.1.4] // should not have been published as v1

require (
	github.com/lmittmann/tint v1.0.5
	github.com/mattn/go-isatty v0.0.20
	github.com/stretchr/testify v1.9.0
	golang.org/x/term v0.24.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
