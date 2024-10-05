package message

type IDate interface {
	Size() (n int)
	Marshal() (dAtA []byte, err error)
	MarshalTo(dAtA []byte) (int, error)
	MarshalToSizedBuffer(dAtA []byte) (int, error)
	Unmarshal(dAtA []byte) error
}
type IMessage interface {
	IDate
	GetId() uint16
}
type IReqMessage interface {
	IMessage
}
type IRepMessage interface {
	IMessage
	GetTag() int32
}
