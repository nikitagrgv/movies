package configs

import _ "embed"

//go:embed watch_servers.yaml
var WatchServersRawConfig []byte
