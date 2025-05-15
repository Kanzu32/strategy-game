package ecs

import (
	"errors"
	"fmt"
	"log"
	"strategy-game/util/ecs/parray"
	"strategy-game/util/ecs/psize"

	"github.com/hajimehoshi/ebiten/v2"
)

const maxEntities = 65536

type P = AnyPool
type AnyPool interface {
	HasEntity(entity Entity) bool
	Entities() []Entity
	EntityCount() int
	RemoveEntity(entity Entity) error
	String() string
}

type Entity struct {
	State   uint8  `json:"state"`
	Id      uint16 `json:"id"`
	Version uint8  `json:"version"`
}

func (e *Entity) Equals(another Entity) bool {
	if !e.isNil() && e.isRegistered() &&
		!another.isNil() && another.isRegistered() &&
		e.Id == another.ID() && e.Version == another.Version {

		return true
	}
	return false
}
func (e *Entity) ID() uint16         { return e.Id }
func (e *Entity) isNil() bool        { return e.State&1 != 0 }
func (e *Entity) isRegistered() bool { return e.State&2 != 0 }
func (e *Entity) setNil() {
	e.State = e.State | 1
	e.State = e.State &^ 2
}
func (e *Entity) setRegistered() {
	e.State = e.State | 2
	e.State = e.State &^ 1
}

func (e *Entity) clear() { e.State = 0 }
func (e Entity) String() string {
	if e.isNil() {
		return fmt.Sprintf("E #%d v%d NIL", e.Id, e.Version)
	}
	if e.isRegistered() {
		return fmt.Sprintf("E #%d v%d REG", e.Id, e.Version)
	}
	return fmt.Sprintf("Ent #%d v%d ", e.Id, e.Version)
}

// SYSTEM

type System interface {
	Run()
}

type RenderSystem interface {
	Run(screen *ebiten.Image)
}

// WORLD

