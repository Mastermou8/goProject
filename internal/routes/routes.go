package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/mastermou8/goProject/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleWorkoutByID)
	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)
	//r.Get("/workout/{id}", app.WorkoutHandler.HandleWorkoutByID)
	//r.Post("/workout", app.WorkoutHandler.HandleCreateWorkout)

	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkoutByID)
	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkoutByID)
	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandleCreateToken)
	return r
}
