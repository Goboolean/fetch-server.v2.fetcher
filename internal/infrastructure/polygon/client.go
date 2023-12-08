package polygon

import (
	"context"
	"sync"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/websocket"
	"github.com/polygon-io/client-go/websocket/models"
)



type client[T models.EquityTrade | models.CryptoTrade] struct {
	conn  *polygonws.Client

	ch chan T

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}


func (c *client[T]) Ping(ctx context.Context) error {

	ch := make(chan error)

	go func(ch chan error) {
		ch <- c.conn.Connect()
	}(ch)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

func (c *client[T]) Close() {
	c.cancel()
	c.conn.Close()
	c.wg.Wait()
	close(c.ch)
}

func newClient[T models.EquityTrade | models.CryptoTrade](c *resolver.ConfigMap) (*client[T], error) {

	key, err := c.GetStringKey("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	_feed, err := c.GetStringKey("FEED")
	if err != nil {
		return nil, err
	}
	feed := polygonws.Feed(_feed)

	_market, err := c.GetStringKey("MARKET")
	if err != nil {
		return nil, err
	}
	market := polygonws.Market(_market)

	conn, err := polygonws.New(polygonws.Config{
		APIKey: key,
		Feed:   feed,
		Market: market,
	})

	buf, err := c.GetIntKey("BUFFER_SIZE")
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &client[T]{
		conn: conn,
		ch: make(chan T, buf),

		ctx: ctx,
		cancel: cancel,
	}, nil
}

func (c *client[T]) Subscribe() (<-chan T, error) {

	if err := c.conn.Subscribe(polygonws.StocksTrades); err != nil {
		return nil, err
	}

	c.wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.conn.Error():
				return
			case out, more := <-c.conn.Output():
				if !more {
					return
				}

				data, ok := out.(T)
				if !ok {
					return
				}

				c.ch <- data
			}
		}
	}(c.ctx, &c.wg)

	return c.ch, nil
}