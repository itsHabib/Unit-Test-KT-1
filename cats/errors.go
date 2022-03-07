package cats

type Error string

func (e Error) Error() string { return string(e) }

const (
	EmptyBreedID Error = "empty breed Id"
)
