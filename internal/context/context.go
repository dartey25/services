package context

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/mdoffice/md-services/internal/log"
	"github.com/mdoffice/md-services/web/views/component"
)

type AppContext struct {
	echo.Context
}

// NewAppContext creates a new AppContext.
func NewAppContext(c echo.Context) *AppContext {
	return &AppContext{
		Context: c,
	}
}

func (c *AppContext) Log() log.Logger {
	return c.Get("logger").(log.Logger)
}

func (c *AppContext) Template(statusCode int, t templ.Component) error {
	c.Response().Writer.WriteHeader(statusCode)
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return t.Render(c.Request().Context(), c.Response().Writer)
}

func ErrorTemplate(message, variant string) templ.Component {
	return component.Alert(component.AlertProps{Text: message, Variant: variant, ClassName: "m-3"})
}

func (c *AppContext) BadRequestTemplate(t templ.Component) error {
	var template templ.Component
	if t != nil {
		template = t
	} else {
		template = ErrorTemplate("Невірний запит", "danger")
	}
	return c.Template(http.StatusBadRequest, template)
}

func (c *AppContext) NotFoundTemplate(t templ.Component) error {
	var template templ.Component
	if t != nil {
		template = t
	} else {
		template = ErrorTemplate("Не знайдено", "info")
	}
	return c.Template(http.StatusNotFound, template)
}

func (c *AppContext) InternalServerErrorTemplate(t templ.Component) error {
	var template templ.Component
	if t != nil {
		template = t
	} else {
		template = ErrorTemplate("Помилка", "danger")
	}
	return c.Template(http.StatusInternalServerError, template)
}
