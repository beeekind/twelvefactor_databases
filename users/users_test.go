package users_test

import (
	"context"
	"errors"
	"log"
	"testing"

	// sqlx in a minimal extension to sql/db
	"github.com/jmoiron/sqlx"
	// postgres driver for sqlx
	_ "github.com/lib/pq"

	// Minimal router middleware that extends net/http
	"encoding/json"
	"fmt"
	"github.com/b3ntly/twelvefactor_databases/users"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http/httptest"
	"os"
)

const (
	usersPathPrefix = "users"
	postgresURI     = "postgresql://postgres@localhost:5432/postgres?sslmode=disable"
	selectManyLimit = 10
)

// Connect to the database client, create the users table if it doesn't exist, delete any preexisting rows,
// and populate it with 10 new users for testing purposes.
func setupDatabase(t testing.TB, ctx context.Context) *sqlx.DB {
	database, err := sqlx.ConnectContext(ctx, "postgres", postgresURI)
	require.Nil(t, err)
	require.Nil(t, bootstrapDatabase(ctx, database))
	require.Nil(t, cleanDatabase(ctx, database))
	require.Nil(t, populateDatabase(ctx, database))
	return database
}

// Create the users table if it doesn't exist.
func bootstrapDatabase(ctx context.Context, database *sqlx.DB) error {
	_, err := database.ExecContext(ctx, users.CreateTableStmt)
	return err
}

// Delete any existing rows in the database.
func cleanDatabase(ctx context.Context, database *sqlx.DB) error {
	_, err := database.ExecContext(ctx, users.DeleteManyStmt)
	return err
}

// Insert users into the database for testing purposes.
func populateDatabase(ctx context.Context, database *sqlx.DB) error {
	for i := 0; i < selectManyLimit; i++ {
		_, err := database.ExecContext(ctx, users.InsertOneStmt, "fred")

		if err != nil {
			return err
		}
	}

	return nil
}

// Test the Get endpoint of the users service.
func TestService_Get(t *testing.T) {
	ctx := context.Background()

	db := setupDatabase(t, ctx)
	router := mux.NewRouter()

	service := users.New(&users.Config{
		Ctx:             ctx,
		Logger:          log.New(os.Stdout, "logger: ", log.Lshortfile),
		DB:              db,
		UsersPathPrefix: usersPathPrefix,
		SelectManyLimit: selectManyLimit,
	})

	service.Mount(router)

	req := httptest.NewRequest("GET", "http://localhost:9090/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode >= 299 {
		require.Nil(t, errors.New(fmt.Sprintf("Invalid status code: %s", res.StatusCode)))
	}

	response, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)

	rows := make([]*users.User, 0)
	err = json.Unmarshal(response, &rows)
	require.Nil(t, err)
	require.Equal(t, 10, len(rows))
}

// YOU MIGHT NEED TO RAISE YOUR ULIMIT ON MACOS TO RUN THIS
func BenchmarkService_Ping(b *testing.B) {
	ctx := context.Background()

	db := setupDatabase(b, ctx)
	router := mux.NewRouter()

	service := users.New(&users.Config{
		Ctx:             ctx,
		Logger:          log.New(os.Stdout, "logger: ", log.Lshortfile),
		DB:              db,
		UsersPathPrefix: usersPathPrefix,
		SelectManyLimit: selectManyLimit,
	})

	service.Mount(router)
	req := httptest.NewRequest("GET", "http://localhost:9090/users/", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
