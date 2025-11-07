package version

import (
	"runtime/debug"
	"time"
)

var (
	Version string
	Revision   = "unknown"
	DirtyTree  = true
	LastCommit time.Time
)

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	Version = info.Main.Version
	for _, kv := range info.Settings {
		if kv.Value == "" {
			continue
		}
		switch kv.Key {
		case "vcs.revision":
			Revision = kv.Value
		case "vcs.time":
			LastCommit, _ = time.Parse(time.RFC3339, kv.Value)
		case "vcs.modified":
			DirtyTree = kv.Value == "true"
		}
	}
}
