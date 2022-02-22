package models

type Request struct {
	IP        string
	Continent string
}

type IPstat struct {
	IP   string
	Hits int
}

type DayStat struct {
	Day       string
	Hits      int
	Addresses []string
}
