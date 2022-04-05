// Package jsonapi provides a very basic JSON REST-like API client.
//
// In any sane world (?!) we could rely entirely on the HTTP status code and
// return (only) that when an error occurs... but there are many APIs,
// SensorPush included, that return useful data in the body regardless of the
// status code.  There are also cases where the API returns the wrong status
// code, and the body is needed for a full understanding.  (For instance
// SensorPush returns 400 with an "access denied" in the body rather than just
// using 401.)
//
// The cleanest way to handle this is to unmarshal into the response object
// *regardless* of the status code, so that it's available even if we end up
// returning an error.  The caller can use an anonymous embedded struct for the
// error fields.  (Alternatively, we can include the body bytes in the returned
// error, for separate unmarshaling?)
package jsonapi
