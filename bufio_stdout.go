package terminal

import (
	"context"
	"io"
	"log"
	"os"
)

type bufioMessage struct {
	data []byte
	err  error
}

func bufioStdout(ctx context.Context, onMessage func(bufioMessage)) {
	bufio := newBufio(
		whenPipe(func(pipedOutput *os.File) {
			os.Stdout = pipedOutput
			log.SetOutput(pipedOutput)
		}),
		whenRestore(func(originalOutput *os.File) {
			os.Stdout = originalOutput
			log.SetOutput(originalOutput)
			onMessage(bufioMessage{err: io.EOF})
		}),
		whenMessage(func(message []byte) {
			onMessage(bufioMessage{data: message})
		}),
	)

	// replace stdout -> buffer
	bufio.pipe(os.Stdout)

	// wait for cancel
	<-ctx.Done()

	// replace it back
	bufio.restore()
}
