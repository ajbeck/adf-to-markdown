//go:build goexperiment.jsonv2

package adfmarkdown

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

//go:embed schema/adf-full-schema-51.5.15.json
var embeddedADFSchema []byte

var (
	schemaOnce     sync.Once
	schemaResolved *jsonschema.Resolved
	schemaInitErr  error
)

func initADFSchema() error {
	schemaOnce.Do(func() {
		var s jsonschema.Schema
		if err := json.Unmarshal(embeddedADFSchema, &s); err != nil {
			schemaInitErr = fmt.Errorf("parse embedded ADF schema: %w", err)
			return
		}
		schemaResolved, schemaInitErr = s.Resolve(nil)
		if schemaInitErr != nil {
			schemaInitErr = fmt.Errorf("resolve embedded ADF schema: %w", schemaInitErr)
		}
	})
	return schemaInitErr
}

func ValidateADFSchema(data []byte) error {
	if err := initADFSchema(); err != nil {
		return err
	}
	var instance any
	if err := json.Unmarshal(data, &instance); err != nil {
		return err
	}
	if err := schemaResolved.Validate(instance); err != nil {
		return fmt.Errorf("ADF schema validation failed: %w", err)
	}
	return nil
}
