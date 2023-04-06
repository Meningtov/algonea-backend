package handler_test

import (
	"github.com/Meningtov/algonea_backend/api"
	"github.com/steinfletcher/apitest"
	"net/http"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	apitest.HandlerFunc(api.Handler).
		Get("/api/health").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"status":"OK"}`).
		End()
}
