package psize

type PageSizes uint16

const (
	Page16 PageSizes = 1 << (iota + 4)
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
