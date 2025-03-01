package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseJunitReport(t *testing.T) {
	t.Run("not existing", func(t *testing.T) {
		tests, err := findFailedTests("not existing", nil, 0)
		assert.Error(t, err)
		assert.Nil(t, tests)
	})
	t.Run("golang", func(t *testing.T) {
		tests, err := findFailedTests("testdata/report.xml", nil, 0)
		assert.NoError(t, err)
		assert.Equal(t, []testCase{
			{
				Name:    "TestDifferentBaseTypes",
				Suite:   "github.com/stackrox/rox/pkg/booleanpolicy/evaluator",
				Message: "Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match: Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_fully_hydrated_object: Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_augmented_object: Failed",
				Error:   "Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match: Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_fully_hydrated_object: \n         evaluator_test.go:96: Error Trace: /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:96 /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:123 Error: Not equal: expected: false actual : true Test: TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_fully_hydrated_object \n    \nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_augmented_object: \n         evaluator_test.go:96: Error Trace: /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:96 /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:145 Error: Not equal: expected: false actual : true Test: TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_augmented_object \n    ",
			},
			{
				Name:    "TestLocalScannerTLSIssuerIntegrationTests",
				Message: "Failed\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh: Failed\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh/no_secrets: Failed",
				Stdout:  "",
				Stderr:  "",
				Error:   "    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca-key.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CERT_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-cert.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-key.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca-key.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CERT_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-cert.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-key.pem\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh: Failed\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh/no_secrets:     tls_issuer_test.go:377:\n        \tError Trace:\t/go/src/github.com/stackrox/stackrox/sensor/kubernetes/localscanner/tls_issuer_test.go:377\n        \t            \t\t\t\t/go/src/github.com/stackrox/stackrox/sensor/kubernetes/localscanner/tls_issuer_test.go:298\n        \t            \t\t\t\t/go/src/github.com/stackrox/stackrox/sensor/kubernetes/localscanner/suite.go:91\n        \tError:      \tcontext deadline exceeded\n        \tTest:       \tTestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh/no_secrets\nkubernetes/localscanner: 2022/10/03 07:32:47.446934 cert_refresher.go:109: Warn: local scanner certificates not found (this is expected on a new deployment), will refresh certificates immediately: 2 errors occurred:\n\t* secrets \"scanner-tls\" not found\n\t* secrets \"scanner-db-tls\" not found\n\n",
				Suite:   "github.com/stackrox/rox/sensor/kubernetes/localscanner",
			},
		}, tests)
	})
	t.Run("golang with threshold", func(t *testing.T) {
		tests, err := findFailedTests("testdata/report.xml", map[string]string{"JOB_NAME": "job-name"}, 1)
		assert.NoError(t, err)
		assert.Equal(t, []testCase{
			{
				Message: `github.com/stackrox/rox/pkg/booleanpolicy/evaluator / TestDifferentBaseTypes FAILED
github.com/stackrox/rox/sensor/kubernetes/localscanner / TestLocalScannerTLSIssuerIntegrationTests FAILED
`,
				JobName: "job-name",
				Suite:   "",
			},
		}, tests)
	})
	t.Run("dir multiple suites with threshold", func(t *testing.T) {
		tests, err := findFailedTests("testdata", map[string]string{"JOB_NAME": "job-name", "BUILD_ID": "1"}, 3)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]testCase{
				{
					Message: `DefaultPoliciesTest / Verify policy Apache Struts  CVE-2017-5638 is triggered FAILED
github.com/stackrox/rox/pkg/booleanpolicy/evaluator / TestDifferentBaseTypes FAILED
github.com/stackrox/rox/sensor/kubernetes/localscanner / TestLocalScannerTLSIssuerIntegrationTests FAILED
github.com/stackrox/rox/central/resourcecollection/datastore/store/postgres / TestCollectionsStore FAILED
`,
					JobName: "job-name",
					BuildId: "1",
				},
			},
			tests,
		)
	})
	t.Run("dir", func(t *testing.T) {
		tests, err := findFailedTests("testdata", map[string]string{"BUILD_ID": "1"}, 0)
		assert.NoError(t, err)

		assert.ElementsMatch(
			t,
			[]testCase{
				{
					Name: "Verify policy Apache Struts: CVE-2017-5638 is triggered",
					Message: "Condition not satisfied:\n" +
						"\n" +
						"waitForViolation(deploymentName,  policyName, 60)\n" +
						"|                |                |\n" +
						"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
						"",
					Stdout: "?[1;30m21:35:15?[0;39m | ?[34mINFO ?[0;39m | DefaultPoliciesTest       | Starting testcase\n" +
						"?[1;30m21:36:16?[0;39m | ?[34mINFO ?[0;39m | Services                  | Failed to trigger Apache Struts: CVE-2017-5638 after waiting 60 seconds\n" +
						"?[1;30m21:36:16?[0;39m | ?[1;31mERROR?[0;39m | Helpers                   | An exception occurred in test\n" +
						"org.spockframework.runtime.ConditionNotSatisfiedError: Condition not satisfied:\n" +
						"\n" +
						"waitForViolation(deploymentName,  policyName, 60)\n" +
						"|                |                |\n" +
						"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
						"\n" +
						"\tat DefaultPoliciesTest.$spock_feature_1_0(DefaultPoliciesTest.groovy:181) [1 skipped]\n" +
						"\tat util.OnFailureInterceptor.intercept(OnFailure.groovy:72) [8 skipped]\n" +
						"\tat util.OnFailureInterceptor.intercept(OnFailure.groovy:72) [10 skipped]\n" +
						" [6 skipped]\n" +
						"?[1;30m21:36:16?[0;39m | ?[39mDEBUG?[0;39m | Helpers                   | 2022-09-30 21:36:16 Will collect various stackrox logs for this failure under /tmp/qa-tests-backend-logs/a57dc4b9-70eb-4391-8a00-c5948fef733d/\n" +
						"?[1;30m21:37:07?[0;39m | ?[39mDEBUG?[0;39m | Helpers                   | Ran: ./scripts/ci/collect-service-logs.sh stackrox /tmp/qa-tests-backend-logs/a57dc4b9-70eb-4391-8a00-c5948fef733d/stackrox-k8s-logs\n" +
						"Exit: 0\n",
					Suite:   "DefaultPoliciesTest",
					BuildId: "1",
					Error: "Condition not satisfied:\n" +
						"\n" +
						"waitForViolation(deploymentName,  policyName, 60)\n" +
						"|                |                |\n" +
						"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
						"\n" +
						"\tat DefaultPoliciesTest.Verify policy #policyName is triggered(DefaultPoliciesTest.groovy:181)\n",
				},
				{
					Name:    "TestDifferentBaseTypes",
					Suite:   "github.com/stackrox/rox/pkg/booleanpolicy/evaluator",
					Message: "Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match: Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_fully_hydrated_object: Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_augmented_object: Failed",
					Error:   "Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match: Failed\nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_fully_hydrated_object: \n         evaluator_test.go:96: Error Trace: /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:96 /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:123 Error: Not equal: expected: false actual : true Test: TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_fully_hydrated_object \n    \nSub test TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_augmented_object: \n         evaluator_test.go:96: Error Trace: /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:96 /go/src/github.com/stackrox/stackrox/pkg/booleanpolicy/evaluator/evaluator_test.go:145 Error: Not equal: expected: false actual : true Test: TestDifferentBaseTypes/base_ts,_query_by_relative,_does_not_match/on_augmented_object \n    ",
					BuildId: "1",
				},
				{
					Name:    "TestLocalScannerTLSIssuerIntegrationTests",
					Message: "Failed\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh: Failed\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh/no_secrets: Failed",
					Suite:   "github.com/stackrox/rox/sensor/kubernetes/localscanner",
					Error:   "    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca-key.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CERT_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-cert.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-key.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CA_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/ca-key.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_CERT_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-cert.pem\n    env_isolator.go:41: EnvIsolator: Setting ROX_MTLS_KEY_FILE to /go/src/github.com/stackrox/stackrox/pkg/mtls/testutils/testdata/central-certs/leaf-key.pem\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh: Failed\nSub test TestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh/no_secrets:     tls_issuer_test.go:377:\n        \tError Trace:\t/go/src/github.com/stackrox/stackrox/sensor/kubernetes/localscanner/tls_issuer_test.go:377\n        \t            \t\t\t\t/go/src/github.com/stackrox/stackrox/sensor/kubernetes/localscanner/tls_issuer_test.go:298\n        \t            \t\t\t\t/go/src/github.com/stackrox/stackrox/sensor/kubernetes/localscanner/suite.go:91\n        \tError:      \tcontext deadline exceeded\n        \tTest:       \tTestLocalScannerTLSIssuerIntegrationTests/TestSuccessfulRefresh/no_secrets\nkubernetes/localscanner: 2022/10/03 07:32:47.446934 cert_refresher.go:109: Warn: local scanner certificates not found (this is expected on a new deployment), will refresh certificates immediately: 2 errors occurred:\n\t* secrets \"scanner-tls\" not found\n\t* secrets \"scanner-db-tls\" not found\n\n",
					BuildId: "1",
				},
				{
					Name:    "TestCollectionsStore",
					Suite:   "github.com/stackrox/rox/central/resourcecollection/datastore/store/postgres",
					Message: "Failed\nSub test TestCollectionsStore/TestStore: Failed",
					Error:   "    env_isolator.go:41: EnvIsolator: Setting ROX_POSTGRES_DATASTORE to true\nSub test TestCollectionsStore/TestStore:     store_test.go:47: collections TRUNCATE TABLE\n    store_test.go:95:\n        \tError Trace:\t/go/src/github.com/stackrox/stackrox/central/resourcecollection/datastore/store/postgres/store_test.go:95\n        \tError:      \tReceived unexpected error:\n        \t            \tERROR: update or delete on table \"collections\" violates foreign key constraint \"fk_collections_embedded_collections_collections_cycle_ref\" on table \"collections_embedded_collections\" (SQLSTATE 23503)\n        \t            \tcould not delete from \"collections\"\n        \t            \tgithub.com/stackrox/rox/pkg/search/postgres.RunDeleteRequestForSchema.func1\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/search/postgres/common.go:833\n        \t            \tgithub.com/stackrox/rox/pkg/postgres/pgutils.Retry.func1\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/postgres/pgutils/retry.go:21\n        \t            \tgithub.com/stackrox/rox/pkg/postgres/pgutils.Retry2[...].func1\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/postgres/pgutils/retry.go:32\n        \t            \tgithub.com/stackrox/rox/pkg/postgres/pgutils.Retry3[...]\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/postgres/pgutils/retry.go:43\n        \t            \tgithub.com/stackrox/rox/pkg/postgres/pgutils.Retry2[...]\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/postgres/pgutils/retry.go:35\n        \t            \tgithub.com/stackrox/rox/pkg/postgres/pgutils.Retry\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/postgres/pgutils/retry.go:23\n        \t            \tgithub.com/stackrox/rox/pkg/search/postgres.RunDeleteRequestForSchema\n        \t            \t\t/go/src/github.com/stackrox/stackrox/pkg/search/postgres/common.go:830\n        \t            \tgithub.com/stackrox/rox/central/resourcecollection/datastore/store/postgres.(*storeImpl).Delete\n        \t            \t\t/go/src/github.com/stackrox/stackrox/central/resourcecollection/datastore/store/postgres/store.go:429\n        \t            \tgithub.com/stackrox/rox/central/resourcecollection/datastore/store/postgres.(*CollectionsStoreSuite).TestStore\n        \t            \t\t/go/src/github.com/stackrox/stackrox/central/resourcecollection/datastore/store/postgres/store_test.go:95\n        \t            \treflect.Value.call\n        \t            \t\t/usr/local/go/src/reflect/value.go:556\n        \t            \treflect.Value.Call\n        \t            \t\t/usr/local/go/src/reflect/value.go:339\n        \t            \tgithub.com/stretchr/testify/suite.Run.func1\n        \t            \t\t/go/pkg/mod/github.com/stretchr/testify@v1.8.0/suite/suite.go:175\n        \t            \ttesting.tRunner\n        \t            \t\t/usr/local/go/src/testing/testing.go:1439\n        \t            \truntime.goexit\n        \t            \t\t/usr/local/go/src/runtime/asm_amd64.s:1571\n        \tTest:       \tTestCollectionsStore/TestStore\n    store_test.go:98:\n        \tError Trace:\t/go/src/github.com/stackrox/stackrox/central/resourcecollection/datastore/store/postgres/store_test.go:98\n        \tError:      \tShould be false\n        \tTest:       \tTestCollectionsStore/TestStore\n    store_test.go:99:\n        \tError Trace:\t/go/src/github.com/stackrox/stackrox/central/resourcecollection/datastore/store/postgres/store_test.go:99\n        \tError:      \tExpected nil, but got: &storage.ResourceCollection{Id:\"a\", Name:\"a\", Description:\"a\", CreatedAt:&types.Timestamp{Seconds: 1,\n        \t            \tNanos: 1,\n        \t            \t}, LastUpdated:&types.Timestamp{Seconds: 1,\n        \t            \tNanos: 1,\n        \t            \t}, CreatedBy:(*storage.SlimUser)(0xc00085fb00), UpdatedBy:(*storage.SlimUser)(0xc00085fb40), ResourceSelectors:[]*storage.ResourceSelector{(*storage.ResourceSelector)(0xc00085fb80)}, EmbeddedCollections:[]*storage.ResourceCollection_EmbeddedResourceCollection{(*storage.ResourceCollection_EmbeddedResourceCollection)(0xc0011e00f0)}, XXX_NoUnkeyedLiteral:struct {}{}, XXX_unrecognized:[]uint8(nil), XXX_sizecache:0}\n        \tTest:       \tTestCollectionsStore/TestStore\n    store_test.go:114:\n        \tError Trace:\t/go/src/github.com/stackrox/stackrox/central/resourcecollection/datastore/store/postgres/store_test.go:114\n        \tError:      \tNot equal:\n        \t            \texpected: 200\n        \t            \tactual  : 201\n        \tTest:       \tTestCollectionsStore/TestStore",
					BuildId: "1",
				},
			},
			tests,
		)
	})
	t.Run("gradle", func(t *testing.T) {
		tests, err := findFailedTests("testdata/TEST-DefaultPoliciesTest.xml", map[string]string{"BUILD_ID": "1"}, 0)
		assert.NoError(t, err)

		assert.Equal(
			t,
			[]testCase{{
				Name: "Verify policy Apache Struts: CVE-2017-5638 is triggered",
				Message: "Condition not satisfied:\n" +
					"\n" +
					"waitForViolation(deploymentName,  policyName, 60)\n" +
					"|                |                |\n" +
					"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
					"",
				Stdout: "?[1;30m21:35:15?[0;39m | ?[34mINFO ?[0;39m | DefaultPoliciesTest       | Starting testcase\n" +
					"?[1;30m21:36:16?[0;39m | ?[34mINFO ?[0;39m | Services                  | Failed to trigger Apache Struts: CVE-2017-5638 after waiting 60 seconds\n" +
					"?[1;30m21:36:16?[0;39m | ?[1;31mERROR?[0;39m | Helpers                   | An exception occurred in test\n" +
					"org.spockframework.runtime.ConditionNotSatisfiedError: Condition not satisfied:\n" +
					"\n" +
					"waitForViolation(deploymentName,  policyName, 60)\n" +
					"|                |                |\n" +
					"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
					"\n" +
					"\tat DefaultPoliciesTest.$spock_feature_1_0(DefaultPoliciesTest.groovy:181) [1 skipped]\n" +
					"\tat util.OnFailureInterceptor.intercept(OnFailure.groovy:72) [8 skipped]\n" +
					"\tat util.OnFailureInterceptor.intercept(OnFailure.groovy:72) [10 skipped]\n" +
					" [6 skipped]\n" +
					"?[1;30m21:36:16?[0;39m | ?[39mDEBUG?[0;39m | Helpers                   | 2022-09-30 21:36:16 Will collect various stackrox logs for this failure under /tmp/qa-tests-backend-logs/a57dc4b9-70eb-4391-8a00-c5948fef733d/\n" +
					"?[1;30m21:37:07?[0;39m | ?[39mDEBUG?[0;39m | Helpers                   | Ran: ./scripts/ci/collect-service-logs.sh stackrox /tmp/qa-tests-backend-logs/a57dc4b9-70eb-4391-8a00-c5948fef733d/stackrox-k8s-logs\n" +
					"Exit: 0\n",
				Stderr:  "",
				Suite:   "DefaultPoliciesTest",
				BuildId: "1",
				Error: "Condition not satisfied:\n" +
					"\n" +
					"waitForViolation(deploymentName,  policyName, 60)\n" +
					"|                |                |\n" +
					"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
					"\n" +
					"\tat DefaultPoliciesTest.Verify policy #policyName is triggered(DefaultPoliciesTest.groovy:181)\n",
			}},
			tests,
		)
	})
}

