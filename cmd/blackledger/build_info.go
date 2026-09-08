package main

import "runtime/debug"

type buildIdentity struct {
	Revision string `json:"revision"`
	Modified *bool  `json:"modified"`
}

func readBuildIdentity() buildIdentity {
	identity := buildIdentity{Revision: "unknown"}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return identity
	}
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			identity.Revision = setting.Value
		case "vcs.modified":
			value := setting.Value == "true"
			identity.Modified = &value
		}
	}
	return identity
}

var runningBuild = readBuildIdentity()
