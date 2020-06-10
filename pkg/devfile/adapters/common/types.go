package common

import (
	devfileParser "github.com/openshift/odo/pkg/devfile/parser"
	"github.com/openshift/odo/pkg/devfile/parser/data/common"
	"github.com/openshift/odo/pkg/envinfo"
)

// AdapterContext is a construct that is common to all adapters
type AdapterContext struct {
	ComponentName string                   // ComponentName is the odo component name, it is NOT related to any devfile components
	Context       string                   // Context is the given directory containing the source code and configs
	Devfile       devfileParser.DevfileObj // Devfile is the object returned by the Devfile parser
}

// DevfileVolume is a struct for Devfile volume that is common to all the adapters
type DevfileVolume struct {
	Name          string
	ContainerPath string
	Size          string
}

// ComponentVolumesPair is a struct to hold the relationship between components and their volumes
// Use splices of this struct instead of a map because it is important to maintain order since iteration order of maps is not guaranteed.
// When maps were used the ordering of the volumes could change in the pod spec causing an unnecessary rollout of the deployment which can be slow.
type ComponentVolumesPair struct {
	ComponentAlias string
	Volumes        []DevfileVolume
}

// VolumePVCPair is a struct to hold the relationship between volumes and their PVCs
// Use splices of this struct instead of a map because it is important to maintain order since iteration order of maps is not guaranteed.
// When maps were used the ordering of the volumes could change in the pod spec causing an unnecessary rollout of the deployment which can be slow.
type VolumePVCPair struct {
	Volume string
	PVC    string
}

// Storage is a struct that is common to all the adapters
type Storage struct {
	Name   string
	Volume DevfileVolume
}

// PushParameters is a struct containing the parameters to be used when pushing to a devfile component
type PushParameters struct {
	Path              string                  // Path refers to the parent folder containing the source code to push up to a component
	WatchFiles        []string                // Optional: WatchFiles is the list of changed files detected by odo watch. If empty or nil, odo will check .odo/odo-file-index.json to determine changed files
	WatchDeletedFiles []string                // Optional: WatchDeletedFiles is the list of deleted files detected by odo watch. If empty or nil, odo will check .odo/odo-file-index.json to determine deleted files
	IgnoredFiles      []string                // IgnoredFiles is the list of files to not push up to a component
	ForceBuild        bool                    // ForceBuild determines whether or not to push all of the files up to a component or just files that have changed, added or removed.
	Show              bool                    // Show tells whether the devfile command output should be shown on stdout
	DevfileInitCmd    string                  // DevfileInitCmd takes the init command through the command line and overwrites devfile init command
	DevfileBuildCmd   string                  // DevfileBuildCmd takes the build command through the command line and overwrites devfile build command
	DevfileRunCmd     string                  // DevfileRunCmd takes the run command through the command line and overwrites devfile run command
	EnvSpecificInfo   envinfo.EnvSpecificInfo // EnvSpecificInfo contains infomation of env.yaml file
}

// SyncParameters is a struct containing the parameters to be used when syncing a devfile component
type SyncParameters struct {
	PushParams      PushParameters
	CompInfo        ComponentInfo
	PodChanged      bool
	ComponentExists bool
}

// ComponentInfo is a struct that holds information about a component i.e.; pod name, container name
type ComponentInfo struct {
	PodName       string
	ContainerName string
}

// PushCommandsMap stores the commands to be executed as per their types.
type PushCommandsMap map[common.DevfileCommandGroupType]common.DevfileCommand

// NewPushCommandMap returns the instance of PushCommandsMap
func NewPushCommandMap() PushCommandsMap {
	return make(map[common.DevfileCommandGroupType]common.DevfileCommand)
}
