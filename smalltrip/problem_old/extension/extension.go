package extension

type Errors struct {
	Errors []Error `json:"errors,omitempty"  xml:"errors,omitempty"`
}

func NewErrorTree(err error) Errors {
	errs := []Error{}
	if us, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range us.Unwrap() {
			errs = append(errs, Error{
				Detail: e.Error(),
				Errors: NewErrorTree(e).Errors,
				error:  e,
			})
		}
	} else if u, ok := err.(interface{ Unwrap() error }); ok {
		e := u.Unwrap()
		errs = append(errs, Error{
			Detail: e.Error(),
			Errors: NewErrorTree(e).Errors,
			error:  err,
		})
	}
	return Errors{Errors: errs}
}

type Error struct {
	Detail  string  `json:"detail,omitempty"  xml:"detail,omitempty"`
	Pointer string  `json:"pointer,omitempty" xml:"pointer,omitempty"`
	Errors  []Error `json:"errors,omitempty"  xml:"errors,omitempty"`
	error   error   `json:"-"                 xml:"-"`
}

func (e Error) Error() string {
	return e.error.Error()
}
