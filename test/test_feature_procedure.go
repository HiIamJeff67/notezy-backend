package test

import "testing"

// the below test procedures should all be implemented
// whatever you'll use it or not, just implement them all
type TestFeatureProcedureInterface interface {
	Main(t *testing.T)
	BeforeAll(t *testing.T)
	BeforeEach(t *testing.T)
	AfterEach(t *testing.T)
	AfterAll(t *testing.T)
}
