package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = doStageWork(in, done, stage)
	}
	return in
}

func doStageWork(in In, done In, stage Stage) Out {
	outCh := make(Bi)
	go func() {
		defer close(outCh)
		stOut := stage(in)

		for {
			select {
			case <-done:
				return
			case value, ok := <-stOut:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case outCh <- value:
				}
			}
		}
	}()
	return outCh
}
