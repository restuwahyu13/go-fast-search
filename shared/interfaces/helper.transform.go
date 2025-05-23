package inf

type ITransform interface {
	SrcToDest(src, dest any) error
	ReqToRes(src, dest any) error
	ResToReq(src, dest any) error
}
