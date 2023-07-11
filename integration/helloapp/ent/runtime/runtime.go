// Code generated by ent, DO NOT EDIT.

package runtime

import (
	"github.com/woocoos/entco/integration/helloapp/ent/schema"
	"github.com/woocoos/entco/integration/helloapp/ent/world"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	worldMixin := schema.World{}.Mixin()
	worldMixinHooks1 := worldMixin[1].Hooks()
	worldMixinHooks2 := worldMixin[2].Hooks()
	world.Hooks[0] = worldMixinHooks1[0]
	world.Hooks[1] = worldMixinHooks2[0]
	worldMixinInters1 := worldMixin[1].Interceptors()
	worldMixinInters2 := worldMixin[2].Interceptors()
	world.Interceptors[0] = worldMixinInters1[0]
	world.Interceptors[1] = worldMixinInters2[0]
	worldFields := schema.World{}.Fields()
	_ = worldFields
	// worldDescPowerBy is the schema descriptor for power_by field.
	worldDescPowerBy := worldFields[1].Descriptor()
	// world.DefaultPowerBy holds the default value on creation for the power_by field.
	world.DefaultPowerBy = worldDescPowerBy.Default.(string)
}

const (
	Version = "v0.12.4-0.20230702151415-1ec75238037c" // Version of ent codegen.
)
