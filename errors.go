package main

const (
	// DecodeError occurs when the request is not valid JSON
	DecodeError = iota
	// BlankHostError occurs when the request's host is blank
	BlankHostError
	// K8sError occurs when an API call to k8s failed
	K8sError
	// NoDeploymentError occurs when no deployment for the given host was found
	NoDeploymentError
	// TooManyDeploymentsError occurs when more than one deployment for the
	// given host were found
	TooManyDeploymentsError
)

// Error wraps another error and has a "Type" field to allow more granular
// handling
type Error struct {
	Type int
	Err  error
}

func (e Error) Error() string {
	return e.Err.Error()
}
