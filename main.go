/*
 * Copyright (C) 2023 Coeus
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License version 3 as published by
 * the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/agpl-3.0.txt>.
 */

package main

import (
	"coeus/controllers"
	"coeus/models"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	globals "coeus/globals"
)

//go:embed views/templates/**/*
var templatesEmbed embed.FS

//go:embed views/static/*
var staticEmbed embed.FS

func main() {

	// Load environment variables from .env file
	err := godotenv.Load("globals/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Running as a single binary: " + strconv.FormatBool(globals.IsBinary()))
	_ = models.NewDB()

	router := gin.Default()

	// Serve static files before the organization is set up
	staticFS, _ := fs.Sub(staticEmbed, "views/static")
	router.StaticFS("/static", http.FS(staticFS))

	templ := template.Must(template.New("").ParseFS(templatesEmbed, "views/templates/coeus/*.html", "views/templates/management/*.html", "views/templates/shared/*.html"))
	router.SetHTMLTemplate(templ)

	router.Use(sessions.Sessions("session", cookie.NewStore(globals.Secret)))

	// Public
	public := router.Group("/")
	public.Use(controllers.OrganizationRequired)
	controllers.PublicRoutes(public)

	// Private
	private := router.Group("/")
	private.Use(controllers.AuthRequired)
	controllers.PrivateRoutes(private)

	// API
	api := router.Group("/")
	controllers.APIRoutes(api)

	// Websockets
	ws := router.Group("/")
	controllers.WSRoutes(ws)

	// Catch all
	router.NoRoute(controllers.NotFoundHandler)

	// Start server
	port := os.Getenv("PORT")
	runOnHeroku := os.Getenv("RUN_ON_HEROKU")

	if runOnHeroku != "true" {
		if port == "" {
			port = "8080"
		}
	} else {
		if port == "" {
			log.Fatal("$PORT must be set")
		}
	}

	router.Run(":" + port)
}