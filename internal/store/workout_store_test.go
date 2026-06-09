package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	//run migrations for our test db
	err = Migrate(db, "../../migrations")
	if err != nil {
		t.Fatalf("migrations error: %v", err)
	}

	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
	return db

}

func TestCreateWorkout(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db)

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body dat",
				DurationMinutes: 60,
				CaloriesBurned:  200,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "bench press",
						Reps:         IntPtr(10),
						Sets:         3,
						Weight:       FloatPtr(100.0),
						Notes:        StringPtr("hekkingly heavy"),
						OrderIndex:   1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "workout with invalid entry",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body dat",
				DurationMinutes: 90,
				CaloriesBurned:  500,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "plank",
						Reps:         IntPtr(10),
						Sets:         3,
						Weight:       FloatPtr(100.0),
						Notes:        StringPtr("hekkingly heavy"),
						OrderIndex:   1,
					},
					{
						ExerciseName:    "bench press",
						Reps:            IntPtr(10),
						Sets:            3,
						DurationSeconds: IntPtr(60),
						Weight:          FloatPtr(185.0),
						Notes:           StringPtr("hekkingly heavy"),
						OrderIndex:      2,
						// This entry is invalid because it has both reps and duration_seconds set, which should not be allowed
						// The CreateWorkout method should return an error when trying to insert this workout into the database
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutByID(int64(createdWorkout.ID))
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Reps, retrieved.Entries[i].Reps)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
			}
		})
	}
}

func IntPtr(i int) *int {
	return &i
}

func FloatPtr(f float64) *float64 {
	return &f
}

func StringPtr(s string) *string {
	return &s
}
