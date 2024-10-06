package ecs

import (
	"errors"
	"fmt"
	"strategy-game/ecs/parray"
	"strategy-game/ecs/psize"
)

const maxEntities = 65536

type P = AnyPool
type AnyPool interface {
	HasEntity(entity Entity) bool
	Entities() []Entity
	EntityCount() int
	RemoveEntity(entity Entity) error
}

type Entity struct {
	state   uint8
	id      uint16
	version uint8
}

func (e *Entity) ID() uint16         { return e.id }
func (e *Entity) isNil() bool        { return e.state&1 != 0 }
func (e *Entity) isRegistered() bool { return e.state&2 != 0 }
func (e *Entity) setNil() {
	e.state = e.state | 1
	e.state = e.state &^ 2
}
func (e *Entity) setRegistered() {
	e.state = e.state | 2
	e.state = e.state &^ 1
}

func (e *Entity) clear() { e.state = 0 }
func (e Entity) String() string {
	if e.isNil() {
		return fmt.Sprintf("E #%d v%d NIL", e.id, e.version)
	}
	if e.isRegistered() {
		return fmt.Sprintf("E #%d v%d REG", e.id, e.version)
	}
	return fmt.Sprintf("Ent #%d v%d ", e.id, e.version)
}

// SYSTEM

type System interface {
	Run(frameCount int)
}

// WORLD

type World struct {
	pools []AnyPool

	systems []System

	next      uint32   // next available entity ID
	entities  []Entity // array to mark registred and destroyed entities
	destroyed Entity   // last entity removed from forld
}

func CreateWorld() *World {
	w := World{}
	w.pools = make([]AnyPool, 0)
	w.entities = make([]Entity, maxEntities)
	e := Entity{}
	e.setNil()
	w.destroyed = e
	return &w
}

func CreateComponentPool[componentType any](w *World, pageSize psize.PageSizes) *ComponentPool[componentType] {
	pool := ComponentPool[componentType]{}

	pool.sparseEntities = parray.CreatePageArray(pageSize)
	pool.world = w
	w.pools = append(w.pools, &pool)
	return &pool
}

func CreateFlagPool(w *World, pageSize psize.PageSizes) *FlagPool {
	pool := FlagPool{}

	pool.sparseEntities = parray.CreatePageArray(pageSize)
	pool.world = w
	w.pools = append(w.pools, &pool)
	return &pool
}

func AddSystem(w *World, system System) {
	w.systems = append(w.systems, system)
}

func (w *World) Update(frameCount int) {
	for _, s := range w.systems {
		s.Run(frameCount)
	}
}

func (w *World) registerNewEntity() (Entity, error) {
	if w.destroyed.isNil() {
		if w.next == maxEntities {
			e := Entity{}
			e.setNil()
			return e, errors.New("too many entities")
		}
		w.entities[w.next].setRegistered()
		e := w.entities[w.next]
		e.id = uint16(w.next)
		w.next += 1
		return e, nil
	}

	ret := w.destroyed
	w.destroyed = w.entities[ret.id]
	w.entities[ret.id].setRegistered()
	ret.setRegistered()
	ret.version = w.entities[ret.id].version
	return ret, nil
}

func (w *World) isRegisteredEntity(entity Entity) bool {
	return w.entities[entity.id].isRegistered()
}

func RemoveEntityFromWorld(w *World, entity Entity) {
	for _, pool := range w.pools {
		if pool.HasEntity(entity) {
			pool.RemoveEntity(entity)
		}
	}
	w.entities[entity.id].id = w.destroyed.id
	w.entities[entity.id].state = w.destroyed.state
	w.entities[entity.id].version++
	w.destroyed.id = entity.id
	w.destroyed.clear()
}

// COMPONENT POOL

type ComponentPool[componentType any] struct {
	denseComponents []componentType

	denseEntities []Entity

	sparseEntities parray.PageArray

	world *World
}

func (pool *ComponentPool[componentType]) AddNewEntity(comp componentType) (Entity, error) {
	entity, err := pool.world.registerNewEntity()
	if err != nil {
		return entity, err
	}
	pool.denseComponents = append(pool.denseComponents, comp)
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *ComponentPool[componentType]) AddExistingEntity(entity Entity, comp componentType) error {
	if !pool.world.isRegisteredEntity(entity) {
		return errors.New("entityID is not registered")
	}
	pool.denseComponents = append(pool.denseComponents, comp)
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return nil
}

