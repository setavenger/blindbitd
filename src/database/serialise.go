package database

type Serialiser interface {
	Serialise() ([]byte, error)
	DeSerialise([]byte) error
}
