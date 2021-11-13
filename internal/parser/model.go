package parser

const (
	MSK CityCode = "msk"
	SPB CityCode = "spb"

	NameMSK = "Москва"
	NameSPB = "Санкт-Петербург"

	sadwaveURL = "https://sadwave.com"
)

var knownCities = []*City{
	{
		Code: MSK,
		Name: NameMSK,
	},
	{
		Code: SPB,
		Name: NameSPB,
	},
}

type CityCode string

type City struct {
	Code CityCode
	Name string
}

type CityEvents struct {
	City   *City
	Events []*Event
}

type Event struct {
	Title           string
	DescriptionHTML string
	ImageURL        string
}
