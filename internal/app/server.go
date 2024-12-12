package app

import (
	v1 "bkh-ecom/internal/app/api/v1"
	"bkh-ecom/internal/logger"
	"bkh-ecom/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	*fiber.App
}

func ServerInit(_ context.Context, clickService service.ClickService) Server {
	srv := Server{}

	srv.App = fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ReadTimeout:   time.Minute * 1,
		WriteTimeout:  time.Second * 20,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		Network:       "tcp",
	})

	//Кастомный recover
	srv.Use(func(ctx *fiber.Ctx) error {
		defer func() {
			err := recover()
			if err != nil {
				p := ctx.Route().Path
				logger.ErrorKV(ctx.Context(), logger.Data{
					Msg:    "panic occured",
					Error:  fmt.Errorf("%v", err),
					Detail: fmt.Sprintf("Path: %v", p),
				})

				ctx.SendStatus(http.StatusInternalServerError)
			}
		}()

		return ctx.Next()
	})

	apiGroup := srv.Group("/api")
	v1.NewRoute(apiGroup, clickService).Routes()

	return srv
}
