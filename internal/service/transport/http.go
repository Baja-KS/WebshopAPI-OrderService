package transport

import (
	"context"
	"errors"
	"github.com/Baja-KS/WebshopAPI-OrderService/internal/service/endpoints"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strings"
)

func GetAuthToken(r *http.Request) (string, error) {
	authHeader := r.Header["Authorization"]
	if len(authHeader) == 0 {
		return "", errors.New("no auth header")
	}
	authHeaderParts := strings.Split(authHeader[0], " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("invalid auth header")
	}
	return authHeaderParts[1], nil
}

func AuthExtractor(ctx context.Context, r *http.Request) context.Context {
	token, err := GetAuthToken(r)
	if err != nil {
		return context.WithValue(ctx, "auth", "")
	}
	return context.WithValue(ctx, "auth", token)
}

func NewHTTPHandler(ep endpoints.EndpointSet) http.Handler {
	router := mux.NewRouter()

	GetByIDHandler := httptransport.NewServer(ep.GetByIDEndpoint, endpoints.DecodeGetByIDRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))
	SearchHandler := httptransport.NewServer(ep.SearchEndpoint, endpoints.DecodeSearchRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))
	CreateHandler := httptransport.NewServer(ep.CreateEndpoint, endpoints.DecodeCreateRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))
	DeleteHandler := httptransport.NewServer(ep.DeleteEndpoint, endpoints.DecodeDeleteRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))
	TotalHandler := httptransport.NewServer(ep.TotalEndpoint, endpoints.DecodeTotalRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))
	TopHandler := httptransport.NewServer(ep.TopEndpoint, endpoints.DecodeTopRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))
	QuantityOrderedHandler := httptransport.NewServer(ep.QuantityOrderedEndpoint, endpoints.DecodeQuantityOrderedRequest, endpoints.EncodeResponse, httptransport.ServerBefore(AuthExtractor))

	router.Handle("/GetByID/{id}", GetByIDHandler).Methods(http.MethodGet)
	router.Handle("/Search", SearchHandler).Methods(http.MethodGet)
	router.Handle("/Create", CreateHandler).Methods(http.MethodPost)
	router.Handle("/Delete/{id}", DeleteHandler).Methods(http.MethodDelete)
	router.Handle("/Total", TotalHandler).Methods(http.MethodGet)
	router.Handle("/Top", TopHandler).Methods(http.MethodGet)
	router.Handle("/QuantityOrdered/{id}", QuantityOrderedHandler).Methods(http.MethodGet)
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	return router
}
