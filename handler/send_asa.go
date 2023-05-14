package handler

import (
	"net/http"

	"github.com/Meningtov/algonea-backend/config"
	"github.com/Meningtov/algonea-backend/utils"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/gin-gonic/gin"
)

// Some testnet ASA
const rewardAssetID = 175771249

type response struct {
	Transactions []utils.ResponseTransaction `json:"transactions"`
}

// SendAsa is an example endpoint which builds an atomic group transaction of size 2
// 1. Transfer ASA from creator account to the address specified in the path param
// 2. User opt-in to ASA
// They either both succeed or both fail
func SendAsa(c *gin.Context) {
	appCtx, err := config.GetAppContext()
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError)
		return
	}

	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, MissingPathParam)
		return
	}

	suggestedParams, err := appCtx.AlgodClient.SuggestedParams().Do(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, MissingPathParam)
		return
	}

	sendAsaTxn, err := future.MakeAssetTransferTxn(
		appCtx.CreatorAccount.Address.String(),
		address,
		1,
		nil,
		suggestedParams,
		"",
		rewardAssetID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError)
		return
	}

	userOptInTxn, err := future.MakeAssetAcceptanceTxn(
		address,
		nil,
		suggestedParams,
		rewardAssetID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError)
		return
	}

	transactions, err := utils.NewGroupTransactionBuilderWithSigner(appCtx.CreatorAccount.PrivateKey).
		Add(sendAsaTxn, true).    // true means that this transaction is signed by the creator
		Add(userOptInTxn, false). // false means that this transaction is to be signed by the user in the browser
		Build()
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError)
		return
	}

	c.JSON(http.StatusOK, response{Transactions: transactions})
}
