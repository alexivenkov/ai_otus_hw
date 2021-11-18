package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		if stage != nil {
			in = proxyStage(done, stage(in))
		}
	}

	return in
}

func proxyStage(done, in In) Out {
	doneCh := make(Bi)

	go func() {
		defer close(doneCh)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok { // handle empty interface type issue
					return
				}
				doneCh <- v
			}
		}
	}()

	return doneCh
}
