package status

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const durKey = "dur_ms"

type Handler interface {
	AvailableServices(w http.ResponseWriter, r *http.Request)
	NotAvailableServices(w http.ResponseWriter, r *http.Request)
	ServicesWithResponseFasterThan(w http.ResponseWriter, r *http.Request)
	ServicesWithResponseSlowerThan(w http.ResponseWriter, r *http.Request)
}

type statusHandler struct {
	statusService Service
}

func NewStatusHandler(ticketService Service) Handler {
	return &statusHandler{
		ticketService,
	}
}

func (h *statusHandler) AvailableServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from, to, err := fromToParams(vars)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, "bad request br979823", http.StatusBadRequest)
		return
	}
	count, err := h.statusService.AvailableServices(*from, *to)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error ise77982093", http.StatusInternalServerError)
		return
	}
	countResponse(count, w)
}

func (h *statusHandler) NotAvailableServices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from, to, err := fromToParams(vars)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, "bad request br97932323", http.StatusBadRequest)
		return
	}
	count, err := h.statusService.NotAvailableServices(*from, *to)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error ise77982321", http.StatusInternalServerError)
		return
	}
	countResponse(count, w)
}

func (h *statusHandler) ServicesWithResponseFasterThan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from, to, err := fromToParams(vars)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, "bad request br97232313", http.StatusBadRequest)
		return
	}
	d, err := durParse(vars)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request br79893300", http.StatusBadRequest)
		return
	}
	count, err := h.statusService.ServicesWithResponseFasterThan(*d, *from, *to)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error ise7098091821", http.StatusInternalServerError)
		return
	}

	countResponse(count, w)
}

func (h *statusHandler) ServicesWithResponseSlowerThan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from, to, err := fromToParams(vars)
	if err != nil {
		logrus.Warn(err)
		http.Error(w, "bad request br79821029", http.StatusBadRequest)
		return
	}
	d, err := durParse(vars)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request br09813130", http.StatusBadRequest)
		return
	}
	count, err := h.statusService.ServicesWithResponseSlowerThan(*d, *from, *to)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error ise170980939", http.StatusInternalServerError)
		return
	}

	countResponse(count, w)
}

func getTs(m map[string]string, key string) (*time.Time, error) {
	tsStr, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("missing param %+v", key)
	}
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	t := time.Unix(ts, 0)
	return &t, nil
}

func countResponse(count int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf(`{"Count": %d}`, count))); err != nil {
		logrus.Error(err)
	}
}

func fromToParams(m map[string]string) (*time.Time, *time.Time, error) {
	from, err := getTs(m, "from_ts")
	if err != nil {
		logrus.Warn(err)
		return nil, nil, err
	}
	to, err := getTs(m, "to_ts")
	if err != nil {
		logrus.Warn(err)
		return nil, nil, err
	}
	return from, to, nil
}

func durParse(m map[string]string) (*time.Duration, error) {
	v, ok := m[durKey]
	if !ok {
		return nil, fmt.Errorf("missing key %+v", durKey)
	}
	ms, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	d := time.Duration(ms) * time.Millisecond
	return &d, nil
}
