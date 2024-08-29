// Code generated by MockGen. DO NOT EDIT.
// Source: location.go
//
// Generated by this command:
//
//	mockgen -source=location.go -destination=mock/location.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	dto "github.com/shredd0r/weather-api/dto"
	gomock "go.uber.org/mock/gomock"
)

// MockLocationService is a mock of LocationService interface.
type MockLocationService struct {
	ctrl     *gomock.Controller
	recorder *MockLocationServiceMockRecorder
}

// MockLocationServiceMockRecorder is the mock recorder for MockLocationService.
type MockLocationServiceMockRecorder struct {
	mock *MockLocationService
}

// NewMockLocationService creates a new mock instance.
func NewMockLocationService(ctrl *gomock.Controller) *MockLocationService {
	mock := &MockLocationService{ctrl: ctrl}
	mock.recorder = &MockLocationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocationService) EXPECT() *MockLocationServiceMockRecorder {
	return m.recorder
}

// FindGeocoding mocks base method.
func (m *MockLocationService) FindGeocoding(ctx context.Context, request *dto.GeocodingRequest) (*[]*dto.Geocoding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindGeocoding", ctx, request)
	ret0, _ := ret[0].(*[]*dto.Geocoding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindGeocoding indicates an expected call of FindGeocoding.
func (mr *MockLocationServiceMockRecorder) FindGeocoding(ctx, request any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindGeocoding", reflect.TypeOf((*MockLocationService)(nil).FindGeocoding), ctx, request)
}

// LocationByCoords mocks base method.
func (m *MockLocationService) LocationByCoords(ctx context.Context, coords *dto.Coords) (*dto.LocationInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LocationByCoords", ctx, coords)
	ret0, _ := ret[0].(*dto.LocationInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LocationByCoords indicates an expected call of LocationByCoords.
func (mr *MockLocationServiceMockRecorder) LocationByCoords(ctx, coords any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LocationByCoords", reflect.TypeOf((*MockLocationService)(nil).LocationByCoords), ctx, coords)
}
