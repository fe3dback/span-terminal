package terminal

func whenPipe(fn onPipe) bufIOOpt {
	return func(io *bufIO) {
		io.onPipe = fn
	}
}

func whenRestore(fn onRestore) bufIOOpt {
	return func(io *bufIO) {
		io.onRestore = fn
	}
}

func whenMessage(fn onMessage) bufIOOpt {
	return func(io *bufIO) {
		io.onMessage = fn
	}
}
