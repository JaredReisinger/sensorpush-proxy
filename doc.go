// sensorpush-proxy is a rate-limiting, authentication-hiding proxy for
// SensorPush data.
//
// See the README (https://github.com/JaredReisinger/sensorpush-proxy#readme)
// for end-user usage information.  From a programmatic standpoint, the
// [jsonapi] and [sensorpush] packages might be interesting.
package main

import (
	_ "github.com/jaredreisinger/sensorpush-proxy/pkg/jsonapi"
	_ "github.com/jaredreisinger/sensorpush-proxy/pkg/sensorpush"
)
