package flags

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Flag struct {
	Name      string
	Shorthand string
	Value     string
	Usage     string
	Env       string
	Required  bool
}


// StringPersistentFlag defines a string persistent flag for a cobra command with a pointer, a specified name, a shorthand letter, default value, usage description, name of env variable incur value of flag. The argument p points to a string variable in which to store the value of the flag.
func StringPersistentFlag(cmd *cobra.Command, p *string, name string, shorthand string, value string, usage string, env string, required bool) {
	if env == "" {
		env = name
	}

	viper.BindEnv(name, env)
	cmd.PersistentFlags().StringVarP(p, name, shorthand, value, usage)
	if required {
		cmd.MarkPersistentFlagRequired(name)
	}
}

