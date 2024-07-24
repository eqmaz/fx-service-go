package rates

import (
	"errors"
	"fx-service/internal/service/providers"
	"fx-service/pkg/e"
	"testing"
)

// MockProvider is a manual mock implementation of the ProviderInterface
type MockProvider struct {
	Name     string
	CallFunc func(from string, to interface{}, isMulti bool) (interface{}, error)
}

func (m *MockProvider) GetName() string {
	return m.Name
}

func (m *MockProvider) CheckApiKey() bool {
	return true
}

func (m *MockProvider) GetRate(from string, to string) (float64, error) {
	return 1, nil
}

func (m *MockProvider) GetRates(from string, to []string) (providers.RateList, error) {
	return make(providers.RateList), nil
}

func (m *MockProvider) Supports(currency string) bool {
	return true
}

func (m *MockProvider) CallProvider(from string, to interface{}, isMulti bool) (interface{}, error) {
	return m.CallFunc(from, to, isMulti)
}

func TestCallProviderRoundRobin(t *testing.T) {
	// Create mock providers with manual setup
	mockProvider1 := &MockProvider{
		Name: "Provider1",
		CallFunc: func(from string, to interface{}, isMulti bool) (interface{}, error) {
			return nil, errors.New("failure")
		},
	}

	mockProvider2 := &MockProvider{
		Name: "Provider2",
		CallFunc: func(from string, to interface{}, isMulti bool) (interface{}, error) {
			return "Success", nil
		},
	}

	mockProvider3 := &MockProvider{
		Name: "Provider3",
		CallFunc: func(from string, to interface{}, isMulti bool) (interface{}, error) {
			return nil, errors.New("failure")
		},
	}

	// Initialize the roundRobinState with mock providers
	rrs = &roundRobinState{
		providers: []providers.ProviderInterface{mockProvider1, mockProvider2, mockProvider3},
		nextIndex: 0,
	}

	// Test successful provider call
	result, providerName, err := callProviderRoundRobin("from", "to", false)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "Success" {
		t.Errorf("expected result to be 'Success', got %v", result)
	}
	if *providerName != "Provider2" {
		t.Errorf("expected provider name to be 'Provider2', got %v", *providerName)
	}

	// Reset nextIndex and test round-robin logic
	rrs.nextIndex = 0
	mockProvider1.CallFunc = func(from string, to interface{}, isMulti bool) (interface{}, error) {
		return "Success", nil
	}

	result, providerName, err = callProviderRoundRobin("from", "to", false)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "Success" {
		t.Errorf("expected result to be 'Success', got %v", result)
	}
	if *providerName != "Provider1" {
		t.Errorf("expected provider name to be 'Provider1', got %v", *providerName)
	}

	// Test all providers failing
	mockProvider1.CallFunc = func(from string, to interface{}, isMulti bool) (interface{}, error) {
		return nil, errors.New("failure")
	}
	mockProvider2.CallFunc = func(from string, to interface{}, isMulti bool) (interface{}, error) {
		return nil, errors.New("failure")
	}
	mockProvider3.CallFunc = func(from string, to interface{}, isMulti bool) (interface{}, error) {
		return nil, errors.New("failure")
	}

	rrs.nextIndex = 0
	result, providerName, err = callProviderRoundRobin("from", "to", false)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	if result != nil {
		t.Errorf("expected result to be nil, got %v", result)
	}
	if providerName != nil {
		t.Errorf("expected provider name to be nil, got %v", *providerName)
	}
	if err != e.FromCode(errAllFailed) {
		t.Errorf("expected error code %v, got %v", e.FromCode(errAllFailed), err)
	}
}
