package urlshortenerctrl

import (
	"github.com/arfan21/backend-test/internal/model"
	"github.com/arfan21/backend-test/internal/urlshortener"
	_ "github.com/arfan21/backend-test/pkg/constant"
	_ "github.com/arfan21/backend-test/pkg/pkgutil"
	"github.com/gofiber/fiber/v2"
)

type HTTPController struct {
	svc urlshortener.Service
}

func NewHTTP(svc urlshortener.Service) *HTTPController {
	return &HTTPController{
		svc: svc,
	}
}

// @Summary Create shorted url
// @Description Create shorted url
// @Tags urlshortener
// @Accept json
// @Produce json
// @Param input body model.CreateURLShortedRequest true "Create URL Shorted Request"
// @Success 200 {object} model.CreateURLShortedResponse
// @Failure 400 {object} pkgutil.HTTPResponse{errors=[]constant.ErrWithCode}
// @Router /shorten [post]
func (h HTTPController) Create(c *fiber.Ctx) error {
	req := model.CreateURLShortedRequest{}
	err := c.BodyParser(&req)
	if err != nil {
		return err
	}

	res, err := h.svc.Create(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// @Summary Redirect to original url
// @Description Redirect to original url
// @Tags urlshorted
// @Accept json
// @Produce json
// @Param key path string true "Shorted URL Key"
// @Success 302 {string} string
// @Router /{key} [get]
func (h HTTPController) Redirect(c *fiber.Ctx) error {
	key := c.Params("key")

	res, err := h.svc.Get(c.Context(), key)
	if err != nil {
		return err
	}

	return c.Redirect(res, fiber.StatusFound)
}
