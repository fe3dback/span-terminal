package terminal

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestTerminal(t *testing.T) {
	ctx := context.Background()
	RegisterTerminal(NewTerminal(ctx, os.Stdout))

	testRoot(ctx)
	Shutdown()
}

func testRoot(ctx context.Context) {
	ctx, span := StartSpan(ctx, "root")
	defer span.End()

	span.WriteMessage("Hello, root in")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		runServiceA(ctx)
		wg.Done()
	}()
	go func() {
		runServiceB(ctx)
		wg.Done()
	}()

	wg.Wait()
}

func runServiceA(ctx context.Context) {
	ctx, span := StartSpan(ctx, "service A")
	defer span.End()

	span.WriteMessage("start A")

	forkDone := make(chan struct{})
	go func() {
		runServiceAFork(ctx)
		close(forkDone)
	}()

	for i := 0; i <= 1000; i++ {
		span.WriteMessage(fmt.Sprintf("making job %d", i))
		time.Sleep(time.Millisecond * 5)
	}

	<-forkDone

	span.WriteMessage("end A")
}

func runServiceAFork(ctx context.Context) {
	ctx, span := StartSpan(ctx, "service A (fork)")
	defer span.End()

	span.WriteMessage("start A fork")

	for i := 0; i <= 300; i++ {
		span.WriteMessage(fmt.Sprintf("making FF job %d", i))
		time.Sleep(time.Millisecond * 5)
	}

	span.WriteMessage("end A fork")
}

func runServiceB(ctx context.Context) {
	ctx, span := StartSpan(ctx, "service B")
	defer span.End()

	span.WriteMessage("b start")

	runServiceC(ctx)

	span.WriteMessage("b done")
}

func runServiceC(ctx context.Context) {
	ctx, span := StartSpan(ctx, "service C")
	defer span.End()

	span.WriteMessage("start C")

	for i := 0; i <= 2500; i++ {
		span.WriteMessage(fmt.Sprintf("making c job %d", i))
		time.Sleep(time.Millisecond * 1)
	}

	span.WriteMessage("end C")
}
