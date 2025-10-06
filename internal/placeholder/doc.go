// Package placeholder provides a temporary compilation target so that CI workflows
// can execute gofmt, go vet, golangci-lint, and go test successfully prior to
// landing the real CloudMoor services. This package will be removed once the
// initial domain modules (connectors, vault, persistence) are implemented in
// Milestone M0.
package placeholder
