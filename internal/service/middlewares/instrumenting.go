package middlewares

import (
	"OrderService/internal/database"
	"OrderService/internal/service"
	"context"
	"fmt"
	"github.com/go-kit/kit/metrics"
	"strconv"
	"time"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           service.Service
}

func (i *InstrumentingMiddleware) QuantityOrdered(ctx context.Context, id uint) (qty uint,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","QuantityOrdered","order_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	qty,err=i.Next.QuantityOrdered(ctx,id)
	return
}

func (i *InstrumentingMiddleware) GetByID(ctx context.Context, id uint) (order []database.OrderItemOut, err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","GetByID","order_id", strconv.Itoa(int(id)),"error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	order,err=i.Next.GetByID(ctx,id)
	return
}

func (i *InstrumentingMiddleware) Search(ctx context.Context, search string, startDate time.Time, endDate time.Time) (orders []database.OrderOut,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Search","order_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	orders,err=i.Next.Search(ctx,search,startDate,endDate)
	return
}

func (i *InstrumentingMiddleware) Create(ctx context.Context, data database.OrderIn) (msg string,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Create","order_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	msg,err=i.Next.Create(ctx,data)
	return
}

func (i *InstrumentingMiddleware) Delete(ctx context.Context, id uint) (msg string,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Delete","order_id", strconv.Itoa(int(id)),"error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	msg,err=i.Next.Delete(ctx,id)
	return
}

func (i *InstrumentingMiddleware) Total(ctx context.Context) (total float32,err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Total","order_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	total,err=i.Next.Total(ctx)
	return
}

func (i *InstrumentingMiddleware) Top(ctx context.Context, count uint) (top []database.ProductOutTop, err error) {
	defer func(begin time.Time) {
		lvs:=[]string{"method","Top","order_id", "none","error",fmt.Sprint(err!=nil)}
		i.RequestCount.With(lvs...).Add(1)
		i.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	top,err=i.Next.Top(ctx,count)
	return
}
