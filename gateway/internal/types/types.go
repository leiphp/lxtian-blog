// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.2

package types

type ArticleListReq struct {
	Cid      uint32 `form:"cid,optional"`
	Types    uint32 `form:"types,optional"`
	Page     uint32 `form:"page"`
	PageSize uint32 `form:"page_size"`
}

type ArticleListResp struct {
	Page     uint32                   `json:"page"`
	PageSize uint32                   `json:"page_size"`
	List     []map[string]interface{} `json:"list"`
	Total    uint64                   `json:"total"`
}

type ArticleReq struct {
	Id uint32 `path:"id"`
}

type ArticleResp struct {
	Data map[string]interface{} `json:"data"`
}

type CategoryListReq struct {
	Page     uint32 `form:"page,optional"`
	PageSize uint32 `form:"page_size,optional"`
}

type CategoryListResp struct {
	Page     uint32                   `json:"page"`
	PageSize uint32                   `json:"page_size"`
	List     []map[string]interface{} `json:"list"`
	Total    uint64                   `json:"total"`
}

type ChatListReq struct {
	Cid      uint32 `form:"cid,optional"`
	Page     uint32 `form:"page"`
	PageSize uint32 `form:"page_size"`
}

type ChatListResp struct {
	Page     uint32                   `json:"page"`
	PageSize uint32                   `json:"page_size"`
	List     []map[string]interface{} `json:"list"`
	Total    uint64                   `json:"total"`
}

type CommentListReq struct {
	Cid      uint32 `form:"cid,optional"`
	Page     uint32 `form:"page"`
	PageSize uint32 `form:"page_size"`
}

type CommentListResp struct {
	Page     uint32                   `json:"page"`
	PageSize uint32                   `json:"page_size"`
	List     []map[string]interface{} `json:"list"`
	Total    uint64                   `json:"total"`
}

type InfoResp struct {
	Data map[string]interface{} `json:"data"`
}

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResp struct {
	AccessToken string                 `json:"access_token"`
	ExpiresIn   uint64                 `json:"expires_in"`
	User        map[string]interface{} `json:"user"`
}

type OrderListReq struct {
	Cid      uint32 `form:"cid,optional"`
	Page     uint32 `form:"page"`
	PageSize uint32 `form:"page_size"`
}

type OrderListResp struct {
	Page     uint32                   `json:"page"`
	PageSize uint32                   `json:"page_size"`
	List     []map[string]interface{} `json:"list"`
	Total    uint64                   `json:"total"`
}

type RegisterReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}

type RegisterResp struct {
	Data map[string]interface{} `json:"data"`
}

type TagsListReq struct {
	Cid      uint32 `form:"cid,optional"`
	Page     uint32 `form:"page"`
	PageSize uint32 `form:"page_size"`
}

type TagsListResp struct {
	Page     uint32                   `json:"page"`
	PageSize uint32                   `json:"page_size"`
	List     []map[string]interface{} `json:"list"`
	Total    uint64                   `json:"total"`
}
