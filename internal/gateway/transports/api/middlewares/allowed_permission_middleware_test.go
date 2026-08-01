package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	contexts "github.com/HiIamJeff67/notezy-backend/internal/gateway/contexts"
	sharedtypes "github.com/HiIamJeff67/notezy-backend/internal/shared/types"
)

func TestAllowedPermissionsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name                string
		middleware          gin.HandlerFunc
		expectedPermissions []sharedtypes.AccessControlPermission
	}{
		{
			name:       "above includes the requested permission",
			middleware: AllowedPermissionsAbove(sharedtypes.AccessControlPermission_Admin),
			expectedPermissions: []sharedtypes.AccessControlPermission{
				sharedtypes.AccessControlPermission_Admin,
				sharedtypes.AccessControlPermission_Owner,
			},
		},
		{
			name:       "below includes the requested permission",
			middleware: AllowedPermissionsBelow(sharedtypes.AccessControlPermission_Write),
			expectedPermissions: []sharedtypes.AccessControlPermission{
				sharedtypes.AccessControlPermission_Read,
				sharedtypes.AccessControlPermission_Write,
			},
		},
		{
			name: "within preserves the explicit permission set",
			middleware: AllowedPermissionsWithin(
				sharedtypes.AccessControlPermission_Owner,
				sharedtypes.AccessControlPermission_Write,
			),
			expectedPermissions: []sharedtypes.AccessControlPermission{
				sharedtypes.AccessControlPermission_Owner,
				sharedtypes.AccessControlPermission_Write,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/", testCase.middleware, func(ctx *gin.Context) {
				allowedPermissions, exception := contexts.GetAllowedPermissions(ctx.Request.Context())
				if exception != nil {
					t.Fatal(exception)
				}

				if len(allowedPermissions) != len(testCase.expectedPermissions) {
					t.Fatalf("expected %d permissions, got %d", len(testCase.expectedPermissions), len(allowedPermissions))
				}
				for index, expectedPermission := range testCase.expectedPermissions {
					if allowedPermissions[index] != expectedPermission {
						t.Fatalf("expected permission %s at index %d, got %s", expectedPermission, index, allowedPermissions[index])
					}
				}

				ctx.Status(http.StatusNoContent)
			})

			responseRecorder := httptest.NewRecorder()
			router.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodGet, "/", nil))

			if responseRecorder.Code != http.StatusNoContent {
				t.Fatalf("expected status %d, got %d", http.StatusNoContent, responseRecorder.Code)
			}
		})
	}
}
