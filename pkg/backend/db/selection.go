package db

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
