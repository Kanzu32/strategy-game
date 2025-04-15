package main_test

import (
	"strategy-game/util/ecs"
	"strategy-game/util/ecs/psize"
	"testing"
)

type AddSystem struct{}

func (s *AddSystem) Run() {
	entities := ecs.PoolFilter([]ecs.AnyPool{NumberPool}, []ecs.AnyPool{})
	for _, entity := range entities {
		c, err := NumberPool.Component(entity)
		if err != nil {
			panic(err)
		}
		c.Num += 5
	}
}

type NumberComponent struct {
	Num int
}

var NumberPool *ecs.ComponentPool[NumberComponent]

func TestECS(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	ent1, err := NumberPool.AddNewEntity(NumberComponent{Num: 10})
	if err != nil {
		t.Error(err)
	}
	ent2, err := NumberPool.AddNewEntity(NumberComponent{Num: 20})
	if err != nil {
		t.Error(err)
	}
	ecs.AddSystem(w, &AddSystem{})
	w.Update()
	c1, err := NumberPool.Component(ent1)
	if err != nil {
		t.Error(err)
	}
	c2, err := NumberPool.Component(ent2)
	if err != nil {
		t.Error(err)
	}

	if c1.Num != 15 || c2.Num != 25 {
		t.Fail()
	}

	w.Update()

	c1, err = NumberPool.Component(ent1)
	if err != nil {
		t.Error(err)
	}
	c2, err = NumberPool.Component(ent2)
	if err != nil {
		t.Error(err)
	}

	if c1.Num != 20 || c2.Num != 30 {
		t.Fail()
	}

}

func TestECSEntityRemove(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	ent, err := NumberPool.AddNewEntity(NumberComponent{Num: 1})
	if err != nil {
		t.Error(err)
	}

	if !NumberPool.HasEntity(ent) {
		t.Fail()
	}

	NumberPool.RemoveEntity(ent)

	if NumberPool.HasEntity(ent) {
		t.Fail()
	}

}
