package server

import (
	urlshortenerctrl "github.com/arfan21/backend-test/internal/urlshortener/controller"
	urlshortenerrepo "github.com/arfan21/backend-test/internal/urlshortener/repository"
	urlshortenersvc "github.com/arfan21/backend-test/internal/urlshortener/service"
	"github.com/gofiber/fiber/v2"
)

func (s *Server) Routes() {

	urlshortedRepoRedis := urlshortenerrepo.NewRedis(s.dbRedis)
	urlshortedRepoPg := urlshortenerrepo.NewPostgres(s.dbPg)
	urlshortedSvc := urlshortenersvc.New(urlshortedRepoRedis, urlshortedRepoPg)
	urlshortedCtrl := urlshortenerctrl.NewHTTP(urlshortedSvc)

	s.app.Get("/:key", urlshortedCtrl.Redirect)
	s.app.Post("/shorten", urlshortedCtrl.Create)

	api := s.app.Group("/api")
	api.Get("/health-check", func(c *fiber.Ctx) error {
		err := s.dbPg.Ping(c.Context())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		err = s.dbRedis.Conn().Ping(c.Context()).Err()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "OK",
		})
	})

}
