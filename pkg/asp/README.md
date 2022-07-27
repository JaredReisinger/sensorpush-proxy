# asp - Automatic Settings Processor

asp, the Automatic Settings Provider, an opinionated companion for viper and
cobra.

## Embedded anonymous structs

It's reasonable to want to compose app configuration out of sub-parts, and embed
anonymous structs to make those values transparently available at runtime.  This
is `asp`'s default behavior with anonymous structs, but there are a few caveats
about which you need to be aware:

  - The config, flag, and environment variable names for an anonymous embedded struct _**do not**_ include the name of the embedded struct itself.  If you want to include the struct name, simply don't make it an anonymous embed, and ignore the rest of this section entirely.

  - When writing the anonymous embedded struct reference, you need to include a `mapstructure` tag to "squash" the members to the parent map for deserialization.  It would be ideal if `asp` could somehow default this for you, but it cannot.  (I wish it were the default for `mapstructure`, but alas, it is not.)  For example:

    ```go
    type CommonFields struct {
        FirstName string
        LastName  string
    }

    type Config struct {
        CommonFields `mapstructure:",squash"` // <==== this is the needed tag!
        More         string
    }
    ```

    Without the `mapstructure` "squash" option, the `viper` configuration file values won't map to the final config object correctly.

  - When you write a config file (in YAML, TOML, or what-have-you), you must write as though the embedded fields exist directly in the parent:

    ```yaml
    # config.yaml
    firstName: John
    lastName: Doe
    more: used for an unknown person

    # --*NOT*--
    # commonFields:
    #   firstName: John
    #   lastName: Doe
    # more: ...
    ```

  - As per standard Go behavior, however, while you will be able to "read" values from your loaded configuration using the embedded struct field shorthand (`config.FirstName`), you _cannot_ programmatically construct your config that way.  In this case, for example to create defaults, you will need to provide the embedded struct explicitly:

    ```go
    var Default = Config{
    	CommonFields: CommonFields{
    		FirstName: "Mia",
    	},
    }
    ```

The examples above will result in:

```text
--first-name string   sets the FirstName value (or use APP_FIRSTNAME) (default "Mia")
--last-name string    sets the LastName value (or use APP_LASTNAME)
--more string         sets the More value (or use APP_MORE)
```

## 2022-07-26

After playing around with some auto viper/cobra generation, I realized that generating a full cobra `Command` was going to be somewhat awkward, as was the idea of wrapping the `cobra.Command.Execute()` so as to pass all the needed stuff to the handler.  Especially in the context of creating one or more utilities that use (some) common config, it really makes more sense to let the tools continue to define the base `cobra.Command`, but to provide helpers that help augment that command with persistent flags and provide a type-safe viper wrapper.  This solves the 80:20 problem.

So, the new model is to let the user/caller define the `cobra.Command` object themselves, and then call an `asp` helper that can take a config struct and process it so as to add some `cobra.Command` flags, and also set up `viper`'s config/flag/env-var settings, finally exposing a simple "load" that will populate/create a config struct.  This "load" method can be called from inside any of the `cobra.Command` "Run" handlers.

Note that the types supported are *heavily* influenced by the pFlags capabilities.


## _OLDER_

Why?  I really like viper and cobra for building 12-factor style applications, but there's a lot of overhead incurred (in lines-of-code) in creating the config, command-line, and environment variable settings for each option.  The goal of `asp` is the (a) reduce the redundant boilerplate by concisely defining all of the necessary information in the config struct itself, (b) to encourage good practices by ensuring that *every* option has config, command-line, *and* environment variable representation, and (c) to avoid possible typos that using string-based configuration lookups can cause--Go can't tell that `viper.Get("sommeSetting")` is misspelled at compile time... but it *can* tell that `config.sommeSetting` is invalid if the struct defines the member as `someSetting`.

This is done by driving the command-line and environment variable settings from a tagged struct that's used to define the config file format.

```
type Config struct {
    Host string `asp:""`
}
```
