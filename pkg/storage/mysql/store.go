package mysql

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"shadowsocks-go/pkg/storage"

	"github.com/golang/glog"
	dbmysql "github.com/jinzhu/gorm"
	"golang.org/x/net/context"
)

type store struct {
	client *dbmysql.DB
}

const (
	//StructTagKey struct tag key
	StructTagKey    = "column"
	ContextTableKey = "table"
)

type result struct {
	sliceValue reflect.Value
	elemType   reflect.Type
}

//New create a mysql store
func New(client *dbmysql.DB) *store {
	return &store{
		client: client,
	}
}

func (s *store) Create(ctx context.Context, key string, obj, out interface{}) error {

	rst := s.client.Table(ctx.Value(ContextTableKey).(string)).Create(obj)

	glog.V(5).Infof("add record line %+v result %v", obj, rst)
	// if record == false {
	// 	return fmt.Errorf("create record failure")
	// }
	return nil
}

//filedsToStructFieldsMap use input fields get struct field name
func (s *store) filedsToStructFieldsMap(fiedls []string, typ reflect.Type) (map[string]string, error) {
	formStructType := typ
	field := typ.NumField()

	structFields := make(map[string]string, len(fiedls))

	for i := 0; i < field; i++ {
		tag := formStructType.Field(i).Tag.Get(StructTagKey)
		structFields[tag] = formStructType.Field(i).Name
	}

	for _, v := range fiedls {
		v, ok := structFields[v]
		if !ok {
			return nil, fmt.Errorf("required export \"%v\" field", v)
		}
	}
	return structFields, nil
}

func (s *store) convertByteArry(kind reflect.Kind, buffer []byte) (interface{}, error) {
	strBuffer := string(buffer)
	var v interface{}
	var err error

	switch kind {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		fval, _ := strconv.Atoi(strBuffer)
		v = int64(fval)
	case reflect.String:
		v = string(buffer)
	default:
		err = fmt.Errorf("not support byte type to %v", kind)
	}

	return v, err
}

// https://github.com/stretchr/testify/blob/master/assert/assertions.go
// func ObjectsAreEqualValues(expected, actual interface{}) bool {
// 	if ObjectsAreEqual(expected, actual) {
// 		return true
// 	}
//
// 	actualType := reflect.TypeOf(actual)
// 	if actualType == nil {
// 		return false
// 	}
// 	expectedValue := reflect.ValueOf(expected)
// 	if expectedValue.IsValid() && expectedValue.Type().ConvertibleTo(actualType) {
// 		// Attempt comparison after type conversion
// 		return reflect.DeepEqual(expectedValue.Convert(actualType).Interface(), actual)
// 	}
//
// 	return false
// }

func (s *store) convertToInterface(kind reflect.Kind, src interface{}) (interface{}, error) {
	var v interface{}
	var err error

	switch t := src.(type) {
	//case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
	case int, int8, int16, int32, int64:
		v = reflect.ValueOf(t).Int()
	case time.Time:
		v = t
	default:
		err = fmt.Errorf("not support field type %v to %v", kind, t)
	}

	return v, err
}

//return sliceValue, elemType in slice, error
func (*store) checkResultInterface(result interface{}) error {
	if result == nil {
		return fmt.Errorf("Cannot restore result from <nil>")
	}

	resultValue := reflect.ValueOf(result)
	if resultValue.IsNil() {
		return fmt.Errorf("Cantnot reflect on a nil pointer")
	}

	resultType := resultValue.Type()
	resultKind := resultType.Kind()

	if resultKind != reflect.Ptr {
		return fmt.Errorf("Cannot reflect into non-poiner")
	}

	sliceValue := resultValue.Elem()
	sliceKind := sliceValue.Kind()

	if sliceKind != reflect.Slice {
		return fmt.Errorf("Pointer must point to a slice")
	}

	return nil
}

