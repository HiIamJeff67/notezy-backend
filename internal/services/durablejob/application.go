package durablejob

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	config "github.com/HiIamJeff67/notezy-backend/internal/platform/config"
	observability "github.com/HiIamJeff67/notezy-backend/internal/platform/observability"
	data "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/data"
	routinetask "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/routinetask"
	yjsmaintenance "github.com/HiIamJeff67/notezy-backend/internal/services/durablejob/yjsmaintenance"
)

func Start() {
	shutdownObservability := observability.Initialize(context.Background())
	defer shutdownObservability()

	db := data.ConnectToDatabase(config.PostgresDatabaseConfig)
	defer data.DisconnectToDatabase(db)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	routineTaskEngine := routinetask.NewEngine(db)
	shutdownRoutineTaskEngine := routineTaskEngine.Start(ctx)
	defer shutdownRoutineTaskEngine()

	yjsMaintenanceEngine := yjsmaintenance.NewEngine(db)
	shutdownYjsMaintenanceEngine := yjsMaintenanceEngine.Start(ctx)
	defer shutdownYjsMaintenanceEngine()

	<-ctx.Done()
}