func (pool *ComponentPool[componentType]) RemoveEntity(entity Entity) error {
	denseRemoveIndex := pool.sparseEntities.Get(entity.id)                                     // индекс для удаления (замены) элемента в dense массивах
	sparseLastIndex := pool.denseEntities[len(pool.denseEntities)-1].id                        // индекс элемента в sparse массиве для последнего dense элемента
	pool.sparseEntities.Set(sparseLastIndex, denseRemoveIndex)                                 // установка нового указателя на dence массив sparce массиве
	pool.denseEntities[denseRemoveIndex] = pool.denseEntities[len(pool.denseEntities)-1]       // перемещение последнего элемента dense массива
	pool.denseComponents[denseRemoveIndex] = pool.denseComponents[len(pool.denseComponents)-1] // на позицию удаления для двух массивов

	pool.sparseEntities.Set(entity.id, -1) // установка sparse эдемента для удаления в -1

	// Уменьшение len без удаления последнего элемента.
	// При необходимости его можно восстановить увиличив len. Append перезапишет скрытый элемент
	pool.denseComponents = pool.denseComponents[:len(pool.denseComponents)-1]
	pool.denseEntities = pool.denseEntities[:len(pool.denseEntities)-1]
	return nil
}

func (pool *ComponentPool[componentType]) HasEntity(entity Entity) bool {
	return pool.sparseEntities.Get(entity.id) != -1
}

func (pool *ComponentPool[componentType]) Entities() []Entity {
	return pool.denseEntities
}

func (pool *ComponentPool[componentType]) Component(entity Entity) (*componentType, error) {
	if !pool.HasEntity(entity) {
		return nil, errors.New("Entity is not in the pool")
	}
	return &pool.denseComponents[pool.sparseEntities.Get(entity.id)], nil
}

func (pool *ComponentPool[poolType]) EntityCount() int {
	return len(pool.denseEntities)
}

func (pool *ComponentPool[componentType]) String() string {
	return fmt.Sprintf("Components: %v\nDense ent: %v\nSparse ent:\n%v", pool.denseComponents, pool.denseEntities, pool.sparseEntities.String())
}

// ENTITY FILTER

func PoolFilter(include []AnyPool, exclude []AnyPool) []Entity {
	if len(include) == 0 {
		panic("include can't be empty")
	}
	shortestIndex := 0
	shortestLen := include[0].EntityCount()
	res := make([]Entity, 0)
	for i, pool := range include {
		if shortestLen > pool.EntityCount() {
			shortestLen = pool.EntityCount()
			shortestIndex = i
		}
	}
EntityLoop:
	for _, entity := range include[shortestIndex].Entities() {
		for _, pool := range include {
			if !pool.HasEntity(entity) {
				continue EntityLoop
			}
		}

		for _, pool := range exclude {
			if pool.HasEntity(entity) {
				continue EntityLoop
			}
		}

		res = append(res, entity)
	}
	return res
}

// FLAG POOL

type FlagPool struct {
	denseEntities []Entity

	sparseEntities parray.PageArray

	world *World
}

func (pool *FlagPool) AddNewEntity() (Entity, error) {
	entity, err := pool.world.registerNewEntity()
	if err != nil {
		return entity, err
	}
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *FlagPool) AddExistingEntity(entity Entity) (Entity, error) {
	if !pool.world.isRegisteredEntity(entity) {
		return entity, errors.New("entityID is not registered")
	}
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *FlagPool) RemoveEntity(entity Entity) error {
	denseRemoveIndex := pool.sparseEntities.Get(entity.id)
	sparseLastIndex := pool.denseEntities[len(pool.denseEntities)-1].id
	pool.sparseEntities.Set(sparseLastIndex, denseRemoveIndex)
	pool.denseEntities[denseRemoveIndex] = pool.denseEntities[len(pool.denseEntities)-1]

	pool.sparseEntities.Set(entity.id, -1)

	pool.denseEntities = pool.denseEntities[:len(pool.denseEntities)-1]
	return nil
}

func (pool *FlagPool) HasEntity(entity Entity) bool {
	return pool.sparseEntities.Get(entity.id) != -1
}

func (pool *FlagPool) Entities() []Entity {
	return pool.denseEntities
}

func (pool *FlagPool) EntityCount() int {
	return len(pool.denseEntities)
}
