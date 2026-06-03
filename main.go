package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/mastermou8/goProject/internal/app"
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
	//initial test passed, now we can start the server
	app.Logger.Println("Application running Successfully")
	http.HandleFunc("/health", HealthCheck)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
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

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
