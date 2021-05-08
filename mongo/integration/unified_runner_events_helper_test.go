// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package integration

import (
	"context"
	"sync"
	"time"

	"github.com/pritunl/mongo-go-driver/event"
	"github.com/pritunl/mongo-go-driver/internal/testutil/assert"
	"github.com/pritunl/mongo-go-driver/mongo/address"
	"github.com/pritunl/mongo-go-driver/mongo/description"
	"github.com/pritunl/mongo-go-driver/mongo/integration/mtest"
	"github.com/pritunl/mongo-go-driver/mongo/readpref"
	"github.com/pritunl/mongo-go-driver/x/mongo/driver/topology"
)

// Helper functions for the operations in the unified spec test runner that require assertions about SDAM and connection
// pool events.

var (
	eventTypesMap = map[string]string{
		"PoolClearedEvent": event.PoolCleared,
	}
	defaultCallbackTimeout = 10 * time.Second
)

// unifiedRunnerEventMonitor monitors connection pool-related events.
type unifiedRunnerEventMonitor struct {
	sync.Mutex
	poolEventCount map[string]int
}

func newUnifiedRunnerEventMonitor() *unifiedRunnerEventMonitor {
	return &unifiedRunnerEventMonitor{
		poolEventCount: make(map[string]int),
	}
}

// handlePoolEvent can be used as the event handler for an connection pool monitor.
func (u *unifiedRunnerEventMonitor) handlePoolEvent(evt *event.PoolEvent) {
	u.Lock()
	defer u.Unlock()

	u.poolEventCount[evt.Type]++
}

// getCount returns the number of events of the given type, or 0 if no events were recorded.
func (u *unifiedRunnerEventMonitor) getCount(eventType string) int {
	u.Lock()
	defer u.Unlock()

	mappedType := eventTypesMap[eventType]
	return u.poolEventCount[mappedType]
}

func waitForEvent(mt *mtest.T, test *testCase, op *operation) {
	eventType := op.Arguments.Lookup("event").StringValue()
	if eventType == "ServerMarkedUnknownEvent" {
		mt.Log("skipping waitForEvent assertion for event type ServerMarkedUnknownEvent")
		return
	}

	expectedCount := int(op.Arguments.Lookup("count").Int32())

	callback := func() {
		for {
			count := test.monitor.getCount(eventType)
			if count >= expectedCount {
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}

	assert.Soon(mt, callback, defaultCallbackTimeout)
}

func assertEventCount(mt *mtest.T, testCase *testCase, op *operation) {
	eventType := op.Arguments.Lookup("event").StringValue()
	if eventType == "ServerMarkedUnknownEvent" {
		mt.Log("skipping waitForEvent assertion for event type ServerMarkedUnknownEvent")
		return
	}

	expectedCount := int(op.Arguments.Lookup("count").Int32())
	gotCount := testCase.monitor.getCount(eventType)
	assert.Equal(mt, expectedCount, gotCount, "expected count %d for event %s, got %d", expectedCount, eventType,
		gotCount)
}

func recordPrimary(mt *mtest.T, testCase *testCase) {
	testCase.recordedPrimary = getPrimaryAddress(mt, testCase.testTopology, true)
}

func waitForPrimaryChange(mt *mtest.T, testCase *testCase, op *operation) {
	callback := func() {
		for {
			if getPrimaryAddress(mt, testCase.testTopology, false) != testCase.recordedPrimary {
				return
			}
		}
	}

	timeout := convertValueToMilliseconds(mt, op.Arguments.Lookup("timeoutMS"))
	assert.Soon(mt, callback, timeout)
}

// getPrimaryAddress returns the address of the current primary. If failFast is true, the server selection fast path
// is used and the function will fail if the fast path doesn't return a server.
func getPrimaryAddress(mt *mtest.T, topo *topology.Topology, failFast bool) address.Address {
	mt.Helper()

	ctx, cancel := context.WithCancel(mtest.Background)
	defer cancel()
	if failFast {
		cancel()
	}

	primary, err := topo.SelectServer(ctx, description.ReadPrefSelector(readpref.Primary()))
	assert.Nil(mt, err, "SelectServer error: %v", err)
	return primary.(*topology.SelectedServer).Description().Addr
}
