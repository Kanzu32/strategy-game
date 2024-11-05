package parray

import (
	"fmt"
	"strategy-game/util/ecs/psize"
)

type PageArray struct {
	data      [][]int
	pageSize  uint16
	arraySize int
}

func CreatePageArray(pageSize psize.PageSize) PageArray {
	p := PageArray{}
	p.pageSize = uint16(pageSize)
	p.arraySize = 65536
	p.data = make([][]int, p.arraySize/int(p.pageSize))
	return p
}

func (p *PageArray) Size() int {
	return p.arraySize
}

func (p *PageArray) Set(index uint16, value int) {
	pageNumber := index / p.pageSize
	pageIndex := index % p.pageSize

	if p.data[pageNumber] == nil {
		p.data[pageNumber] = make([]int, p.pageSize)

		p.data[pageNumber][0] = -1

		for j := 1; j < len(p.data[pageNumber]); j *= 2 {
			copy(p.data[pageNumber][j:], p.data[pageNumber][:j])
		}
	}

	p.data[pageNumber][pageIndex] = value
}

func (p *PageArray) Get(index uint16) int {
	pageNumber := index / p.pageSize
	pageIndex := index % p.pageSize

	if p.data[pageNumber] == nil {
		return -1
	}

	return p.data[pageNumber][pageIndex]
}

func (p *PageArray) String() string {
	return fmt.Sprintf("Page size: %d\nArray size: %d\n%v", p.pageSize, p.arraySize, p.data)
}
