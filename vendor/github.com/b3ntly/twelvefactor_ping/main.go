package ping

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/TV4/env"
	"github.com/b3ntly/twelvefactor_ping/ping"
)

// Using middleware to compose route definitions can be extremely powerful. Here we define a middleware
// that injects a context with a timeout. If that timeout is exceeded the request returns a error code
// indicating timeout.
func injectContextWithTimeout(ctx context.Context, reqTimeout time.Duration, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, reqTimeout)

		// defer cancel to prevent context leaking
		defer cancel()

		next.ServeHTTP(w, r.WithContext(ctxWithTimeout))
	})
}

// break individual initialization steps into small sequential pieces
func buildServer(serverReadTimeout time.Duration, serverWriteTimeout time.Duration, port string, mux http.Handler) *http.Server {
	return &http.Server{
		// Defaults to no timeouts, which is really bad
		// See: https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
		// And: https://blog.cloudflare.com/exposing-go-on-the-internet/
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
	}
}

func main() {
	// use the background context from which we will derive all future contexts
	// contexts can be used for request-scoped variables (like user ids) or cancellation (like request timeouts)
	ctx := context.Background()

	// you can replace this with something like Logrus or Zap in production
	logger := log.New(os.Stdout, "", log.Lshortfile)

	// replace with gorilla mux etc. if desired
	router := http.NewServeMux()

	// read configuration from our environment with default values
	// ...expose as much configuration as possible to your users/ops team
	// PORT to serve the application on
	port := env.String("PORT", "9090")

	// PATH to return a response from, defaults to /ping
	endpoint := filepath.Join("/", env.String("ENDPOINT", "ping"))

	// DEFAULT_RESPONSE to send
	pingResponse := env.String("DEFAULT_RESPONSE", "PONG")

	// REQ_TIMEOUT timeout for http handler
	reqTimeout := time.Duration(env.Int("REQ_TIMEOUT", 500)) * time.Millisecond

	// SERVER_READ_TIMEOUT timeout for the server, REQ_TIMEOUT should occur before a server timeout in most circumstances
	serverReadTimeout := time.Duration(env.Int("SERVER_READ_TIMEOUT", 1000)) * time.Millisecond

	// SERVER_WRITE_TIMEOUT timeout for the server, REQ_TIMEOUT should occur before a server timeout
	serverWriteTimeout := time.Duration(env.Int("SERVER_WRITE_TIMEOUT", 2000)) * time.Millisecond

	// instantiate the service
	service := ping.New(pingResponse, logger)

	// build the router with desired middleware
	router.Handle(endpoint, injectContextWithTimeout(ctx, reqTimeout, service))

	// instantiate the http.Server with our router
	server := buildServer(serverReadTimeout, serverWriteTimeout, port, router)

	// start the server
	logger.Fatal(server.ListenAndServe())
}
