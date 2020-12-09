package Week03

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"golang.org/x/sync/errgroup"
)



func LinuxSignal(ctx context.Context) error {
	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println("Start Signal")
	for {
		select {
		case s := <-c:
			return fmt.Errorf("%v Signal end ", s)
		case <-ctx.Done():
			return fmt.Errorf("return channel")
		}
	}
}

func startServer(ctx context.Context, addr string, h http.Handler) error {
	srv := http.Server{
		Addr:    addr,
		Handler: h,
	}
	go func(ctx context.Context) {
		<-ctx.Done()
		fmt.Println("http return channel")
		err := srv.Shutdown(context.Background())
		if err != nil {
			fmt.Println(err)
		}
	}(ctx)
	fmt.Println("listening at " + addr)
	return srv.ListenAndServe()
}

type HelloHandlerStruct struct {
	content string
}

func (handler *HelloHandlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, handler.content)
}


func main() {
	ctx := context.Background()
	group, ct := errgroup.WithContext(ctx)
	group.Go(func() error {
		return LinuxSignal(ct)
	})
	group.Go(func() error {
		return startServer(ct, ":8080", &HelloHandlerStruct{})
	})
	if err := group.Wait(); err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("FINISH")
}