package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"
	"gofreezer/pkg/selection"
	"gofreezer/pkg/storage"
	storagehelper "gofreezer/pkg/storage/helper"
	"gofreezer/pkg/util/cache"

	"github.com/golang/glog"
	"github.com/jinzhu/gorm"
)

type TableTag struct {
	column      string   //it is a column name in db
	tableName   string   //a table name
	keyword     []string //the keyword for sql like as resourceKey unique and so on
	structField string
}

type Table struct {
	name string
	obj  reflect.Value

	//json tag value as key,valus is TableTag: it contains column or talbe or  other keyword for sql
	freezerTag map[string]TableTag

	//freezer column as key,the key in struct freezerTag  as value
	columnToFreezerTagKey map[string]string

	//resourcekey hold a column name,this key as a resouce name in restful url
	resoucekey string
}

const (
	freezerTag = "freezer"
	jsonTag    = "json"

	//StructTagKey struct tag key for mysql
	tagColumn = "column"
	//tagTableKey table name in tag
	tagTable = "table"
	//primary_key tag
	tagResourceKey = "resoucekey"
)

//example struct
/*
type DBRoot struct {
	embedded DBResource `freezer:"table:dbresource"`
}

type DBResource struct {
	//use column as extend, because of use gorm, so append gorm tag(gorm and sql tag)
	Name string `freezer:"column:name;resoucekey" gorm:"column:name" sql:"type:varchar(100);unique"`
}

you must give out a resource key for rest requst
*/

const (
	maxTableCache int = 128
)

var tableCache *cache.LRUExpireCache

func init() {
	tableCache = cache.NewLRUExpireCache(maxTableCache)
}

func FindTableTag(typ reflect.Type, index int, t *Table) bool {
	field := typ.Field(index)

	value, ok := field.Tag.Lookup(freezerTag)
	findTable := false

	if ok {
		tagMap := parseTag(value)

		var jsonKey string
		if value, ok := field.Tag.Lookup(jsonTag); ok {
			jsonKey = stripJsonTagValue(value)
			tagMap.structField = field.Name
			t.freezerTag[jsonKey] = tagMap
		}

		if len(tagMap.tableName) != 0 {
			t.columnToFreezerTagKey[tagTable] = jsonKey
			t.name = tagMap.tableName
			findTable = true
		}

		if len(tagMap.column) != 0 {
			t.columnToFreezerTagKey[tagMap.column] = jsonKey
		}

		for _, v := range tagMap.keyword {
			if strings.Compare(tagResourceKey, v) == 0 {
				t.resoucekey = tagMap.column
			}
		}

		// var column string
		// for k, v := range tagMap {
		// 	glog.V(9).Infof("in %s range tag map %v:%v", jsonKey, k, v)
		// 	switch k {
		// 	case tagTable:
		// 		t.columnToFreezerTagKey[k] = field.Name
		// 		t.name = v
		// 		findTable = true
		// 	case tagColumn:
		// 		column = v
		// 		t.columnToFreezerTagKey[v] = field.Name
		// 	case tagResourceKey:
		// 		//we used column value in  this tag line
		// 		t.resoucekey = column
		// 		glog.V(9).Infof("in %s find resource key %s", jsonKey, t.resoucekey)
		// 	}

	}

	return findTable
}

//BuildTable search tag in obj
//return the reflect.value of tag
//return error if has a error
func BuildTable(obj reflect.Value, t *Table) error {

	vType := obj.Type()
	for i := 0; i < vType.NumField(); i++ {
		embV := obj.Field(i)

		if FindTableTag(vType, i, t) {
			t.obj = reflect.Indirect(reflect.New(embV.Type()))
			t.obj.Set(embV)
		}

		switch embV.Kind() {
		case reflect.Struct:
			if err := BuildTable(embV, t); err != nil {
				return err
			}
		}
	}

	return nil
}

func stripJsonTagValue(origin string) string {
	vals := strings.Split(origin, ",")
	return vals[0]
}

