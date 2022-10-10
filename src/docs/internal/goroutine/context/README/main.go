package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()

	preCtx, preCancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer preCancel()

	childCtx, childCancel := context.WithTimeout(preCtx, 300*time.Millisecond)
	defer childCancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-childCtx.Done()
		fmt.Println("childCtx is done")
	}()

	<-preCtx.Done()
	fmt.Println("preCtx is done")

	wg.Wait()
}
