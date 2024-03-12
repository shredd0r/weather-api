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

type ApiNinjasReverseGeocodingRequestDto struct {
	Latitude  float32
	Longitude float32
}

type ApiNinjasReverseGeocodingResponseDto struct {
	Country string
	Name    string
	State   string
}
