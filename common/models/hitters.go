package models

var HitterNames = []string{
	"Died", // 0
	"Crushing",
	"Fall",
	"Water",
	"Water Stun",
	"Water Stun Force",
	"Drowning",
	"Fire",
	"Burn",
	"Flying",
	"Stomp",
	"Suicide", // 11
	"Bite",
	"Pickaxe",
	"Sword",
	"Shield",
	"Bomb",
	"Stab",
	"Arrow",
	"Bomb Arrow",
	"Ballista Bolt",
	"Catapult Stones",
	"Catapult Boulder",
	"Boulder",
	"Ram",
	"Explosion",
	"Keg",
	"Mine",
	"Mine Special",
	"Spikes",
	"Saw",
	"Drill",
	"Muscles",
	"Sudden Gib",
}

func HitterName(id int64) string {
	if id >= int64(len(HitterNames)) {
		return HitterNames[0]
	}
	return HitterNames[id]
}
