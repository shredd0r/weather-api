package dto

type GeocodingRequestDto struct {
	city    string
	country string
	state   string
}

type GeocodingResponseDto struct {
	Name      string
	Latitude  float32
	Longitude float32
	Country   string
}
