package handler_test

import (
	"net/http"
	"testing"

	"github.com/Meningtov/algonea-backend/api"
	"github.com/Meningtov/algonea-backend/testdata"
	"github.com/steinfletcher/apitest"
	jsonpath "github.com/steinfletcher/apitest-jsonpath"
)

func TestSendAsa(t *testing.T) {
	apitest.HandlerFunc(api.Handler).
		Getf("/api/account/%s/send-asa", testdata.UserAddress).
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len("transactions", 2)).
		Assert(jsonpath.Equal("transactions[0].requires_signing", false)).
		Assert(jsonpath.Equal("transactions[1].requires_signing", true)).
		End()
}
