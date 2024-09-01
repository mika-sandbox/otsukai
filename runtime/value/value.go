package value

import (
	"github.com/mika-sandbox/otsukai/parser"
	re "github.com/mika-sandbox/otsukai/runtime/errors"
	"strconv"
)

const (
	VALUE_INT64 = iota
	VALUE_FLOAT64
	VALUE_STRING
	VALUE_BOOLEAN
	VALUE_HASH_OBJECT
)

type IValueObject interface {
	Type() int
	ToInt64() (*int64, error)
	ToFloat64() (*float64, error)
	ToString() (*string, error)
	ToBoolean() (*bool, error)
	ToHashObject() (map[string]IValueObject, error)
}

type Int64ValueObject struct {
	Val int64
}

func (i Int64ValueObject) Type() int {
	return VALUE_INT64
}

func (i Int64ValueObject) ToInt64() (*int64, error) {
	return &i.Val, nil
}

func (i Int64ValueObject) ToFloat64() (*float64, error) {
	return nil, re.CAST_ERROR
}

func (i Int64ValueObject) ToString() (*string, error) {
	return nil, re.CAST_ERROR
}

func (i Int64ValueObject) ToBoolean() (*bool, error) {
	return nil, re.CAST_ERROR
}

func (i Int64ValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, re.CAST_ERROR
}

var _ IValueObject = (*Int64ValueObject)(nil)

//

type Float64ValueObject struct {
	Val float64
}

func (i Float64ValueObject) Type() int {
	return VALUE_FLOAT64
}

func (i Float64ValueObject) ToInt64() (*int64, error) {
	return nil, re.CAST_ERROR
}

func (i Float64ValueObject) ToFloat64() (*float64, error) {
	return &i.Val, nil

}

func (i Float64ValueObject) ToString() (*string, error) {
	return nil, re.CAST_ERROR
}

func (i Float64ValueObject) ToBoolean() (*bool, error) {
	return nil, re.CAST_ERROR
}

func (i Float64ValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, re.CAST_ERROR
}

var _ IValueObject = (*Float64ValueObject)(nil)

//

type StringValueObject struct {
	Val string
}

func (i StringValueObject) Type() int {
	return VALUE_STRING
}

func (i StringValueObject) ToInt64() (*int64, error) {
	return nil, re.CAST_ERROR
}

func (i StringValueObject) ToFloat64() (*float64, error) {
	return nil, re.CAST_ERROR
}

func (i StringValueObject) ToString() (*string, error) {
	return &i.Val, nil
}

func (i StringValueObject) ToBoolean() (*bool, error) {
	return nil, re.CAST_ERROR
}

func (i StringValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, re.CAST_ERROR
}

var _ IValueObject = (*StringValueObject)(nil)

//

type BooleanValueObject struct {
	Val bool
}

func (i BooleanValueObject) Type() int {
	return VALUE_BOOLEAN
}

func (i BooleanValueObject) ToInt64() (*int64, error) {
	return nil, re.CAST_ERROR
}

func (i BooleanValueObject) ToFloat64() (*float64, error) {
	return nil, re.CAST_ERROR
}

func (i BooleanValueObject) ToString() (*string, error) {
	return nil, re.CAST_ERROR
}

func (i BooleanValueObject) ToBoolean() (*bool, error) {
	return &i.Val, nil
}

func (i BooleanValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, re.CAST_ERROR
}

var _ IValueObject = (*BooleanValueObject)(nil)

//

type HashValueObject struct {
	Val map[string]IValueObject
}

func (i HashValueObject) Type() int {
	return VALUE_HASH_OBJECT
}

func (i HashValueObject) ToInt64() (*int64, error) {
	return nil, re.CAST_ERROR
}

func (i HashValueObject) ToFloat64() (*float64, error) {
	return nil, re.CAST_ERROR
}

func (i HashValueObject) ToString() (*string, error) {
	return nil, re.CAST_ERROR
}

func (i HashValueObject) ToBoolean() (*bool, error) {
	return nil, re.CAST_ERROR
}

func (i HashValueObject) ToHashObject() (map[string]IValueObject, error) {
	return i.Val, nil
}

var _ IValueObject = (*HashValueObject)(nil)

func ToValueObject(v parser.Value) (IValueObject, error) {
	if v.HashSymbol != nil {
		return &StringValueObject{Val: v.HashSymbol.Identifier}, nil
	}

	if v.Hash != nil {
		items := map[string]IValueObject{}

		for _, pair := range v.Hash.Pairs {
			items[pair.Identifier.Identifier], _ = ToValueObject(pair.Value)
		}

		return &HashValueObject{Val: items}, nil
	}

	if v.Literal != nil {
		if v.Literal.String != nil {
			return &StringValueObject{Val: *v.Literal.String}, nil
		}

		if v.Literal.Number != nil {
			float, err := strconv.ParseFloat(*v.Literal.Number, 64)
			if err != nil {
				return nil, re.SYNTAX_ERROR
			}

			return &Float64ValueObject{Val: float}, nil
		}

		if v.Literal.True != nil {
			return &BooleanValueObject{Val: true}, nil
		}

		if v.Literal.False != nil {
			return &BooleanValueObject{Val: false}, nil
		}

		if v.Literal.Null != nil {
			return nil, nil
		}

	}

	return nil, re.RUNTIME_ERROR
}
