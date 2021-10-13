package main

import (
	"context"
	ctx "github.com/gorilla/context"
	"github.com/justinas/alice"
	"github/gvozdev08101993/db"
	"github/gvozdev08101993/pkg"
	"log"
	"net/http"
)

func main() {
	var err error

	emptyCtx := context.Background()
	db.ClientDB, err = db.ConnectAndPingToClientDB()
	if err != nil {
		log.Fatalf("failed to connect with client database: %v", err)
	}
	defer db.ClientDB.Disconnect(emptyCtx)

	commonHandlers := alice.New(ctx.ClearHandler, pkg.LoggingHandler, pkg.RecoverHandler, pkg.AcceptHandler)
	router := pkg.NewRouter()

	log.Println("teas-api is running...")
	router.Get("/teas/:name", commonHandlers.ThenFunc(pkg.GetTeaHandler))
	router.Post("/teas", commonHandlers.Append(pkg.ContentTypeHandler, pkg.BodyParserHandler(pkg.TeaResource{})).ThenFunc(pkg.CreateTeaHandler))
	router.Delete("/teas/:name", commonHandlers.ThenFunc(pkg.DeleteTeaHandler))
	router.Put("/teas/:name", commonHandlers.Append(pkg.ContentTypeHandler, pkg.BodyParserHandler(pkg.TeaResource{})).ThenFunc(pkg.UpdateTeaHandler))
	router.Get("/teas", commonHandlers.ThenFunc(pkg.GetAllTeasHandler))
	http.ListenAndServe(":8080", router)
}
