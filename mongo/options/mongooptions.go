// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"fmt"
	"reflect"

	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/bsoncodec"
	"github.com/pritunl/mongo-go-driver/x/bsonx/bsoncore"
)

// Collation allows users to specify language-specific rules for string comparison, such as
// rules for lettercase and accent marks.
type Collation struct {
	Locale          string `bson:",omitempty"` // The locale
	CaseLevel       bool   `bson:",omitempty"` // The case level
	CaseFirst       string `bson:",omitempty"` // The case ordering
	Strength        int    `bson:",omitempty"` // The number of comparision levels to use
	NumericOrdering bool   `bson:",omitempty"` // Whether to order numbers based on numerical order and not collation order
	Alternate       string `bson:",omitempty"` // Whether spaces and punctuation are considered base characters
	MaxVariable     string `bson:",omitempty"` // Which characters are affected by alternate: "shifted"
	Normalization   bool   `bson:",omitempty"` // Causes text to be normalized into Unicode NFD
	Backwards       bool   `bson:",omitempty"` // Causes secondary differences to be considered in reverse order, as it is done in the French language
}

// ToDocument converts the Collation to a bson.Raw.
func (co *Collation) ToDocument() bson.Raw {
	idx, doc := bsoncore.AppendDocumentStart(nil)
	if co.Locale != "" {
		doc = bsoncore.AppendStringElement(doc, "locale", co.Locale)
	}
	if co.CaseLevel {
		doc = bsoncore.AppendBooleanElement(doc, "caseLevel", true)
	}
	if co.CaseFirst != "" {
		doc = bsoncore.AppendStringElement(doc, "caseFirst", co.CaseFirst)
	}
	if co.Strength != 0 {
		doc = bsoncore.AppendInt32Element(doc, "strength", int32(co.Strength))
	}
	if co.NumericOrdering {
		doc = bsoncore.AppendBooleanElement(doc, "numericOrdering", true)
	}
	if co.Alternate != "" {
		doc = bsoncore.AppendStringElement(doc, "alternate", co.Alternate)
	}
	if co.MaxVariable != "" {
		doc = bsoncore.AppendStringElement(doc, "maxVariable", co.MaxVariable)
	}
	if co.Normalization {
		doc = bsoncore.AppendBooleanElement(doc, "normalization", true)
	}
	if co.Backwards {
		doc = bsoncore.AppendBooleanElement(doc, "backwards", true)
	}
	doc, _ = bsoncore.AppendDocumentEnd(doc, idx)
	return doc
}

// CursorType specifies whether a cursor should close when the last data is retrieved. See
// NonTailable, Tailable, and TailableAwait.
type CursorType int8

const (
	// NonTailable specifies that a cursor should close after retrieving the last data.
	NonTailable CursorType = iota
	// Tailable specifies that a cursor should not close when the last data is retrieved and can be resumed later.
	Tailable
	// TailableAwait specifies that a cursor should not close when the last data is retrieved and
	// that it should block for a certain amount of time for new data before returning no data.
	TailableAwait
)

// ReturnDocument specifies whether a findAndUpdate operation should return the document as it was
// before the update or as it is after the update.
type ReturnDocument int8

const (
	// Before specifies that findAndUpdate should return the document as it was before the update.
	Before ReturnDocument = iota
	// After specifies that findAndUpdate should return the document as it is after the update.
	After
)

// FullDocument specifies whether a change stream should include a copy of the entire document that was changed from
// some time after the change occurred.
type FullDocument string

const (
	// Default does not include a document copy
	Default FullDocument = "default"
	// UpdateLookup includes a delta describing the changes to the document and a copy of the entire document that
	// was changed
	UpdateLookup FullDocument = "updateLookup"
)

// ArrayFilters is used to hold filters for the array filters CRUD option. If a registry is nil, bson.DefaultRegistry
// will be used when converting the filter interfaces to BSON.
type ArrayFilters struct {
	Registry *bsoncodec.Registry // The registry to use for converting filters. Defaults to bson.DefaultRegistry.
	Filters  []interface{}       // The filters to apply
}

// ToArray builds a []bson.Raw from the provided ArrayFilters.
func (af *ArrayFilters) ToArray() ([]bson.Raw, error) {
	registry := af.Registry
	if registry == nil {
		registry = bson.DefaultRegistry
	}
	filters := make([]bson.Raw, 0, len(af.Filters))
	for _, f := range af.Filters {
		filter, err := bson.MarshalWithRegistry(registry, f)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filter)
	}
	return filters, nil
}

// MarshalError is returned when attempting to transform a value into a document
// results in an error.
type MarshalError struct {
	Value interface{}
	Err   error
}

// Error implements the error interface.
func (me MarshalError) Error() string {
	return fmt.Sprintf("cannot transform type %s to a bson.Raw", reflect.TypeOf(me.Value))
}

var defaultRegistry = bson.DefaultRegistry
