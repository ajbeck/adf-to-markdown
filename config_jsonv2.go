//go:build goexperiment.jsonv2

package adfmarkdown

import "fmt"

type HardBreakStyle string

const (
	HardBreakTwoSpaces HardBreakStyle = "two-spaces"
	HardBreakBackslash HardBreakStyle = "backslash"
)

type CodeFenceStyle string

const (
	CodeFenceBackticks CodeFenceStyle = "```"
	CodeFenceTildes    CodeFenceStyle = "~~~"
)

type config struct {
	StrictSchema          bool
	AllowUnsupportedNodes bool
	HardBreakStyle        HardBreakStyle
	CodeFenceStyle        CodeFenceStyle
	SchemaValidator       func([]byte) error
	BuiltInSchemaCheck    bool
	UnsupportedBlock      UnsupportedBlockHandler
	UnsupportedInline     UnsupportedInlineHandler
	ExtensionBlock        ExtensionBlockHandler
	ExtensionInline       ExtensionInlineHandler
}

type UnsupportedBlockHandler func(nodeType, path string) (markdown string, handled bool, err error)
type UnsupportedInlineHandler func(nodeType, path string) (markdown string, handled bool, err error)
type ExtensionBlockHandler func(nodeType, key, path string) (markdown string, handled bool, err error)
type ExtensionInlineHandler func(nodeType, key, path string) (markdown string, handled bool, err error)

func defaultConfig() config {
	return config{
		StrictSchema:          true,
		AllowUnsupportedNodes: false,
		HardBreakStyle:        HardBreakTwoSpaces,
		CodeFenceStyle:        CodeFenceBackticks,
		BuiltInSchemaCheck:    true,
	}
}

func (c config) validate() error {
	switch c.HardBreakStyle {
	case HardBreakTwoSpaces, HardBreakBackslash:
	default:
		return fmt.Errorf("invalid hard break style: %q", c.HardBreakStyle)
	}
	switch c.CodeFenceStyle {
	case CodeFenceBackticks, CodeFenceTildes:
	default:
		return fmt.Errorf("invalid code fence style: %q", c.CodeFenceStyle)
	}
	return nil
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(c *config) { f(c) }

func WithStrictSchema(v bool) Option {
	return optionFunc(func(c *config) { c.StrictSchema = v })
}

func WithAllowUnsupportedNodes(v bool) Option {
	return optionFunc(func(c *config) { c.AllowUnsupportedNodes = v })
}

func WithHardBreakStyle(v HardBreakStyle) Option {
	return optionFunc(func(c *config) { c.HardBreakStyle = v })
}

func WithCodeFenceStyle(v CodeFenceStyle) Option {
	return optionFunc(func(c *config) { c.CodeFenceStyle = v })
}

func WithSchemaValidator(fn func([]byte) error) Option {
	return optionFunc(func(c *config) { c.SchemaValidator = fn })
}

func WithBuiltInSchemaValidation(v bool) Option {
	return optionFunc(func(c *config) { c.BuiltInSchemaCheck = v })
}

func WithUnsupportedBlockHandler(fn UnsupportedBlockHandler) Option {
	return optionFunc(func(c *config) { c.UnsupportedBlock = fn })
}

func WithUnsupportedInlineHandler(fn UnsupportedInlineHandler) Option {
	return optionFunc(func(c *config) { c.UnsupportedInline = fn })
}

func WithExtensionBlockHandler(fn ExtensionBlockHandler) Option {
	return optionFunc(func(c *config) { c.ExtensionBlock = fn })
}

func WithExtensionInlineHandler(fn ExtensionInlineHandler) Option {
	return optionFunc(func(c *config) { c.ExtensionInline = fn })
}
