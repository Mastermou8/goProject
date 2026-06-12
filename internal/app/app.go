package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mastermou8/goProject/internal/api"
	"github.com/mastermou8/goProject/internal/middleware"
	"github.com/mastermou8/goProject/internal/store"
	"github.com/mastermou8/goProject/migrations"
)

// Application struct to hold application-wide dependencies
type Application struct {
	Logger         *log.Logger
	WorkoutHandler *api.WorkoutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	Middleware     *middleware.UserMiddleware
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
	userStore := store.NewPostgresUserStore(pgDB)
	tokenStore := store.NewPostgresTokenStore(pgDB)

	// our handlers will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)
	userHandler := api.NewUserHandler(userStore, logger)
	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)
	middlewareHandler := &middleware.UserMiddleware{UserStore: userStore}

	app := &Application{
		Logger:         logger,
		WorkoutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		Middleware:     middlewareHandler,
		DB:             pgDB,
	}
	//returns application and nil (not significant errors)
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
