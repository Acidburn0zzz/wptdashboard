// Copyright 2017 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package compute

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/w3c/wptdashboard/metrics"
	"github.com/w3c/wptdashboard/shared"
)

var timeA = time.Unix(0, 0)
var timeB = time.Unix(0, 1)

func TestGatherResultsById_TwoRuns_SameTest(t *testing.T) {
	runA := shared.TestRun{
		"ABrowser",
		"1.0",
		"MyOS",
		"1.0",
		"abcd",
		"http://example.com/a_run.json",
		timeA,
	}
	runB := shared.TestRun{
		"BBrowser",
		"1.0",
		"MyOS",
		"1.0",
		"dcba",
		"http://example.com/b_run.json",
		timeB,
	}
	testName := "Do a thing"
	results := &[]metrics.TestRunResults{
		{
			&runA,
			&metrics.TestResults{
				"A test",
				"OK",
				&testName,
				[]metrics.SubTest{},
			},
		},
		{
			&runB,
			&metrics.TestResults{
				"A test",
				"ERROR",
				&testName,
				[]metrics.SubTest{},
			},
		},
	}
	gathered := GatherResultsById(results)
	assert.Equal(t, len(gathered), 1) // Merged to single TestId: {"A test",""}.
	for testId, runStatusMap := range gathered {
		assert.Equal(t, testId, metrics.TestId{"A test", ""})
		for testRun, testStatus := range runStatusMap {
			if testRun == runA {
				assert.Equal(t, testStatus, metrics.CompleteTestStatus{
					metrics.TestStatus_fromString("OK"),
					metrics.SubTestStatus_fromString("STATUS_UNKNOWN"),
				})
			} else if testRun == runB {
				assert.Equal(t, testStatus, metrics.CompleteTestStatus{
					metrics.TestStatus_fromString("ERROR"),
					metrics.SubTestStatus_fromString("STATUS_UNKNOWN"),
				})
			} else {
				assert.Fail(t, "Invalid test run in GatherResultsById output")
			}
		}
	}
}