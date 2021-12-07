package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}

	// Дополнительные типы.
	TestType struct {
		UnsupportedTypeField float64 `validate:"min:3"`
		WrongTagField        string  `validate:"min:3"`
		LessMin              int     `validate:"min:3"`
	}

	MyStruct struct {
		Password     string `validate:"password"`
		PositiveInt  int    `validate:"pos"`
		NegativeInts []int  `validate:"amount:5,10|neg"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{in: "string", expectedErr: ErrIsNotStruct},

		{
			in: TestType{
				UnsupportedTypeField: 3.14,
				WrongTagField:        "test wrong tag",
				LessMin:              2,
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "UnsupportedTypeField", Err: ErrUnsupportedType},
				ValidationError{Field: "WrongTagField", Err: ErrWrongTag},
				ValidationError{Field: "LessMin", Err: ErrLessThanMin},
			},
		},

		{in: User{
			ID:     "1234567890",
			Name:   "DIMAN",
			Age:    52,
			Email:  " string@mail.com",
			Role:   "user",
			Phones: []string{"12345678901", "0987654321", "123456789012"},
		}, expectedErr: ValidationErrors{
			ValidationError{Field: "ID", Err: ErrWrongLen},
			ValidationError{Field: "Age", Err: ErrGreaterThanMax},
			ValidationError{Field: "Email", Err: ErrNotMatchPattern},
			ValidationError{Field: "Role", Err: ErrNotInTegList},
			ValidationError{Field: "Phones", Err: ErrWrongLen},
			ValidationError{Field: "Phones", Err: ErrWrongLen},
		}},

		{in: MyStruct{
			Password:     "qwerty",
			PositiveInt:  -5,
			NegativeInts: []int{5, -15, -5, -6},
		}, expectedErr: ValidationErrors{
			ValidationError{Field: "Password", Err: ErrPasswordRequirements},
			ValidationError{Field: "PositiveInt", Err: ErrMustBePositive},
			ValidationError{Field: "NegativeInts", Err: ErrWrongAmountElem},
			ValidationError{Field: "NegativeInts", Err: ErrMustBeNegative},
		}},

		{in: App{
			Version: "1.0.5",
		}, expectedErr: nil},

		{in: Token{
			Header:    []byte("general"),
			Payload:   []byte("100"),
			Signature: []byte("cAtwa1kkEy"),
		}, expectedErr: nil},

		{in: Response{
			Code: 404,
			Body: "Not found",
		}, expectedErr: nil},

		{in: User{
			ID:     "123456789012345678901234567890123456",
			Name:   "Dmitry",
			Age:    26,
			Email:  "string@mail.com",
			Role:   "admin",
			Phones: []string{"+7912234567", "+7987654321"},
			meta:   []byte("test"),
		}, expectedErr: nil},

		{in: MyStruct{
			Password:     "Qwerty1!",
			PositiveInt:  5,
			NegativeInts: []int{-5, -15, -5, -6, -8},
		}, expectedErr: nil},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			var result bool
			if errors.As(err, &ValidationErrors{}) {
				result = err.Error() == tt.expectedErr.Error()
			} else {
				result = errors.Is(err, tt.expectedErr)
			}
			require.True(t, result, "actual err - %v", err)
			_ = tt
		})
	}
}
