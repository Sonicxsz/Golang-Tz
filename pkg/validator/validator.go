package validator

import (
	"fmt"
	"github.com/google/uuid"
	"unicode/utf8"
)

type Validator struct {
	errors []string
	count  int
}

// Проверяет есть ли ошибки
func (v *Validator) HasErrors() bool {
	return len(v.errors) > 0
}

// Возвращает количество провалидированных полей
func (v *Validator) ValidatedFieldsCount() int {
	return v.count
}

func (v *Validator) AddError(msg string) {
	v.errors = append(v.errors, msg)
}

func (v *Validator) GetErrors() []string {
	return v.errors
}

type StringValidator struct {
	value     string
	validator *Validator
	name      string
}

type NumberValidator struct {
	value     any
	validator *Validator
	name      string
}

func New() *Validator {
	return &Validator{
		errors: make([]string, 0),
		count:  0,
	}
}

// Создаем структуру с методами для проверки строк
func (v *Validator) CheckString(value string, name string) *StringValidator {
	v.count += 1
	return &StringValidator{
		value:     value,
		validator: v,
		name:      name,
	}
}

// // Создаем структуру с методами для проверки чисел
func (v *Validator) CheckNumber(value any, name string) *NumberValidator {
	v.count += 1
	return &NumberValidator{
		value:     value,
		validator: v,
		name:      name,
	}
}

// Проверяет что длина строки не больше указанного
func (v *StringValidator) IsMax(max int) *StringValidator {
	length := utf8.RuneCountInString(v.value)
	if length > max {
		v.validator.AddError(fmt.Sprintf("[%s] - Max aviable length is %d, Provided: %d", v.name, max, length))
	}
	return v
}

// Проверяет что длина строки не больше указанного
func (v *StringValidator) IsUuid() *StringValidator {
	_, err := uuid.Parse(v.value)
	if err != nil {
		v.validator.AddError(fmt.Sprintf("[%s] - Invalid uuid", v.name))
	}
	return v
}

// Проверяет что длина строки не меньше указанного
func (v *StringValidator) IsMin(min int) *StringValidator {
	length := utf8.RuneCountInString(v.value)
	if length < min {
		v.validator.AddError(fmt.Sprintf("[%s] - Min required length is %d, Provided: %d", v.name, min, length))
	}
	return v
}

// Проверяет что число не меньше указанного
func (v *NumberValidator) IsMin(min float64) *NumberValidator {
	value, ok := v.toInt64()

	if !ok {
		v.validator.AddError(fmt.Sprintf("[%s] - Unsupported type: %T", v.name, v.value))
		return v
	}
	if value < min {
		v.validator.AddError(fmt.Sprintf("[%s] - Min required: %g, Provided: %g", v.name, min, value))

	}
	return v
}

// Проверяет что число не больше указанного
func (v *NumberValidator) IsMax(max float64) *NumberValidator {
	value, ok := v.toInt64()

	if !ok {
		v.validator.AddError(fmt.Sprintf("[%s] - Unsupported type: %T", v.name, v.value))
		return v
	}

	if value > max {
		v.validator.AddError(fmt.Sprintf("[%s] - Max aviable: %g, Provided: %g", v.name, max, value))
	}
	return v
}

// Преобразует любой числовой тип в int64 для единого сравнения
func (v *NumberValidator) toInt64() (float64, bool) {
	switch val := v.value.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		// Проверяем переполнение
		if val > 9223372036854775807 { // math.MaxInt64
			return 0, false
		}
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}
