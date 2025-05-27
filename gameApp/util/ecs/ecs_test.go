package ecs_test

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
var TestFlag *ecs.FlagPool

func TestECSEntityCount(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	NumberPool.AddNewEntity(NumberComponent{Num: 10})
	NumberPool.AddNewEntity(NumberComponent{Num: 10})
	NumberPool.AddNewEntity(NumberComponent{Num: 10})

	if NumberPool.EntityCount() != 3 {
		t.Fail()
	}

	NumberPool.AddNewEntity(NumberComponent{Num: 10})
	NumberPool.AddNewEntity(NumberComponent{Num: 10})

	if NumberPool.EntityCount() != 5 {
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

func TestECSFilters(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	TestFlag = ecs.CreateFlagPool(w, psize.Page1)
	ent1, err := NumberPool.AddNewEntity(NumberComponent{Num: 10})
	if err != nil {
		t.Error(err)
	}
	_, err = NumberPool.AddNewEntity(NumberComponent{Num: 20})
	if err != nil {
		t.Error(err)
	}
	TestFlag.AddExistingEntity(ent1)

	if NumberPool.EntityCount() != 2 {
		t.Fail()
	}

	if TestFlag.EntityCount() != 1 {
		t.Fail()
	}

	entites := ecs.PoolFilter([]ecs.AnyPool{NumberPool}, []ecs.AnyPool{TestFlag})

	if len(entites) != 1 {
		t.Fail()
	}
}

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

func TestECSEntityEquality(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	ent1, err := NumberPool.AddNewEntity(NumberComponent{Num: 10})
	if err != nil {
		t.Error(err)
	}
	ent2, err := NumberPool.AddNewEntity(NumberComponent{Num: 10})
	if err != nil {
		t.Error(err)
	}

	if ent1.Equals(ent2) {
		t.Error(err)
	}

	if !ent1.Equals(ent1) {
		t.Fail()
	}

	if !ent2.Equals(ent2) {
		t.Fail()
	}
}

func TestECSRemoveFromWorld(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	ent1, err := NumberPool.AddNewEntity(NumberComponent{Num: 10})
	if err != nil {
		t.Error(err)
	}

	w.RemoveEntityFromWorld(ent1)

	if NumberPool.EntityCount() != 0 {
		t.Fail()
	}
}

func TestECSRemoveFlag(t *testing.T) {
	w := ecs.CreateWorld()
	TestFlag = ecs.CreateFlagPool(w, psize.Page32)
	ent1, err := TestFlag.AddNewEntity()
	if err != nil {
		t.Error(err)
	}

	w.RemoveEntityFromWorld(ent1)

	if len(TestFlag.Entities()) != 0 {
		t.Fail()
	}
}

func TestECSString(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	TestFlag = ecs.CreateFlagPool(w, psize.Page32)

	NumberPool.String()

	TestFlag.String()
}

func TestECSComponentAdd(t *testing.T) {
	w := ecs.CreateWorld()
	NumberPool = ecs.CreateComponentPool[NumberComponent](w, psize.Page32)
	TestFlag = ecs.CreateFlagPool(w, psize.Page32)

	ent, err := TestFlag.AddNewEntity()

	if err != nil {
		t.Error(err)
	}

	NumberPool.AddExistingEntity(ent, NumberComponent{Num: 1})

	comp, err := NumberPool.Component(ent)
	if err != nil {
		t.Error(err)
	}

	if comp.Num != 1 {
		t.Fail()
	}

}
