package http

import (
	"fmt"
	"net/http"

	"github.com/abstractbreazy/rate-limited-api/internal/repository/redis"
	"github.com/abstractbreazy/rate-limited-api/pkg/limiter"
	"github.com/abstractbreazy/rate-limited-api/pkg/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	redis   *redis.Redis
	limiter *limiter.Config
}

func NewHandler(limiter *limiter.Config, redis *redis.Redis) (h *Handler) {
	h = new(Handler)
	h.redis = redis
	h.limiter = limiter
	return
}

func (h *Handler) Ping(c echo.Context) error {
	subnet, err := h.limiter.Subnet(c.RealIP())
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, err.Error()))
	}

	value, err := h.redis.Get(c.Request().Context(), subnet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"),
		)
	}

	if value.Value > h.limiter.Limit {
		return c.JSON(http.StatusTooManyRequests,
			utils.NewHTTPError(http.StatusTooManyRequests, "Too many request"),
		)
	}

	requestsCount, err := h.redis.Incr(c.Request().Context(), subnet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"),
		)
	}

	if requestsCount > h.limiter.Limit {
		return c.JSON(http.StatusTooManyRequests,
			utils.NewHTTPError(http.StatusTooManyRequests, "Too many request"),
		)
	}

	return c.String(http.StatusOK, "Pong\n")
}

func (h *Handler) Reset(c echo.Context) error {
	subnet, err := h.limiter.Subnet(c.RealIP())
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewHTTPError(http.StatusBadRequest, err.Error()))
	}

	err = h.redis.Del(c.Request().Context(), subnet)
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			utils.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"),
		)
	}

	return c.String(http.StatusOK, fmt.Sprintf("Reset rate limit for subnet : %s\n", subnet))
}
