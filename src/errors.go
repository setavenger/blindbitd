package src

type LabelAlreadyExistsErr struct{}

func (l LabelAlreadyExistsErr) Error() string {
	return "label already exists"
}
