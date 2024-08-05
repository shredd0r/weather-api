package marshals

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/shredd0r/weather-api/dto"
)

func MarshalGeocodingRequest(v dto.GeocodingRequest) graphql.Marshaler {
	return graphql.MarshalAny(v)
}

func UnmarshalGeocodingRequest(v any) (dto.GeocodingRequest, error) {
	return v.(dto.GeocodingRequest), nil
}

func MarshalGeocoding(v dto.Geocoding) graphql.Marshaler {
	return graphql.MarshalAny(v)
}

func UnmarshalGeocoding(v any) (dto.Geocoding, error) {
	return v.(dto.Geocoding), nil
}
