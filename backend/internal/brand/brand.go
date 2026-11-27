// brand.go defines the placeholder product identity used when site settings do
// not supply a value.

package brand

const (
	// Name is the default site name shown in titles, emails, and setup flows.
	Name = "AI Relay"
	// TitleSuffix is appended to the browser title after the site name.
	TitleSuffix = "AI API Gateway"
	// Subtitle is the default marketing subtitle used by admin settings.
	Subtitle = "Subscription to API Conversion Platform"
	// TOTPIssuer is the fallback issuer label for authenticator apps.
	TOTPIssuer = Name
	// MCPName is the implementation name exposed by the MCP server.
	MCPName = Name
	// WebAuthnDisplayName is the default passkey prompt label.
	WebAuthnDisplayName = Name
	// ReleaseRepository is the GitHub repository used for release checks.
	ReleaseRepository = "Sworol/TokenRouter"
	// DockerImage is the container image shown by admin version details.
	DockerImage = "ghcr.io/sworol/tokenrouter"
	// DocsURL is the documentation link shown by admin payment settings.
	DocsURL = "https://github.com/Sworol/TokenRouter/tree/main/docs"
)
