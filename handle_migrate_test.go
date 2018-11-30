package goose

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSchema struct {
}

func (s *mockSchema) SetDebug(b bool) {
}

func (s *mockSchema) Reset() (m []string, d []string, err error) {
	return
}

func (s *mockSchema) Install() (err error) {
	return
}

func (s *mockSchema) Drop() (d []string, err error) {
	return
}

func (s *mockSchema) Migrate() (m []string, err error) {
	return
}

func newMockSchema() SchemaInterface {
	return &mockSchema{}
}

func resetTestEnv() {
	productionOperations = []string{}
	productionOperationsSet = false
	isProduction = false
	envSet = false
}

func Test_SetEnvProduction_True(t *testing.T) {
	resetTestEnv()

	SetEnvProduction(true)

	assert.True(t, isProduction, "isProduction flag not set correctly")
	assert.True(t, envSet, "envSet flag not set correctly")
}

func Test_SetEnvProduction_False(t *testing.T) {
	resetTestEnv()

	SetEnvProduction(false)

	assert.False(t, isProduction, "isProduction flag not set correctly")
	assert.True(t, envSet, "envSet flag not set correctly")
}

func Test_SetProductionOperations(t *testing.T) {
	resetTestEnv()

	allowed := []string{"qwerty", "asdf"}
	SetProductionOperations(allowed)
	ops := getProductionOperations()

	assert.Equal(t, allowed, ops, "production operations not set correctly")
}

func Test_SetProductionOperations_Defaults(t *testing.T) {
	resetTestEnv()

	ops := getProductionOperations()

	assert.Equal(t, []string{"up", "install"}, ops, "unexpected default production operations")
}

func Test_OperationAllowed_Production_Defaults_Drop(t *testing.T) {
	resetTestEnv()

	SetEnvProduction(true)
	allowed := operationAllowed("drop")

	require.True(t, envSet, "failed to set envSet")
	require.True(t, isProduction, "failed to set isProduction")
	assert.False(t, allowed, "failed to block 'drop' operation in production by default")
}

func Test_OperationAllowed_Production_Defaults_Reset(t *testing.T) {
	resetTestEnv()

	SetEnvProduction(true)
	allowed := operationAllowed("reset")

	require.True(t, envSet, "failed to set envSet")
	require.True(t, isProduction, "failed to set isProduction")
	assert.False(t, allowed, "failed to block 'reset' operation in production by default")
}

func Test_OperationAllowed_Production_Defaults_Install(t *testing.T) {
	resetTestEnv()

	SetEnvProduction(true)
	allowed := operationAllowed("install")

	require.True(t, envSet, "failed to set envSet")
	require.True(t, isProduction, "failed to set isProduction")
	assert.True(t, allowed, "failed to allow 'install' operation in production by default")
}

func Test_OperationAllowed_Production_Defaults_Up(t *testing.T) {
	resetTestEnv()

	SetEnvProduction(true)
	allowed := operationAllowed("up")

	require.True(t, envSet, "failed to set envSet")
	require.True(t, isProduction, "failed to set isProduction")
	assert.True(t, allowed, "failed to allow 'up' operation in production by default")
}

func Test_HandleMigrate_RequireEnv(t *testing.T) {
	resetTestEnv()
	s := newMockSchema()
	var a []string
	r := true

	m, d, err := HandleMigrate(s, a, &r)

	assert.NotNil(t, err, "HandleMigrate failed to return an error")
	assert.Len(t, m, 0, "HandleMigrate returned migrated tables")
	assert.Len(t, d, 0, "HandleMigrate returned dropped tables")
}

func Test_HandleMigrate_ProductionOperations_Defaults_Drop(t *testing.T) {
	resetTestEnv()
	s := newMockSchema()
	a := []string{"drop"}
	r := true

	SetEnvProduction(true)
	m, d, err := HandleMigrate(s, a, &r)

	assert.NotNil(t, err, "HandleMigrate failed to return an error")
	assert.Len(t, m, 0, "HandleMigrate returned migrated tables")
	assert.Len(t, d, 0, "HandleMigrate returned dropped tables")
}

func Test_HandleMigrate_ProductionOperations_Defaults_Reset(t *testing.T) {
	resetTestEnv()
	s := newMockSchema()
	a := []string{"reset"}
	r := true

	SetEnvProduction(true)
	m, d, err := HandleMigrate(s, a, &r)

	assert.NotNil(t, err, "HandleMigrate failed to return an error")
	assert.Len(t, m, 0, "HandleMigrate returned migrated tables")
	assert.Len(t, d, 0, "HandleMigrate returned dropped tables")
}

func Test_HandleMigrate_ProductionOperations_Defaults_Up(t *testing.T) {
	resetTestEnv()
	s := newMockSchema()
	a := []string{"up"}
	r := true

	SetEnvProduction(true)
	_, _, err := HandleMigrate(s, a, &r)

	assert.Nil(t, err, "HandleMigrate returned an error")
}

func Test_HandleMigrate_ProductionOperations_Defaults_Install(t *testing.T) {
	resetTestEnv()
	s := newMockSchema()
	a := []string{"install"}
	r := true

	SetEnvProduction(true)
	_, _, err := HandleMigrate(s, a, &r)

	assert.Nil(t, err, "HandleMigrate returned an error")
}

// @TODO investigate Table-Driven Testing
//func Test_HandleMigrate(t *testing.T) {
//	type handleMigrateAssertion struct {
//		MigratedLen   int
//		DroppedLen    int
//		ReturnedError bool
//	}
//	type handleMigrateTest struct {
//		Production bool
//		Args       []string
//		DryRun     bool
//		Require    *handleMigrateAssertion
//		Assert     *handleMigrateAssertion
//	}
//	tests := map[string]handleMigrateTest{
//		"Test_HandleMigrate_ProductionOperations_Defaults_Drop": {
//			Production: true,
//			Args:       []string{"drop"},
//			DryRun:     true,
//			Assert: &handleMigrateAssertion{
//				MigratedLen:   0,
//				DroppedLen:    0,
//				ReturnedError: true,
//			},
//		},
//	}
//	for _, test := range tests {
//		resetTestEnv()
//		s := newMockSchema()
//
//		SetEnvProduction(test.Production)
//		m, d, err := HandleMigrate(s, test.Args, &test.DryRun)
//
//		if test.Require != nil {
//			require.Len(t, m, test.Assert.MigratedLen)
//			require.Len(t, d, test.Assert.DroppedLen)
//			if test.Assert.ReturnedError {
//				require.NotNil(t, err)
//			} else {
//				require.Nil(t, err)
//			}
//		}
//
//		if test.Assert != nil {
//			assert.Len(t, m, test.Assert.MigratedLen)
//			assert.Len(t, d, test.Assert.DroppedLen)
//			if test.Assert.ReturnedError {
//				assert.NotNil(t, err)
//			} else {
//				assert.Nil(t, err)
//			}
//		}
//	}
//}
