package handler_test

import (
	"net/http"
	"testing"

	"github.com/Meningtov/algonea-backend/api"
	"github.com/steinfletcher/apitest"
)

func TestHealthcheck(t *testing.T) {
	apitest.HandlerFunc(api.Handler).
		Get("/api/health").
		Expect(t).
		Status(http.StatusOK).
		Body(`{"status":"OK"}`).
		End()
}
