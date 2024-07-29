package dto

type ApiNinjasGeocodingRequestDto struct {
	City    string
	State   string
	Country string
}

type ApiNinjasGeocodingResponseDto struct {
	Name      string
	Latitude  float64
	Longitude float64
	Country   string
	State     string
}

type ApiNinjasReverseGeocodingRequestDto struct {
	Latitude  float64
	Longitude float64
}

type ApiNinjasReverseGeocodingResponseDto struct {
	Country string
	Name    string
	State   string
}

type ApiNinjasError struct {
	Error string `json:"error"`
}
