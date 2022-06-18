package partner

type Partner struct {
	ID              int
	Rating          int
	OperatingRadius int
	Experiences     []string
	Location
}

type Location struct {
	Latitude  float64
	Longitude float64
}
