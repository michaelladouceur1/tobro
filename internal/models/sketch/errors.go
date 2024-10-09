package sketch

type StepNotFoundError struct{}

func (e *StepNotFoundError) Error() string {
	return "Step not found"
}
