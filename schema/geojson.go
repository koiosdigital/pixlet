package schema

import (
	"fmt"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

type GeoJSON struct {
	SchemaField
}

func newGeoJSON(
	thread *starlark.Thread,
	_ *starlark.Builtin,
	args starlark.Tuple,
	kwargs []starlark.Tuple,
) (starlark.Value, error) {
	var (
		id           starlark.String
		name         starlark.String
		desc         starlark.String
		icon         starlark.String
		def          starlark.String
		collectPoint starlark.Bool
	)

	if err := starlark.UnpackArgs(
		"GeoJSON",
		args, kwargs,
		"id", &id,
		"name", &name,
		"desc", &desc,
		"icon", &icon,
		"default?", &def,
		"collect_point?", &collectPoint,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for GeoJSON: %s", err)
	}

	s := &GeoJSON{}
	s.SchemaField.Type = "geojson"
	s.ID = id.GoString()
	s.Name = name.GoString()
	s.Description = desc.GoString()
	s.Icon = icon.GoString()
	s.Default = def.GoString()
	s.CollectPoint = bool(collectPoint)

	return s, nil
}

func (s *GeoJSON) AsSchemaField() SchemaField {
	return s.SchemaField
}

func (s *GeoJSON) AttrNames() []string {
	return []string{
		"id", "name", "desc", "icon", "default", "collect_point",
	}
}

func (s *GeoJSON) Attr(name string) (starlark.Value, error) {
	switch name {
	case "id":
		return starlark.String(s.ID), nil
	case "name":
		return starlark.String(s.Name), nil
	case "desc":
		return starlark.String(s.Description), nil
	case "icon":
		return starlark.String(s.Icon), nil
	case "default":
		return starlark.String(s.Default), nil
	case "collect_point":
		return starlark.Bool(s.CollectPoint), nil
	default:
		return nil, nil
	}
}

func (s *GeoJSON) String() string       { return "GeoJSON(...)" }
func (s *GeoJSON) Type() string         { return "GeoJSON" }
func (s *GeoJSON) Freeze()              {}
func (s *GeoJSON) Truth() starlark.Bool { return true }

func (s *GeoJSON) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(s, hashstructure.FormatV2, nil)
	return uint32(sum), err
}
