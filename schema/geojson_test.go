package schema_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"tidbyt.dev/pixlet/runtime"
)

var geojsonSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.GeoJSON(
	id = "area",
	name = "Area",
	desc = "A polygon area on the map.",
	icon = "map",
	default = '{"type":"Polygon","coordinates":[[[-73.95,40.68],[-73.94,40.68],[-73.94,40.69],[-73.95,40.69],[-73.95,40.68]]]}',
)

assert(s.id == "area")
assert(s.name == "Area")
assert(s.desc == "A polygon area on the map.")
assert(s.icon == "map")

def main():
	return []
`

var geojsonNoDefaultSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.GeoJSON(
	id = "zone",
	name = "Zone",
	desc = "Select a zone.",
	icon = "globe",
)

assert(s.id == "zone")
assert(s.name == "Zone")
assert(s.desc == "Select a zone.")
assert(s.icon == "globe")

def main():
	return []
`

func TestGeoJSON(t *testing.T) {
	app, err := runtime.NewApplet("geojson.star", []byte(geojsonSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

func TestGeoJSONNoDefault(t *testing.T) {
	app, err := runtime.NewApplet("geojson_no_default.star", []byte(geojsonNoDefaultSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}

var geojsonCollectPointSource = `
load("schema.star", "schema")

def assert(success, message=None):
    if not success:
        fail(message or "assertion failed")

s = schema.GeoJSON(
	id = "geo",
	name = "Geo",
	desc = "Draw a polygon and mark your location.",
	icon = "map",
	collect_point = True,
)

assert(s.id == "geo")
assert(s.name == "Geo")
assert(s.collect_point == True)

def main():
	return []
`

func TestGeoJSONCollectPoint(t *testing.T) {
	app, err := runtime.NewApplet("geojson_collect_point.star", []byte(geojsonCollectPointSource))
	assert.NoError(t, err)

	screens, err := app.Run(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, screens)
}
