package config

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"reflect"
)

func bindFlags(cfg any, cmd *cobra.Command, prefix string) {
	val := reflect.ValueOf(cfg).Elem()
	typ := reflect.TypeOf(cfg).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Tag.Get("mapstructure")
		desc := field.Tag.Get("desc")

		if desc == "" {
			continue
		}

		flagName := fieldName
		if prefix != "" {
			flagName = prefix + "." + fieldName
		}

		fieldVal := val.Field(i)
		if field.Type.Kind() == reflect.Struct {
			bindFlags(fieldVal.Addr().Interface(), cmd, flagName)
			continue
		}
		switch field.Type.Kind() {
		case reflect.String:
			defaultVal := fieldVal.String()
			cmd.PersistentFlags().String(flagName, defaultVal, desc)
		case reflect.Int:
			defaultVal := int(fieldVal.Int())
			cmd.PersistentFlags().Int(flagName, defaultVal, desc)
		case reflect.Bool:
			defaultVal := fieldVal.Bool()
			cmd.PersistentFlags().Bool(flagName, defaultVal, desc)
		default:
			fmt.Printf("[BindCobraFlags] Unsupported type: %s\n", field.Type.Kind())
			continue
		}
		flag := cmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			fmt.Printf("[BindCobraFlags] Flag for %s is nil, cannot bind to viper\n", flagName)
			continue
		}
		if err := viper.BindPFlag(flagName, flag); err != nil {
			fmt.Printf("[BindCobraFlags] Orrured an error when BindPFlag to viper: %v\n", err)
		}
	}
}
