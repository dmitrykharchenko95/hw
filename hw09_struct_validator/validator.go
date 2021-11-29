package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrIsNotStruct     = errors.New("interface is not a struct") //+
	ErrUnsupportedType = errors.New("unsupported type of field")
	ErrWrongTag        = errors.New("wrong tag")

	// Possible string validation errors.
	ErrWrongLen             = errors.New("field length does not match the tag")
	ErrNotInTegList         = errors.New("field value is not in the tag list")
	ErrNotMatchPattern      = errors.New("field value does not match the pattern")
	ErrPasswordRequirements = errors.New("field value must contain a symbol, number, uppercase and lowercase letters ")

	// Possible int validation errors.
	ErrLessThanMin    = errors.New("field value is less than min")
	ErrGreaterThanMax = errors.New("field value is greater than max")
	ErrMustBePositive = errors.New("field value is not positive")
	ErrMustBeNegative = errors.New("field value is not negative")

	// Possible slice validation errors.
	ErrWrongAmountElem = errors.New("fields value has wrong amount of elements")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var s string
	for _, err := range v {
		s += fmt.Sprintf("%v:%v\n", err.Field, err.Err)
	}
	return s
}

func Validate(v interface{}) error {
	var VE ValidationErrors
	var subVE ValidationErrors

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Struct {
		return ErrIsNotStruct
	}

	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := rv.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		var err error
		switch fv.Kind() { //nolint:exhaustive
		case reflect.String:
			err = stringValidate(field.Name, fv, tag)
		case reflect.Int:
			err = intValidate(field.Name, fv, tag)
		case reflect.Slice:
			err = sliceValidate(fv.Index(0).Kind(), field.Name, fv, tag)
		default:
			err = ValidationErrors{ValidationError{Field: field.Name, Err: ErrUnsupportedType}}
		}
		if errors.As(err, &subVE) {
			VE = append(VE, subVE...)
		} else {
			return err
		}
	}
	if len(VE) == 0 {
		return nil
	}
	return VE
}

func stringValidate(fn string, fv reflect.Value, tag string) error {
	var SVE ValidationErrors

	tags := strings.Split(tag, "|")
mainLoop:
	for _, t := range tags {
		args := strings.Split(t, ":")
		switch args[0] {
		case "len":
			tagLen, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			if len(fv.String()) != tagLen {
				SVE = append(SVE, ValidationError{Field: fn, Err: ErrWrongLen})
				continue
			}

		case "in":
			ptr := strings.Split(args[1], ",")
			for _, s := range ptr {
				if fv.String() == s {
					continue mainLoop
				}
			}
			SVE = append(SVE, ValidationError{Field: fn, Err: ErrNotInTegList})
			continue

		case "regexp":
			ok, err := regexp.MatchString(args[1], fv.String())
			if err != nil {
				return err
			}
			if !ok {
				SVE = append(SVE, ValidationError{Field: fn, Err: ErrNotMatchPattern})
				continue
			}

		// password validator checks that the string contains a character, number, lowercase letter, and uppercase letter
		case "password":
			if stringPasswordValidate(fv.String()) {
				SVE = append(SVE, ValidationError{Field: fn, Err: ErrPasswordRequirements})
				continue
			}

		case "":
			continue

		default:
			SVE = append(SVE, ValidationError{Field: fn, Err: ErrWrongTag})
			continue
		}
	}
	return SVE
}

func stringPasswordValidate(fv string) bool {
	var A, a, i, s bool
	for _, r := range fv {
		switch {
		case unicode.IsDigit(r):
			i = true
		case unicode.IsUpper(r):
			A = true
		case unicode.IsLower(r):
			a = true
		case !unicode.IsDigit(r) && !unicode.IsLetter(r):
			s = true
		}
	}
	return !A || !a || !i || !s
}

func intValidate(fn string, fv reflect.Value, tag string) error {
	var IVE ValidationErrors

	tags := strings.Split(tag, "|")
mainLoop:
	for _, t := range tags {
		args := strings.Split(t, ":")
		switch args[0] {
		case "min":
			min, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			if int(fv.Int()) < min {
				IVE = append(IVE, ValidationError{Field: fn, Err: ErrLessThanMin})
				continue
			}

		case "max":
			max, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			if int(fv.Int()) > max {
				IVE = append(IVE, ValidationError{Field: fn, Err: ErrGreaterThanMax})
				continue
			}

		case "in":
			r := strings.Split(args[1], ",")
			for _, n := range r {
				if fmt.Sprint(fv.Int()) == n {
					continue mainLoop
				}
			}
			IVE = append(IVE, ValidationError{Field: fn, Err: ErrNotInTegList})

		// pos validator checks that the value is positive
		case "pos":
			if fv.Int() < 0 {
				IVE = append(IVE, ValidationError{Field: fn, Err: ErrMustBePositive})
				continue
			}

		// neg validator checks that the value is positive
		case "neg":
			if fv.Int() >= 0 {
				IVE = append(IVE, ValidationError{Field: fn, Err: ErrMustBeNegative})
				continue
			}
		case "":
			continue
		default:
			return ErrWrongTag
		}
	}
	return IVE
}

func sliceValidate(k reflect.Kind, fn string, fv reflect.Value, tag string) error {
	var (
		validateFunc func(fn string, fv reflect.Value, tag string) error
		err          error
		SSVE         ValidationErrors
		subSSVE      ValidationErrors
	)

	// amount:min,max validator checks that the number of elements in the slice is in the range [min, max]
	tags := strings.Split(tag, "|")
	for _, t := range tags {
		args := strings.Split(t, ":")
		if args[0] == "amount" {
			tag = strings.ReplaceAll(tag, t, "")
			n := strings.Split(args[1], ",")
			min, err := strconv.Atoi(n[0])
			if err != nil {
				return err
			}
			max, err := strconv.Atoi(n[1])
			if err != nil {
				return err
			}

			if fv.Len() < min || fv.Len() > max {
				SSVE = append(SSVE, ValidationError{Field: fn, Err: ErrWrongAmountElem})
				break
			}
		}
	}

	switch k { //nolint:exhaustive
	case reflect.String:
		validateFunc = stringValidate
	case reflect.Int:
		validateFunc = intValidate
	default:
		return ErrUnsupportedType
	}

	for i := 0; i < fv.Len(); i++ {
		err = validateFunc(fn, fv.Index(i), tag)
		if errors.As(err, &subSSVE) {
			SSVE = append(SSVE, subSSVE...)
		} else {
			return err
		}
	}

	return SSVE
}
