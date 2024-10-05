package tzDataI

import (
	"tzgit.kaixinxiyou.com/engine/message/byteBuf"
)

type ITZData interface {
	Init(bb byteBuf.IByteBuf) error
	GetTable(tableId int32) ITable
	GetVer() int64
}
type ITable interface {
	GetRowCount() int32
	GetRows() []IRow
	InsertRow(idx int32, row IRow)
}
type IRow interface {
	GetValue(field int32) int64
	GetValueStr(field int32) string
}
type ITZManage interface {
	Create() ITZData
}
