package dns

const (
	dataNotFoundErrorDescription = "data not found"
	nameErrorDescription         = "name does not exist"
)

type Error struct {
	Description string
	IsTemporary bool
}

func NewDataNotFoundError() Error {
	return Error{
		Description: dataNotFoundErrorDescription,
		IsTemporary: false,
	}
}

func NewNameError() Error {
	return Error{
		Description: nameErrorDescription,
		IsTemporary: false,
	}
}
