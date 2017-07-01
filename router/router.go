package router

import (
	"github.com/gorilla/mux"
	"github.com/holwech/learnxyz-backend/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}

var routes = Routes{
	Route{
		"addTopics",
		"POST",
		"/topics/add",
		handlers.AddTopic,
	},
	Route{
		"searchTopics",
		"GET",
		"/topics/search",
		handlers.SearchTopics,
	},
	Route{
		"getTopics",
		"GET",
		"/topics/get",
		handlers.GetTopics,
	},
	Route{
		"addResource",
		"POST",
		"/resources/add",
		handlers.AddResource,
	},
	Route{
		"getResources",
		"GET",
		"/resources/get",
		handlers.GetResources,
	},
}
