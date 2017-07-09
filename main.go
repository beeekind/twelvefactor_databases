package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// sqlx in a minimal extension to sql/db
	"github.com/jmoiron/sqlx"
	// postgres driver for sqlx
	_ "github.com/lib/pq"

	// Gathers and castes environmental variables
	"github.com/kelseyhightower/envconfig"
	// Minimal router middleware that extends net/http
	"github.com/gorilla/mux"
	// Simple Ping service
	"github.com/b3ntly/twelvefactor_databases/ping"
	// Users service: GetAll
	// Authentication service: Register, Login
	"github.com/b3ntly/twelvefactor_databases/users"
)

// This application is composed of submodules which mount routes to a router.
type Service interface {
	Mount(*mux.Router)
}

// Environment declares variables gathered from our environment using kelseyhightower/envconfig.
// The goal is to expose as much configuration of your application as possible so knobs can be easily turned
// by your devops team.
type Environment struct {
	Port               string        `envconfig:"PORT" default:"9090"`
	PingPath           string        `envconfig:"PING_PATH" default:"/ping"`
	PingResponse       string        `envconfig:"PING_RESPONSE" default:"PONG"`
	ReqTimeout         time.Duration `envconfig:"REQ_TIMEOUT" default:"500ms"`
	ServerReadTimeout  time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"1000ms"`
	ServerWriteTimeout time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"2000ms"`
	DBConnMaxLifetime  time.Duration `envconfig:"DB_CONN_MAX_LIFETIME" default:"0"`
	DBMaxOpen          int           `envconfig:"DB_MAX_OPEN" default:"0"`
	DBMaxIdle          int           `envconfig:"DB_MAX_IDLE" default:"2"`
	PostgresURI        string        `envconfig:"POSTGRES_URI" default:"postgresql://postgres@localhost:5432/postgres?sslmode=disable"`
	// Expose the users service at this path, defaults to /users.
	UsersPathPrefix    string        `envconfig:"USERS_PATH" default:"users"`
	// The number of users returned by the /users endpoint.
	SelectManyLimit    int           `envconfig:"USERS_SELECT_LIMIT" default:"10"`
}

// Here we define a middleware that injects a context with a timeout.
// If that timeout is exceeded the request returns a error code indicating timeout.
func injectContextWithTimeout(ctx context.Context, reqTimeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, reqTimeout)

		// Defer cancel to prevent context leaking.
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctxWithTimeout))
	})
}

// Return a sqlx database client.
func getDatabaseConnection(ctx context.Context, env *Environment) (*sqlx.DB, error) {
	database, err := sqlx.ConnectContext(ctx, "postgres", env.PostgresURI)

	if err != nil {
		return nil, err
	}

	// If d <= 0, connections are reused forever. (default=0)
	database.SetConnMaxLifetime(env.DBConnMaxLifetime)
	// If SetMaxIdleConns <= 0, then there is no limit on the number of open connections (default=0).
	database.SetMaxOpenConns(env.DBMaxOpen)
	// If n <= 0, no idle connections are retained. (default=2).
	database.SetMaxIdleConns(env.DBMaxIdle)
	return database, nil
}

// Return a customized http.Server. Many defaults within net/http are bad for protection so it's important to do this.
func buildServer(env *Environment, mux http.Handler) *http.Server {
	return &http.Server{
		// Defaults to no timeouts, which is really bad
		// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		// And: https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout:  env.ServerReadTimeout,
		WriteTimeout: env.ServerWriteTimeout,
		Addr:         fmt.Sprintf(":%s", env.Port),
		Handler:      mux,
	}
}

func main() {
	// Contexts can be used for request-scoped variables (like user ids) or cancellation (like request timeouts). This will be
	// the root context for our application.
	ctx := context.Background()

	// You can replace this with something like Logrus or Zap in production.
	logger := log.New(os.Stdout, "", log.Lshortfile)

	// Initialize kelseyhightower/envconfig library.
	env := &Environment{}
	if err := envconfig.Process("", env); err != nil {
		logger.Fatal(err)
	}

	// Connect with and ping our database client.
	database, err := getDatabaseConnection(ctx, env)
	if err != nil {
		logger.Fatal(err)
	}

	// This is our root routing component provided by gorilla/mux. All service routes and subrouters will be mounted
	// to it. Note services are fully capable of overriding each-other if they have identical paths.
	router := mux.NewRouter()

	// Instantiate the service(s) with requisite configurations.
	services := []Service{
		ping.New(&ping.Config{
			PingPath:     env.PingPath,
			PingResponse: env.PingResponse,
			Logger:       logger,
		}),

		users.New(&users.Config{
			Ctx:             ctx,
			Logger:          logger,
			DB:              database,
			UsersPathPrefix: env.UsersPathPrefix,
			SelectManyLimit: env.SelectManyLimit,
		}),
	}

	for _, service := range services {
		service.Mount(router)
	}

	// instantiate the http.Server with our router
	server := buildServer(env, injectContextWithTimeout(ctx, env.ReqTimeout, router))

	// start the server
	logger.Fatal(server.ListenAndServe())
}
