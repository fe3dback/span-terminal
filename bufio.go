package terminal

import (
	"bufio"
	"fmt"
	stdIo "io"
	"os"
	"sync"
)

type (
	bufIO struct {
		onPipe    onPipe
		onRestore onRestore
		onMessage onMessage

		original    *os.File
		pipeReader  *os.File
		pipeWriter  *os.File
		redirecting bool

		sync.RWMutex
	}

	onPipe    = func(pipedOutput *os.File)
	onRestore = func(originalOutput *os.File)
	onMessage = func(message []byte)

	bufIOOpt = func(io *bufIO)
)

func newBufio(opts ...bufIOOpt) *bufIO {
	io := &bufIO{}

	for _, opt := range opts {
		opt(io)
	}

	return io
}

func (io *bufIO) pipe(original *os.File) {
	io.Lock()
	defer io.Unlock()

	// utils
	var err error

	// save original
	io.original = original

	// create pipe
	io.pipeReader, io.pipeWriter, err = os.Pipe()
	if err != nil {
		panic(fmt.Sprintf("failed buf io: %v", err))
	}

	// user pipe
	io.onPipe(io.pipeWriter)

	// redirect all original data to custom stream
	io.redirecting = true
	go io.streamFrom(io.pipeReader)
}

func (io *bufIO) streamFrom(r stdIo.Reader) {
	scan := bufio.NewScanner(r)

	for io.redirecting && scan.Scan() {
		io.onMessage(scan.Bytes())
	}

	// when scan stopped
	// possible here can be error, but we just ignore it
	// and no more messages come to stream, it`s ok
}

func (io *bufIO) restore() {
	io.Lock()
	defer io.Unlock()

	// stop redirecting
	io.redirecting = false

	// restore pipe
	io.onRestore(io.original)
}
