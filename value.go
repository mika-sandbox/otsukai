package otsukai

import "errors"

const (
	VALUE_INT64 = iota
	VALUE_FLOAT64
	VALUE_STRING
	VALUE_BOOLEAN
	VALUE_HASH_OBJECT
)

var CAST_ERROR = errors.New("failed to cast value")

type IValueObject interface {
	Type() int
	ToInt64() (*int64, error)
	ToFloat64() (*float64, error)
	ToString() (*string, error)
	ToBoolean() (*bool, error)
	ToHashObject() (map[string]IValueObject, error)
}

type Int64ValueObject struct {
	val int64
}

func (i *Int64ValueObject) Type() int {
	return VALUE_INT64
}

func (i *Int64ValueObject) ToInt64() (*int64, error) {
	return &i.val, nil
}

func (i *Int64ValueObject) ToFloat64() (*float64, error) {
	return nil, CAST_ERROR
}

func (i *Int64ValueObject) ToString() (*string, error) {
	return nil, CAST_ERROR
}

func (i *Int64ValueObject) ToBoolean() (*bool, error) {
	return nil, CAST_ERROR
}

func (i *Int64ValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, CAST_ERROR
}

var _ IValueObject = (*Int64ValueObject)(nil)

//

type Float64ValueObject struct {
	val float64
}

func (i *Float64ValueObject) Type() int {
	return VALUE_FLOAT64
}

func (i *Float64ValueObject) ToInt64() (*int64, error) {
	return nil, CAST_ERROR
}

func (i *Float64ValueObject) ToFloat64() (*float64, error) {
	return &i.val, nil

}

func (i *Float64ValueObject) ToString() (*string, error) {
	return nil, CAST_ERROR
}

func (i *Float64ValueObject) ToBoolean() (*bool, error) {
	return nil, CAST_ERROR
}

func (i *Float64ValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, CAST_ERROR
}

var _ IValueObject = (*Float64ValueObject)(nil)

//

type StringValueObject struct {
	val string
}

func (i *StringValueObject) Type() int {
	return VALUE_STRING
}

func (i *StringValueObject) ToInt64() (*int64, error) {
	return nil, CAST_ERROR
}

func (i *StringValueObject) ToFloat64() (*float64, error) {
	return nil, CAST_ERROR
}

func (i *StringValueObject) ToString() (*string, error) {
	return &i.val, nil
}

func (i *StringValueObject) ToBoolean() (*bool, error) {
	return nil, CAST_ERROR
}

func (i *StringValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, CAST_ERROR
}

var _ IValueObject = (*StringValueObject)(nil)

//

type BooleanValueObject struct {
	val bool
}

func (i *BooleanValueObject) Type() int {
	return VALUE_BOOLEAN
}

func (i *BooleanValueObject) ToInt64() (*int64, error) {
	return nil, CAST_ERROR
}

func (i *BooleanValueObject) ToFloat64() (*float64, error) {
	return nil, CAST_ERROR
}

func (i *BooleanValueObject) ToString() (*string, error) {
	return nil, CAST_ERROR
}

func (i *BooleanValueObject) ToBoolean() (*bool, error) {
	return &i.val, nil
}

func (i *BooleanValueObject) ToHashObject() (map[string]IValueObject, error) {
	return nil, CAST_ERROR
}

var _ IValueObject = (*BooleanValueObject)(nil)

//

type HashValueObject struct {
	val map[string]IValueObject
}

func (i *HashValueObject) Type() int {
	return VALUE_HASH_OBJECT
}

func (i *HashValueObject) ToInt64() (*int64, error) {
	return nil, CAST_ERROR
}

func (i *HashValueObject) ToFloat64() (*float64, error) {
	return nil, CAST_ERROR
}

func (i *HashValueObject) ToString() (*string, error) {
	return nil, CAST_ERROR
}

func (i *HashValueObject) ToBoolean() (*bool, error) {
	return nil, CAST_ERROR
}

func (i *HashValueObject) ToHashObject() (map[string]IValueObject, error) {
	return i.val, nil
}

var _ IValueObject = (*HashValueObject)(nil)
