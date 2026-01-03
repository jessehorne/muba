package internal

import (
	"encoding/json"
	"io"
	"os"
)

// MapData holds all the details of a map in memory
type MapData struct {
	Name      string        `json:"name"`
	Author    string        `json:"author"`
	Version   string        `json:"version"`
	TickRate  float64       `json:"tickRate"`
	Map       [][]string    `json:"map"`
	StartTime int           `json:"startTime"`
	Teams     []MapDataTeam `json:"teams"`
	Visions   struct {
		Definitions []struct {
			ID          string       `json:"id"`
			Name        string       `json:"name"`
			Health      string       `json:"health"`
			Location    string       `json:"location"`
			RespawnTime int          `json:"respawnTime"`
			Gives       MapDataGives `json:"gives"`
		} `json:"definitions"`
	} `json:"visions"`
	Towers struct {
		Health      int          `json:"health"`
		Red         []string     `json:"red"`
		Blue        []string     `json:"blue"`
		AttackSpeed int          `json:"attackSpeed"`
		Damage      int          `json:"damage"`
		Gives       MapDataGives `json:"gives"`
	} `json:"towers"`
	Bosses struct {
		Locations   []string `json:"locations"`
		Definitions []struct {
			ID          string       `json:"id"`
			Name        string       `json:"name"`
			Health      int          `json:"health"`
			Damage      int          `json:"damage"`
			AttackSpeed int          `json:"attackSpeed"`
			Location    string       `json:"location"`
			SpawnAt     int          `json:"spawnAt"`
			Gives       MapDataGives `json:"gives"`
		} `json:"definitions"`
	} `json:"bosses"`
	Camps struct {
		Red         []string `json:"red"`
		Blue        []string `json:"blue"`
		Definitions []struct {
			ID          string   `json:"id"`
			Name        string   `json:"name"`
			DescShort   string   `json:"descShort"`
			DescLong    string   `json:"descLong"`
			Damage      int      `json:"damage"`
			AttackSpeed int      `json:"attackSpeed"`
			Health      int      `json:"health"`
			Locations   []string `json:"locations"`
			RespawnTime int      `json:"respawnTime"`
			HealthGive  int      `json:"healthGive"`
			CoinGive    int      `json:"coinGive"`
		} `json:"definitions"`
	} `json:"camps"`
	Minions struct {
		Health           int           `json:"health"`
		Speed            int           `json:"speed"`
		Gives            *MapDataGives `json:"gives"`
		RespawnTime      int           `json:"respawnTime"`
		SmallCount       int           `json:"smallCount"`
		SmallDesc        string        `json:"smallDesc"`
		SmallDamage      int           `json:"smallDamage"`
		SmallAttackSpeed int           `json:"smallAttackSpeed"`
		BigCount         int           `json:"bigCount"`
		BigDesc          string        `json:"bigDesc"`
		BigDamage        int           `json:"bigDamage"`
		BigAttackSpeed   int           `json:"bigAttackSpeed"`
	} `json:"minions"`
}

type MapDataTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Base string `json:"base"`
	Size int    `json:"size"`
}

type MapDataGives struct {
	LastHitter struct {
		Coin      int `json:"coin"`
		Health    int `json:"health"`
		MoveSpeed int `json:"moveSpeed"`
	} `json:"lastHitter"`
	Team struct {
		Vision      []string `json:"vision"`
		AttackSpeed int      `json:"attackSpeed"`
		Damage      int      `json:"damage"`
	} `json:"team"`
}

type Map struct {
	Path  string
	Data  *MapData
	Rooms [7][7]*Room
}

func NewMap(path string) (*Map, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var mapData *MapData
	err = json.Unmarshal(data, &mapData)
	if err != nil {
		return nil, err
	}

	return &Map{
		Path:  path,
		Data:  mapData,
		Rooms: [7][7]*Room{},
	}, nil
}

func (m *Map) CoordsToRoomID(v Vector2) string {
	if v.X < 0 || v.X > 8 || v.Y < 0 || v.Y > 8 {
		return ""
	}

	if len(m.Data.Map) < v.Y+1 {
		return ""
	}
	if len(m.Data.Map[v.Y]) < v.X+1 {
		return ""
	}
	return m.Data.Map[v.Y][v.X]
}

func (m *Map) RoomIDToCoords(id string) Vector2 {
	return NewVector2(int(id[1]-'1'), int(id[0]-'A'))
}

func (m *Map) CoordsToRoom(coords Vector2) *Room {
	return m.Rooms[coords.Y][coords.X]
}
