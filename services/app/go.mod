module github.com/itscheems/golang-template/services/app

go 1.26

toolchain go1.26.2

require github.com/itscheems/golang-template/libs/shared v0.0.0-00010101000000-000000000000

// Resolve workspace-local modules without depending on a published version.
// Keep these `replace` directives in sync with go.work's `use` list so the
// module also builds standalone (e.g. `cd services/app && go build` with
// GOWORK=off, or any tool that ignores go.work like `go mod tidy`).
replace github.com/itscheems/golang-template/libs/shared => ../../libs/shared