func (s *store) GetToList(ctx context.Context, filter storage.Filter, result interface{}) error {
	err := s.checkResultInterface(result)
	if err != nil {
		return err
	}

	resultValue := reflect.ValueOf(result)

	sliceValue := resultValue.Elem()
	sliceType := sliceValue.Type()
	elemType := sliceType.Elem()

	structFiledMap, err := s.filedsToStructFieldsMap(filter.Field(), elemType)
	if err != nil {
		return err
	}

	query, queryArgs := filter.Condition()

	rows, err := s.client.Table(ctx.Value(ContextTableKey).(string)).Model(elemType).Select(filter.Field()).Where(query, queryArgs...).Rows()

	if err != nil {
		glog.Errorf("query row failure err %v\r\n", err)
		return err
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuesPtrs := make([]interface{}, count)

	for rows.Next() {
		objVal := reflect.Indirect(reflect.New(elemType))

		for i := range columns {
			valuesPtrs[i] = &values[i]
		}
		err := rows.Scan(valuesPtrs...)
		if err != nil {
			glog.Errorf("scan rows errors %v\r\n", err)
			continue
		}

		for i := range columns {
			val := values[i]

			fieldVal := objVal.FieldByName(structFiledMap[columns[i]])
			b, ok := val.([]byte)
			if ok {
				v, err := s.convertByteArry(fieldVal.Kind(), b)
				if err != nil {
					return err
				}
				fieldVal.Set(reflect.ValueOf(v))
			} else {
				v, err := s.convertToInterface(fieldVal.Kind(), val)
				if err != nil {
					return err
				}
				fieldVal.Set(reflect.ValueOf(v))
			}
		}

		sliceValue.Set(reflect.Append(sliceValue, objVal))
	}
	return nil
}

//return sliceValue, elemType in slice, error
func (*store) checkObjInterface(result interface{}) error {
	if result == nil {
		return fmt.Errorf("Cannot restore result from <nil>")
	}

	resultValue := reflect.ValueOf(result)
	if resultValue.IsNil() {
		return fmt.Errorf("Cantnot reflect on a nil pointer")
	}

	resultType := resultValue.Type()
	resultKind := resultType.Kind()

	if resultKind != reflect.Ptr {
		return fmt.Errorf("Cannot reflect into non-poiner")
	}

	objValue := resultValue.Elem()
	objKind := objValue.Kind()

	if objKind != reflect.Struct {
		return fmt.Errorf("Pointer must point to a slice")
	}

	return nil
}

func (s *store) GuaranteedUpdate(ctx context.Context, keyField string, updateFields []string, obj interface{}) error {
	err := s.checkObjInterface(obj)
	if err != nil {
		return err
	}

	resultValue := reflect.ValueOf(obj)

	elem := resultValue.Elem()
	elemType := elem.Type()

	structFiledMap, err := s.filedsToStructFieldsMap(updateFields, elemType)
	if err != nil {
		return err
	}

	formStructValue := elem

	cond := make(map[string]interface{})
	cond[keyField] = formStructValue.FieldByName(structFiledMap[keyField]).Interface()

	update := make(map[string]interface{})
	for _, v := range updateFields {
		update[v] = formStructValue.FieldByName(structFiledMap[v]).Interface()
	}

	err = s.client.Table(ctx.Value(ContextTableKey).(string)).Where(cond).Updates(update).Error
	if err != nil {
		glog.Errorf("update err %v\r\n", err)
		return err
	}

	return nil

}

func (s *store) Delete(ctx context.Context, filter storage.Filter, out interface{}) error {
	if out == nil {
		return fmt.Errorf("Cannot restore result to <nil>")
	}

	resultValue := reflect.ValueOf(out)
	if resultValue.IsNil() {
		return fmt.Errorf("Cantnot reflect on a nil pointer")
	}

	resultType := resultValue.Type()
	resultKind := resultType.Kind()

	if resultKind != reflect.Ptr {
		return fmt.Errorf("Cannot reflect into non-poiner")
	}

	// if resultKind != reflect.Struct {
	// 	return fmt.Errorf("Cannot reflect into non-struct")
	// }

	elemType := resultType.Elem()
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("Cannot reflect into non-struct")
	}

	obj := reflect.Indirect(reflect.New(elemType))

	query, queryArgs := filter.Condition()

	err := s.client.Table(ctx.Value(ContextTableKey).(string)).Where(query, queryArgs).Delete(obj).Error
	if err != nil {
		glog.Errorf("delete err %v\r\n", err)
		return err
	}
	return nil
}
