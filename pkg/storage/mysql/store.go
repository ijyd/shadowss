package mysql

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/golang/glog"
	dbmysql "github.com/jinzhu/gorm"
)

type store struct {
	client *dbmysql.DB
}

const (
	//StructTagKey struct tag key
	StructTagKey = "column"
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
			return nil, fmt.Errorf("required export %s field", v)
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
		v, err = strconv.Atoi(strBuffer)
	case reflect.String:
		v = string(buffer)
	default:
		err = fmt.Errorf("not support field type %v", kind)
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

func (s *store) GetToList(table string, fields []string, result interface{}) error {
	err := s.checkResultInterface(result)
	if err != nil {
		return err
	}

	resultValue := reflect.ValueOf(result)

	sliceValue := resultValue.Elem()
	sliceType := sliceValue.Type()
	elemType := sliceType.Elem()

	structFiledMap, err := s.filedsToStructFieldsMap(fields, elemType)
	if err != nil {
		return err
	}

	glog.Infof("elem type:%+v", elemType)
	rows, err := s.client.Table(table).Model(elemType).Select(fields).Rows()
	if err != nil {
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
				return fmt.Errorf("not a []byte value(%v) unexcept error\r\n", b)
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

func (s *store) GuaranteedUpdate(table string, conditionFields []string, updateFields []string, obj interface{}) error {
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
	for _, v := range conditionFields {
		glog.V(5).Infof("condition %v:%v in  fieldname[%s]", v, formStructValue.FieldByName(structFiledMap[v]).Interface(), structFiledMap[v])
		cond[v] = formStructValue.FieldByName(structFiledMap[v]).Interface()
	}

	update := make(map[string]interface{})
	for _, v := range updateFields {
		glog.V(5).Infof("update %v:%v in  fieldname[%s]", v, formStructValue.FieldByName(structFiledMap[v]).Interface(), structFiledMap[v])
		update[v] = formStructValue.FieldByName(structFiledMap[v]).Interface()
	}

	err = s.client.Table(table).Where(cond).Updates(update).Error
	if err != nil {
		glog.Errorf("update err %v\r\n", err)
		return err
	}

	return nil

}
