package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/alveary/overseer/announce"
	"github.com/alveary/user-registration/registration"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
)

// Result is a combination of a given registration and the related errors
type Result struct {
	Registration registration.Registration
	Errors       map[string]string
}

func errorMap(errors binding.Errors) map[string]string {
	errorMap := map[string]string{}

	for _, error := range errors {
		if len(error.FieldNames) > 0 {
			errorMap[error.FieldNames[0]] = error.Message
		}
	}

	return errorMap
}

// AppEngine for web engine setup
func AppEngine() *martini.ClassicMartini {
	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", nil)
	})

	m.Head("/alive", func(resp http.ResponseWriter) {
		resp.WriteHeader(http.StatusOK)
	})

	m.Post("/", binding.Form(registration.Registration{}), func(r render.Render, errors binding.Errors, registration registration.Registration, resp http.ResponseWriter) {
		if errors.Len() > 0 {

			errorMap := errorMap(errors)
			r.HTML(200, "index", Result{Registration: registration, Errors: errorMap})

		} else {

			target, err := registration.RequestRegistration()
			fmt.Println(err)

			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				r.HTML(200, "failure", err)
				return
			}

			r.HTML(200, "success", target)
		}
	})

	return m
}

func init() {
	rootURL := os.Getenv("ROOT_URL")
	aliveURL := os.Getenv("ALIVE_URL")

	if rootURL == "" || aliveURL == "" {
		fmt.Println("ROOT_URL and/or ALIVE_URL not set")
		return
	}

	announce.NewService(
		"user-registration",
		rootURL,
		aliveURL)
}

func main() {
	var port int
	flag.IntVar(&port, "p", 9000, "the port number")
	flag.Parse()

	m := AppEngine()
	m.RunOnAddr(":" + strconv.Itoa(port))
}
