package handler

import (
	"github.com/Meningtov/algonea_backend/config"
	"github.com/Meningtov/algonea_backend/utils"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Some testnet ASA
const rewardAssetID = 175771249

type response struct {
	Transactions []utils.ResponseTransaction `json:"transactions"`
}

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
		Add(sendAsaTxn, true).
		Add(userOptInTxn, false).
		Build()
	if err != nil {
		c.JSON(http.StatusInternalServerError, InternalServerError)
		return
	}

	c.JSON(http.StatusOK, response{Transactions: transactions})
}
