package httphandler

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type HTTPHandler struct {
	loan LoanHandler
	Mux  *mux.Router
}

func InitHTTPHandler(loan LoanHandler) *HTTPHandler {
	return &HTTPHandler{
		loan: loan,
		Mux:  mux.NewRouter().StrictSlash(false),
	}
}

type LoanHandler interface {
	GetUserLoansHandler(w http.ResponseWriter, r *http.Request)
	CheckUserDelinquentHandler(w http.ResponseWriter, r *http.Request)
	GetUserOutstandingHandler(w http.ResponseWriter, r *http.Request)
	PayLoanHandler(w http.ResponseWriter, r *http.Request)
	RequestLoanHandler(w http.ResponseWriter, r *http.Request)
	TriggerDelinquentCheckHandler(w http.ResponseWriter, r *http.Request)
}

// RegisterRoutes register all routes
func (handlers *HTTPHandler) RegisterRoutes() {
	var routes = handlers.Mux

	// GET Routes
	routes.Methods("GET").Path("/").HandlerFunc(indexHandler)
	routes.Methods("GET").Path("/user/loans").HandlerFunc(handlers.loan.GetUserLoansHandler)
	routes.Methods("GET").Path("/user/delinquent").HandlerFunc(handlers.loan.CheckUserDelinquentHandler)
	routes.Methods("GET").Path("/user/outstanding").HandlerFunc(handlers.loan.GetUserOutstandingHandler)

	// POST Routes
	routes.Methods("POST").Path("/request/loan").HandlerFunc(handlers.loan.RequestLoanHandler)
	routes.Methods("POST").Path("/pay/loan").HandlerFunc(handlers.loan.PayLoanHandler)
	routes.Methods("POST").Path("/internal/trigger/delinquent_check").HandlerFunc(handlers.loan.TriggerDelinquentCheckHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Service is running...")
}
