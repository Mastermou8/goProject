package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mastermou8/goProject/internal/app"
	"github.com/mastermou8/goProject/internal/routes"
)

func main() {
	var port int //name of flag
	flag.IntVar(&port, "port", 8080, "Go backend server port")
	flag.Parse()

	//start of the application
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
		//cricial errors to stop app
	}
	defer app.DB.Close()
	//initial test passed, now we can start the server

	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.Logger.Printf("We are running on port %d", port)

	//start the server
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}

}
