package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gitlab.com/g6834/team41/analytics/internal/ports"
)

type DeclinedResponse struct {
	Count int `json:"count" example:"3"`
}

type CountDeclinedTask struct {
	Statistics ports.Statistics
}

// @Summary Get count declined task
// @Description Get count declined task
// @Produce json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /analytics/count-declined [get]
func (i CountDeclinedTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var response DeclinedResponse
	var err error
	response.Count, err = i.Statistics.GetCountDeclinedTask()
	if err != nil {
		logrus.Error(err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(response)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "{}", http.StatusInternalServerError)
		return
	}
}
