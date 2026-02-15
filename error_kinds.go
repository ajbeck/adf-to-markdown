package adfmarkdown

type ErrorKind string

const (
	ErrKindInvalidRoot       ErrorKind = "invalid_root"
	ErrKindInvalidJSON       ErrorKind = "invalid_json"
	ErrKindMissingType       ErrorKind = "missing_type"
	ErrKindUnsupportedNode   ErrorKind = "unsupported_node"
	ErrKindUnsupportedInline ErrorKind = "unsupported_inline"
	ErrKindUnsupportedMark   ErrorKind = "unsupported_mark"
	ErrKindInvalidAttr       ErrorKind = "invalid_attr"
	ErrKindInvalidMark       ErrorKind = "invalid_mark"
	ErrKindInvalidMarkCombo  ErrorKind = "invalid_mark_combo"
	ErrKindInvalidStructure  ErrorKind = "invalid_structure"
	ErrKindInvalidText       ErrorKind = "invalid_text"
)
