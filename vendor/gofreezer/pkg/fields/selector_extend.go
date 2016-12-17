package fields

import (
	"fmt"

	"gofreezer/pkg/selection"
	"gofreezer/pkg/util/sets"
)

type existTerm struct {
	field string
}

func (t *existTerm) Matches(ls Fields) bool {
	return ls.Has(t.field)
}

func (t *existTerm) Empty() bool {
	return false
}

//RequiresExactMatch check filed exist, ignore field value,
//so return value always return empty string.
func (t *existTerm) RequiresExactMatch(field string) (value string, found bool) {
	if t.field == field {
		return "", true
	}
	return "", false
}

func (t *existTerm) Transform(fn TransformFunc) (Selector, error) {
	tValue := ""
	field, _, err := fn(t.field, tValue)
	if err != nil {
		return nil, err
	}
	return &existTerm{field}, nil
}

func (t *existTerm) Requirements() Requirements {
	return []Requirement{{
		Field:    t.field,
		Operator: selection.Exists,
	}}
}

func (t *existTerm) String() string {
	return fmt.Sprintf("%v", t.field)
}

type notExistTerm struct {
	field string
}

func (t *notExistTerm) Matches(ls Fields) bool {
	return ls.Has(t.field)
}

func (t *notExistTerm) Empty() bool {
	return false
}

//RequiresExactMatch check filed exist, ignore field value,
//so return value always return empty string.
func (t *notExistTerm) RequiresExactMatch(field string) (value string, found bool) {
	if t.field == field {
		return "", true
	}
	return "", false
}

func (t *notExistTerm) Transform(fn TransformFunc) (Selector, error) {
	field, _, err := fn(t.field, "")
	if err != nil {
		return nil, err
	}
	return &notExistTerm{field}, nil
}

func (t *notExistTerm) Requirements() Requirements {
	return []Requirement{{
		Field:    t.field,
		Operator: selection.DoesNotExist,
	}}
}

func (t *notExistTerm) String() string {
	return fmt.Sprintf("%s%s", selection.DoesNotExist, t.field)
}

type inTerm struct {
	field     string
	operator  selection.Operator
	strValues sets.String
}

func (t *inTerm) Matches(ls Fields) bool {
	switch t.operator {
	case selection.In:
		if !ls.Has(t.field) {
			return false
		}
		return t.strValues.Has(ls.Get(t.field))
	case selection.NotIn:
		if !ls.Has(t.field) {
			return true
		}
		return !t.strValues.Has(ls.Get(t.field))
	default:
		return false
	}
}

func (t *inTerm) Empty() bool {
	return false
}

//RequiresExactMatch check filed exist, ignore field value,
//so return value always return empty string.
func (t *inTerm) RequiresExactMatch(field string) (value string, found bool) {
	if t.field == field {
		return "", true
	}
	return "", false
}

func (t *inTerm) Transform(fn TransformFunc) (Selector, error) {
	field, _, err := fn(t.field, "")
	if err != nil {
		return nil, err
	}
	return &notExistTerm{field}, nil
}

func (t *inTerm) Requirements() Requirements {
	return []Requirement{{
		Field:    t.field,
		Operator: selection.DoesNotExist,
	}}
}

func (t *inTerm) String() string {
	return fmt.Sprintf("%s%s", selection.DoesNotExist, t.field)
}
