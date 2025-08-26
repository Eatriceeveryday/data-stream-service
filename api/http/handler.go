package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Eatriceeveryday/data-stream-service/internal/service"
	"github.com/labstack/echo/v4"
)

type SensorHandler struct {
	s *service.EMQXService
}

func NewHandler(s *service.EMQXService) *SensorHandler {
	return &SensorHandler{
		s: s,
	}
}

func (h *SensorHandler) ChangeFrequency(c echo.Context) error {
	duration := c.QueryParam("d")
	if duration == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing duration param"})
	}

	d, err := strconv.Atoi(duration)
	if err != nil || d <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid duration"})
	}

	h.s.ChangeInterval(int32(d))

	return c.JSON(http.StatusOK, map[string]string{
		"message": fmt.Sprintf("interval updated to %d seconds", d),
	})
}
