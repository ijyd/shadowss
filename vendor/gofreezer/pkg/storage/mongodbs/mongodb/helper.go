package mongodb

import (
	"fmt"
	"sort"
	"strings"

	"gofreezer/pkg/pagination"
	storage "gofreezer/pkg/storage"
	"gofreezer/pkg/storage/mongodbs/client"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
)

type operator string

const (
	equalsOperator       operator = "="
	doubleEqualsOperator operator = "=="
	inOperator           operator = "in"
	notEqualsOperator    operator = "!="
	notInOperator        operator = "notin"
	existsOperator       operator = "exists"
)

type selectorItem struct {
	key, value string
	opCode     operator
}

func labelsSelectorToCondition(labels string, condition *client.QueryMetaData) {
	glog.V(5).Infof("mongo driver list with filter(%s):labelsSelectorToCondition", labels)
}

func try(selectorPiece, op string) (lhs, rhs string, ok bool) {
	pieces := strings.Split(selectorPiece, op)
	if len(pieces) == 2 {
		keyslice := strings.Split(pieces[0], ".")
		key := keyslice[len(keyslice)-1]
		return key, pieces[1], true
	}
	return "", "", false
}

func parseFields(selector string, selectorItems *[]selectorItem) error {
	parts := strings.Split(selector, ",")
	sort.StringSlice(parts).Sort()
	for _, part := range parts {
		if part == "" {
			continue
		}
		glog.V(5).Infof("Parse filed:%v", part)
		if lhs, rhs, ok := try(part, string(notEqualsOperator)); ok {
			*selectorItems = append(*selectorItems, selectorItem{key: lhs, value: rhs, opCode: notEqualsOperator})
		} else if lhs, rhs, ok := try(part, string(doubleEqualsOperator)); ok {
			*selectorItems = append(*selectorItems, selectorItem{key: lhs, value: rhs, opCode: doubleEqualsOperator})
		} else if lhs, rhs, ok := try(part, string(equalsOperator)); ok {
			*selectorItems = append(*selectorItems, selectorItem{key: lhs, value: rhs, opCode: equalsOperator})
		} else {
			return fmt.Errorf("invalid selector: '%s'; can't understand '%s'", selector, part)
		}
	}
	return nil
}

func fieldsSelectorToCondition(fields string, selector []selectorItem, condition *client.QueryMetaData) {
	glog.V(5).Infof("mongo driver list with filter(%s):fieldsSelectorToCondition", fields)

	err := parseFields(fields, &selector)
	if err != nil {
		glog.Errorf("parse fields err:%v", err)
		return
	}

	items := selector
	if len(items) > 0 {
		condition.Condition["$and"] = []bson.M{}
		glog.V(5).Infof("Convert selector items(%+v) to condition", items)
		for _, item := range items {
			equalRegex := fmt.Sprintf("\"%v\":\"%v\"", item.key, item.value)
			equalRegexBson := bson.M{"value": bson.M{"$regex": bson.RegEx{equalRegex, ""}}}

			notEqualRegexBson := bson.M{"value": bson.M{"$not": bson.RegEx{equalRegex, ""}}}

			switch item.opCode {
			case equalsOperator:
				fallthrough
			case doubleEqualsOperator:
				condition.Condition["$and"] = append(condition.Condition["$and"].([]bson.M), equalRegexBson)
			case notEqualsOperator:
				glog.V(5).Infof("Convert selector item(%+v) to condition", item)
				condition.Condition["$and"] = append(condition.Condition["$and"].([]bson.M), notEqualRegexBson)
			default:
				glog.Warningln("invalid selector operator")
			}
		}
	}
}

func pagerToCondition(meta *client.RequestMeta, pager pagination.Pager, condition *client.QueryMetaData) {
	glog.V(5).Infof("mongo driver list with filter:pagerToCondition")

	itemSum, err := client.MongoQueryCount(meta, condition)
	if err != nil {
		glog.Errorf("Request Document Count err:%v", err)
		return
	}
	glog.V(5).Infof("Query Count is:%v", itemSum)
	//update current item sum
	pager.SetItemTotal(uint64(itemSum))

	//if there have not present page do nothing
	has, _, perPage := pager.PresentPage()
	if !has {
		return
	}

	var skip int
	hasPrev, prevPage, prevPerPage := pager.PreviousPage()
	if hasPrev {
		skip = int(prevPage * prevPerPage)
	} else {
		skip = 0
	}

	condition.Limit = int(perPage)
	condition.Skip = skip
	condition.Sort = append(condition.Sort, "lastmodifytime")
}

func Condition(meta *client.RequestMeta, condition *client.QueryMetaData, p storage.SelectionPredicate) error {
	var selector []selectorItem

	if p.Label != nil && !p.Label.Empty() {
		labelsSelectorToCondition(p.Label.String(), condition)
	}

	if p.Field != nil && !p.Field.Empty() {
		fieldsSelectorToCondition(p.Field.String(), selector, condition)
	}

	if p.Pager != nil && !p.Pager.Empty() {
		pagerToCondition(meta, p.Pager, condition)
	}
	return nil
}
