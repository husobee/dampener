// Copyright 2015 - husobee associates, llc; all rights reserved.

// Package dampener - A toolkit to apply throttling to web applications
package dampener

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ThrottlerOptions - Interface that defines what a throttler option should be able to do
type ThrottlerOptions interface {
	GetStatus() int
	GetMessage() []byte
}

// dampenerPolicyOptions - implementation of ThrottlerOptions, to store options for the throttler
type dampenerPolicyOptions struct {
	StatusCode int
	Message    []byte
}

// NewThrottlerOptions - patchable entry point for getting new throttler options
var NewThrottlerOptions = newDampenerPolicyOptions

// newDampenerPolicyOptions - create a new throttler options instance
func newDampenerPolicyOptions(statusCode int, message []byte) ThrottlerOptions {
	return &dampenerPolicyOptions{
		StatusCode: statusCode,
		Message:    message,
	}
}

// GetStatus - implementation of ThrottlerOptions, to get the status to return on throttling
func (d *dampenerPolicyOptions) GetStatus() int {
	return d.StatusCode
}

// GetMessage - implementation of ThrottlerOptions, to get the message to return on throttling
func (d *dampenerPolicyOptions) GetMessage() []byte {
	return d.Message
}

// Throttler - interface that extends the http.handler with a list of throttles that will be processed
type Throttler interface {
	http.Handler
	GetThrottles() []Throttle
	GetOptions() ThrottlerOptions
}

// dapenerPolicy - An Implementation of throttler
type dampenerPolicy struct {
	throttles []Throttle
	options   ThrottlerOptions
	next      http.Handler
}

// NewThrottler - patchable entry point for getting a new throttler, to assist with unit testing
var NewThrottler = newDampenerPolicy

// newDampenerPolicy - this will create a new throttler
func newDampenerPolicy(next http.Handler, options ThrottlerOptions, throttles ...Throttle) Throttler {
	return &dampenerPolicy{
		throttles: throttles,
		options:   options,
		next:      next,
	}
}

// ServeHTTP - implementation of a Throttler for damenerPolicy
func (d *dampenerPolicy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// foreach of the throttles we are watching, check if this request
	// should be throttled.
	for _, t := range d.throttles {
		if yes, _ := t.ShouldThrottle(r); yes {
			log.Printf("should throttle, sending response")
			w.WriteHeader(d.GetOptions().GetStatus())
			w.Write(d.GetOptions().GetMessage())
			return
		}
		t.AppendEvent(r)
	}
	d.next.ServeHTTP(w, r)
}

// GetThrottles - implementation of a Throttler for DamenerPolicy
func (d *dampenerPolicy) GetThrottles() []Throttle {
	return d.throttles
}

// GetOptions - implementation of a Throttler for DamenerPolicy
func (d *dampenerPolicy) GetOptions() ThrottlerOptions {
	return d.options
}

// ThrottleOptions - interface that describes throttling options
type ThrottleOptions interface {
	GetPrefix() string
	GetMaxCount() int
	GetDuration() time.Duration
	MatchCriteria() func(*http.Request) bool
	GetStorage() Storage
}

// NewThrottleOptions - patchable entry point for getting new throttle options
var NewThrottleOptions = newDampenerThrottleOptions

// newDampenerThrottleOptions - create a new dampenerThrottleOptions
func newDampenerThrottleOptions(prefix string, count int, duration time.Duration, f func(*http.Request) bool, storage Storage) ThrottleOptions {
	return &dampenerThrottleOptions{
		prefix:   prefix,
		duration: duration,
		f:        f,
		storage:  storage,
		max:      count,
	}
}

// dampenerThrottleOptions - implementation of ThrottleOptions
type dampenerThrottleOptions struct {
	prefix   string
	duration time.Duration
	f        func(*http.Request) bool
	storage  Storage
	max      int
}

// GetPrefix - implementation of ThrottleOptions
func (d *dampenerThrottleOptions) GetPrefix() string {
	return d.prefix
}

// GetDuration - implementation of ThrottleOptions
func (d *dampenerThrottleOptions) GetDuration() time.Duration {
	return d.duration
}

