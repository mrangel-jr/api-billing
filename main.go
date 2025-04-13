package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.cm/mrangel-jr/api-billing/internals/app"
	"github.cm/mrangel-jr/api-billing/internals/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8081, "Go backend server port")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	app.Logger.Printf("We are running out app on port %d", port)

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
