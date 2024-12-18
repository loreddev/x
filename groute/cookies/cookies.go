package cookies

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"forge.capytal.company/loreddev/x/groute/router/rerrors"
)

type Marshaler interface {
	MarshalCookie() (*http.Cookie, error)
}

type Unmarshaler interface {
	UnmarshalCookie(*http.Cookie) error
}

func Marshal(v any) (*http.Cookie, error) {
	if m, ok := v.(Marshaler); ok {
		return m.MarshalCookie()
	}

	c, err := marshalValue(v)
	if err != nil {
		return c, err
	}

	if err := setCookieProps(c, v); err != nil {
		return c, err
	}

	return c, err
}

func MarshalToWriter(v any, w http.ResponseWriter) error {
	if ck, err := Marshal(v); err != nil {
		return err
	} else {
		http.SetCookie(w, ck)
	}
	return nil
}

func Unmarshal(c *http.Cookie, v any) error {
	if m, ok := v.(Unmarshaler); ok {
		return m.UnmarshalCookie(c)
	}

	value := c.Value
	b, err := base64.URLEncoding.DecodeString(value)
	if err != nil {
		return errors.Join(ErrDecodeBase64, err)
	}

	if err := json.Unmarshal(b, v); err != nil {
		return errors.Join(ErrUnmarshal, err)
	}

	return nil
}

func UnmarshalRequest(r *http.Request, v any) error {
	name, err := getCookieName(v)
	if err != nil {
		return err
	}

	c, err := r.Cookie(name)
	if errors.Is(err, http.ErrNoCookie) {
		return ErrNoCookie{name}
	} else if err != nil {
		return err
	}

	return Unmarshal(c, v)
}

func UnmarshalIfRequest(r *http.Request, v any) (bool, error) {
	if err := UnmarshalRequest(r, v); err != nil {
		if _, ok := err.(ErrNoCookie); ok {
			return false, nil
		} else {
			return true, err
		}
	} else {
		return true, nil
	}
}

func RerrUnmarshalCookie(err error) rerrors.RouteError {
	if e, ok := err.(ErrNoCookie); ok {
		return rerrors.MissingCookies([]string{e.name})
	} else {
		return rerrors.InternalError(err)
	}
}

func marshalValue(v any) (*http.Cookie, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return &http.Cookie{}, errors.Join(ErrMarshal, err)
	}

	s := base64.URLEncoding.EncodeToString(b)

	return &http.Cookie{
		Value: s,
	}, nil
}

var COOKIE_EXPIRE_VALID_FORMATS = []string{
	time.DateOnly, time.DateTime,
	time.RFC1123, time.RFC1123Z,
}

func setCookieProps(c *http.Cookie, v any) error {
	tag, err := getCookieTag(v)
	if err != nil {
		return err
	}

	c.Name, err = getCookieName(v)
	if err != nil {
		return err
	}

	tvs := strings.Split(tag, ",")

	if len(tvs) == 1 {
		return nil
	}

	tvs = tvs[1:]

	for _, tv := range tvs {
		var k, v string
		if strings.Contains(tv, "=") {
			s := strings.Split(tv, "=")
			k = s[0]
			v = s[1]
		} else {
			k = tv
			v = ""
		}

		switch k {
		case "SECURE":
			c.Name = "__Secure-" + c.Name
			c.Secure = true

		case "HOST":
			c.Name = "__Host" + c.Name
			c.Secure = true
			c.Path = "/"

		case "path":
			c.Path = v

		case "domain":
			c.Domain = v

		case "httponly":
			if v == "" {
				c.HttpOnly = true
			} else if v, err := strconv.ParseBool(v); err != nil {
				c.HttpOnly = false
			} else {
				c.HttpOnly = v
			}

		case "samesite":
			if v == "" {
				c.SameSite = http.SameSiteDefaultMode
			} else if v == "strict" {
				c.SameSite = http.SameSiteStrictMode
			} else if v == "lax" {
				c.SameSite = http.SameSiteLaxMode
			} else {
				c.SameSite = http.SameSiteNoneMode
			}
		case "secure":
			if v == "" {
				c.Secure = true
			} else if v, err := strconv.ParseBool(v); err != nil {
				c.Secure = false
			} else {
				c.Secure = v
			}

		case "max-age", "age":
			if v == "" {
				c.MaxAge = 0
			} else if v, err := strconv.Atoi(v); err != nil {
				c.MaxAge = 0
			} else {
				c.MaxAge = v
			}

		case "expires":
			if v == "" {
				c.Expires = time.Now()
			} else if v, err := timeParseMultiple(v, COOKIE_EXPIRE_VALID_FORMATS...); err != nil {
				c.Expires = time.Now()
			} else {
				c.Expires = v
			}
		}
	}

	return nil
}

func getCookieName(v any) (name string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(ErrReflectPanic, fmt.Errorf("Panic recovered: %#v", r))
		}
	}()

	tag, err := getCookieTag(v)
	if err != nil {
		return name, err
	}

	tvs := strings.Split(tag, ",")
	if len(tvs) == 0 {
		t := reflect.TypeOf(v)
		name = t.Name()
	} else {
		name = tvs[0]
	}

	if name == "" {
		return name, ErrMissingName
	}

	return name, nil
}

func getCookieTag(v any) (t string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Join(ErrReflectPanic, fmt.Errorf("Panic recovered: %#v", r))
		}
	}()

	rt := reflect.TypeOf(v)

	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		if t := ft.Tag.Get("cookie"); t != "" {
			return t, nil
		}
	}

	return "", nil
}

func timeParseMultiple(v string, formats ...string) (time.Time, error) {
	errs := []error{}
	for _, f := range formats {
		t, err := time.Parse(v, f)
		if err != nil {
			errs = append(errs, err)
		} else {
			return t, nil
		}
	}

	return time.Time{}, errs[len(errs)-1]
}

var (
	ErrDecodeBase64 = errors.New("Failed to decode base64 string from cookie value")
	ErrMarshal      = errors.New("Failed to marhal JSON value for cookie value")
	ErrUnmarshal    = errors.New("Failed to unmarshal JSON value from cookie value")
	ErrReflectPanic = errors.New("Reflect panic while trying to get tag from value")
	ErrMissingName  = errors.New("Failed to get name of cookie")
)

type ErrNoCookie struct {
	name string
}

func (e ErrNoCookie) Error() string {
	return fmt.Sprintf("Cookie \"%s\" missing from request", e.name)
}
