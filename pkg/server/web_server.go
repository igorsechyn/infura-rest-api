package server

import (
	"context"
	"ethereum-api/pkg/app"
	"ethereum-api/pkg/server/routes"
	"net/http"

	"github.com/go-chi/chi"
)

func Run(ctx context.Context, app *app.App) {
	server := &server{}
	server.Run(ctx, app)
}

type server struct{}

func (webServer *server) Run(ctx context.Context, app *app.App) {
	router := chi.NewRouter()

	router.Use(CreateHttResponseLogger(app.Reporter))

	router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	routes.RegisterGetTransactionByBlockNumberAndIndex(router, app)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.Reporter.Error("server.startup.error", err, map[string]interface{}{})
		}
	}()

	<-ctx.Done()
	err := srv.Shutdown(context.Background())
	if err != nil {
		app.Reporter.Error("server.shutdown", err, map[string]interface{}{})
	}
}
