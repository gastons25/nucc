package handlers

import(
	"gastonstec/nuricc/models"
	"gastonstec/nuricc/memdb"
	"github.com/gofiber/fiber/v2"

)

const(
	STATUS_FAIL = "fail"
	STATUS_SUCCESS = "success"
)

func GetBlock(ctx *fiber.Ctx) error {
	var err error

	// Get parameters
	networkCode := ctx.Params("network_code")
	blockHash := ctx.Params("block_hash")

	// Validate network code
	netwname, _ := memdb.GetNetwork(networkCode)

	// Validate parameters
	if netwname == "" || blockHash == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": "Invalid parameters",
		})
	}

	// Get block information
	var data string
	data, err = models.GetBlockInfo(networkCode, blockHash)

	// Check and return error
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": STATUS_FAIL,
			"data": err.Error(),
		})	
	}

	// Return OK
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": STATUS_SUCCESS,
		"data": data,
	})
}