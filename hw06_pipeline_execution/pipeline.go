package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		in = stage(doneControl(in, done))
	}
	return in
}

func doneControl(in In, done In) Out {
	outChan := make(Bi)

	go func() {
		defer close(outChan)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					return // больше элементов в канале нет
				}
				outChan <- v
			case <-done:
				return // остановка пайпалана по сигналу
			}
		}
	}()

	return outChan
}
