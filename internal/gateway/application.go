package gateway

import (
	"fmt"
	"os"
	"strings"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"

	developmentroutes "github.com/HiIamJeff67/notezy-backend/internal/gateway/transports/api/routes/developmentroutes"
)

func Start() {
	developmentroutes.DevelopmentRouter = gin.Default()
	proxies := strings.Split(os.Getenv("GIN_TRUSTED_PROXIES"), ",")
	if err := developmentroutes.DevelopmentRouter.SetTrustedProxies(proxies); err != nil {
		fmt.Println("Failed to set trusted proxies for router: ", err)
		return
	}

	developmentroutes.ConfigureDevelopmentRoutes()
	ginAddress := os.Getenv("GIN_DOMAIN") + ":" + os.Getenv("GIN_PORT")
	if err := endless.ListenAndServe(ginAddress, developmentroutes.DevelopmentRouter); err != nil {
		fmt.Println("Failed to connect to the server")
	}
}
