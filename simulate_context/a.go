package  main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	TerminationSignal = syscall.SIGINT
)

type CsContext struct {
	context.Context
	cancel context.CancelFunc
	singl  chan os.Signal
	msg    string
	out    chan struct{}
	close  chan interface{}
}

func (ctx *CsContext) Done() <-chan struct{} {
	signal.Notify(ctx.singl, TerminationSignal)
	go func() {
		select {
		case <-ctx.Context.Done():
			ctx.msg = "timeout"
		case <-ctx.singl:
			ctx.msg = "context canceled"
		case <-ctx.close:
		}
		ctx.out <- struct{}{}
	}()
	return ctx.out
}
func (ctx *CsContext) Err() error {
	return errors.New(ctx.msg)
}
func WithSignal(parent context.Context, duration time.Duration) (context.Context, context.CancelFunc) {
	timeCtx, timeCancel := context.WithTimeout(parent, duration)
	cc := CsContext{
		Context: timeCtx,
		cancel:  timeCancel,
		singl:   make(chan os.Signal),
		msg:     "",
		out:     make(chan struct{}),
		close:   make(chan interface{}),
	}
	return &cc, func() {
		select {
		case cc.close <- 666:
		default:
		}
		cc.cancel()
		signal.Stop(cc.singl)
	}
}

func func1(ctx context.Context) {
	hctx, hcancel := context.WithTimeout(ctx, time.Second*4)
	defer hcancel()
	resp := make(chan struct{}, 1)
	go func() {
		time.Sleep(time.Second * 10)
		resp <- struct{}{}
	}()
	// 超时机制
	select {
	case <-hctx.Done():
		fmt.Println("ctx done, ", hctx.Err(), ", ", time.Now().Format("15:04:05"))
	case <-resp:
		fmt.Println("fun1 handle done, ", time.Now().Format("15:04:05"))
	}
	fmt.Println("func1 finish, ", time.Now().Format("15:04:05"))
	return
}

func main() {
	ctx := context.Background()
	ctx, cancel := WithSignal(ctx, 2*time.Second)
	defer cancel()
	fmt.Println("start, ", time.Now().Format("15:04:05"))
	go func1(ctx)
	time.Sleep(4 * time.Second)
	fmt.Println("main exit...")
}
