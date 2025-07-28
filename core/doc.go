// Package core contains the core logic for parsing command-line arguments
// into user-defined structs using reflection.
//
// It provides the primary parsing function and defines marker types used
// to annotate struct fields with CLI metadata such as flags, required fields,
// and descriptions.
//
// This package is intended to be used internally by higher-level packages,
// but some core functions may be exposed for advanced use cases.
package core
