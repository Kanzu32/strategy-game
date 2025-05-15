package damagetype

type DamageType uint8

const (
	Hit DamageType = iota + 1
	Status
)