// MatchCriteria - implementation of ThrottleOptions
func (d *dampenerThrottleOptions) MatchCriteria() func(*http.Request) bool {
	return d.f
}

// GetStorage - implementation of ThrottleOptions
func (d *dampenerThrottleOptions) GetStorage() Storage {
	return d.storage
}

// GetMaxCount - implementation of ThrottleOptions
func (d *dampenerThrottleOptions) GetMaxCount() int {
	return d.max
}

// Throttle - Interface to describe capabilities of a "throttle"
type Throttle interface {
	GetOptions() ThrottleOptions
	ShouldThrottle(*http.Request) (bool, error)
	AppendEvent(*http.Request)
}

// NewThrottle - patchable entry point for getting new throttle options
var NewThrottle = newDampenerThrottle

// dampenerThrottle - implementation of throttle
type dampenerThrottle struct {
	options ThrottleOptions
}

// newDampenerThrottle - new throttle implemenation
func newDampenerThrottle(options ThrottleOptions) Throttle {
	dt := &dampenerThrottle{
		options: options,
	}
	go func() {
		// clean up the MemoryStorage
		// should include a way to gracefully stop this
		for {
			time.Sleep(1 * time.Minute)
			dt.GetOptions().GetStorage().Clean(
				dt.GetOptions().GetPrefix(), time.Now())
		}
	}()
	return dt
}

// GetOptions - get options from the throttle
func (d *dampenerThrottle) GetOptions() ThrottleOptions {
	return d.options
}

// ShouldThrottle - should throttle? from the throttle
func (d *dampenerThrottle) ShouldThrottle(r *http.Request) (bool, error) {
	// is this something that matches our criteria?
	if d.GetOptions().MatchCriteria()(r) {
		// get count of events for this given prefix
		count, err := d.GetOptions().GetStorage().EventsInDuration(d.GetOptions().GetPrefix(), d.GetOptions().GetDuration())
		if err != nil {
			return false, err
		}
		fmt.Println("\n\nmax count: ", d.GetOptions().GetMaxCount(), count)
		if count > int64(d.GetOptions().GetMaxCount()) {
			return true, nil
		}
	}
	return false, nil
}

// AppendEvent - append event to storage? from the throttle
func (d *dampenerThrottle) AppendEvent(*http.Request) {
	d.GetOptions().GetStorage().AppendEvent(d.GetOptions().GetPrefix(), time.Now())
}

// Storage - Interface to implement to use various backends for storage
// of throttling data.  Comes down to the ability to store new events, clean up
// based on time elapsed, and addition of events
// and produce aggragate data for the collections
type Storage interface {
	EventsInDuration(string, time.Duration) (int64, error)
	AppendEvent(string, time.Time) error
	Clean(string, time.Time) error
}

// MemoryStorage - Dead Simple in memory storage example implementation of Storage
// would recommend making a storage implementation using a real backend data store
type MemoryStorage struct {
	m *sync.RWMutex
	s map[string][]time.Time
}

// NewMemoryStorage - Create a new memory storage
func NewMemoryStorage() Storage {
	s := &MemoryStorage{
		m: new(sync.RWMutex),
		s: make(map[string][]time.Time),
	}
	return s
}

// AppendEvent - implementing Storage Interface, adding a new event timestamp to the
// memorystorage
func (s *MemoryStorage) AppendEvent(k string, t time.Time) error {
	s.m.Lock()
	defer s.m.Unlock()
	if _, ok := s.s[k]; !ok {
		s.s[k] = []time.Time{}
	}
	s.s[k] = append(s.s[k], t)
	return nil
}

// EventsInDuration - get a cound of the events within the duration stipulated
func (s *MemoryStorage) EventsInDuration(k string, d time.Duration) (int64, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	var counter int64 = 0
	for _, v := range s.s[k] {
		if v.After(time.Now().Add(-1 * d)) {
			counter++
		}
	}
	return counter, nil
}

// Clean - clean up every earlier than "to"
func (s *MemoryStorage) Clean(k string, to time.Time) error {
	s.m.Lock()
	defer s.m.Unlock()
	for i, v := range s.s[k] {
		if v.Before(to) {
			continue
		}
		s.s[k] = s.s[k][i-1:]
		return nil
	}
	return nil
}
