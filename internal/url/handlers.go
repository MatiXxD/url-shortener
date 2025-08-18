package url

import "net/http"

type Handlers interface {
	Router(w http.ResponseWriter, r *http.Request)
	ReduceURL(w http.ResponseWriter, r *http.Request)
	BatchReduceURL(w http.ResponseWriter, r *http.Request)
	GetURL(w http.ResponseWriter, r *http.Request)
}
