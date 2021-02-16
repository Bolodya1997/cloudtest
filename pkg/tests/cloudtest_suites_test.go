// Copyright (c) 2020 Doc.ai and/or its affiliates.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tests

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/networkservicemesh/cloudtest/pkg/commands"
	"github.com/networkservicemesh/cloudtest/pkg/config"
	"github.com/networkservicemesh/cloudtest/pkg/reporting"
	"github.com/networkservicemesh/cloudtest/pkg/utils"
)

func TestCloudtestCanWorkWithSuites(t *testing.T) {
	testConfig := config.NewCloudTestConfig()

	testConfig.Timeout = 300

	tmpDir, err := ioutil.TempDir(os.TempDir(), "cloud-test-temp")
	require.NoError(t, err)
	defer utils.ClearFolder(tmpDir, false)
	testConfig.ConfigRoot = tmpDir
	createProvider(testConfig, "a_provider")
	testConfig.MinSuiteSize = 3

	testConfig.Executions = append(testConfig.Executions, &config.Execution{
		Name:        "simple",
		Timeout:     15,
		PackageRoot: "./sample/suites",
	})

	testConfig.Reporting.JUnitReportFile = JunitReport

	report, err := commands.PerformTesting(testConfig, &TestValidationFactory{}, &commands.Arguments{})
	require.NoError(t, err)
	require.NotNil(t, report)

	var providerSuite *reporting.Suite
	for providerSuite = report.Suites[0]; providerSuite.Name != "a_provider"; providerSuite = providerSuite.Suites[0] {
	}

	const testName = "TestRunSuite"

	var providerSuiteTestCount int
	for _, suite := range providerSuite.Suites {
		if suite.Name == testName {
			require.Equal(t, 4, suite.Tests)
			providerSuiteTestCount += suite.Tests
		}
	}
	require.Equal(t, 4, providerSuiteTestCount)
}

func TestCloudtestCanWorkWithSuitesSplit(t *testing.T) {
	testConfig := config.NewCloudTestConfig()

	testConfig.Timeout = 300
	testConfig.MinSuiteSize = 2

	tmpDir, err := ioutil.TempDir(os.TempDir(), "cloud-test-temp")
	require.NoError(t, err)
	defer utils.ClearFolder(tmpDir, false)

	testConfig.ConfigRoot = tmpDir
	createProvider(testConfig, "a_provider")

	testConfig.Executions = append(testConfig.Executions, &config.Execution{
		Name:        "simple",
		Timeout:     15,
		PackageRoot: "./sample/suites",
	})

	testConfig.Reporting.JUnitReportFile = JunitReport

	report, err := commands.PerformTesting(testConfig, &TestValidationFactory{}, &commands.Arguments{})
	require.NoError(t, err)
	require.NotNil(t, report)

	var providerSuite *reporting.Suite
	for providerSuite = report.Suites[0]; providerSuite.Name != "a_provider"; providerSuite = providerSuite.Suites[0] {
	}

	const testName = "TestRunSuite"

	var providerSuiteTestCount int
	for _, suite := range providerSuite.Suites {
		if suite.Name == testName {
			require.Equal(t, 2, suite.Tests)
			providerSuiteTestCount += suite.Tests
		}
	}
	require.Equal(t, 4, providerSuiteTestCount)
}
