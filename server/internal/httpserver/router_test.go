package httpserver

import (
	"net/http"
	"testing"
)

func TestMount_noRouteConflicts(t *testing.T) {
	t.Helper()
	mux := http.NewServeMux()
	Mount(mux, &App{})
}
