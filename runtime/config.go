package runtime

import (
	"fmt"
	"strconv"

	"github.com/mitchellh/hashstructure/v2"
	"go.starlark.net/starlark"
)

// DefaultDisplayWidth is the default width for display rendering
const DefaultDisplayWidth = 64

// DefaultDisplayHeight is the default height for display rendering
const DefaultDisplayHeight = 32

type AppletConfig map[string]string

// AppletConfigWithDimensions wraps AppletConfig and adds display dimension accessors.
// This allows apps to access width/height as integers directly instead of parsing strings.
type AppletConfigWithDimensions struct {
	config        AppletConfig
	displayWidth  int
	displayHeight int
}

// NewAppletConfigWithDimensions creates a config wrapper with explicit dimensions.
// If width or height is 0, defaults to 64x32.
func NewAppletConfigWithDimensions(config map[string]string, width, height int) *AppletConfigWithDimensions {
	if width <= 0 {
		width = DefaultDisplayWidth
	}
	if height <= 0 {
		height = DefaultDisplayHeight
	}
	return &AppletConfigWithDimensions{
		config:        AppletConfig(config),
		displayWidth:  width,
		displayHeight: height,
	}
}

func (c *AppletConfigWithDimensions) AttrNames() []string {
	return []string{
		"get",
		"str",
		"bool",
		"width",
		"height",
	}
}

func (c *AppletConfigWithDimensions) Attr(name string) (starlark.Value, error) {
	switch name {
	case "get", "str":
		return starlark.NewBuiltin("str", c.getString), nil
	case "bool":
		return starlark.NewBuiltin("bool", c.getBoolean), nil
	case "width":
		return starlark.NewBuiltin("width", c.getWidth), nil
	case "height":
		return starlark.NewBuiltin("height", c.getHeight), nil
	default:
		return nil, nil
	}
}

func (c *AppletConfigWithDimensions) Get(key starlark.Value) (starlark.Value, bool, error) {
	return c.config.Get(key)
}

func (c *AppletConfigWithDimensions) String() string       { return "AppletConfig(...)" }
func (c *AppletConfigWithDimensions) Type() string         { return "AppletConfig" }
func (c *AppletConfigWithDimensions) Freeze()              {}
func (c *AppletConfigWithDimensions) Truth() starlark.Bool { return true }

func (c *AppletConfigWithDimensions) Hash() (uint32, error) {
	return c.config.Hash()
}

func (c *AppletConfigWithDimensions) getString(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return c.config.getString(thread, nil, args, kwargs)
}

func (c *AppletConfigWithDimensions) getBoolean(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return c.config.getBoolean(thread, nil, args, kwargs)
}

func (c *AppletConfigWithDimensions) getWidth(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.MakeInt(c.displayWidth), nil
}

func (c *AppletConfigWithDimensions) getHeight(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.MakeInt(c.displayHeight), nil
}

func (a AppletConfig) AttrNames() []string {
	return []string{
		"get",
		"str",
		"bool",
	}
}

func (a AppletConfig) Attr(name string) (starlark.Value, error) {
	switch name {

	case "get", "str":
		return starlark.NewBuiltin("str", a.getString), nil

	case "bool":
		return starlark.NewBuiltin("bool", a.getBoolean), nil

	default:
		return nil, nil
	}
}

func (a AppletConfig) Get(key starlark.Value) (starlark.Value, bool, error) {
	switch v := key.(type) {
	case starlark.String:
		val, found := a[v.GoString()]
		return starlark.String(val), found, nil
	default:
		return nil, false, nil
	}
}

func (a AppletConfig) String() string       { return "AppletConfig(...)" }
func (a AppletConfig) Type() string         { return "AppletConfig" }
func (a AppletConfig) Freeze()              {}
func (a AppletConfig) Truth() starlark.Bool { return true }

func (a AppletConfig) Hash() (uint32, error) {
	sum, err := hashstructure.Hash(a, hashstructure.FormatV2, nil)
	return uint32(sum), err
}

func (a AppletConfig) getString(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String
	var def starlark.Value
	def = starlark.None

	if err := starlark.UnpackPositionalArgs(
		"str", args, kwargs, 1,
		&key, &def,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for config.str: %v", err)
	}

	val, ok := a[key.GoString()]
	if !ok {
		return def, nil
	} else {
		return starlark.String(val), nil
	}
}

func (a AppletConfig) getBoolean(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var key starlark.String
	var def starlark.Value
	def = starlark.None

	if err := starlark.UnpackPositionalArgs(
		"bool", args, kwargs, 1,
		&key, &def,
	); err != nil {
		return nil, fmt.Errorf("unpacking arguments for config.bool: %v", err)
	}

	val, ok := a[key.GoString()]
	if !ok {
		return def, nil
	} else {
		b, _ := strconv.ParseBool(val)
		return starlark.Bool(b), nil
	}
}
