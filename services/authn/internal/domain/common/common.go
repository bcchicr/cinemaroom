package common

import (
	"fmt"

	"github.com/bcchicr/cinemaroom/services/authn/internal/domain"
	"github.com/google/uuid"
)

const (
	PasswordMinLength     = 8
	PasswordMaxByteLength = 72
)

type ValueObject interface {
	fmt.Stringer
	Equals(ValueObject) bool
	ToMap() map[string]any
}

type ID interface {
	ValueObject
	Value() uuid.UUID
}

type TypedId[T any] struct {
	value uuid.UUID
}

func (id *TypedId[T]) Value() uuid.UUID {
	return id.value
}

func (id *TypedId[T]) String() string {
	return id.value.String()
}

func (id *TypedId[T]) ToMap() map[string]any {
	return map[string]any{
		"value": id.value.String(),
	}
}

func (id *TypedId[T]) Equals(v ValueObject) bool {
	other, ok := v.(*TypedId[T])
	if !ok {
		return false
	}

	return id.Value() == other.Value()
}

func NewTypeIdFromString[T any](rawId string) (*TypedId[T], error) {
	value, err := uuid.Parse(rawId)
	if err != nil {
		return nil, domain.NewInvalidArgumentError(
			fmt.Sprintf("raw id must be uuid string, got %s", rawId),
		)
	}

	return &TypedId[T]{value: value}, nil
}
