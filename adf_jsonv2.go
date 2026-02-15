//go:build goexperiment.jsonv2

package adfmarkdown

import "io"

func UnmarshalADF(data []byte, opts ...Option) ([]byte, error) {
	cfg := defaultConfig()
	for _, opt := range opts {
		opt.apply(&cfg)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	if cfg.SchemaValidator != nil {
		if err := cfg.SchemaValidator(data); err != nil {
			return nil, err
		}
	} else if cfg.StrictSchema && cfg.BuiltInSchemaCheck {
		if err := ValidateADFSchema(data); err != nil {
			return nil, err
		}
	}

	dec := newDecoder(cfg)
	doc, err := dec.decodeDocument(data)
	if err != nil {
		return nil, err
	}

	em := newEmitter(cfg)
	return em.renderDocument(doc)
}

func UnmarshalADFTo(w io.Writer, data []byte, opts ...Option) error {
	out, err := UnmarshalADF(data, opts...)
	if err != nil {
		return err
	}
	_, err = w.Write(out)
	return err
}