func TestDescription(t *testing.T) {
	tc := testCase{
		Name: "Verify policy Apache Struts: CVE-2017-5638 is triggered",
		Message: "Condition not satisfied:\n" +
			"\n" +
			"waitForViolation(deploymentName,  policyName, 60)\n" +
			"|                |                |\n" +
			"false            qadefpolstruts   Apache Struts: CVE-2017-5638\n" +
			"",
		Stdout: "?[1;30m21:35:15?[0;39m | ?[34mINFO ?[0;39m | DefaultPoliciesTest       | Starting testcase\n" +
			"?[1;30m21:36:16?[0;39m | ?[34mINFO ?[0;39m | Services                  | Failed to trigger Apache Struts: CVE-2017-5638 after waiting 60 seconds\n" +
			"?[1;30m21:36:16?[0;39m | ?[1;31mERROR?[0;39m | Helpers                   | An exception occurred in test\n" +
			"org.spockframework.runtime.ConditionNotSatisfiedError: Condition not satisfied:\n",
		Stderr:  "",
		Suite:   "DefaultPoliciesTest",
		BuildId: "1",
	}
	actual, err := tc.description()
	assert.NoError(t, err)
	assert.Equal(t, `
{code:title=Message|borderStyle=solid}
Condition not satisfied:

waitForViolation(deploymentName,  policyName, 60)
|                |                |
false            qadefpolstruts   Apache Struts: CVE-2017-5638

{code}
{code:title=STDOUT|borderStyle=solid}
?[1;30m21:35:15?[0;39m | ?[34mINFO ?[0;39m | DefaultPoliciesTest       | Starting testcase
?[1;30m21:36:16?[0;39m | ?[34mINFO ?[0;39m | Services                  | Failed to trigger Apache Struts: CVE-2017-5638 after waiting 60 seconds
?[1;30m21:36:16?[0;39m | ?[1;31mERROR?[0;39m | Helpers                   | An exception occurred in test
org.spockframework.runtime.ConditionNotSatisfiedError: Condition not satisfied:

{code}

||    ENV     ||      Value           ||
| BUILD ID     | [1|https://prow.ci.openshift.org/view/gs/origin-ci-test/logs//1]|
| BUILD TAG    | [|]|
| JOB NAME     ||
| CLUSTER      ||
| ORCHESTRATOR ||
`, actual)
	s, err := tc.summary()
	assert.NoError(t, err)
	assert.Equal(t, `DefaultPoliciesTest / Verify policy Apache Struts  CVE-2017-5638 is triggered FAILED`, s)
}
