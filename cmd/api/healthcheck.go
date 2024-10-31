//cmd/api/healthcheck.go
package main

import (
	"net/http"
)

<<<<<<< HEAD
func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter,
=======
func (a *applicationDependencies) healthcheckHandler(w http.ResponseWriter,
>>>>>>> 0cc1270f48216b9318fcc1ef24b827397488e322
	r *http.Request) {
	//panic("Apples & Oranges")
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.environment,
			"version":     appVersion,
		},
	}
	err := a.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)

	}
}