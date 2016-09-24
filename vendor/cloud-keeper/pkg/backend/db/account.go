package db

import (
	"fmt"
	"time"

	"cloud-keeper/pkg/api"
	"cloud-keeper/pkg/pagination"

	"golib/pkg/storage"

	"github.com/golang/glog"
)

func GetAccounts(handle storage.Interface, page pagination.Pager) ([]api.AccountDetail, error) {

	ctx := createContextWithValue(apiKeyTableName)

	var notPage bool
	if page == nil {
		notPage = true
	} else {
		notPage = page.Empty()
	}

	var selection storage.RetrieveFilter
	fileds := []string{"id", "name", "operators", "api_key", "descryption", "lables", "expire_time", "credit_ceilings", "created_time"}
	if notPage {
		//query all records from db
		query := string("id > ?")
		queryArgs := []interface{}{0}
		selection = NewSelection(fileds, query, queryArgs)
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
			query := string("id > ?")
			queryArgs := []interface{}{0}
			skipVal := skip
			sortVal := string("id")
			limitVal := perPage
			selection = NewPageSelection(fileds, query, queryArgs, sortVal, limitVal, skipVal)
		} else {
			query := string("id > ?")
			queryArgs := []interface{}{0}
			selection = NewSelection(fileds, query, queryArgs)
		}
	}

	var acc []api.AccountDetail
	err := handle.GetToList(ctx, selection, &acc)
	if err != nil {
		return nil, err
	}

	if len(acc) > 0 {
		return acc, nil
	} else {
		return nil, fmt.Errorf("not found")
	}

}

func GetAccountByname(handle storage.Interface, name string) (*api.AccountDetail, error) {

	fileds := []string{"id", "name", "operators", "api_key", "descryption", "lables", "expire_time", "credit_ceilings", "created_time"}

	//query all records from db
	query := string("name = ?")
	queryArgs := []interface{}{name}
	selection := NewSelection(fileds, query, queryArgs)

	ctx := createContextWithValue(apiKeyTableName)

	var acc []api.AccountDetail
	err := handle.GetToList(ctx, selection, &acc)
	if err != nil {
		return nil, err
	}

	if len(acc) > 0 {
		return &acc[0], nil
	} else {
		return nil, fmt.Errorf("not found")
	}

}

func CreateAccount(handle storage.Interface, detail api.AccountDetail) error {

	ctx := createContextWithValue(apiKeyTableName)

	acc := &api.AccountDetail{
		Name:           detail.Name,
		Operators:      detail.Operators,
		Key:            detail.Key,
		Descryption:    detail.Descryption,
		CreditCeilings: detail.CreditCeilings,
		Lables:         detail.Lables,
		ExpireDBTime:   detail.ExpireTime.Time,
		CreateDBTime:   time.Now(),
	}

	glog.Infof("create a api key record %+v\r\n", acc)

	err := handle.Create(ctx, detail.Name, acc, acc)
	if err != nil {
		glog.Errorf("create a api key record failure %v\r\n", err)
	}
	return err
}

func DeleteAccount(handle storage.Interface, name string) error {

	ctx := createContextWithValue(apiKeyTableName)

	var acc api.AccountDetail
	query := string("name = ?")
	queryArgs := []interface{}{name}
	selection := NewSelection(nil, query, queryArgs)

	err := handle.Delete(ctx, selection, &acc)
	if err != nil {
		glog.Errorf("delete a apikey record failure %v\r\n", err)
	}
	return err
}
