package stream

func orDone(done <-chan struct{}, input <-chan interface{}) <-chan interface{} {
	outch := make(chan interface{})

	go func() {
		defer close(outch)

	loop:
		for {
			select {
			case <-done:
				break loop

			case v, ok := <-input:
				if !ok {
					break loop
				}

				select {
				case <-done:
					break loop

				case outch <- v:
					// continue
				}
			}
		}
	}()

	return outch
}
