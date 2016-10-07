package db

import (
	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/pagination"
	"golib/pkg/storage"

	"golang.org/x/net/context"

	"github.com/golang/glog"
)

type Selection struct {
	ResultFields []string
	//plain sql like as : ("name = ? AND age >= ?", "jinzhu", "22")
	//Query := string("name = ? AND age >= ?")
	Query interface{}
	//QueryArgs := make([]interface{}, 2)
	//QueryArgs[0] = string("jinzhu")
	//QueryArgs[1] = string("22")
	QueryArgs []interface{}
}

func NewSelection(selectedFields []string, query interface{}, queryArgs []interface{}) *Selection {
	return &Selection{
		ResultFields: selectedFields,
		Query:        query,
		QueryArgs:    queryArgs,
	}
}

func (s *Selection) Field() []string {
	return s.ResultFields
}

func (s *Selection) Condition() (query interface{}, args []interface{}) {
	return s.Query, s.QueryArgs
}

func (s *Selection) Sort() interface{} {
	return nil
}
func (s *Selection) Limit() interface{} {
	return nil
}
func (s *Selection) Skip() interface{} {
	return nil
}

type PageSelection struct {
	ResultFields []string
	//plain sql like as : ("name = ? AND age >= ?", "jinzhu", "22")
	//Query := string("name = ? AND age >= ?")
	Query interface{}
	//QueryArgs := make([]interface{}, 2)
	//QueryArgs[0] = string("jinzhu")
	//QueryArgs[1] = string("22")
	QueryArgs []interface{}

	SortVal  interface{}
	LimitVal interface{}
	SkipVal  interface{}
}

func NewPageSelection(selectedFields []string, query interface{}, queryArgs []interface{}, sort, limit, skip interface{}) *PageSelection {
	return &PageSelection{
		ResultFields: selectedFields,
		Query:        query,
		QueryArgs:    queryArgs,
		SortVal:      sort,
		LimitVal:     limit,
		SkipVal:      skip,
	}
}

func (s *PageSelection) Field() []string {
	return s.ResultFields
}

func (s *PageSelection) Condition() (query interface{}, args []interface{}) {
	return s.Query, s.QueryArgs
}

func (s *PageSelection) Sort() interface{} {
	return s.SortVal
}
func (s *PageSelection) Limit() interface{} {
	return s.LimitVal
}
func (s *PageSelection) Skip() interface{} {
	return s.SkipVal
}

func buildListSelecttion(ctx context.Context, handle storage.Interface, page pagination.Pager, fileds []string) (storage.RetrieveFilter, error) {
	var selection storage.RetrieveFilter

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	if notPage {
		selection = NewSelection(fileds, nil, nil)
	} else {
		var count uint64
		err := handle.GetCount(ctx, nil, &count)
		if err != nil {
			return nil, err
		}
		glog.V(5).Infof("Got Total count %v \r\n", count)
		hasPage, perPage, skip := api.PagerToCondition(page, count)
		glog.V(5).Infof("Got page has %v  perpage %v skip %v\r\n", hasPage, perPage, skip)
		if hasPage {
			skipVal := skip
			sortVal := string("id")
			limitVal := perPage
			selection = NewPageSelection(fileds, nil, nil, sortVal, limitVal, skipVal)
		} else {
			selection = NewSelection(fileds, nil, nil)
		}
	}
	return selection, nil
}
