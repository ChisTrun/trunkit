package publish

import (
	"fmt"
	"os"

	"trunkit/internal/config"
	"trunkit/internal/emitter/generate/js"
	osutil "trunkit/internal/util/os"

	"github.com/briandowns/spinner"
)

func publishJs(cfg *config.GenerateConfig, npmCommands []string, s *spinner.Spinner) {
	outPaths := js.GenerateJs(cfg)
	for _, outPath := range outPaths {
		for _, cmd := range npmCommands {
			s.Prefix = cmd
			s.Start()
			_, err := osutil.Exec([]string{
				fmt.Sprintf("cd %s", outPath),
				cmd,
			})
			if err != nil {
				fmt.Println(cmd, "failed", err)
				os.Exit(1)
			}
		}
	}
}
