package dampener

import (
	"net/http"
	"time"
)

// MockStorage - Mockable implementation of the storage interface.  Use this in your unit tests
// so that you can mock out the storage layer, like this:
// storage := &MockStorage{
//     MockSetKeyWithDuration: func(string, interface{}, time.Duration) error {
//         return errors.New("mocking out a failure condition")
//     },
//}
// See this blog for rationale: https://husobee.github.io/golang/testing/unit-test/2015/06/08/golang-unit-testing.html
type MockStorage struct {
	MockEventsInDuration func(string, time.Duration) (int64, error)
	MockAppendEvent      func(string, time.Time) error
	MockClean            func(string, time.Time) error
}

// EventsInDuration - implementation of Storage interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockStorage) EventsInDuration(k string, d time.Duration) (int64, error) {
	if s.MockEventsInDuration != nil {
		return s.MockEventsInDuration(k, d)
	}
	return 0, nil
}

// AppendEvent - implementation of Storage interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockStorage) AppendEvent(k string, t time.Time) error {
	if s.MockAppendEvent != nil {
		return s.MockAppendEvent(k, t)
	}
	return nil

}

// Clean - implementation of Storage interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockStorage) Clean(k string, t time.Time) error {
	if s.MockClean != nil {
		return s.MockClean(k, t)
	}
	return nil
}

// MockThrottleOptions - Mockable implementation of the throttle options interface.  Use this in your unit tests
// so that you can mock out the throttle options layer, like this:
// throttleOptions := &MockThrottleOptions{
//     MockGetStorage: func() Storage{
//	       return &MockStorage{}
//     },
//}
// See this blog for rationale: https://husobee.github.io/golang/testing/unit-test/2015/06/08/golang-unit-testing.html
type MockThrottleOptions struct {
	MockGetPrefix     func() string
	MockGetDuration   func() time.Duration
	MockMatchCriteria func() func(*http.Request) bool
	MockGetStorage    func() Storage
	MockGetMaxCount   func() int
}

// GetPrefix - implementation of ThrottleOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottleOptions) GetPrefix() string {
	if s.MockGetPrefix != nil {
		return s.MockGetPrefix()
	}
	return ""
}

// GetMaxCount - implementation of ThrottleOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottleOptions) GetMaxCount() int {
	if s.MockGetPrefix != nil {
		return s.MockGetMaxCount()
	}
	return 100
}

// GetDuration - implementation of ThrottleOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottleOptions) GetDuration() time.Duration {
	if s.MockGetDuration != nil {
		return s.MockGetDuration()
	}
	return 1 * time.Hour
}

// MatchCriteria - implementation of ThrottleOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottleOptions) MatchCriteria() func(*http.Request) bool {
	if s.MockMatchCriteria != nil {
		return s.MockMatchCriteria()
	}
	return func(*http.Request) bool { return true }
}

// GetStorage - implementation of ThrottleOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottleOptions) GetStorage() Storage {
	if s.MockGetStorage != nil {
		return s.MockGetStorage()
	}
	return &MockStorage{}
}

// MockThrottler - Mockable implementation of the throttler interface.  Use this in your unit tests
// so that you can mock out the throttle options layer, like this:
// throttler := &MockThrottler{
//     MockGetThrottles: func() []Throttle{
//	       return []Throttles{&MockThrottle{}}
//     },
//}
// See this blog for rationale: https://husobee.github.io/golang/testing/unit-test/2015/06/08/golang-unit-testing.html
type MockThrottler struct {
	MockServeHTTP    func(http.ResponseWriter, *http.Request)
	MockGetThrottles func() []Throttle
	MockGetOptions   func() ThrottlerOptions
}

// ServeHTTP - implementation of Throttler interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.MockServeHTTP != nil {
		s.MockServeHTTP(w, r)
		return
	}
	return
}

// GetThrottles - implementation of Throttler interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottler) GetThrottles() []Throttle {
	if s.MockGetThrottles != nil {
		return s.MockGetThrottles()
	}
	return []Throttle{}
}

// GetOptions - implementation of Throttler interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottler) GetOptions() ThrottlerOptions {
	if s.MockGetOptions != nil {
		return s.MockGetOptions()
	}
	return &MockThrottlerOptions{}
}

// MockThrottlerOptions - Mockable implementation of the throttleroptions interface.  Use this in your unit tests
// so that you can mock out the throttle options layer, like this:
// throttlerOptions := &MockThrottlerOptions{
//     MockGetStatus: func() int {
//         return http.StatusOK
//     },
//}
// See this blog for rationale: https://husobee.github.io/golang/testing/unit-test/2015/06/08/golang-unit-testing.html
type MockThrottlerOptions struct {
	MockGetStatus  func() int
	MockGetMessage func() []byte
}

// GetStatus - implementation of ThrottlerOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottlerOptions) GetStatus() int {
	if s.MockGetStatus != nil {
		return s.MockGetStatus()
	}
	return http.StatusOK
}

// GetMessage - implementation of ThrottlerOptions interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottlerOptions) GetMessage() []byte {
	if s.MockGetMessage != nil {
		return s.MockGetMessage()
	}
	return []byte{}
}

// MockThrottle - Mockable implementation of the throttle interface.  Use this in your unit tests
// so that you can mock out the throttle options layer, like this:
// throttle := &MockThrottle{
//     MockShouldThrottle: func(*http.Request) bool{
//         return true
//     },
//}
// See this blog for rationale: https://husobee.github.io/golang/testing/unit-test/2015/06/08/golang-unit-testing.html
type MockThrottle struct {
	MockGetOptions     func() ThrottleOptions
	MockShouldThrottle func(*http.Request) bool
	MockAppendEvent    func(*http.Request)
}

// GetOptions - implementation of Throttle interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottle) GetOptions() ThrottleOptions {
	if s.MockGetOptions != nil {
		return s.MockGetOptions()
	}
	return &MockThrottleOptions{}
}

// ShouldThrottle - implementation of Throttle interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottle) ShouldThrottle(r *http.Request) bool {
	if s.MockShouldThrottle != nil {
		return s.MockShouldThrottle(r)
	}
	return false
}

// AppendEvent - implementation of Throttle interface, allowing for a custom mock
// function to be specified for unit testing
func (s *MockThrottle) AppendEvent(r *http.Request) {
	if s.MockAppendEvent != nil {
		s.MockAppendEvent(r)
		return
	}
	return
}
