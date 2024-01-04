package models

// swagger:model error
type BaseError struct {
	Code string `json:"code"`
}

var InternalError = BaseError{Code: "INTERNAL_ERROR"}

type (
	BadRequestErr     BaseError
	ForbiddenErr      BaseError
	NotFoundErr       BaseError
	NotAllowedErr     BaseError
	NoContentErr      BaseError
	TooManyRequestErr BaseError
	AuthorizationErr  BaseError
	ConflictErr       BaseError
)

func (e BaseError) Error() string {
	return e.Code
}

func (e BadRequestErr) Error() string {
	return e.Code
}

func (e ForbiddenErr) Error() string {
	return e.Code
}

func (e NotFoundErr) Error() string {
	return e.Code
}

func (e NotAllowedErr) Error() string {
	return e.Code
}

func (e NoContentErr) Error() string {
	return e.Code
}

func (e TooManyRequestErr) Error() string {
	return e.Code
}

func (e AuthorizationErr) Error() string {
	return e.Code
}

func (e ConflictErr) Error() string {
	return e.Code
}
