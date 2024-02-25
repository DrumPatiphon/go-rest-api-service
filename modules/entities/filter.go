package entities

/*Pagination ก็คือการแบ่งข้อมูลที่ระบบของเรานั้นทำการตอบกลับให้กับผู้ใช้งานเป็นหน้า ๆ เหมือนกับหน้าหนังสือ (Page)
จะทำให้สามารถโหลดข้อมูลได้เร็วขึ้นแล้วยังทำให้ระบบใช้งานได้ง่ายมากขึ้น*/
type PaginationReq struct {
	Page      int `query:"page"`
	Limit     int `query:"limit"`
	TotalPage int `query:"total_page" json:"total_page"`
	TotalItem int `query:"total_item" json:"total_item"`
}

type SortReq struct {
	OrderBy string `query"order_by"`
	Sort    string `query"sort"` //DESC ASC
}
