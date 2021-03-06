package ping_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/b3ntly/twelvefactor_ping/ping"
	"github.com/stretchr/testify/require"
)

const defaultPingResponse = "PONG"

func TestService_Endpoint(t *testing.T) {
	service := ping.New(defaultPingResponse, log.New(os.Stdout, "logger: ", log.Lshortfile))

	req := httptest.NewRequest("GET", "http://localhost:9090/ping", nil)
	w := httptest.NewRecorder()
	service.ServeHTTP(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode >= 299 {
		require.Nil(t, errors.New("Invalid status code"))
	}

	response, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	res.Body.Close()

	// json.Marshal will surround the string with quotes so we marshal it here to get an identical representation for comparison
	repr, err := json.Marshal(defaultPingResponse)
	require.Nil(t, err)
	require.True(t, bytes.Equal(repr, response))
}

// YOU MIGHT NEED TO RAISE YOUR ULIMIT ON MACOS TO RUN THIS
func BenchmarkService_Endpoint(b *testing.B) {
	service := ping.New(defaultPingResponse, log.New(os.Stdout, "logger: ", log.Lshortfile))
	req := httptest.NewRequest("GET", "http://localhost:9090/ping", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		service.ServeHTTP(w, req)
	}
}
