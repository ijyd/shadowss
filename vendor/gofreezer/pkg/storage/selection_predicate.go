package storage

import (
	"gofreezer/pkg/fields"
	"gofreezer/pkg/labels"
	"gofreezer/pkg/pagination"
	"gofreezer/pkg/runtime"
)

// AttrFunc returns label and field sets for List or Watch to match.
// In any failure to parse given object, it returns error.
type AttrFunc func(obj runtime.Object) (labels.Set, fields.Set, error)

// SelectionPredicate is used to represent the way to select objects from api storage.
type SelectionPredicate struct {
	Label       labels.Selector
	Field       fields.Selector
	GetAttrs    AttrFunc
	IndexFields []string

	//extension selection predicate with pager if needed
	Pager pagination.Pager

	//strange define mysqls selection at here.
	//may be abstracts backend selection
	MysqlsSelectionPredicate
}

// SelectionPredicate is used to represent the way to select objects from mysql storage.
type MysqlsSelectionPredicate struct {
	Query     interface{}
	QueryArgs interface{}
	SortField string
	LimitCon  uint64
	SkipCon   uint64
}

// Matches returns true if the given object's labels and fields (as
// returned by s.GetAttrs) match s.Label and s.Field. An error is
// returned if s.GetAttrs fails.
func (s *SelectionPredicate) Matches(obj runtime.Object) (bool, error) {
	if s.Label.Empty() && s.Field.Empty() {
		return true, nil
	}
	labels, fields, err := s.GetAttrs(obj)
	if err != nil {
		return false, err
	}
	matched := s.Label.Matches(labels)
	if s.Field != nil {
		matched = (matched && s.Field.Matches(fields))
	}
	return matched, nil
}

// MatchesSingle will return (name, true) if and only if s.Field matches on the object's
// name.
func (s *SelectionPredicate) MatchesSingle() (string, bool) {
	// TODO: should be namespace.name
	if name, ok := s.Field.RequiresExactMatch("metadata.name"); ok {
		return name, true
	}
	return "", false
}

// For any index defined by IndexFields, if a matcher can match only (a subset)
// of objects that return <value> for a given index, a pair (<index name>, <value>)
// wil be returned.
// TODO: Consider supporting also labels.
func (s *SelectionPredicate) MatcherIndex() []MatchValue {
	var result []MatchValue
	for _, field := range s.IndexFields {
		if value, ok := s.Field.RequiresExactMatch(field); ok {
			result = append(result, MatchValue{IndexName: field, Value: value})
		}
	}
	return result
}

//BuildPagerCondition use total count parse condition for list
//return value:have pager, perPagecount,skipitem
func (s *SelectionPredicate) BuildPagerCondition(count uint64) (bool, uint64, uint64) {
	page := s.Pager
	var needPager, hasPage bool
	if page == nil {
		needPager = false
	} else {
		needPager = !(page.Empty())
	}

	var perPage, skip uint64
	if needPager {
		hasPage, perPage, skip = page.Condition(count)
	}
	return hasPage, perPage, skip
}

// Filter returns skip condition
func (s *MysqlsSelectionPredicate) Skip() uint64 {
	return s.SkipCon
}

// Filter returns limit condition
func (s *MysqlsSelectionPredicate) Limit() uint64 {
	return s.LimitCon
}

// Filter returns sort condition
func (s *MysqlsSelectionPredicate) Sort() string {
	return s.SortField
}

// Filter returns select condition
func (s *MysqlsSelectionPredicate) SelectField() []string {
	//return s.IndexFields
	return []string{}
}

// Filter returns where condition
func (s *MysqlsSelectionPredicate) Where() (interface{}, interface{}) {
	return s.Query, s.QueryArgs
}
