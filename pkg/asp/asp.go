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
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrConfigMustBeStruct     = errors.New("config must be a struct or pointer to struct")
	ErrConfigFieldUnsupported = errors.New("config struct field is of an unsupported type (pointer, array, channel or size-specific number)")
)

// IncomingConfig is a placeholder generic type that exists only to allow us to
// define our innner and exposed value as strongly-typed to the originating
// configuration struct.
type IncomingConfig interface {
	interface{}
}

type contextKey struct{}

var ContextKey = contextKey{}

// Asp is an interface that represents the "callable" interface for
// settings/options.  After creating/initializing with a configuration structure
// (with default values), the methods on the interface allow for loading from
// command-line/config/environment, as well as lower-level access to the created
// viper instance and cobra command.  (In most cases these should not be needed,
// though!)
type Asp[T IncomingConfig] interface {
	Config() *T
	// Command() *cobra.Command
	// Viper() *viper.Viper

	// Execute(handler func(config T, args []string)) error

	Debug()
}

// maybe the incoming config should be *explicitly* the defaults, and a new
// struct created instead of loading in-place?
func Attach[T IncomingConfig](cmd *cobra.Command, config T, envPrefix string) (Asp[T], error) {
	vip := viper.New()
	// cmd := &cobra.Command{
	// 	// Run: func(cmd *cobra.Command, args []string) {
	// 	// 	log.Printf("INSIDE COMMAND! %q", args)
	// 	// },
	// }

	a := &asp[T]{
		// config: config,
		envPrefix: envPrefix,
		vip:       vip,
		cmd:       cmd,
	}
	// log.Printf("initializing config for: %#v", config)

	var err error
	a.baseType, err = a.processStruct(config, "", envPrefix)
	if err != nil {
		return nil, err
	}

	vip.SetConfigName("config")
	// viper.SetConfigType("yaml") // setting the config type takes precedence
	// over the extension, which seems wrong!
	appName := "TEMPDUMMY"
	vip.AddConfigPath(fmt.Sprintf("/etc/%s", appName))
	vip.AddConfigPath(fmt.Sprintf("$HOME/.config/%s", appName))
	vip.AddConfigPath(fmt.Sprintf("$HOME/.%s", appName))
	vip.AddConfigPath(".")

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
func (a *asp[T]) processStruct(s interface{}, parentCanonical string, parentEnv string) (baseType reflect.Type, err error) {
	vip, flags := a.vip, a.cmd.PersistentFlags()

	// log.Printf("initializing struct for: %#v", s)

	// We expect the incoming value to be a struct or a pointer to a struct.
	// Anything else is invalid.
	structVal := reflect.Indirect(reflect.ValueOf(s))
	if structVal.Kind() != reflect.Struct {
		err = ErrConfigMustBeStruct
		return
	}

	baseType = structVal.Type()
	fields := reflect.VisibleFields(baseType)
	// log.Printf("fields: %#v", fields)

	for _, f := range fields {
		// We deal with anonymous (embedded) structs by *not* updating the
		// parentCanonical/parentEnv strings when recursing.  We also need to
		// *not* attempt to process the mirrored sub-elements directly, because
		// we need the canonical structure to get serialized properly.  We can
		// tell if a field is a mirrored embedded field because its "Index"
		// value isn't a length-1 array, it's length 2+.
		if len(f.Index) > 1 {
			continue
		}

		// The attrNoDesc/attrDesc bifurcation exists because *originally* there
		// was no string->map[string]string decoding support, and thus no way to
		// represent maps in an environment variable.  This has been fixed, but
		// it's still a useful concept to have the field description with and
		// without the env-var notation, just in case.
		canonicalName, attrLong, attrShort, attrEnv, attrDescNoEnv := getAttributes(f, parentCanonical, parentEnv)
		attrDesc := fmt.Sprintf("%s (or use %s)", attrDescNoEnv, attrEnv)

		// Rather than setting handled to true in out myriad cases, we default
		// to true, and make sure to set it to false in our default/unhandled
		// cases.
		handled := true
		addBindings := true
		addEnvBinding := true

		// log.Printf("handling field %q : anonymous? %v, index: %v", canonicalName, f.Anonymous, f.Index)

		// This is some very repetitive code!  The flags helpers are typesafe
		// (but not "fluent"), and thus `attrLong`, `attrShort`, and `atrDesc`
		// have to get specified again and again.  Perhaps there's a better
		// library for this, or an opportunity for a new one?
		//
		// ```
		// flag := flags.AsP(attrLong, attrShort, attrDesc)
		// flag.IntP(v2.Int())
		// ```

		// WAIT!!!!!! can we use flags.VarP()?

		// use shortened names purely for concision...
		l, s, d := attrLong, attrShort, attrDesc
		// fieldVal := structVal.Field(i)
		fieldVal := structVal.FieldByIndex(f.Index)
		intf := fieldVal.Interface()

		// switch it := intf.(type) {
		// default:
		// 	log.Printf("interface type is: %v\n\n%v\n\n", it, v2.InterfaceData())
		// }

		// There are special-case types that we handle up-front, falling back to
		// low-level "kinds" only if we need to...
		switch val := intf.(type) {
		case time.Time:
			// create our own time-parsing flag
			flags.VarP(newTimeValue(time.Time{}, new(time.Time)), l, s, d)

		case time.Duration:
			flags.DurationP(l, s, val, d)

		case []time.Duration:
			flags.DurationSliceP(l, s, val, d)

		case bool:
			flags.BoolP(l, s, val, d)

		case int:
			flags.IntP(l, s, val, d)

		case uint:
			flags.UintP(l, s, val, d)

		case string:
			// FUTURE: should we handle "rich" parsing for things like IP
			// addresses, Durations, etc?
			flags.StringP(l, s, val, d)

		case []bool:
			flags.BoolSliceP(l, s, val, d)

		case []int:
			flags.IntSliceP(l, s, val, d)

		case []uint:
			flags.UintSliceP(l, s, val, d)

		// pFlags supports []byte, but the parsing gets confused?
		// maybe that's viper?

		case []byte:
			// This is really []byte!... we'd double-check, but at runtime,
			// all we see is []uint8.
			flags.BytesHexP(l, s, val, d)

		case []string:
			flags.StringSliceP(l, s, val, d)

		case map[string]int:
			// note that viper doesn't parse environment variables for maps
			// correctly (but it *does* seem to do slices!)... we also use
			// `attrDescNoEnv` here!
			flags.StringToIntP(l, s, val, attrDescNoEnv)
			addEnvBinding = false

		case map[string]string:
			// note that viper doesn't parse environment variables for maps
			// correctly (but it *does* seem to do slices!)... we also use
			// `attrDescNoEnv` here!
			flags.StringToStringP(l, s, val, attrDescNoEnv)
			// addEnvBinding = false

		default:
			if f.Type.Kind() == reflect.Struct {
				nestedParentCanonical := parentCanonical
				nestedParentEnv := parentEnv

				if !f.Anonymous {
					nestedParentCanonical = fmt.Sprintf("%s.", canonicalName)
					nestedParentEnv = fmt.Sprintf("%s_", attrEnv)
				}

				a.processStruct(
					intf,
					nestedParentCanonical,
					nestedParentEnv)

				addBindings = false // prevent default flag/config additions!
			} else {
				handled = false
			}
		}

		if !handled {
			log.Printf("unsupported type? %q %#v", f.Type.Kind(), f)
			err = ErrConfigFieldUnsupported
			return
		}

		if addBindings {
			log.Printf("%q, %v, CLI: %q / %q, env: %q, desc: %q",
				canonicalName, f.Type.Kind(),
				attrLong, attrShort, attrEnv, attrDesc)

			// Start pushing into viper?  Note that we're going to need to handle
			// parent paths pretty quickly!
			vip.SetDefault(canonicalName, intf)
			vip.BindPFlag(canonicalName, flags.Lookup(attrLong))
			if addEnvBinding {
				vip.BindEnv(canonicalName, attrEnv)
			}
		}
	}

	return
}

// func (a *asp[T]) Execute(handler func(config T, args []string)) error {
// 	// Set up run-handler for the cobra command...
// 	a.cmd.Run = func(cmd *cobra.Command, args []string) {
// 		log.Printf("BEFORE (INSIDE): %v", a.vip.AllSettings())
// 		// TODO: unmarshal the settings into the expected config type!
// 		cfgVal := reflect.New(a.baseType)
// 		handler(cfgVal.Interface().(T), args)
// 		log.Printf("AFTER (INSIDE): %v", a.vip.AllSettings())
// 	}

// 	log.Printf("BEFORE: %v", a.vip.AllSettings())

// 	// a.cmd.ParseFlags()
// 	err := a.cmd.Execute()
// 	log.Printf("error? %v", err)

// 	log.Printf("AFTER: %v", a.vip.AllSettings())
// 	return err
// }

// func (a *asp[T]) Command() *cobra.Command {
// 	return a.cmd
// }

// func (a *asp[T]) Viper() *viper.Viper {
// 	return a.vip
// }

func (a *asp[T]) Debug() {
	log.Printf("asp.Debug: %#v", a.vip.AllSettings())
}

func (a *asp[T]) Config() *T {
	val := reflect.New(a.baseType)
	log.Printf("created config: %+v", val.Interface())
	cfg := val.Interface().(*T)

	err := a.vip.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
		case *viper.ConfigFileNotFoundError:
			log.Printf("no config file found... perhaps there are environment variables")
		default:
			log.Fatalf("read config error: (%T) %s", err, err.Error())
		}
	}

	err = a.vip.Unmarshal(
		cfg,
		viper.DecodeHook(
			mapstructure.ComposeDecodeHookFunc(
				mapstructure.StringToTimeDurationHookFunc(),
				betterStringToTime(),
				stringToByteSlice(),
				stringToMapStringString(),
				betterStringToSlice(","),
			)))

	if err != nil {
		log.Fatalf("unmarshal config error: %+v", err)
	}

	log.Printf("returning merged config: %+v", cfg)
	return cfg
}

