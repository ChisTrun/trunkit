package publish

import (
	"fmt"
	"os"

	"mykit/internal/config"
	"mykit/internal/emitter/generate/js"
	osutil "mykit/internal/util/os"

	"github.com/briandowns/spinner"
)

func publishConnectJs(cfg *config.GenerateConfig, npmCommands []string, s *spinner.Spinner) {
	connectOutPaths := js.GenerateConnectJs(cfg)
	for _, connectOutPath := range connectOutPaths {
		for _, cmd := range npmCommands {
			s.Prefix = cmd
			s.Start()
			_, err := osutil.Exec([]string{
				fmt.Sprintf("cd %s", connectOutPath),
				cmd,
			})
			if err != nil {
				fmt.Println(cmd, "failed", err)
				os.Exit(1)
			}
		}
	}

}
