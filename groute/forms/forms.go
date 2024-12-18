package forms

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"forge.capytal.company/loreddev/x/groute/router/rerrors"
)

type Unmarshaler interface {
	UnmarshalForm(r *http.Request) error
}

func Unmarshal(r *http.Request, v any) (err error) {
	if u, ok := v.(Unmarshaler); ok {
		return u.UnmarshalForm(r)
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(ErrReflectPanic, fmt.Errorf("Panic recovered: %#v", r))
		}
	}()

	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	rt := rv.Type()

	for i := 0; i < rv.NumField(); i++ {
		ft := rt.Field(i)
		fv := rv.FieldByName(ft.Name)

		log.Print(ft.Name)

		if !fv.CanSet() {
			continue
		}

		// TODO: Support embedded fields
		if ft.Anonymous {
			continue
		}

		var tv string
		if t := ft.Tag.Get("form"); t != "" {
			tv = t
		} else if t = ft.Tag.Get("query"); t != "" {
			tv = t
		} else {
			tv = ft.Name
		}

		tvs := strings.Split(tv, ",")

		name := tvs[0]
		required := false
		defaultv := ""

		for _, v := range tvs {
			if v == "required" {
				required = true
			} else if strings.HasPrefix(v, "default=") {
				defaultv = strings.TrimPrefix(v, "default=")
			}
		}

		qv := r.FormValue(name)
		if qv == "" {
			if defaultv != "" {
				qv = defaultv
			} else if required {
				return &ErrMissingRequiredValue{name}
			} else {
				continue
			}
		}

		if err := setFieldValue(fv, qv); errors.Is(err, &ErrInvalidValueType{}) {
			e, _ := err.(*ErrInvalidValueType)
			e.value = name
			return e
		} else if errors.Is(err, &ErrUnsuportedValueType{}) {
			e, _ := err.(*ErrUnsuportedValueType)
			e.value = name
			return e
		} else if err != nil {
			return err
		}
	}

	return nil
}

func RerrUnsmarshal(err error) rerrors.RouteError {
	if e, ok := err.(*ErrMissingRequiredValue); ok {
		return rerrors.MissingParameters([]string{e.value})
	} else if e, ok := err.(*ErrInvalidValueType); ok {
		return rerrors.BadRequest(e.Error())
	} else {
		return rerrors.InternalError(err)
	}
}

func setFieldValue(rv reflect.Value, v string) error {
	switch rv.Kind() {

	case reflect.Pointer:
		return setFieldValue(rv.Elem(), v)

	case reflect.String:
		rv.SetString(v)

	case reflect.Bool:
		if cv, err := strconv.ParseBool(v); err != nil {
			return &ErrInvalidValueType{"bool", err, ""}
		} else {
			rv.SetBool(cv)
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if cv, err := strconv.Atoi(v); err != nil {
			return &ErrInvalidValueType{"int", err, ""}
		} else {
			rv.SetInt(int64(cv))
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if cv, err := strconv.Atoi(v); err != nil {
			return &ErrInvalidValueType{"uint", err, ""}
		} else {
			rv.SetUint(uint64(cv))
		}

	case reflect.Float32, reflect.Float64:
		if cv, err := strconv.ParseFloat(v, 64); err != nil {
			return &ErrInvalidValueType{"float64", err, ""}
		} else {
			rv.SetFloat(cv)
		}

	case reflect.Complex64, reflect.Complex128:
		if cv, err := strconv.ParseComplex(v, 128); err != nil {
			return &ErrInvalidValueType{"complex128", err, ""}
		} else {
			rv.SetComplex(cv)
		}

	// TODO: Support strucys
	// TODO: Support slices
	// TODO: Support maps
	default:
		return &ErrUnsuportedValueType{
			[]string{
				"string",
				"bool",
				"int", "int8", "int16", "int32", "int64",
				"uint", "uint8", "uint16", "uint32", "uint64",
				"float32", "float64",
				"complex64", "complex64",
			},
			"",
		}

	}

	return nil
}

type ErrInvalidValueType struct {
	expected string
	err      error
	value    string
}

func (e ErrInvalidValueType) Error() string {
	return fmt.Sprintf(
		"Value \"%s\" is a invalid type, expected type \"%s\". Got err: %s",
		e.value,
		e.expected,
		e.err.Error(),
	)
}

type ErrUnsuportedValueType struct {
	supported []string
	value     string
}

func (e ErrUnsuportedValueType) Error() string {
	return fmt.Sprintf(
		"Value \"%s\" is a unsupported type, supported types are: \"%s\"",
		e.value,
		strings.Join(e.supported, ", "),
	)
}

type ErrMissingRequiredValue struct {
	value string
}

func (e ErrMissingRequiredValue) Error() string {
	return fmt.Sprintf("Required value \"%s\" missing from query", e.value)
}

var (
	ErrParseForm    = errors.New("Failed to parse form from body or query parameters")
	ErrReflectPanic = errors.New("Reflect panic while trying to parse request")
)
