package authe2etest

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	models "notezy-backend/app/models"
	testroutes "notezy-backend/app/routes/test_routes"
	shared "notezy-backend/shared"
	test "notezy-backend/test"
)

const (
	testTargetPath         = "notezy-backend/app/routes/test_routes/auth_route.go"
	testAuthRouteNamespace = "/testRoute/auth"
)

type testAuthFeatureProcedure struct {
	testDB          *gorm.DB
	testRouter      *gin.Engine
	testRouterGroup *gin.RouterGroup
}

func (p *testAuthFeatureProcedure) BeforeAll(t *testing.T) {
	p.testDB = models.ConnectToDatabase(shared.PostgresDatabaseConfig)
	gin.SetMode(gin.TestMode)
	p.testRouter = gin.New()
	p.testRouterGroup = p.testRouter.Group(testAuthRouteNamespace)
	testroutes.ConfigureTestAuthRoutes(p.testDB, p.testRouterGroup)
}

func (p *testAuthFeatureProcedure) BeforeEach(t *testing.T) { /* Do Nothing */ }

func (p *testAuthFeatureProcedure) AfterEach(t *testing.T) { /* Do Nothing */ }

func (p *testAuthFeatureProcedure) AfterAll(t *testing.T) {
	models.DisconnectToDatabase(p.testDB)
}

func (p *testAuthFeatureProcedure) Main(t *testing.T) {
	t.Run(fmt.Sprintf("E2E-Test---Auth-(%s):", testTargetPath), func(t *testing.T) { // feature level
		t.Run("Test-Register-Route", func(t *testing.T) { // spec level
			var registerE2ETester = NewRegisterE2ETester(p.testRouter)
			if registerE2ETester == nil {
				t.Fatal("NewRegisterE2ETester returned nil, router may be nil")
			}

			t.Run("[Valid-Test-Account]", func(t *testing.T) { // case level
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterValidTestAccount(t)
				p.AfterEach(t)
			})
			t.Run("[Valid-User-Account]", func(t *testing.T) { // case level
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterValidUserAccount(t)
				p.AfterEach(t)
			})
			t.Run("[No-Name]", func(t *testing.T) { // case level
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterNoName(t)
				p.AfterEach(t)
			})
			t.Run("[Name-Without-Number]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterNameWithoutNumber(t)
				p.AfterEach(t)
			})
			t.Run("[Short-Name]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterShortName(t)
				p.AfterEach(t)
			})
			t.Run("[Invalid-Email]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterInvalidEmail(t)
				p.AfterEach(t)
			})
			t.Run("[Short-Password]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterShortPassword(t)
				p.AfterEach(t)
			})
			t.Run("[Password-Without-Lower-Case-Letter]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterPasswordWithoutLowerCaseLetter(t)
				p.AfterEach(t)
			})
			t.Run("[Password-Without-Upper-Case-Letter]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterPasswordWithoutUpperCaseLetter(t)
				p.AfterEach(t)
			})
			t.Run("[Password-Without-Number]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterPasswordWithoutNumber(t)
				p.AfterEach(t)
			})
			t.Run("[Password-Without-Sign]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				registerE2ETester.TestRegisterPasswordWithoutSign(t)
				p.AfterEach(t)
			})
		})
		t.Run("Test-Login-Route", func(t *testing.T) { // spec
			var loginE2ETester = NewLoginE2ETester(p.testRouter)
			if loginE2ETester == nil {
				t.Fatal("NewLoginE2ETester returned nil, router may be nil")
			}

			t.Run("[Valid-Test-Account-By-Name]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				loginE2ETester.TestLoginValidTestAccountByName(t)
				p.AfterEach(t)
			})
			t.Run("[Valid-Test-Account-By-Email]", func(t *testing.T) {
				t.Parallel()
				p.BeforeEach(t)
				loginE2ETester.TestLoginValidTestAccountByEmail(t)
				p.AfterEach(t)
			})
		})
		// logout...
	})
}

func TestMain(t *testing.T) {
	var testRegisterFeatureProcedure test.TestFeatureProcedureInterface = &testAuthFeatureProcedure{}

	testRegisterFeatureProcedure.BeforeAll(t)
	defer testRegisterFeatureProcedure.AfterAll(t)
	testRegisterFeatureProcedure.Main(t)
}
