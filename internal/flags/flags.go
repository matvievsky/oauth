package flags

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// []Name/Usage/DefaultValue
type Flags []string

func Init() Flags {
	return Flags{}
}

func (flags Flags) WithFlag(name, value, usage string) Flags {
	return append(flags, strings.Join([]string{name, value, usage}, "/"))
}

// Bind добавляет флаги и привязывает их к Viper
func (flags Flags) Bind(cmd *cobra.Command) error {
	for _, flag := range flags {
		nameDefaultValueUsage := strings.Split(flag, "/")
		if len(nameDefaultValueUsage) != 3 {
			return errors.New("invalid flag format: expected 'name/value/usage'")
		}

		prefix := viper.GetEnvPrefix()
		if prefix != "" {
			prefix += "_"
		}

		value := viper.GetString(strings.Replace(fmt.Sprintf("%s%s", prefix, nameDefaultValueUsage[0]), "-", "_", -1))
		if value == "" {
			value = nameDefaultValueUsage[1]
		}

		cmd.Flags().String(nameDefaultValueUsage[0], value, nameDefaultValueUsage[2])
		viper.BindPFlag(nameDefaultValueUsage[0], cmd.Flags().Lookup(nameDefaultValueUsage[0]))
	}

	return nil
}
