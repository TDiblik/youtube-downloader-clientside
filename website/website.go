package main

import (
	"errors"
	"html/template"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	engine.AddFunc(
		"htmlsafe", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/assets", "./assets")

	app.Get("/", func(c fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"PageType": "index",
		})
	})

	// takes "id" / "url" query parameters
	app.Get("/api/v1/youtube-download-urls", func(ctx fiber.Ctx) error {
		params := ctx.Queries()
		query_id := params["id"]
		query_url := params["url"]

		if query_id == "" && query_url == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Bad Request",
				"message": "You have to supply either a video id or an url to a video",
			})
		}

		// Get id (which cannot be trusted yet) either from query or url
		var possibly_invalid_id string
		if query_id != "" {
			possibly_invalid_id = query_id
		} else {
			var err error
			possibly_invalid_id, err = GetIdFromYoutubeUrl(query_url)
			if err != nil {
				if errors.Is(err, ErrUnableToParseUrl) {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "Unable to parse the url supplied",
					})
				} else if errors.Is(err, ErrNotYoutubeUrl) {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "Provided url is not officialy from youtube",
					})
				} else if errors.Is(err, ErrNoIdProvidedInsideUrl) {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "The url you provided does not include an id",
					})
				} else {
					return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"message": "Error type not handled on the server. Please report: https://github.com/TDiblik/youtube-downloader-clientside/issues/new",
					})
				}
			}
		}

		// additional cleanup, probably unnecessary
		possibly_invalid_id = strings.TrimSpace(possibly_invalid_id)
		possibly_invalid_id = strings.TrimRight(possibly_invalid_id, "/")

		if !IsYoutubeIdValid(possibly_invalid_id) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "The youtube video ID provided is NOT valid.",
			})
		}
		valid_youtube_id := possibly_invalid_id // reassign for clarity

		return ctx.JSON(fiber.Map{
			"message": "valid id: " + valid_youtube_id,
		})
	})

	app.Listen(":3000")
}