func parseTag(origin string) TableTag {
	vals := strings.Split(origin, ";")

	tag := TableTag{}
	for _, query := range vals {
		key := query
		if i := strings.IndexAny(key, ":"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else if len(key) != 0 {
			singleTag := key[:]
			key, query = singleTag, singleTag
		} else {
			key, query = "", ""
		}

		switch key {
		case tagTable:
			tag.tableName = query
		case tagColumn:
			tag.column = query
		case tagResourceKey:
			//we used column value in  this tag line
			tag.keyword = append(tag.keyword, query)
		}
	}

	return tag
}

func GetTable(ctx context.Context, obj runtime.Object) (*Table, error) {
	v, err := conversion.EnforcePtr(obj)
	if err != nil {
		return nil, err
	}

	kind := storagehelper.GetObjKind(obj)
	cacheVal, found := tableCache.Get(kind)
	if found {
		table := cacheVal.(*Table)
		return table, nil
	}

	table := &Table{
		freezerTag:            make(map[string]TableTag),
		columnToFreezerTagKey: make(map[string]string),
	}

	err = BuildTable(v, table)
	if err != nil {
		return nil, err
	}

	if len(table.name) == 0 {
		return nil, fmt.Errorf("not find tag('table') in struct")
	}

	if len(table.resoucekey) == 0 {
		return nil, fmt.Errorf("not find resource key in struct")
	}

	glog.V(5).Infof("find table %+v", table)
	tableCache.Add(kind, table, 24*time.Hour)

	WithTable(ctx, table)

	return table, err
}

type AfterFindTable func(tableObj reflect.Value) error

func (t *Table) ExtractTableObj(obj runtime.Object, afterFunc AfterFindTable) error {
	//tObj := reflect.Value{}
	//var tObj uintptr
	v, err := conversion.EnforcePtr(obj)
	if err != nil {
		return err
	}

	if !v.CanInterface() {
		return fmt.Errorf("object(%v) cannt interface by reflect", v.Type())
	}

	//find in struct root
	tableFiledName := t.freezerTag[t.columnToFreezerTagKey[tagTable]].structField
	val := v.FieldByName(tableFiledName)
	if val.IsValid() {
		return afterFunc(val)
	}

	vType := v.Type()
	for i := 0; i < vType.NumField(); i++ {
		embV := v.Field(i)

		glog.V(9).Infof("search table %v table:%v", embV.Kind(), tableFiledName)
		switch embV.Kind() {
		case reflect.Struct:
			val := embV.FieldByName(tableFiledName)
			if val.IsValid() {
				return afterFunc(val)
			}
			vType = embV.Type()
		}
	}

	return fmt.Errorf("runtime object(%v) not found table", v.Type())
}

func (t *Table) SetTable(obj runtime.Object, table reflect.Value) error {
	err := t.ExtractTableObj(obj, func(tObj reflect.Value) error {
		if !tObj.CanSet() {
			return fmt.Errorf("object(%v) cannt set by reflect", tObj.Type())
		}
		tObj.Set(table)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (t *Table) CovertRowsToObject(row *RowResult, obj runtime.Object, table reflect.Value) error {
	err := t.SetTable(obj, table)
	if err != nil {
		return err
	}

	row.data, err = json.Marshal(obj)
	if err != nil {
		return err
	}

	return nil
}

func (t *Table) GetColumnByField(filed string) (column string) {
	i := strings.LastIndexAny(filed, ".")
	if i >= 0 {
		column = filed[i+1:]
	} else {
		column = filed
	}

	//check is a invalid column
	_, ok := t.freezerTag[column]
	if !ok {
		column = ""
		glog.Warningf("field %v is not a valid field", filed)
	} else {
		tag := t.freezerTag[column]
		column = tag.column
	}

	glog.V(9).Infof("Get fileds %v", column)

	return
}

func (t *Table) Fields(dbHandle *gorm.DB, p storage.SelectionPredicate) *gorm.DB {
	if p.Field.Empty() {
		return dbHandle
	}

	fieldsCondition := p.Field.Requirements()
	for _, v := range fieldsCondition {
		column := t.GetColumnByField(v.Field)
		switch v.Operator {
		case selection.Equals:
			fallthrough
		case selection.DoubleEquals:
			query := fmt.Sprintf("%s = ?", column)
			queryArgs := v.Value
			dbHandle = dbHandle.Where(query, queryArgs)
		case selection.NotEquals:
			query := fmt.Sprintf("%s != ?", column)
			queryArgs := v.Value
			dbHandle = dbHandle.Where(query, queryArgs)
		}
	}

	return dbHandle
}

func (t *Table) BaseCondition(dbHandle *gorm.DB, p storage.SelectionPredicate) *gorm.DB {
	dbHandle = t.Fields(dbHandle, p)
	// query, args := p.Where()
	// if query != nil && args != nil {
	// 	dbHandle = dbHandle.Where(query, args)
	// }

	// selectField := p.SelectField()
	// if len(selectField) != 0 {
	// 	dbHandle = dbHandle.Select(selectField)
	// }

	//always use resoucekey as sort field
	dbHandle = dbHandle.Order(t.resoucekey)

	return dbHandle
}

func (t *Table) PageCondition(dbHandle *gorm.DB, p storage.SelectionPredicate, totalCount uint64) *gorm.DB {

	hasPage, perPage, skip := p.BuildPagerCondition(uint64(totalCount))
	if hasPage {
		limitVal := perPage
		if limitVal != 0 {
			dbHandle = dbHandle.Limit(int(limitVal))
		}

		skipVal := skip
		if skipVal != 0 {
			dbHandle = dbHandle.Offset(int(skipVal))
		}
	}

	return dbHandle
}
