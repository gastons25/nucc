// Package handlers contains the routes handlers
package handlers

import(
	"time"
	"sync"
	"gastonstec/nuricc/memdb"
	"gastonstec/nuricc/models"
	"github.com/gofiber/fiber/v2"

)

const(
	STATUS_FAIL = "fail"
	STATUS_SUCCESS = "success"
)

// Response struct	
type TxResponse struct {
	Txid      string `json:"txid"`
	DateTime  string `json:"datetime"`
	Fee       string `json:"fee"`
	SentValue string `json:"sent_value"`
}

// Response struct	
type BlockResponse struct {
	Network       		string 	`json:"network"`
	BlockNo       		int    	`json:"block_no"`
	DateTime      		string  `json:"datetime"`
	PreviousBlockhash 	string 	`json:"previous_blockhash"`
	NextBlockhash     	string  `json:"next_blockhash"`
	Size              	int    	`json:"size"`
	Txs                 []TxResponse `json:"txs"`
}



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
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": err.Error(),
		})	
	}

	// Get block transactions
	numTxs := 10
	if len(block.Data.Txs) < 10 {
		numTxs = len(block.Data.Txs)
	}

	var txsMap sync.Map
	wg := sync.WaitGroup{}

	for i := 0; i < numTxs; i++ {
		wg.Add(1)
		go func(idx int) {
			tx, err := models.GetTxInfo(networkCode, block.Data.Txs[idx])
			if err == nil {
				txsMap.Store(idx, tx)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	// Set response
	blockResp := new(BlockResponse)
	blockResp.Network = block.Data.Network
	blockResp.BlockNo = block.Data.BlockNo
	blockResp.DateTime = time.Unix(block.Data.Time, 0).Format("02-Jan-2006 15:04")
	blockResp.PreviousBlockhash = block.Data.PreviousBlockhash
	blockResp.NextBlockhash = block.Data.NextBlockhash
	blockResp.Size = block.Data.Size

	// Add transactions to response
	for i := 0; i < numTxs; i++ {
		result, ok := txsMap.Load(i)
		if ok {
			txResp := new(TxResponse)
			txResp.Txid = result.(*models.TxInfo).Data.Txid
			txResp.DateTime = time.Unix(result.(*models.TxInfo).Data.Time, 0).Format("02-Jan-2006 15:04")
			txResp.Fee = result.(*models.TxInfo).Data.Fee
			txResp.SentValue = result.(*models.TxInfo).Data.SentValue			

			blockResp.Txs = append(blockResp.Txs, *txResp)
		}
	}

	// Return OK
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": blockResp,
		"status": STATUS_SUCCESS,
	})
}

// Get transaction handler function
func GetTx(ctx *fiber.Ctx) error {
	var err error

	// Get parameters
	networkCode := ctx.Params("network_code")
	txId := ctx.Params("transaction_id")

	// Validate parameters
	if networkCode == "" || txId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": "Invalid parameters",
		})
	}

	// Check network code on memdb
	netwname, _ := memdb.GetNetwork(networkCode)
	if netwname == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": "Invalid network code",
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

	// Fill and send response
	txResp := new(TxResponse)
	txResp.Txid = tx.Data.Txid
	txResp.DateTime = time.Unix(tx.Data.Time, 0).Format("02-Jan-2006 15:04")
	txResp.Fee = tx.Data.Fee
	txResp.SentValue = tx.Data.SentValue

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": STATUS_SUCCESS,
		"data": txResp,
	})
	
}