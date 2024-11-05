package psize

type PageSize uint16

const (
	Page1 PageSize = 1 << iota
	Page2
	Page4
	Page8
	Page16
	Page32
	Page64
	Page128
	Page256
	Page512
	Page1024
	Page2048
	Page4096
	Page8192
)
