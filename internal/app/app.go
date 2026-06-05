package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mastermou8/goProject/internal/api"
	"github.com/mastermou8/goProject/internal/store"
	"github.com/mastermou8/goProject/migrations"
)

// Application struct to hold application-wide dependencies
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	DB             *sql.DB
}

// NewApplication initializes the application and its dependencies
func NewApplication() (*Application, error) {
	pgDB, err := store.Open()
	if err != nil {
		return nil, err
	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//out stores will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore)
	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		DB:             pgDB,
	}
	//returns application and nil (not significant errors)
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
