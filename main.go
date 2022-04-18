package main

import (
	"myapp/data"
	"myapp/handlers"
	"myapp/middleware"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/brucebotes/celeritas"
)

type application struct {
	App        *celeritas.Celeritas
	Handlers   *handlers.Handlers
	Models     data.Models
	Middleware *middleware.Middleware
	wg         sync.WaitGroup
}

func main() {
	c := initApplication()
	go c.listenForShutdown()
	err := c.App.ListenAndServe()
	c.App.ErrorLog.Println(err)
}

func (a *application) shutdown() {
	// TODO put any clean up tasks here!!

	// this will block until waitgroup is empty
	a.wg.Wait()
}

func (a *application) listenForShutdown() {
	// listen for SIGINT and SIGTERM and pipe it into a channel
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit

	a.App.InfoLog.Println("Received signal", s.String())
	a.shutdown()

	os.Exit(0)
}
