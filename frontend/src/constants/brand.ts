// brand.ts centralizes the placeholder product identity shown by the frontend.
//
// Replace the values here when the final brand is chosen. Technical identifiers
// such as provider keys, storage keys, and protocol headers stay unchanged so
// existing clients and deployments keep working.

export const BRAND = {
  name: 'AI Relay',
  titleSuffix: 'AI API Gateway',
  subtitle: 'Subscription to API Conversion Platform',
  repositoryUrl: 'https://github.com/Sworol/TokenRouter',
  dockerImage: 'ghcr.io/sworol/tokenrouter',
} as const
