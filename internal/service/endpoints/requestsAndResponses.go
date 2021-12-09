package endpoints

import (
	"OrderService/internal/database"
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)


func ParseIDFromURL(r *http.Request) (uint, error) {
	params:=mux.Vars(r)
	idStr:=params["id"]
	id,err:=strconv.ParseUint(idStr,10,32)
	if err != nil {
		return 0,err
	}
	return uint(id),nil
}

type GetByIDRequest struct {
	ID uint `json:"id,omitempty"`
}

type GetByIDResponse struct {
	OrderItems []database.OrderItemOut `json:"orderItems"`
}
type SearchRequest struct {
	Search string `json:"search"`
	StartDate time.Time `json:"startDate"`
	EndDate time.Time `json:"endDate"`
}

type SearchResponse struct {
	Orders []database.OrderOut `json:"orders"`
}
type CreateRequest struct {
	Data database.OrderIn `json:"data"`
}

type CreateResponse struct {
	Message string `json:"message"`
}
type DeleteRequest struct {
	ID uint `json:"id,omitempty"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}
type TotalRequest struct {

}

type TotalResponse struct {
	Total float32 `json:"total"`
}
type TopRequest struct {
 	Count uint `json:"count"`
}

type TopResponse struct {
	Products []database.ProductOutTop `json:"products"`
}

type QuantityOrderedRequest struct {
	ID uint `json:"id"`
}

type QuantityOrderedResponse struct {
	Quantity uint `json:"quantity"`
}


func DecodeGetByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request GetByIDRequest
	id,err:=ParseIDFromURL(r)
	if err != nil {
		return request,err
	}
	request.ID=id
	return request,nil
}
func DecodeSearchRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request SearchRequest
	request.Search=r.URL.Query().Get("search")
	startDateParam:=r.URL.Query().Get("startDate")
	endDateParam:=r.URL.Query().Get("endDate")
	if startDateParam!="" {
		request.StartDate,_=time.Parse(time.RFC3339,startDateParam)
	}
	if endDateParam!="" {
		request.EndDate,_=time.Parse(time.RFC3339,endDateParam)
	}
	return request,nil
}
func DecodeCreateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request CreateRequest
	err:=json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request,nil
}
func DecodeDeleteRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request DeleteRequest
	id,err:=ParseIDFromURL(r)
	if err != nil {
		return request,err
	}
	request.ID=id
	return request,nil
}
func DecodeTotalRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request TotalRequest
	return request,nil
}
func DecodeTopRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request TopRequest
	countParam:=r.URL.Query().Get("count")
	if  countParam=="" {
		countParam="3"
	}
	count,err:=strconv.ParseUint(countParam,10,32)
	if err != nil {
		return request,err
	}
	request.Count= uint(count)
	return request,nil
}

func DecodeQuantityOrderedRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var request QuantityOrderedRequest
	id,err:=ParseIDFromURL(r)
	if err != nil {
		return request,err
	}
	request.ID=id
	return request,nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type","application/json; charset=UTF-8")
	return json.NewEncoder(w).Encode(response)
}