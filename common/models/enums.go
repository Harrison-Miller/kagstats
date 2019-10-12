package models

const (
	Nothing int64 = iota

	//env
	Crush //(required to be 1 for engine reasons)
	Fall
	Water            //just fire
	Water_stun       //splash
	Water_stun_force //splash
	Drown
	Fire //initial burst (should ignite things)
	Burn //burn damage
	Flying

	//common actor
	Stomp
	Suicide //(required to be 11 for engine reasons)

	//natural
	Bite

	//builders
	Builder

	//knight
	Sword
	Shield
	Bomb

	//archer
	Stab

	//arrows and similar projectiles
	Arrow
	Bomb_arrow
	Ballista

	//cata
	Cata_stones
	Cata_boulder
	Boulder

	//siege
	Ram

	// explosion
	Explosion
	Keg
	Mine
	Mine_special

	//traps
	Spikes

	//machinery
	Saw
	Drill

	//barbarian
	Muscles

	// scrolls
	Suddengib
)
