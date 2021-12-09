package middlewares

import (
	"OrderService/internal/database"
	"OrderService/internal/service"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

func CheckAuth(ctx context.Context,authServiceURL string) bool {
	client:=&http.Client{}
	req,err:=http.NewRequest("GET",authServiceURL+"/User",nil)
	if err != nil {
		return false
	}
	token:=ctx.Value("auth").(string)
	authHeader:=fmt.Sprintf("Bearer %s",token)
	req.Header.Add("Authorization",authHeader)
	res, err := client.Do(req)
	if err != nil || res.StatusCode!=200 {
		return false
	}

	return true
}


type AuthenticationMiddleware struct {
	Next service.Service
}

func (a *AuthenticationMiddleware) QuantityOrdered(ctx context.Context, id uint) (uint, error) {
	if !CheckAuth(ctx,os.Getenv("AUTH_SERVICE"))  {
		return 0, errors.New("unauthorized")
	}
	return a.Next.QuantityOrdered(ctx,id)
}

func (a *AuthenticationMiddleware) GetByID(ctx context.Context, id uint) ([]database.OrderItemOut, error) {
	if !CheckAuth(ctx,os.Getenv("AUTH_SERVICE"))  {
		return nil, errors.New("unauthorized")
	}
	return a.Next.GetByID(ctx,id)
}

func (a *AuthenticationMiddleware) Search(ctx context.Context, search string, startDate time.Time, endDate time.Time) ([]database.OrderOut, error) {
	if !CheckAuth(ctx,os.Getenv("AUTH_SERVICE"))  {
		return nil,errors.New("unauthorized")
	}
	return a.Next.Search(ctx,search,startDate,endDate)
}

func (a *AuthenticationMiddleware) Create(ctx context.Context, data database.OrderIn) (string, error) {
	return a.Next.Create(ctx,data)
}

func (a *AuthenticationMiddleware) Delete(ctx context.Context, id uint) (string, error) {
	if !CheckAuth(ctx,os.Getenv("AUTH_SERVICE"))  {
		return "Unauthorized",errors.New("unauthorized")
	}
	return a.Next.Delete(ctx,id)
}

func (a *AuthenticationMiddleware) Total(ctx context.Context) (float32, error) {
	if !CheckAuth(ctx,os.Getenv("AUTH_SERVICE"))  {
		return 0,errors.New("unauthorized")
	}
	return a.Next.Total(ctx)
}

func (a *AuthenticationMiddleware) Top(ctx context.Context, count uint) ([]database.ProductOutTop, error) {
	if !CheckAuth(ctx,os.Getenv("AUTH_SERVICE"))  {
		return nil,errors.New("unauthorized")
	}
	return a.Next.Top(ctx,count)
}
