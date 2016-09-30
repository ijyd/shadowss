package api

import (
	"cloud-keeper/pkg/pagination"
	"fmt"

	"github.com/golang/glog"

	restful "github.com/emicklei/go-restful"
)

func PageParse(req *restful.Request) (pagination.Pager, error) {
	pagParam := req.QueryParameter(string("pagination"))

	return pagination.ParsePaginaton(pagParam)

}

func PagerToCondition(pager pagination.Pager, total uint64) (bool, uint64, uint64) {

	//update current item sum
	pager.SetItemTotal(total)

	//if there have not present page do nothing
	has, _, perPage := pager.PresentPage()
	if !has {
		return false, 0, 0
	}

	var skip uint64
	hasPrev, prevPage, prevPerPage := pager.PreviousPage()
	if hasPrev {
		skip = prevPage * prevPerPage
	} else {
		skip = 0
	}

	return true, perPage, skip
}

func SetPageLink(baseLink string, response *restful.Response, pager pagination.Pager) (string, error) {

	glog.V(5).Infof("Got base link %v\r\n", baseLink)
	var link string

	if pager == nil {
		return "", nil
	}

	if !pager.Empty() {
		var prevPageLink, nextPageLink, lastPageLink string

		has, page, perPage := pager.PreviousPage()
		if has {
			prevPageLink = fmt.Sprintf("%v?pagination=page=%v,perPage=%v; rel=prev", baseLink, page, perPage)
		}

		has, page, perPage = pager.NextPage()
		if has {
			nextPageLink = fmt.Sprintf("%v?pagination=page=%v,perPage=%v; rel=next", baseLink, page, perPage)
		}

		has, page, perPage = pager.LastPage()
		if has {
			lastPageLink = fmt.Sprintf("%v?pagination=page=%v,perPage=%v; rel=last", baseLink, page, perPage)
		}

		separator := string("")
		if len(prevPageLink) > 0 {
			link += prevPageLink
			separator = string(",")
		}

		if len(nextPageLink) > 0 {
			link += separator + nextPageLink
			separator = string(",")
		}

		if len(lastPageLink) > 0 {
			link += separator + lastPageLink
		}
	}
	if response != nil {
		response.AddHeader("Link", link)
	}

	return link, nil
}
