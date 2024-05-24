package exception

import (
	"errors"
	"net/http"

	"github.com/arfan21/backend-test/pkg/constant"
	"github.com/arfan21/backend-test/pkg/logger"
	"github.com/arfan21/backend-test/pkg/pkgutil"

	"github.com/gofiber/fiber/v2"
)

func FiberErrorHandler(ctx *fiber.Ctx, err error) error {
	defer func() {
		logger.Log(ctx.UserContext()).Error().Msg(err.Error())
	}()

	defaultRes := pkgutil.HTTPResponse{
		Code: fiber.StatusInternalServerError,
	}

	var withCodeErrorArr constant.ErrsWithCode

	var withCodeErrors constant.ErrsWithCode
	if errors.As(err, &withCodeErrors) && len(withCodeErrors) > 0 {
		defaultRes.Code = withCodeErrors[0].HTTPCode
		withCodeErrorArr = withCodeErrors
	}

	var withCodeError constant.ErrWithCode
	if errors.As(err, &withCodeError) {
		defaultRes.Code = withCodeError.HTTPCode
		withCodeErrorArr = append(withCodeErrorArr, withCodeError)
	}

	var fiberError *fiber.Error
	if errors.As(err, &fiberError) {
		defaultRes.Code = fiberError.Code
		errwithcode := constant.ErrWithCode{
			HTTPCode: fiberError.Code,
			Message:  fiberError.Message,
		}
		withCodeErrorArr = append(withCodeErrorArr, errwithcode)
	}

	if defaultRes.Code >= 500 {
		defaultRes.Errors = []constant.ErrWithCode{
			{
				HTTPCode: http.StatusInternalServerError,
				Message:  http.StatusText(http.StatusInternalServerError),
			},
		}
	} else {
		if len(withCodeErrorArr) > 0 {
			defaultRes.Errors = withCodeErrorArr
		}
	}

	return ctx.Status(defaultRes.Code).JSON(defaultRes)
}
