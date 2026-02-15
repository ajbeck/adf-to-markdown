//go:build !goexperiment.jsonv2

package adfmarkdown

import (
	"errors"
	"io"
)

var errJSONV2Required = errors.New("adfmarkdown requires GOEXPERIMENT=jsonv2")

type HardBreakStyle string
type CodeFenceStyle string
type UnsupportedBlockHandler func(nodeType, path string) (markdown string, handled bool, err error)
type UnsupportedInlineHandler func(nodeType, path string) (markdown string, handled bool, err error)
type ExtensionBlockHandler func(nodeType, key, path string) (markdown string, handled bool, err error)
type ExtensionInlineHandler func(nodeType, key, path string) (markdown string, handled bool, err error)

type config struct{}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

type Error struct {
	Path   string
	Kind   ErrorKind
	Detail string
	Cause  error
}

func (e *Error) Error() string { return e.Detail }
func (e *Error) Unwrap() error { return e.Cause }

func WithStrictSchema(bool) Option          { return optionFunc(func(*config) {}) }
func WithAllowUnsupportedNodes(bool) Option { return optionFunc(func(*config) {}) }
func WithHardBreakStyle(HardBreakStyle) Option {
	return optionFunc(func(*config) {})
}
func WithCodeFenceStyle(CodeFenceStyle) Option {
	return optionFunc(func(*config) {})
}
func WithSchemaValidator(func([]byte) error) Option { return optionFunc(func(*config) {}) }
func WithBuiltInSchemaValidation(bool) Option       { return optionFunc(func(*config) {}) }
func WithUnsupportedBlockHandler(UnsupportedBlockHandler) Option {
	return optionFunc(func(*config) {})
}
func WithUnsupportedInlineHandler(UnsupportedInlineHandler) Option {
	return optionFunc(func(*config) {})
}
func WithExtensionBlockHandler(ExtensionBlockHandler) Option {
	return optionFunc(func(*config) {})
}
func WithExtensionInlineHandler(ExtensionInlineHandler) Option {
	return optionFunc(func(*config) {})
}

func UnmarshalADF([]byte, ...Option) ([]byte, error) {
	return nil, errJSONV2Required
}

func UnmarshalADFTo(io.Writer, []byte, ...Option) error {
	return errJSONV2Required
}
