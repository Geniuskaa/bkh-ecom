package v1

import (
	"bkh-ecom/internal/domain"
	"bkh-ecom/internal/dto"
	"bkh-ecom/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

type Router struct {
	fiber.Router
	clickService service.ClickService
}

func NewRoute(router fiber.Router, clickService service.ClickService) *Router {
	r := Router{
		Router:       router,
		clickService: clickService,
	}

	return &r
}

func (r *Router) Routes() {
	r.Get("/counter/:bannerID", r.counter)
	r.Post("/stats/:bannerID", r.stats)
}

// Сохраняет полученное событие
func (r *Router) counter(c *fiber.Ctx) error {

	bannerIDString := c.Params("bannerID")
	bannerID, err := strconv.Atoi(bannerIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bannerID parameter",
		})
	}

	clickInfo := domain.Click{
		BannerID:  bannerID,
		ClickTime: time.Now(),
	}

	r.clickService.SaveClick(clickInfo)

	return c.SendStatus(fiber.StatusAccepted)
}

// Возвращает события полученные поле применения фильтра
func (r *Router) stats(c *fiber.Ctx) error {

	bannerIDString := c.Params("bannerID")
	bannerID, err := strconv.Atoi(bannerIDString)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid bannerID parameter",
		})
	}

	filter := dto.ClicksStatRequest{
		BannerID: bannerID,
	}
	if err := c.BodyParser(&filter); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON body for interval",
		})
	}

	data, err := r.clickService.ListClicks(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Неизвестная ошибка",
		})
	}

	return c.JSON(data)
}
