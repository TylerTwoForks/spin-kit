package sf

// SFResponse is the generic envelope returned by sf CLI with --json.
type SFResponse[T any] struct {
	Status  int    `json:"status"`
	Result  T      `json:"result"`
	Message string `json:"message"`
}

type OrgListResult struct {
	Other          []OrgInfo `json:"other"`
	Sandboxes      []OrgInfo `json:"sandboxes"`
	NonScratchOrgs []OrgInfo `json:"nonScratchOrgs"`
	DevHubs        []OrgInfo `json:"devHubs"`
	ScratchOrgs    []OrgInfo `json:"scratchOrgs"`
}

type OrgInfo struct {
	Alias           string `json:"alias"`
	Username        string `json:"username"`
	OrgID           string `json:"orgId"`
	ConnectedStatus string `json:"connectedStatus"`
	InstanceURL     string `json:"instanceUrl"`
	IsDevHub        bool   `json:"isDevHub"`
}

type SandboxProcessRecord struct {
	SandboxName  string `json:"SandboxName"`
	Status       string `json:"Status"`
	CopyProgress int    `json:"CopyProgress"`
}

type DataQueryResult struct {
	Records   []SandboxProcessRecord `json:"records"`
	TotalSize int                    `json:"totalSize"`
	Done      bool                   `json:"done"`
}

type RefreshResult struct {
	SandboxProcessID string `json:"Id"`
	Status           string `json:"Status"`
	SandboxName      string `json:"SandboxName"`
}

type LoginResult struct {
	OrgID    string `json:"orgId"`
	URL      string `json:"url"`
	Username string `json:"username"`
}

// TUI message types produced by sf commands.

type OrgListMsg struct {
	Orgs []OrgInfo
	Err  error
}

type SandboxStatusMsg struct {
	Records []SandboxProcessRecord
	Err     error
}

type SandboxRefreshMsg struct {
	Name   string
	Result string
	Err    error
}

type WebLoginMsg struct {
	Alias string
	URL   string
	Err   error
}

type OrgLogoutMsg struct {
	Alias string
	Err   error
}

type CommandErrorMsg struct {
	Err error
}
