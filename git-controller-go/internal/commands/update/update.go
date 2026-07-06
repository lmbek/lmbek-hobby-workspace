package update

import (
	"fmt"
	"workspace/git-controller/internal/commands/fetch"
	"workspace/git-controller/internal/commands/pull"
	"workspace/git-controller/internal/commands/status"
	"workspace/git-controller/internal/ui"
)

// Run performs a full update cycle: fetch all remotes, pull latest changes,
// and display the repository status dashboard.
func Run() error {
	ui.Header("Update Repositories")

	if err := fetch.Run(); err != nil {
		return fmt.Errorf("update failed during fetch: %w", err)
	}

	if err := pull.Run(); err != nil {
		return fmt.Errorf("update failed during pull: %w", err)
	}

	if err := status.Run(); err != nil {
		return fmt.Errorf("update failed during status: %w", err)
	}

	ui.Success("Update complete!")
	return nil
}
