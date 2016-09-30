package pagination

import (
	"fmt"
	"github.com/golang/glog"
)

func PagerToCondition(pager Pager, total uint64) (bool, uint64, uint64) {

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

//BuildPageLink build a string 'Link' that like as :
// "Link: /api/v1beta1/namespace/default/users?pagination=page=1,perPage=1; rel=prev,
// "/api/v1beta1/users?pagination=page=3,perPage=1; rel= next,"
// "/api/v1beta1/users?pagination=page=5,perPage=1; rel=last"
func BuildDefPageLink(pager Pager, baseLink string) (string, error) {

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

	return link, nil
}
