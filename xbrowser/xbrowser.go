package xbrowser

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/pkg/browser"

	"oss.terrastruct.com/util-go/xos"
)

func Open(ctx context.Context, env *xos.Env, url string) error {
	browserEnv := env.Getenv("BROWSER")
	if browserEnv == "0" || browserEnv == "false" {
		return nil
	}
	if browserEnv != "" && browserEnv != "1" && browserEnv != "true" {
		browserSh := fmt.Sprintf(`%s "$1"`, browserEnv)
		cmd := exec.CommandContext(ctx, "sh", "-sc", browserSh, "--", url)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to run %v (out: %q): %w", cmd.Args, out, err)
		}
		return nil
	}
	return browser.OpenURL(url)
}
