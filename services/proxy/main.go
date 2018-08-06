package proxy

import (
	"github.com/ironstar-io/tokaido/services/unison"
	"github.com/ironstar-io/tokaido/system/ssl"
)

// Setup ...
func Setup() {
	buildDirectories()

	ssl.Configure(getProxyClientTLSDir())

	GenerateProxyDockerCompose()
	DockerComposeUp()

	unison.CreateOrUpdatePrf(UnisonPort(), "proxy", getProxyClientDir())

	// RebuildNginxConfigFile() // Every run
	// RestartProxyContainer()  // Every run

	// AppendHostsfile() // Check then append
	// StartUnisonSyncSvc() // If not already started?
}
