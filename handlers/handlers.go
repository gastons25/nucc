package handlers

import(
	"gastonstec/nuricc/models"
	"github.com/gofiber/fiber/v2"

)

const(
	STATUS_FAIL = "fail"
	STATUS_SUCCESS = "success"
)

func GetBlock(ctx *fiber.Ctx) error {
	var err error

	// comment
	networkCode := ctx.Params("network_code")
	blockHash := ctx.Params("block_hash")


	var data string
	data, err = models.GetBlockInfo(networkCode, blockHash)

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