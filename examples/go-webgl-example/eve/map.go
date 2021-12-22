package eve

import (
	"evemap/examples/go-webgl-example/fetch"
)

const (
	ScaleFactor     = 10e-15
	solarSystemsURL = "https://www.fuzzwork.co.uk/dump/latest/mapSolarSystems.csv.bz2"
	jumpsURL        = "https://www.fuzzwork.co.uk/dump/latest/mapSolarSystemJumps.csv.bz2"
)

type SolarSystem struct {
	RegionID        int     `csv:"regionID"`
	ConstellationID int     `csv:"constellationID"`
	SolarSystemID   int     `csv:"solarSystemID"`
	SolarSystemName string  `csv:"solarSystemName"`
	X               float64 `csv:"x"`
	Y               float64 `csv:"y"`
	Z               float64 `csv:"z"`
	Border          bool    `csv:"border"`
	Fringe          bool    `csv:"fringe"`
	Corridor        bool    `csv:"corridor"`
	Hub             bool    `csv:"hub"`
	International   bool    `csv:"international"`
	Regional        bool    `csv:"regional"`
	Security        float64 `csv:"security"`
}

type JumpInfo struct {
	FromRegionID        int `csv:"fromRegionID"`
	FromConstellationID int `csv:"fromConstellationID"`
	FromSolarSystemID   int `csv:"fromSolarSystemID"`
	ToSolarSystemID     int `csv:"toSolarSystemID"`
	ToConstellationID   int `csv:"toConstellationID"`
	ToRegionID          int `csv:"toRegionID"`
}

//nolint:gochecknoglobals
var (
	SolarSystems       map[int]*SolarSystem
	SolarSystemsByName map[string]*SolarSystem
	solarSystemsData   []SolarSystem
	JumpsData          []JumpInfo
)

func LoadSolarSystems(log func(...interface{})) {
	if err := fetch.CSV(
		solarSystemsURL,
		&solarSystemsData,
		log,
	); err != nil {
		log("ERROR:", err.Error())
	}

	SolarSystems = make(map[int]*SolarSystem, len(solarSystemsData))
	SolarSystemsByName = make(map[string]*SolarSystem, len(solarSystemsData))

	log("fetched stars:", len(solarSystemsData))

	for i := range solarSystemsData {
		s := &solarSystemsData[i]
		s.X *= ScaleFactor
		s.Y *= ScaleFactor
		s.Z *= ScaleFactor

		SolarSystems[s.SolarSystemID] = s
		SolarSystemsByName[s.SolarSystemName] = s
	}

	log("stars preparation complete")
}

func LoadJumps(log func(...interface{})) {
	if err := fetch.CSV(
		jumpsURL,
		&JumpsData,
		log,
	); err != nil {
		log("ERROR:", err.Error())
	}

	log("jumps load complete")
}
