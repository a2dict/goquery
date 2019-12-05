package goquery

type PagedQueryFunc func(*QReq) (*PageWrap, error)

func (f PagedQueryFunc) Query(qReq *QReq) (*PageWrap, error) {
	return f(qReq)
}

// QReq ...
type QReq struct {
	// 页码
	Page int64
	// 每页显示条数
	Size int64
	// 显示的字段列表. 例如: ["name","age","created_at"],为空,则默认显示所有0
	Select []string
	// 排序, 例如: ["-created_at"] 则表示按照 created_at 降序排列, 默认按照: ["-created_at"] 排序
	Sort []string
	// 查询条件，例如： GET /v1/logs?q["level"]=DEBUG
	Q map[string]string
}

// PageWrap ...
type PageWrap struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Size  int64       `json:"size"`
	Page  int64       `json:"page"`
	Pages int64       `json:"pages"`
}
