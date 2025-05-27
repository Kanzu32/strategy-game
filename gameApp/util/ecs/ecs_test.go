package ecs

import (
	"strategy-game/util/ecs/psize"
	"testing"
)

type AddNumSystem struct{}

func (s *AddNumSystem) Run() {
	entities := PoolFilter([]AnyPool{NumberPool}, []AnyPool{})
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

var NumberPool *ComponentPool[NumberComponent]
var TestFlag *FlagPool

func TestECSEntityCount(t *testing.T) {
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
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
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
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
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
	TestFlag = CreateFlagPool(w, psize.Page1)
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

	entites := PoolFilter([]AnyPool{NumberPool}, []AnyPool{TestFlag})

	if len(entites) != 1 {
		t.Fail()
	}
}

func TestECS(t *testing.T) {
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
	ent1, err := NumberPool.AddNewEntity(NumberComponent{Num: 10})
	if err != nil {
		t.Error(err)
	}
	ent2, err := NumberPool.AddNewEntity(NumberComponent{Num: 20})
	if err != nil {
		t.Error(err)
	}
	AddSystem(w, &AddNumSystem{})
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
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
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
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
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
	w := CreateWorld()
	TestFlag = CreateFlagPool(w, psize.Page32)
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
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
	TestFlag = CreateFlagPool(w, psize.Page32)

	NumberPool.String()

	TestFlag.String()
}

func TestECSComponentAdd(t *testing.T) {
	w := CreateWorld()
	NumberPool = CreateComponentPool[NumberComponent](w, psize.Page32)
	TestFlag = CreateFlagPool(w, psize.Page32)

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

// Тесты для Entity
func TestEntity(t *testing.T) {
	t.Run("StateManagement", func(t *testing.T) {
		e := &Entity{}

		// Проверка начального состояния
		if e.isNil() {
			t.Error("New entity should not be nil")
		}
		if e.isRegistered() {
			t.Error("New entity should not be registered")
		}

		// Регистрация
		e.setRegistered()
		if e.isNil() {
			t.Error("Registered entity should not be nil")
		}
		if !e.isRegistered() {
			t.Error("Entity should be registered")
		}

		// Помечаем как nil
		e.setNil()
		if !e.isNil() {
			t.Error("Entity should be nil after setNil")
		}
		if e.isRegistered() {
			t.Error("Entity should not be registered after setNil")
		}
	})

	t.Run("Equals", func(t *testing.T) {
		e1 := &Entity{Id: 1, Version: 1, State: 2} // registered
		e2 := &Entity{Id: 1, Version: 1, State: 2}

		if !e1.Equals(*e2) {
			t.Error("Identical entities should be equal")
		}

		e3 := &Entity{Id: 1, Version: 2, State: 2}
		if e1.Equals(*e3) {
			t.Error("Different versions should not be equal")
		}
	})
}

// Тесты для World
func TestWorld(t *testing.T) {
	t.Run("CreateWorld", func(t *testing.T) {
		w := CreateWorld()

		if len(w.pools) != 0 {
			t.Error("New world should have no pools")
		}
		if len(w.entities) != maxEntities {
			t.Errorf("Expected %d entities, got %d", maxEntities, len(w.entities))
		}
	})

	t.Run("EntityLifecycle", func(t *testing.T) {
		w := CreateWorld()

		// Создаем новую сущность
		e, err := w.registerNewEntity()
		if err != nil {
			t.Fatalf("Failed to register entity: %v", err)
		}

		if !w.isRegisteredEntity(e) {
			t.Error("Entity should be registered")
		}

		// Удаляем сущность
		w.RemoveEntityFromWorld(e)
		if w.isRegisteredEntity(e) {
			t.Error("Entity should be unregistered after removal")
		}

		// Проверяем версию
		if w.entities[e.Id].Version != 1 {
			t.Error("Version should increment after removal")
		}
	})
}

// Тесты для ComponentPool
func TestComponentPool(t *testing.T) {
	w := CreateWorld()
	pool := CreateComponentPool[int](w, psize.Page1)

	t.Run("AddRemove", func(t *testing.T) {
		e, err := pool.AddNewEntity(42)
		if err != nil {
			t.Fatalf("AddNewEntity failed: %v", err)
		}

		if !pool.HasEntity(e) {
			t.Error("Pool should contain the added entity")
		}

		// Проверяем внутреннее состояние
		if len(pool.denseComponents) != 1 || pool.denseComponents[0] != 42 {
			t.Error("Component not stored correctly")
		}

		// Удаляем
		err = pool.RemoveEntity(e)
		if err != nil {
			t.Fatalf("RemoveEntity failed: %v", err)
		}

		if pool.HasEntity(e) {
			t.Error("Pool should not contain removed entity")
		}
	})
}

// Вспомогательные функции для тестов
func TestSparseArrayIntegration(t *testing.T) {
	w := CreateWorld()
	pool := CreateComponentPool[string](w, psize.Page1)

	e1, _ := pool.AddNewEntity("first")
	e2, _ := pool.AddNewEntity("second")

	// Проверяем sparse индексы
	idx1 := pool.sparseEntities.Get(e1.Id)
	idx2 := pool.sparseEntities.Get(e2.Id)

	if idx1 != 0 || idx2 != 1 {
		t.Errorf("Sparse indices incorrect: %d, %d", idx1, idx2)
	}

	// Удаляем первый элемент и проверяем перестановку
	pool.RemoveEntity(e1)
	idx2 = pool.sparseEntities.Get(e2.Id)
	if idx2 != 0 {
		t.Error("Sparse index not updated after removal")
	}
}
