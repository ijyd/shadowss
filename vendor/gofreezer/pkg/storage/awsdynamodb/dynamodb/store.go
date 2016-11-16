package dynamodb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"gofreezer/pkg/conversion"
	"gofreezer/pkg/runtime"
	storage "gofreezer/pkg/storage"
	"gofreezer/pkg/storage/awsdynamodb"
	storagehelper "gofreezer/pkg/storage/helper"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsdb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/golang/glog"

	"golang.org/x/net/context"
)

type store struct {
	codec     runtime.Codec
	versioner APIObjectVersioner
	dbHandler *awsdb.DynamoDB
	table     string
}

//New create a mongo store
func New(sess *session.Session, table string, codec runtime.Codec) *store {
	versioner := APIObjectVersioner{}
	db := awsdb.New(sess)

	if len(table) == 0 {
		table = defaultTable
	}
	desc, err := CreateTable(db, table)
	if err != nil {
		glog.Fatalf("table not active, error : %v", err)
		return nil
	}
	glog.V(5).Infof("Got table(%v) description: %v", table, desc)

	return &store{
		codec:     codec,
		versioner: versioner,
		dbHandler: db,
		table:     table,
	}
}

func (s *store) Type() string {
	return string("dynamodb")
}

// Versioner implements storage.Interface.Versioner.
func (s *store) Versioner() storage.Versioner {
	return s.versioner
}

func (s *store) Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error {

	//check item with this key exist,we need replace this by aws PutItem
	err := s.queryObjByKey(key, out, false)
	if err != nil {
		if storage.IsNotFound(err) == false {
			return storage.NewInternalErrorf("key %v, object search error %v", err.Error())
		}
	} else if err == nil {
		return storage.NewItemExistsError(key, string("object exist"))
	}

	data, err := runtime.Encode(s.codec, obj)
	if err != nil {
		return storage.NewInternalErrorf("key %v, object encode error %v", key, err.Error())
	}

	mapObj, err := ConvertByteToMap(data)
	if err != nil {
		return storage.NewInternalErrorf("key %v, object encode error %v", key, err.Error())
	}

	item, err := dynamodbattribute.MarshalMap(mapObj)
	if err != nil {
		return err
	}
	item[primaryKey] = &awsdb.AttributeValue{S: aws.String(key)}
	item[sortKey] = &awsdb.AttributeValue{S: aws.String(time.Now().String())}

	_, err = s.dbHandler.PutItem(&awsdb.PutItemInput{
		Item:      item,
		TableName: aws.String(s.table),
		//ReturnValues: aws.String("ALL_OLD"),
	})
	if err != nil {
		return storage.NewInternalErrorf("key %v, put error %v\r\n", key, err)
	}

	return decode(s.codec, s.versioner, data, out)
}

func (s *store) Delete(ctx context.Context, key string, out runtime.Object, preconditions *storage.Preconditions) error {
	params := &awsdb.DeleteItemInput{
		Key: map[string]*awsdb.AttributeValue{
			"key": &awsdb.AttributeValue{
				S: aws.String(key),
			},
		},
		ReturnValues: aws.String("ALL_OLD"),
		TableName:    aws.String(s.table),
	}

	resp, err := s.dbHandler.DeleteItem(params)
	if err != nil {
		return storage.NewInternalErrorf("key %v delete error %v\r\n", err.Error())
	}
	glog.V(5).Infof("got result %v err %v\r\n", resp.Attributes, err)

	return s.getObject(key, out, false, resp.Attributes)
}

func (s *store) Get(ctx context.Context, key string, out runtime.Object, ignoreNotFound bool) error {
	return s.queryObjByKey(key, out, ignoreNotFound)
}

func (s *store) GetToList(ctx context.Context, key string, p storage.SelectionPredicate, listObj runtime.Object) error {
	listPtr, _, err := storagehelper.GetListItemObj(listObj)
	if err != nil {
		return storage.NewInvalidObjError(key, err.Error())
	}

	//always need to get list count,
	//prevent return a lots of item from backend
	scanParam := &awsdb.ScanInput{
		TableName: aws.String(s.table),
		Select:    aws.String("COUNT"),
	}
	output, err := s.dbHandler.Scan(scanParam)
	if err != nil {
		return storage.NewInternalErrorf("key %v, scan list count error %v", err.Error())
	}

	hasPage, perPage, _ := p.BuildPagerCondition(uint64(*output.Count))

	scanParam = &awsdb.ScanInput{
		TableName: aws.String(s.table),
	}
	if hasPage && perPage != 0 {
		limit := int64(perPage)
		scanParam.Limit = &(limit)
	}
	output, err = s.dbHandler.Scan(scanParam)
	if err != nil {
		return storage.NewInternalErrorf("key %v, scan list  error %v", err.Error())
	}

	glog.V(5).Infof("Get query output %+v\r\n", output)

	jsonData, cnt, err := ConvertTOJson(&output.Items)
	if cnt == 0 {
		return nil
	}

	return decodeList(jsonData, listPtr, s.codec, s.versioner)
}

