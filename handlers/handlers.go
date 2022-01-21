// Package handlers contains the routes handlers
package handlers

import(
	"time"
	"gastonstec/nuricc/models"
	"gastonstec/nuricc/memdb"
	"github.com/gofiber/fiber/v2"

)

const(
	STATUS_FAIL = "fail"
	STATUS_SUCCESS = "success"
)

// Get block handler function
func GetBlock(ctx *fiber.Ctx) error {
	var err error

	// Get parameters
	networkCode := ctx.Params("network_code")
	blockHash := ctx.Params("block_hash")

	// Validate parameters
	netwname, _ := memdb.GetNetwork(networkCode) // check network code on memdb
	if netwname == "" || blockHash == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": "Invalid parameters",
		})
	}

	// Get block information
	var block *models.BlockInfo
	block, err = models.GetBlockInfo(networkCode, blockHash)

	// Check and return error
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": err.Error(),
		})	
	}

	// Return OK
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": block,
		"status": STATUS_SUCCESS,
	})
}

// Get transaction handler function
func GetTx(ctx *fiber.Ctx) error {
	var err error

	// Get parameters
	networkCode := ctx.Params("network_code")
	txId := ctx.Params("transaction_id")

	// Validate network code
	netwname, _ := memdb.GetNetwork(networkCode)

	// Validate parameters
	if netwname == "" || txId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": "Invalid parameters",
		})
	}

	// Get transaction information
	var tx *models.TxInfo
	tx, err = models.GetTxInfo(networkCode, txId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": err.Error(),
		})	
	}

		
	type TxResponse struct {
		Txid      string `json:"txid"`
		DateTime  string `json:"datetime"`
		Fee       string `json:"fee"`
		SentValue string `json:"sent_value"`
	}

	var txResp TxResponse
	txResp.Txid = tx.Data.Txid
	txResp.DateTime = time.Unix(tx.Data.Time, 0).Format("dd-MM-yyyy HH:mm")
	txResp.Fee = tx.Data.Fee
	txResp.SentValue = tx.Data.SentValue

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": STATUS_SUCCESS,
		"data": txResp,
	})
	
}