// betterStringToTime handles empty strings as zero time
func betterStringToTime() mapstructure.DecodeHookFuncType {
	return func(f reflect.Type, t reflect.Type, data interface{}) (
		interface{}, error) {
		// log.Printf("attempting to convert string to time? %v --> %v", f, t)
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		// log.Printf("attempting to convert string to time!! %v --> %v (%q)", f, t, data.(string))
		// Convert it by parsing
		return timeConv(data.(string))
	}
}

func stringToByteSlice() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to byte slice? %v", from.Interface())
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf([]byte{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to byte slice! %v, %v", from.Interface(), to.Interface())

		return hex.DecodeString(from.String())
	}
}

func stringToMapStringString() mapstructure.DecodeHookFuncValue {
	return func(from reflect.Value, to reflect.Value) (interface{}, error) {
		// log.Printf("attempting to convert string to byte slice? %v", from.Interface())
		if from.Kind() != reflect.String ||
			to.Type() != reflect.TypeOf(map[string]string{}) {
			return from.Interface(), nil
		}

		// log.Printf("attempting to convert string to map[string]string! %v, %v", from.Interface(), to.Interface())

		// dest := to.Interface().(map[string]string)
		dest := make(map[string]string)
		entries := strings.Split(from.String(), ",")
		for _, entry := range entries {
			log.Printf("converting %q", entry)
			keyVal := strings.SplitN(entry, "=", 2)
			if len(keyVal) != 2 {
				return nil, fmt.Errorf("unexpected map entry %q", entry)
			}
			dest[keyVal[0]] = keyVal[1]
		}

		return dest, nil
	}
}

// betterStringToSlice improves on mapstructure's StringToSliceHookFunc
// by checking for a wrapping "[" and "]" which sometimes happens during flag
// serialization.
func betterStringToSlice(sep string) mapstructure.DecodeHookFuncKind {
	return func(f reflect.Kind, t reflect.Kind, data interface{}) (
		interface{}, error) {
		if f != reflect.String || t != reflect.Slice {
			return data, nil
		}

		raw := data.(string)
		if raw == "" {
			return []string{}, nil
		}

		// check for "[]" around the string...
		if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
			raw = strings.TrimSuffix(strings.TrimPrefix(raw, "["), "]")
		}

		return strings.Split(raw, sep), nil
	}
}