func (s *store) GuaranteedUpdate(ctx context.Context, key string, out runtime.Object, ignoreNotFound bool, precondtions *storage.Preconditions, tryUpdate awsdynamodb.UpdateFunc) error {
	//check item with this key exist,we need replace this by aws PutItem
	err := s.queryObjByKey(key, out, false)
	if err != nil {
		return storage.NewInternalErrorf("key %s, search error %v", err.Error())
	}

	attrValue := make(map[string]interface{})
	ret, _, err := userUpdate(out, tryUpdate, attrValue)
	if err != nil {
		return storage.NewInternalErrorf("key %s, update by user error:%v", key, err.Error())
	}

	updateExpression, expressionAttributeNames, expressionAttributeValues, err := BuildUpdateAttr(ret, out, attrValue)
	if err != nil {
		return storage.NewInternalErrorf("key %v, update error %v\r\n", key, err.Error())
	}
	glog.V(5).Infof("build attr UpdateExpression: %v expressionAttributeNames:%v expressionAttributeValues:%v",
		updateExpression, expressionAttributeNames, expressionAttributeValues)

	params := &awsdb.UpdateItemInput{
		Key: map[string]*awsdb.AttributeValue{
			"key": &awsdb.AttributeValue{
				S: aws.String(key),
			},
		},
		ReturnValues:              aws.String("ALL_NEW"),
		TableName:                 aws.String(s.table),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeNames:  expressionAttributeNames,
		ExpressionAttributeValues: expressionAttributeValues,
	}
	resp, err := s.dbHandler.UpdateItem(params)
	if err != nil {
		return storage.NewInternalErrorf("key %v update error %v\r\n", err.Error())
	}
	glog.V(5).Infof("got result %v err %v\r\n", resp.Attributes, err)

	return s.getObject(key, out, false, resp.Attributes)
}

// decode decodes value of bytes into object. It will also set the object resource version to rev.
// On success, objPtr would be set to the object.
func decode(codec runtime.Codec, versioner storage.Versioner, value []byte, objPtr runtime.Object) error {
	if _, err := conversion.EnforcePtr(objPtr); err != nil {
		panic("unable to convert output object to pointer")
	}
	_, _, err := codec.Decode(value, nil, objPtr)
	if err != nil {
		return err
	}
	// being unable to set the version does not prevent the object from being extracted
	//versioner.UpdateObject(objPtr, uint64(rev))
	return nil
}

// decodeList decodes a list of values into a list of objects, with resource version set to corresponding rev.
// On success, ListPtr would be set to the list of objects.
func decodeList(elems []map[string]interface{}, ListPtr interface{}, codec runtime.Codec, versioner storage.Versioner) error {
	v, err := conversion.EnforcePtr(ListPtr)
	if err != nil || v.Kind() != reflect.Slice {
		panic("need ptr to slice")
	}
	for _, elem := range elems {
		data, err := json.Marshal(elem)
		if err != nil {
			return storage.NewInternalError(err.Error())
		}
		obj, _, err := codec.Decode(data, nil, reflect.New(v.Type().Elem()).Interface().(runtime.Object))
		if err != nil {
			return err
		}
		// being unable to set the version does not prevent the object from being extracted
		// versioner.UpdateObject(obj, elem.rev)
		// if filter(obj) {
		v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
		// }
	}
	return nil
}

func userUpdate(input runtime.Object, userUpdate awsdynamodb.UpdateFunc, attributeValues map[string]interface{}) (output runtime.Object, ttl *uint64, err error) {
	ret, ttl, err := userUpdate(input, attributeValues)
	if err != nil {
		return nil, nil, err
	}
	return ret, ttl, nil
}

func (s *store) queryObjByKey(key string, out runtime.Object, ignoreNotFound bool) error {
	scanParam := &awsdb.ScanInput{
		ScanFilter: map[string]*awsdb.Condition{
			"key": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*awsdb.AttributeValue{
					{
						S: aws.String(key),
					},
				},
			},
		},
		TableName: aws.String(s.table),
	}

	output, err := s.dbHandler.Scan(scanParam)
	if err != nil {
		return err
	}

	return s.getObject(key, out, ignoreNotFound, output.Items...)
}

func (s *store) getObject(key string, out runtime.Object, ignoreNotFound bool, attrs ...map[string]*awsdb.AttributeValue) error {

	jsonData, count, err := ConvertTOJson(&attrs)
	if count == 0 {
		if ignoreNotFound {
			return runtime.SetZeroValue(out)
		}
		return storage.NewItemNotFoundError(key)
	} else if count > 1 {
		return storage.NewTooManyItemError(key, fmt.Sprint("too many item found by key"))
	}

	firstObj := jsonData[0]
	data, err := json.Marshal(firstObj)
	if err != nil {
		return storage.NewInternalError(err.Error())
	}
	return decode(s.codec, s.versioner, data, out)
}
