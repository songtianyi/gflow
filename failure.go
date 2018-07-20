package gflow

type FailureFunc func(err error, step Step, context Context) error
type StepFailureFunc func(err error, context Context) error

func RetryFailure(tries int) FailureFunc {
	return func(err error, step Step, context Context) error {
		var currError error
		for tries > 0 {
			currError = step.Run(context)
			if currError == nil {
				return nil
			}
			tries--
		}
		return currError
	}
}
