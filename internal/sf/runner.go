package sf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func runSF(args ...string) ([]byte, error) {
	cmd := exec.Command("sf", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// If we got stdout, return it for JSON parsing even on non-zero exit.
	// sf CLI commonly returns valid JSON alongside a non-zero exit code
	// (e.g., warnings about expired auth tokens).
	if stdout.Len() > 0 {
		return stdout.Bytes(), nil
	}

	if err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return nil, fmt.Errorf("sf %s: %s", strings.Join(args, " "), errMsg)
		}
		return nil, fmt.Errorf("sf %s: %w", strings.Join(args, " "), err)
	}

	return nil, nil
}

func ListOrgs() tea.Cmd {
	return func() tea.Msg {
		out, err := runSF("org", "list", "--json")
		if err != nil {
			return OrgListMsg{Err: err}
		}

		var resp SFResponse[OrgListResult]
		if err := json.Unmarshal(out, &resp); err != nil {
			return OrgListMsg{Err: fmt.Errorf("parse org list: %w", err)}
		}

		// Deduplicate across categories -- the same org appears under
		// multiple keys (e.g., a sandbox appears in both "sandboxes"
		// and "nonScratchOrgs").
		seen := make(map[string]bool)
		var orgs []OrgInfo
		for _, list := range [][]OrgInfo{
			resp.Result.DevHubs,
			resp.Result.Other,
			resp.Result.Sandboxes,
			resp.Result.NonScratchOrgs,
			resp.Result.ScratchOrgs,
		} {
			for _, o := range list {
				if o.OrgID != "" && !seen[o.OrgID] {
					seen[o.OrgID] = true
					orgs = append(orgs, o)
				}
			}
		}

		return OrgListMsg{Orgs: orgs}
	}
}

func SandboxStatus(prodAlias string) tea.Cmd {
	return func() tea.Msg {
		query := "SELECT Id, SandboxName, Status, CopyProgress FROM SandboxProcess WHERE Status IN ('0', '2', '3', '4', '5') ORDER BY Status DESC"
		out, err := runSF(
			"data", "query",
			"--query", query,
			"--target-org", prodAlias,
			"--use-tooling-api",
			"--json",
		)
		if err != nil {
			return SandboxStatusMsg{Err: err}
		}

		var resp SFResponse[DataQueryResult]
		if err := json.Unmarshal(out, &resp); err != nil {
			return SandboxStatusMsg{Err: fmt.Errorf("parse sandbox status: %w", err)}
		}

		return SandboxStatusMsg{Records: resp.Result.Records}
	}
}

// RefreshSandbox suspends the TUI and runs the refresh command interactively.
// The sf CLI shows sandbox details and prompts for confirmation before starting.
func RefreshSandbox(name, prodAlias string) tea.Cmd {
	c := exec.Command("sf", "org", "refresh", "sandbox",
		"--name", name,
		"--target-org", prodAlias,
	)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		if err != nil {
			return SandboxRefreshMsg{Name: name, Err: err}
		}
		return SandboxRefreshMsg{
			Name:   name,
			Result: fmt.Sprintf("Sandbox %q refresh completed.", name),
		}
	})
}

func WebLogin(alias string) tea.Cmd {
	return func() tea.Msg {
		out, err := runSF("org", "login", "web", "--alias", alias, "--json")
		if err != nil {
			return WebLoginMsg{Alias: alias, Err: err}
		}

		var resp SFResponse[LoginResult]
		if err := json.Unmarshal(out, &resp); err != nil {
			return WebLoginMsg{Alias: alias, Err: fmt.Errorf("parse login result: %w", err)}
		}

		return WebLoginMsg{Alias: alias, URL: resp.Result.URL}
	}
}

func OpenOrg(alias string) tea.Cmd {
	return func() tea.Msg {
		_, err := runSF("org", "open", "--target-org", alias)
		if err != nil {
			return OrgOpenMsg{Alias: alias, Err: err}
		}
		return OrgOpenMsg{Alias: alias}
	}
}

func LogoutOrg(alias string) tea.Cmd {
	return func() tea.Msg {
		_, err := runSF("org", "logout", "--target-org", alias, "--no-prompt")
		if err != nil {
			return OrgLogoutMsg{Alias: alias, Err: err}
		}
		return OrgLogoutMsg{Alias: alias}
	}
}

func ReconnectSandbox(alias string) tea.Cmd {
	return func() tea.Msg {
		out, err := runSF(
			"org", "login", "web",
			"--instance-url", "https://test.salesforce.com",
			"--alias", alias,
			"--json",
		)
		if err != nil {
			return WebLoginMsg{Alias: alias, Err: err}
		}

		var resp SFResponse[LoginResult]
		if err := json.Unmarshal(out, &resp); err != nil {
			return WebLoginMsg{Alias: alias, Err: fmt.Errorf("parse reconnect result: %w", err)}
		}

		return WebLoginMsg{Alias: alias, URL: resp.Result.URL}
	}
}
