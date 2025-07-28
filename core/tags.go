package core

// These empty structs serve as declarative annotations embedded within user-defined
// structs to indicate CLI metadata such as whether a field is a short or long flag,
// required, or holds descriptive information.
//
// The parsing and help-generation logic uses reflection to detect these markers and
// adjust behavior accordingly.
//
// They carry no data or methods themselves and are not meant to be tested directly,
// but rather through the functionality that consumes them.

// === META TAGS ===

type Clifford struct{}
type Version struct{}
type Help struct{}

// === TAGGING ===

type ShortTag struct{}
type LongTag struct{}
type Required struct{}
type Desc struct{}
