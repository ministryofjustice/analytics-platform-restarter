package main

const (
	DecodeError = iota
	BlankHostError
	K8sError
	NoDeploymentError
	TooManyDeploymentsError
)

type Error struct {
	Type int
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
}
