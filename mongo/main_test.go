// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"os"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pritunl/mongo-go-driver/bson"
	"github.com/pritunl/mongo-go-driver/bson/bsoncodec"
	"github.com/pritunl/mongo-go-driver/internal/testutil/assert"
	"github.com/pritunl/mongo-go-driver/mongo/readconcern"
	"github.com/pritunl/mongo-go-driver/mongo/readpref"
	"github.com/pritunl/mongo-go-driver/mongo/writeconcern"
	"github.com/pritunl/mongo-go-driver/x/bsonx/bsoncore"
)

func TestMain(m *testing.M) {
	// register equality options
	assert.RegisterOpts(reflect.TypeOf(&Client{}), cmp.Comparer(func(c1, c2 *Client) bool {
		return c1 == c2
	}))
	assert.RegisterOpts(reflect.TypeOf(&bsoncodec.Registry{}), cmp.Comparer(func(r1, r2 *bsoncodec.Registry) bool {
		return r1 == r2
	}))

	assert.RegisterOpts(reflect.TypeOf(&readconcern.ReadConcern{}), cmp.AllowUnexported(readconcern.ReadConcern{}))
	assert.RegisterOpts(reflect.TypeOf(&writeconcern.WriteConcern{}), cmp.AllowUnexported(writeconcern.WriteConcern{}))
	assert.RegisterOpts(reflect.TypeOf(&readpref.ReadPref{}), cmp.AllowUnexported(readpref.ReadPref{}))
	assert.RegisterOpts(reflect.TypeOf(bson.D{}), cmp.AllowUnexported(bson.E{}, bsoncore.Value{}))
	assert.RegisterOpts(reflect.TypeOf(bson.A{}), cmp.AllowUnexported(bsoncore.Value{}))

	os.Exit(m.Run())
}
