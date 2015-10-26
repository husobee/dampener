// Copyright 2015 - husobee associates, llc; all rights reserved.

// Package dampener - A toolkit to apply throttling to web applications
package dampener_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/husobee/dampener"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewThrottler(t *testing.T) {
	Convey("Test creating a NewThrottle", t, func() {
		throttle := dampener.NewThrottle(
			dampener.NewThrottleOptions(
				"IPAddress", 100, 1*time.Minute,
				func(r *http.Request) bool {
					return true
				},
				dampener.NewMemoryStorage()))
		Convey("set add event ot storage", func() {
			req, _ := http.NewRequest("GET", "/", nil)
			for i := 0; i < 100; i++ {
				throttle.AppendEvent(req)
			}
			beingThrottled, _ := throttle.ShouldThrottle(req)
			So(beingThrottled, ShouldBeFalse)
			throttle.AppendEvent(req)
			beingThrottled, _ = throttle.ShouldThrottle(req)
			So(beingThrottled, ShouldBeTrue)
		})
	})
}