type World struct {
	pools []AnyPool

	systems       []System
	renderSystems []RenderSystem

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

func CreateComponentPool[componentType any](w *World, pageSize psize.PageSize) *ComponentPool[componentType] {
	pool := ComponentPool[componentType]{}

	pool.sparseEntities = parray.CreatePageArray(pageSize)
	pool.world = w
	w.pools = append(w.pools, &pool)
	return &pool
}

func CreateFlagPool(w *World, pageSize psize.PageSize) *FlagPool {
	pool := FlagPool{}

	pool.sparseEntities = parray.CreatePageArray(pageSize)
	pool.world = w
	w.pools = append(w.pools, &pool)
	return &pool
}

func AddSystem(w *World, system System) {
	w.systems = append(w.systems, system)
}

func AddRenderSystem(w *World, renderSystem RenderSystem) {
	w.renderSystems = append(w.renderSystems, renderSystem)
}

func (w *World) Update() {
	for _, s := range w.systems {
		s.Run()
	}
}

func (w *World) Draw(screen *ebiten.Image) {
	for _, s := range w.renderSystems {
		s.Run(screen)
	}
}

func (w *World) registerNewEntity() (Entity, error) {
	if w.destroyed.isNil() {
		if w.next == maxEntities {
			e := Entity{}
			e.setNil()
			log.Println("world have too many entities")
			return e, errors.New("too many entities")
		}
		w.entities[w.next].setRegistered()
		e := w.entities[w.next]
		e.Id = uint16(w.next)
		w.next += 1
		return e, nil
	}

	ret := w.destroyed
	w.destroyed = w.entities[ret.Id]
	w.entities[ret.Id].setRegistered()
	ret.setRegistered()
	ret.Version = w.entities[ret.Id].Version
	return ret, nil
}

func (w *World) isRegisteredEntity(entity Entity) bool {
	return w.entities[entity.Id].isRegistered()
}

// func RemoveEntityFromWorld(w *World, entity Entity) {
// 	for _, pool := range w.pools {
// 		if pool.HasEntity(entity) {
// 			pool.RemoveEntity(entity)
// 		}
// 	}
// 	w.entities[entity.Id].Id = w.destroyed.Id
// 	w.entities[entity.Id].State = w.destroyed.State
// 	w.entities[entity.Id].Version++
// 	w.destroyed.Id = entity.Id
// 	w.destroyed.clear()
// }

func (w *World) RemoveEntityFromWorld(entity Entity) {
	for _, pool := range w.pools {
		if pool.HasEntity(entity) {
			pool.RemoveEntity(entity)
		}
	}
	w.entities[entity.Id].Id = w.destroyed.Id
	w.entities[entity.Id].State = w.destroyed.State
	w.entities[entity.Id].Version++
	w.destroyed.Id = entity.Id
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
	pool.sparseEntities.Set(entity.Id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *ComponentPool[componentType]) AddExistingEntity(entity Entity, comp componentType) error {
	if pool.HasEntity(entity) {
		log.Println("entity already in the pool")
		return errors.New("entity already in the pool")
	}
	if !pool.world.isRegisteredEntity(entity) {
		log.Println("entityID is not registered")
		return errors.New("entityID is not registered")
	}
	pool.denseComponents = append(pool.denseComponents, comp)
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.Id, len(pool.denseEntities)-1)
	return nil
}

func (pool *ComponentPool[componentType]) RemoveEntity(entity Entity) error {
	denseRemoveIndex := pool.sparseEntities.Get(entity.Id)                                     // индекс для удаления (замены) элемента в dense массивах
	sparseLastIndex := pool.denseEntities[len(pool.denseEntities)-1].Id                        // индекс элемента в sparse массиве для последнего dense элемента
	pool.sparseEntities.Set(sparseLastIndex, denseRemoveIndex)                                 // установка нового указателя на dence массив sparce массиве
	pool.denseEntities[denseRemoveIndex] = pool.denseEntities[len(pool.denseEntities)-1]       // перемещение последнего элемента dense массива
	pool.denseComponents[denseRemoveIndex] = pool.denseComponents[len(pool.denseComponents)-1] // на позицию удаления для двух массивов

	pool.sparseEntities.Set(entity.Id, -1) // установка sparse эдемента для удаления в -1

	// Уменьшение len без удаления последнего элемента.
	// При необходимости его можно восстановить увиличив len. Append перезапишет скрытый элемент
	pool.denseComponents = pool.denseComponents[:len(pool.denseComponents)-1]
	pool.denseEntities = pool.denseEntities[:len(pool.denseEntities)-1]
	return nil
}

func (pool *ComponentPool[componentType]) HasEntity(entity Entity) bool {
	return pool.sparseEntities.Get(entity.Id) != -1
}

func (pool *ComponentPool[componentType]) Entities() []Entity {
	return pool.denseEntities
}

func (pool *ComponentPool[componentType]) Component(entity Entity) (*componentType, error) {
	if !pool.HasEntity(entity) {
		log.Println("entity is not in the pool")
		return nil, errors.New("entity is not in the pool")
	}
	return &pool.denseComponents[pool.sparseEntities.Get(entity.Id)], nil
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
	pool.sparseEntities.Set(entity.Id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *FlagPool) AddExistingEntity(entity Entity) (Entity, error) {
	if pool.HasEntity(entity) {
		log.Println("entity already in the pool")
		return entity, errors.New("entity already in the pool")
	}
	if !pool.world.isRegisteredEntity(entity) {
		log.Println("entityID is not registered")
		return entity, errors.New("entityID is not registered")
	}
	pool.denseEntities = append(pool.denseEntities, entity)
	pool.sparseEntities.Set(entity.Id, len(pool.denseEntities)-1)
	return entity, nil
}

func (pool *FlagPool) RemoveEntity(entity Entity) error {
	denseRemoveIndex := pool.sparseEntities.Get(entity.Id)
	sparseLastIndex := pool.denseEntities[len(pool.denseEntities)-1].Id
	pool.sparseEntities.Set(sparseLastIndex, denseRemoveIndex)
	pool.denseEntities[denseRemoveIndex] = pool.denseEntities[len(pool.denseEntities)-1]

	pool.sparseEntities.Set(entity.Id, -1)

	pool.denseEntities = pool.denseEntities[:len(pool.denseEntities)-1]
	return nil
}

func (pool *FlagPool) HasEntity(entity Entity) bool {
	return pool.sparseEntities.Get(entity.Id) != -1
}

func (pool *FlagPool) Entities() []Entity {
	return pool.denseEntities
}

func (pool *FlagPool) EntityCount() int {
	return len(pool.denseEntities)
}

func (pool *FlagPool) String() string {
	return fmt.Sprintf("Dense ent: %v\nSparse ent:\n%v", pool.denseEntities, pool.sparseEntities.String())
}
