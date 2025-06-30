package authe2etest

import (
	"fmt"
	"testing"
)

const (
	testTargetPath = "notezy-backend/app/routes/test_routes/auth_route.go"
)

// - how to initialize the router :
// 		gin.SetMode(gin.TestMode)
// 		router := gin.New()
// 		testRouterGroup := router.Group("/testRegisterRoute")
// 		testroutes.ConfigureTestAuthRoutes(testRouterGroup)

func TestObjectInParallel(t *testing.T) {
	t.Run(fmt.Sprintf("E2E-Test---Auth-(%s):", testTargetPath), func(t *testing.T) { // object level
		t.Run("Test-Register-Route", func(t *testing.T) {
			t.Parallel()
			// TestRegisterRoute(t)
		})
		// login
		// logout...
	})
}
