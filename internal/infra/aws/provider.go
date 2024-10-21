package aws

import "github.com/google/wire"

// ProviderSetAwsBucket is providers.
var ProviderSetAwsBucket = wire.NewSet(
	NewS3Client,
)
