package habitica

// Common structures that can be reused across multiple domains.

// UUID is an alias for IDs represented as UUID strings.
type UUID string

// Timestamp is the ISO8601 timestamp as typically returned by the Habitica API.
// We keep it as a string to avoid making assumptions about the exact format.
type Timestamp string

