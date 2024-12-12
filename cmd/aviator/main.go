package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/Ckefa/aviator/internals/services"
	"github.com/Ckefa/aviator/models"
	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func NewRenderer() *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/home/*.html")),
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	println("Starting the program")

	flight := models.NewFlight()
	flight.Run()

	// Initialize Socket.IO server
	go services.RunSocketIO()

	// Initialize Echo
	e := echo.New()

	e.Renderer = NewRenderer()
	e.Static("/static", "static")

	// e.GET("/socket.io/", echo.WrapHandler(server))
	// e.POST("/socket.io/", echo.WrapHandler(server))

	e.GET("/", func(c echo.Context) error {

		return c.Render(http.StatusOK, "index", nil)
	})

	e.Logger.Fatal(e.Start(":7000"))
	flight.Wait.Wait()
}
