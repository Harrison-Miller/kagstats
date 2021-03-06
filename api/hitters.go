package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/Harrison-Miller/kagstats/common/models"
	"github.com/gorilla/mux"
)

type Hitter struct {
	PlayerID int64  `json:"-" db:"playerID"`
	Hitter   int64  `json:"hitter"`
	Kills    int64  `json:"kills"`
	Name     string `json:"name"`
}

type HittersStats struct {
	PlayerID int64 `json:"-" db:"playerID"`
	models.Player `json:"player" db:"player,prefix=player."`

	// hitters
	Died int64 `json:"died"`
	Crushing int64 `json:"crushing"`
	Fall int64 `json:"fall"`
	Water int64 `json:"water"`
	WaterStun int64 `json:"water_stun" db:"water_stun"`
	WaterStunForce int64 `json:"water_stun_force" db:"water_stun_force"`
	Drowning int64 `json:"drowning"`
	Fire int64 `json:"fire"`
	Burn int64 `json:"burn"`
	Flying int64 `json:"flying"`
	Stomp int64 `json:"stomp"`
	Suicide int64 `json:"suicide"`
	Bite int64 `json:"bite"`
	Pickaxe int64 `json:"pickaxe"`
	Sword int64 `json:"sword"`
	Shield int64 `json:"shield"`
	Bomb int64 `json:"bomb"`
	Stab int64 `json:"stab"`
	Arrow int64 `json:"arrow"`
	BombArrow int64 `json:"bomb_arrow" db:"bomb_arrow"`
	BallistaBolt int64 `json:"ballista_bolt" db:"ballista_bolt"`
	CatapultStones int64 `json:"catapult_stones" db:"catapult_stones"`
	CatapultBoulder int64 `json:"catapult_boulder" db:"catapult_boulder"`
	Boulder int64 `json:"boulder"`
	Ram int64 `json:"ram"`
	Explosion int64 `json:"explosion"`
	Keg int64 `json:"keg"`
	Mine int64 `json:"mine"`
	MineSpecial int64 `json:"mine_special" db:"mine_special"`
	Spikes int64 `json:"spikes"`
	Saw int64 `json:"saw"`
	Drill int64 `json:"drill"`
	Muscles int64 `json:"muscles"`
	SuddenGib int64 `json:"sudden_gib" db:"sudden_gib"`
}

type MonthlyHittersStats struct {
	Year int64 `json:"year"`
	Month int64 `json:"month"`
	HittersStats
}

func HitterRoutes(r *mux.Router) {
	r.HandleFunc("/players/{id:[0-9]+}/hitters", GetHitters).Methods("GET")
	r.HandleFunc("/players/{id:[0-9]+}/hitters/monthly", GetMonthlyHitters).Methods("GET")
	r.HandleFunc("/leaderboard/monthly/hitter/{id:[0-9]+}", GetMonthlyHitterLeaderboard).Methods("GET")
}

type HittersList struct {
	MyPlayer models.Player `json:"player"`
	Size     int           `json:"size"`
	Hitters  []Hitter      `json:"hitters"`
}

// GetHitters godoc
// @Tags Detailed Stats
// @Summary Returns the top five weapons (hitters) used by the player
// @Produce json
// @Param id path int true "PlayerID"
// @Success 200 {object} HittersList
// @Router /players/{id}/hitters [get]
func GetHitters(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var player models.Player
	err = db.Get(&player, "SELECT * FROM players WHERE ID=?", playerID)
	if err != nil {
		playerNotFoundError(w, err)
		return
	}

	h := []Hitter{}
	err = db.Select(&h, `SELECT * FROM top_hitters AS hitters WHERE hitters.playerID=? ORDER BY hitters.kills DESC LIMIT 5`, playerID)
	if err != nil {
		log.Printf("Could not find hitters for player: %v\n", err)
		http.Error(w, "Could not find hitters for player", http.StatusInternalServerError)
		return
	}

	for i, hitter := range h {
		h[i].Name = models.HitterName(hitter.Hitter)
	}

	JSONResponse(w, HittersList{
		MyPlayer: player,
		Size:     len(h),
		Hitters:  h,
	})
}

func GetMonthlyHitters(w http.ResponseWriter, r *http.Request) {
	playerID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	var stats []MonthlyHittersStats
	err = db.Select(&stats, `SELECT * FROM monthly_hitters WHERE playerID=? ORDER BY monthly_hitters.year DESC LIMIT 12`, playerID)
	if err != nil {
		log.Printf("Could not find hitters for player: %v\n", err)
		http.Error(w, "Could not find hitters for player", http.StatusInternalServerError)
		return
	}

	JSONResponse(w, &stats)
}

type MonthlyHittersLeaderboardList struct {
	Size int `json:"size"`
	Year int `json:"year"`
	Month int `json:"month"`
	Hitter int `json:"hitter"`
	HitterName string `json:"hitterName"`
	Leaderboard []MonthlyHittersStats `json:"leaderboard"`
}

func GetMonthlyHitterLeaderboard(w http.ResponseWriter, r *http.Request) {
	year, month := getYearMonth(r)

	hitterID, err := GetIntURLArg("id", r)
	if err != nil {
		http.Error(w, "could not get id", http.StatusBadRequest)
		return
	}

	hitterName := models.HitterName(int64(hitterID))
	hitterName = strings.ToLower(hitterName)
	hitterName = strings.ReplaceAll(hitterName, " ", "_")

	var stats []MonthlyHittersStats
	err = db.Select(&stats, `SELECT monthly_hitters.*, p.ID "player.ID", p.username "player.username",
		p.charactername "player.charactername", p.clantag "player.clantag", p.clanID "player.clanID", p.oldgold "player.oldgold",
		p.registered "player.registered", p.role "player.role", p.avatar "player.avatar", p.tier "player.tier",
		p.gold "player.gold", p.silver "player.silver", p.bronze "player.bronze", p.participation "player.participation",
		p.github "player.github", p.community "player.community", p.mapmaker "player.mapmaker", p.moderation "player.moderation",
		p.leaderboardBan "player.leaderboardBan", p.statsBan "player.statsBan", c.name "player.clan_info.name"
		FROM monthly_hitters
		INNER JOIN players as p ON monthly_hitters.playerID=p.ID 
		LEFT JOIN clan_info as c ON p.clanID=c.ID 
		WHERE NOT p.monthlyLeaderboardBan AND NOT p.statsBan AND monthly_hitters.year=? AND monthly_hitters.month=? ORDER BY monthly_hitters.` + hitterName + ` DESC LIMIT 20`, year, month)
	if err != nil {
		leaderboardError(w, err)
		return
	}

	JSONResponse(w, MonthlyHittersLeaderboardList{
		Size:        len(stats),
		Year:        year,
		Month:       month,
		Hitter:      hitterID,
		HitterName:  hitterName,
		Leaderboard: stats,
	})
}