package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/src/main/app/services"
)

type IConsumerHandler interface {
	GetStatus(ctx *fiber.Ctx) error
	Start(ctx *fiber.Ctx) error
	Stop(ctx *fiber.Ctx) error
}

type ConsumerHandler struct {
	consumerService services.IConsumerService
}

func NewConsumerHandler(consumerService services.IConsumerService) *ConsumerHandler {
	return &ConsumerHandler{
		consumerService: consumerService,
	}
}

// GetStatus godoc
//
// @Summary		Get status for current consumer
// @Description	Started or stopped
// @Tags		Consumer
// @Success		200
// @Accept 		json
// @Produce		json
// @Success     200 {object} model.AppStatusDTO
// @Router		/consumer/status [get].
func (h ConsumerHandler) GetStatus(ctx *fiber.Ctx) error {
	result := h.consumerService.GetAppStatus()
	return ctx.JSON(result)
}

// Start godoc
//
// @Summary		Start consumer
// @Description	Starts the consumer
// @Tags		Consumer
// @Success		200
// @Accept 		json
// @Produce		json
// @Success     200 {object} model.AppStatusDTO
// @Router		/consumer/start [put].
func (h ConsumerHandler) Start(ctx *fiber.Ctx) error {
	err := h.consumerService.Start()
	if err != nil {
		return err
	}

	result := h.consumerService.GetAppStatus()
	return ctx.JSON(result)
}

// Stop godoc
//
// @Summary		Get status for current consumer
// @Description	Started or stopped
// @Tags		Consumer
// @Success		200
// @Accept 		json
// @Produce		json
// @Success     200 {object} model.AppStatusDTO
// @Router		/consumer/stop [put].
func (h ConsumerHandler) Stop(ctx *fiber.Ctx) error {
	err := h.consumerService.Stop()
	if err != nil {
		return err
	}

	result := h.consumerService.GetAppStatus()
	return ctx.JSON(result)
}
