package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/justinas/alice"
	"github.com/neustar/httprouter"
	"github.com/rs/xlog"
)

// App application
type App struct {
	Name string
}

// InternalHandler internal
type InternalHandler struct {
	h func(w http.ResponseWriter, r *http.Request)
}

// ServeHTTP serve
func (ih InternalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ih.h(w, r)
	return
}

func (app *App) hello(w http.ResponseWriter, r *http.Request) {
	l := xlog.FromRequest(r)
	l.Info("hello handler")
	log.Println("this is usual logger")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "hello")
	return
}

func main() {
	host, _ := os.Hostname()
	conf := xlog.Config{
		Fields: xlog.F{
			"role": "my-service",
			"host": host,
		},
		Output: xlog.NewOutputChannel(xlog.NewConsoleOutput()),
	}

	c := alice.New(
		xlog.NewHandler(conf),
		xlog.MethodHandler("method"),
		xlog.URLHandler("url"),
		xlog.UserAgentHandler("user_agent"),
		xlog.RefererHandler("referer"),
		xlog.RequestIDHandler("req_id", "Request-Id"),
		accessLoggingMiddleware,
	)
	app := App{Name: "my-service"}
	r := httprouter.New()
	// r.GET("/hello", c.Then(InternalHandler{h: app.hello}))
	r.GET("/hello", http.HandlerFunc(app.hello))

	xlog.Info("xlog")
	xlog.Infof("chain: %+v", c)
	log.Println("start server")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
