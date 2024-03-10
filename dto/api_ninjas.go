package dto

type ApiNinjasGeocodingRequestDto struct {
	City    string
	State   string
	Country string
}

type ApiNinjasGeocodingResponseDto struct {
	Name      string
	Latitude  float32
	Longitude float32
	Country   string
	State     string
}
