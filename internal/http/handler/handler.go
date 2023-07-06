package handler

import (
	"github.com/bamdadam/rubbit/internal/http/request"
	"github.com/bamdadam/rubbit/internal/rabbit"
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
)

type Handler struct {
	RS *rabbit.RabbitServer
}

func (h *Handler) PublishMessage(c *fiber.Ctx) error {
	body := new(request.PublishRequest)
	err := c.BodyParser(body)
	if err != nil {
		log.Error("error while parsing request body: ", err)
		return fiber.ErrBadRequest
	}
	err = h.RS.PublishMessage(body.Topic, body.Message)
	if err != nil {
		log.Error("error while publishing message: ", err)
		return fiber.ErrInternalServerError
	}
	return nil
}

func (h *Handler) GetSubjectMessages(c *fiber.Ctx) error {
	return nil
}

func (h *Handler) RegisterHandler(g fiber.Router) {
	g.Post("publish/", h.PublishMessage)
	g.Get("subject/", h.GetSubjectMessages)
}