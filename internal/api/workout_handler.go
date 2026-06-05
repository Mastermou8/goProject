package api

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"encoding/json"
	"strconv"

	"github.com/mastermou8/goProject/internal/store"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore //this is interface
}

func NewWorkoutHandler(workoutStore store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
	}
}

func (wh *WorkoutHandler) HandleWorkoutByID(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutID := chi.URLParam(r, "id")
	if paramsWorkoutID == "" {
		http.NotFound(w, r)
		return
	}

	workoutID, err := strconv.ParseInt(paramsWorkoutID, 10, 64)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Workout ID: %d", workoutID)
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		fmt.Println("Error decoding workout:", err)                  //this is for us
		http.Error(w, "Invalid request body", http.StatusBadRequest) // this is for the client
		return
	}

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		fmt.Println("Error creating workout:", err)                               //this is for us
		http.Error(w, "Failed to create workout", http.StatusInternalServerError) // this is for the client
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdWorkout)

}
