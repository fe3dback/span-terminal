# CLI logs and progress with spans

Library allow to track nested progress of all CLI
running tasks and write theirs logs in one _container_

Also it automatically track execution time for each span

This working on `span` concept (like in tracing)

## install

Inside main.go or other setup
```
terminal.SetGlobalTerminal(
    terminal.NewTerminal(terminal.WithWriter(os.Stdout)),
)

terminal.CaptureOutput()

// -- all code

terminal.ReleaseOutput()
```

### Spans

Between capture and release calls, we can start spans
and write some logs to it, like in tracing

__somewhere in code:__
```
func Check(ctx context.Context) (error) {
  ctx, span := terminal.StartSpan(ctx, "Some task")
  defer span.End()

  // we can write some message to logs
  span.Write(fmt.Sprintf("working on %s", path))

  go forkA()
  go forkB()
}

func forkA(ctx context.Context) (error) {
  ctx, span := terminal.StartSpan(ctx, "fork A (level 2)")
  defer span.End()

  // we can write some message to logs
  span.Write("log in A")
  
  // and/or update span progress
  span.UpdateProgress(0.4) // 40%
}
```

### Output

```
[+] Some task
| working on dasd
| working on ewqeqweqweqw                      
| log in A                      
| working on eqweqweqw                   
                                               
 > 53ms fork B
 >  40% fork A (level 2)
    23ms | span level 3
         | ..
      5s | http
     18s | system
      5s | generated-models
      1s | config-http
     98% | operations
      5s | generated-clients
   212ms | version
   314ms | config-cli
   648ms | container-cli
     10s | interaction-kafka
      5s | generated-restapi
      5s | models
    54ms | container-http
   428ms | cli
      4s | interaction
   163ms | generated-proto-models
      7s | repository
   592ms | inmem
```
