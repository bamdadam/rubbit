package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/bamdadam/rubbit/internal/http/request"
	"github.com/bamdadam/rubbit/internal/rabbit"
	"github.com/bamdadam/rubbit/internal/store/rdb"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
)

type Handler struct {
	RH  *rabbit.RabbitHandler
	RDB *rdb.RedisStore
}

func (h *Handler) PublishMessage(c *fiber.Ctx) error {
	body := new(request.PublishRequest)
	err := c.BodyParser(body)
	if err != nil {
		log.Error("error while parsing request body: ", err)
		return fiber.ErrBadRequest
	}
	if body.IsDelayed {
		if !strings.Contains(body.PublishDelay, "ms") {
			return c.Status(fiber.ErrBadRequest.Code).SendString("wrong publish delay format. should be: 100ms")
		}
		pubDelayDuration, err := time.ParseDuration(body.PublishDelay)
		if pubDelayDuration.Milliseconds() < 5000 {
			return c.Status(fiber.ErrBadRequest.Code).SendString("cant have delay less than 5 seconds")
		}
		if err != nil {
			log.Error("error while parsing request publish delay: ", err)
			return fiber.ErrBadRequest
		}
		err = h.RH.PublishDelayedMessage(body.Topic, body.Message, pubDelayDuration.Milliseconds())
		if err != nil {
			log.Error("error while publishing message: ", err)
			return fiber.ErrInternalServerError
		}
	} else {
		err := h.RH.PublishMessage(body.Topic, body.Message)
		if err != nil {
			log.Error("error while publishing message: ", err)
			return fiber.ErrInternalServerError
		}
	}
	return nil
}

func (h *Handler) GetSubjectMessages(c *fiber.Ctx) error {
	body := new(request.GetMessageRequest)
	err := c.BodyParser(body)
	if err != nil {
		log.Error("error while parsing request body: ", err)
		return fiber.ErrBadRequest
	}
	messages, err := h.RDB.GetMessage(c.Context(), body.Topic)
	if err != nil {
		log.Error("error while reading message from redis: ", err)
		return fiber.ErrBadRequest
	}
	return c.Status(http.StatusOK).JSON(messages)
}

func (h *Handler) RegisterHandler(g fiber.Router) {
	g.Post("publish/", h.PublishMessage)
	g.Get("subject/", h.GetSubjectMessages)
}
