package ping

import (
	"encoding/json"
	"log"
	"net/http"
)

// Service contains an HTTP endpoint for ping/pong functionality
type Service struct {
	// don't export variables you don't need to
	pingResponse string
	logger       *log.Logger
}

// New: inject dependencies via an explicit constructor. Though sometimes people will read environmental variables or
// initialize defaults here I prefer to do so explicitly within the program entry-point.
func New(pingResponse string, logger *log.Logger) *Service {
	return &Service{pingResponse: pingResponse, logger: logger}
}

// Endpoint, Alternatively you could return a mux or sub-router from this subpackage.
// Note that we name it ServeHTTP so that Service fulfills the http.Handler interface.
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := json.Marshal(s.pingResponse)

	if err != nil {
		s.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// logic for logging and writing an error, log your errors!
func (s *Service) writeError(w http.ResponseWriter, err error) {
	s.logger.Println(err)

	// here you can decide what type of error to return to the user, you should usually refrain from the actual error
	// in case it contains sensitive information. That said you should try to tell the user something helpful.
	http.Error(w, "", http.StatusInternalServerError)
}
