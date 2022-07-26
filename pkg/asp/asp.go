package asp

// We want to encourage the best behavior by making it as easy as possible to
// provide all of the setting variants:
//   * config file (the config struct itself)
//   * command line - long and short versions
//   * environment variable(s)
// ... and also description info.
//
// for example:
//   type Config struct {
//     Host string `asp:"host,h,APP_HOST,The host to use."`
//   }

import (
	// "builtin"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrConfigMustBeStruct = errors.New("config must be a struct or pointer to struct")
)

// IncomingConfig is a placeholder generic type that exists only to allow us to
// define our innner and exposed value as strongly-typed to the originating
// configuration struct.
type IncomingConfig interface {
	interface{}
}

// Asp is an interface that represents the "callable" interface for
// settings/options.  After creating/initializing with a configuration structure
// (with default values), the methods on the interface allow for loading from
// command-line/config/environment, as well as lower-level access to the created
// viper instance and cobra command.  (In most cases these should not be needed,
// though!)
type Asp[T IncomingConfig] interface {
	// Config() T
	Command() *cobra.Command
	Viper() *viper.Viper

	Execute(handler func(config T, args []string)) error
}

// maybe the incoming config should be *explicitly* the defaults, and a new
// struct created instead of loading in-place?
func New[T IncomingConfig](config T, envPrefix string) (Asp[T], error) {
	vip := viper.New()
	cmd := &cobra.Command{
		// Run: func(cmd *cobra.Command, args []string) {
		// 	log.Printf("INSIDE COMMAND! %q", args)
		// },
	}

	a := &asp[T]{
		// config: config,
		envPrefix: envPrefix,
		vip:       vip,
		cmd:       cmd,
	}
	log.Printf("initializing config for: %#v", config)

	var err error
	a.baseType, err = a.processStruct(config, "")
	if err != nil {
		return nil, err
	}

	return a, nil

}

type asp[T IncomingConfig] struct {
	baseType reflect.Type
	// config T
	envPrefix string
	vip       *viper.Viper
	cmd       *cobra.Command
}

// processStruct is the workhorse that adds a (sub-)struct config into the viper
// config and cobra command.
func (a *asp[T]) processStruct(s interface{}, parentCanonical string) (baseType reflect.Type, err error) {
	vip, cmd := a.vip, a.cmd

	log.Printf("initializing struct for: %#v", s)

	// We expect the incoming value to be a struct or a pointer to a struct.
	// Anything else is invalid.
	val := reflect.Indirect(reflect.ValueOf(s))
	if val.Kind() != reflect.Struct {
		err = ErrConfigMustBeStruct
		return
	}

	baseType = val.Type()

	fields := reflect.VisibleFields(val.Type())
	// log.Printf("fields: %#v", fields)

	for i, f := range fields {
		v2 := val.Field(i)

		canonicalName, attrLong, attrShort, attrEnv, attrDesc := getAttributes(f, parentCanonical, a.envPrefix)

		// *some* kinds (pointers, array, slices, maps, structs) fundamentally
		// adjust the command-line and env-var options... other kinds (value
		// types) just change which CLI mapping/parsing are really used.
		k := f.Type.Kind()
		handled := false
		isSlice := false
		// isArray := false
		// isPointer := false
		// isMap := false

		for !handled {
			handled = true

			switch k {
			case reflect.Pointer, reflect.Array, reflect.Slice, reflect.Map:
				k = f.Type.Elem().Kind()
				log.Printf("%q is pointer, array, slice, or map... to %v", canonicalName, k)
				isSlice = true
				handled = false
			case reflect.Struct:
				log.Printf("%q is struct...", canonicalName)
				a.processStruct(v2.Interface(), fmt.Sprintf("%s.", canonicalName))
			default:
				log.Printf("%q, %v, CLI: %q / %q, env: %q, desc: %q", canonicalName, f.Type.Kind(), attrLong, attrShort, attrEnv, attrDesc)
				// Start pushing into viper?  Note that we're going to need to handle
				// parent paths pretty quickly!
				vip.SetDefault(canonicalName, v2.Interface())
				vip.BindEnv(canonicalName, attrEnv)
				// cmd.Flags().StringP(attrLong, attrShort, v2.String(), attrDesc)

				switch k {
				case reflect.Bool:
					cmd.Flags().BoolP(attrLong, attrShort, v2.Bool(), attrDesc)
				case reflect.Int:
					if !isSlice {
						cmd.Flags().IntP(attrLong, attrShort, int(v2.Int()), attrDesc)
					} else {
						cmd.Flags().IntSliceP(attrLong, attrShort, v2.Interface().([]int), attrDesc)
					}
				case reflect.Uint:
					cmd.Flags().UintP(attrLong, attrShort, uint(v2.Uint()), attrDesc)
				case reflect.String:
					cmd.Flags().StringP(attrLong, attrShort, v2.String(), attrDesc)
				}

				vip.BindPFlag(canonicalName, cmd.Flags().Lookup(attrLong))
			}
		}
	}

	return
}

func (a *asp[T]) Execute(handler func(config T, args []string)) error {
	// Set up run-handler for the cobra command...
	a.cmd.Run = func(cmd *cobra.Command, args []string) {
		log.Printf("BEFORE (INSIDE): %v", a.vip.AllSettings())
		// TODO: unmarshal the settings into the expected config type!
		cfgVal := reflect.New(a.baseType)
		handler(cfgVal.Interface().(T), args)
		log.Printf("AFTER (INSIDE): %v", a.vip.AllSettings())
	}

	log.Printf("BEFORE: %v", a.vip.AllSettings())

	// a.cmd.ParseFlags()
	err := a.cmd.Execute()
	log.Printf("error? %v", err)

	log.Printf("AFTER: %v", a.vip.AllSettings())
	return err
}

// func (a *asp[T]) Config() T {
// 	return a.config
// }

func (a *asp[T]) Command() *cobra.Command {
	return a.cmd
}

func (a *asp[T]) Viper() *viper.Viper {
	return a.vip
}
