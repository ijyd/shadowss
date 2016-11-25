# development

## pagination
格式:  `{URL}?pageSelector=page=1,perPage=2`

`paga` 指定输出第几页。

`perPage` 指定每一页包含的资源条数。

当访问的页面超出范围，则返回全部资源

*Http Header*

在分页请求的HTTP应答包中， 系统会在头域填加`Link`，作为前后一页和最后一页做操作提示：
rel=prev提示前一页页码，rel=next提示后一页页码，rel=last提示最后一页的页码。

例如:
> Link: /api/v1beta1/namespace/default/users?pageSelector=page=1,perPage=1; rel= **prev** ,/api/v1beta1/users?pageSelector=page=3,perPage=1; rel= **next** ,/api/v1beta1/users?pageSelector=page=5,perPage=1; rel= **last**


## fieldSelector
格式:  `{URL}?fieldSelector=key=value,key1=value1`



## lableSelector
格式:  `{URL}?lableSelector=key=value,key1=value1`




*in url use '=' may be like as {URL}?fieldSelector=key%3Dvalue,key2!%3Dvalue2*
