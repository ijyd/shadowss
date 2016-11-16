package dynamodb

import (
	"fmt"
	"gofreezer/pkg/runtime"
	"strings"

	awsdb "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func BuildUpdateAttr(newObj runtime.Object, oldObj runtime.Object, attr map[string]interface{}) (updateExpression string, expressionAttributeNames map[string]*string, expressionAttributeValues map[string]*awsdb.AttributeValue, err error) {

	updateExpression = string("SET ")
	expressionAttributeNames = make(map[string]*string)
	expressionAttributeValues = make(map[string]*awsdb.AttributeValue)
	for k, v := range attr {
		fieldName := k
		var tagName string
		if i := strings.LastIndexAny(fieldName, "."); i >= 0 {
			fieldName = fieldName[i+1:]
			if i := strings.IndexAny(fieldName, "#"); i >= 0 {
				tagName = fieldName[i+1:]
			} else {
				err = fmt.Errorf("not found '#' with %v", k)
				return
			}
		} else {
			tagName = fieldName
		}

		if updateExpression != string("SET ") {
			updateExpression += string(",")
		}
		updateExpression += fmt.Sprintf("%s= :%s ", k, tagName)

		items := strings.Split(k, ".")

		for _, query := range items[1:] {
			if i := strings.IndexAny(query, "#"); i >= 0 {
				name := fmt.Sprintf("%s", query[i+1:])
				expressionAttributeNames[query] = &name
			} else {
				err = fmt.Errorf("not found '#' with %v", query)
				return
			}
		}

		dynaAttr, dynaErr := dynamodbattribute.ConvertTo(v)
		if dynaErr != nil {
			err = dynaErr
			return
		}
		expressionAttributeValues[fmt.Sprintf(":%s", tagName)] = dynaAttr
	}

	return
}
