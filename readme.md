# span terminal

This lib allows to write text/progress in stdout in multiple
log containers

Logs processed in tracing manner (with spans)

Safe for concurrent use

## Init

Initialize terminal logger in main.go or some setup file:

terminal require exclusive use of os.Stdout, so we need
to buffer any fmt.Print and logs, instead of immediate output

```go
// create root context
mainCtx, cancel := context.WithCancel(context.Background())
defer cancel()

// create terminal with opts
// context used for cancellation check
frontendTerminal := terminal.NewTerminal(mainCtx, os.Stdout)

// save real stdOut
defaultStdout := os.Stdout
logBufferReader, logBufferWriter, _ := os.Pipe()

// forward all logs and output to buffer
// (replace real stdout with our fake pipe)
os.Stdout = logBufferWriter
log.SetOutput(logBufferWriter)

// ------------------

// register frontend terminal
// also we start capture output right after this line
terminal.RegisterTerminal(frontendTerminal)

// do something
// this is main program call
// we can use `terminal.StartSpan` inside
err := rootCmd.ExecuteContext(mainCtx)

// end terminal session and clear screen
terminal.Shutdown()
time.Sleep(time.Millisecond * 50) // wait some time for output

// ------------------

// return real stdout
// and write all collected logs into it
_ = logBufferWriter.Close()
bufferedOutput, _ := ioutil.ReadAll(logBufferReader)
os.Stdout = defaultStdout

fmt.Printf("%s", bufferedOutput)

// handle errors and other staff
// ..
```

## Usage

```go
func (c *Service) Check(ctx context.Context) error {
  ctx, span := terminal.StartSpan(ctx, fmt.Sprintf("some operation"))
  defer span.End()

  // optional state update
  span.WriteMessage(fmt.Sprintf("checking component '%s'..", component.Name))
  span.UpdateProgress(0.4) // 40% (0 .. 1)

  // ..
  
  go forkA(ctx)
  go forkB(ctx)
  
  // ..
}

func forkA(ctx context.Context) error {
  ctx, span := terminal.StartSpan(ctx, fmt.Sprintf("forkA"))
  defer span.End()

  span.WriteMessage("fork A started")
  // ..
}
```

## Example of output

In example logs writes in 4 containers in ||:
- interaction-kafka
- models
- system
- generated-models

When span is ended, all logs in container will be collapsed

```
Processing:
[    0% ] Checking project
[  67ms ] imports check
[   42% ] deepscan (4 workers)
[  29ms ] scan project files
[   70% ] - interaction-kafka
  |  .. update.SkipErrorDecorator
  |  .. update.SkipErrorDecorator
  |  .. upsert.Processor
  |  .. bulkupsert.Processor
[ 464ms ] - cli
[    1s ] - config-http
[   12s ] - interaction
[   92% ] - models
  | in internal/models/a
  | in internal/models/b
  | in internal/models/c
  | in internal/models/d
[    2s ] - config-cli
[   63% ] - system
  | in internal/system/psql
  | in internal/system/rand
  | in internal/system/slice
  | in internal/system/storage
[    0% ] - generated-models

```
