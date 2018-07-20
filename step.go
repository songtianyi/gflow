package gflow

type Step interface {
	Label() string
	OnFailure(error, Context) error
	Run(context Context) error
}
