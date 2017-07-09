package users

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type (
	// Config for the users service.
	Config struct {
		Ctx             context.Context
		// The path prefix to expose the subrouter provided by this service, defaults to /users.
		UsersPathPrefix string
		Logger          *log.Logger
		DB              *sqlx.DB
		// The number of users to return from the Get endpoint. Defaults to 10.
		SelectManyLimit int
	}

	// Service: users.
	Service struct {
		ctx             context.Context
		db              *sqlx.DB
		pathPrefix      string
		logger          *log.Logger
		selectManyLimit int
	}

	// User model for the table defined in sql.go .
	User struct {
		ID        int64  `json:"ID" db:"id"`
		Username  string `json:"username" db:"username"`
		CreatedAt string `json:"createdAt" db:"created_at"`
	}
)

// New: Instantiate a new users service. Fail hard if errors occur (our application shouldn't run without this service).
func New(config *Config) *Service {
	_, err := config.DB.Exec(CreateTableStmt)

	if err != nil {
		config.Logger.Fatal(err)
	}

	return &Service{
		ctx:             config.Ctx,
		db:              config.DB,
		pathPrefix:      config.UsersPathPrefix,
		logger:          config.Logger,
		selectManyLimit: config.SelectManyLimit,
	}
}

// Mount the subRouter of this service to the root router.
func (s *Service) Mount(r *mux.Router) {
	subRouter := r.PathPrefix(filepath.Join("/", s.pathPrefix)).Subrouter()
	subRouter.HandleFunc("", s.Get)
}

// Get endpoint returns an array of users with JSON encoding.
func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	results := []*User{}
	err := s.db.Select(&results, SelectManyStmt, s.selectManyLimit)

	if err != nil {
		s.writeError(w, err)
		return
	}

	response, err := json.Marshal(results)
	if err != nil {
		s.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Error handling logic for this service.
func (s *Service) writeError(w http.ResponseWriter, err error) {
	// for tracing purposes you may also access s.ctx here...
	s.logger.Println(err)

	// here you can decide what type of error to return to the user, you should usually refrain from the actual error
	// in case it contains sensitive information. That said you should try to tell the user something helpful.
	http.Error(w, "", http.StatusInternalServerError)
}
