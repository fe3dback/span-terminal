# CLI logs and progress with spans

Library allow to track nested progress of all CLI
running tasks and write theirs logs in one _container_

Also it automatically track execution time for each span

This based on `span`-like concept from tracing

## install

Inside main.go or other setup
```go
terminal.SetGlobalTerminal(
    terminal.NewTerminal(), // can be customized here
)

terminal.CaptureOutput()

// -- all code

terminal.ReleaseOutput()
```

### Spans

Between `CaptureOutput` and `ReleaseOutput` calls, we can start spans
and write some logs to it, like in tracing

__somewhere in code:__
```go
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

### Example of output

[![asciicast](https://asciinema.org/a/lAWXPqIZfii8p01zOpDrW76Pr.svg)](https://asciinema.org/a/lAWXPqIZfii8p01zOpDrW76Pr)
