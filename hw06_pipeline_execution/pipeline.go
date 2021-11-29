package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func doneStage(in In, done In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for range in {
			}
		}()
		for {
			select {
			case outVal, ok := <-in:
				if !ok {
					return
				}
				select {
				case out <- outVal:
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()
	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := doneStage(in, done)

	for _, stage := range stages {
		if stage != nil {
			out = stage(doneStage(out, done))
		}
	}

	return doneStage(out, done)
}
