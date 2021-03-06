package server

import (
	"covid19-update-service/model"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (ws *Covid19UpdateWebServer) registerEventRoutes(r *mux.Router) {
	eventsRouter := r.Path(eventsBaseRoute).Subrouter().StrictSlash(strictSlash)
	eventsRouter.HandleFunc("", ws.checkAcceptType(ws.getEvents)).Methods(http.MethodGet)
	eventsRouter.HandleFunc("", ws.optionHandler(eventsRouter)).Methods(http.MethodOptions)
	eventsRouter.Use(newCorsHandler(eventsRouter))
	eventsRouter.Use(ws.authorizationAndIdentification())
	eventsRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(eventsRouter)

	eventRouter := r.Path(eventRoute).Subrouter().StrictSlash(strictSlash)
	eventRouter.HandleFunc("", ws.checkAcceptType(ws.getEvent)).Methods(http.MethodGet)
	eventRouter.HandleFunc("", ws.optionHandler(eventRouter)).Methods(http.MethodOptions)
	eventRouter.Use(newCorsHandler(eventRouter))
	eventRouter.Use(ws.authorizationAndIdentification())
	eventRouter.MethodNotAllowedHandler = ws.createNotAllowedHandler(eventRouter)
}

func (ws *Covid19UpdateWebServer) getEvents(w http.ResponseWriter, r *http.Request) {
	t, ok := findTopic(w, r)
	if !ok {
		return
	}
	query := r.URL.Query()
	var limit uint
	var err error
	limitRequested := false
	if rawLimit, ok := query["limit"]; ok {
		limit, err = toUInt(rawLimit[0])
		limitRequested = true
		if err != nil {
			writeHTTPResponse(model.NewError("Limit has to be an unsigned integer."), http.StatusBadRequest, w, r)
			return
		}
	}
	var e model.EventCollection
	if limitRequested {
		e, err = model.GetEventsWithLimit(t.ID, limit)
	} else {
		e, err = model.GetEvents(t.ID)
	}
	if err != nil {
		writeHTTPResponse(model.NewError(fmt.Sprintf("Could not load events: %v.", err)), http.StatusInternalServerError, w, r)
		return
	}

	writeHTTPResponse(e, http.StatusOK, w, r)
}

func (ws *Covid19UpdateWebServer) getEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eID, err := toUInt(vars[eventId])
	if err != nil {
		writeHTTPResponse(model.NewError("Event ID has to be an unsigned integer."), http.StatusBadRequest, w, r)
		return
	}

	t, ok := findTopic(w, r)
	if !ok {
		return
	}
	e, err := model.GetEvent(eID, t.ID)
	if e == nil {
		writeHTTPResponse(model.NewError("Could not find event."), http.StatusNotFound, w, r)
		return
	}
	if err != nil {
		writeHTTPResponse(model.NewError(fmt.Sprintf("Could not load event: %v.", err)), http.StatusInternalServerError, w, r)
		return
	}

	writeHTTPResponse(*e, http.StatusOK, w, r)
}
