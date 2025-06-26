package containerd

import (
	"os"

	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/api"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/daemon"
	"github.com/awslabs/amazon-eks-ami/nodeadm/internal/util"
	"github.com/pelletier/go-toml/v2"
)

const ContainerdDaemonName = "containerd"
const SociSnapshotterSocketName = "soci-snapshotter.socket"

var _ daemon.Daemon = &containerd{}

type containerd struct {
	daemonManager daemon.DaemonManager
}

func NewContainerdDaemon(daemonManager daemon.DaemonManager) daemon.Daemon {
	return &containerd{
		daemonManager: daemonManager,
	}
}

func (cd *containerd) Configure(cfg *api.NodeConfig) error {
	if err := writeBaseRuntimeSpec(cfg); err != nil {
		return err
	}
	if err := writeSnapshotterConfig(cfg); err != nil {
		return err
	}
	return writeContainerdConfig(cfg)
}

func (cd *containerd) EnsureRunning(cfg *api.NodeConfig) error {
	if err := cd.daemonManager.StartDaemon(ContainerdDaemonName); err != nil {
		return err
	}
	if api.IsFeatureEnabled(api.FastContainerImagePull, cfg.Spec.FeatureGates) {
		return cd.daemonManager.StartDaemon(SociSnapshotterSocketName)
	}
	return nil
}

func (cd *containerd) PostLaunch(cfg *api.NodeConfig) error {
	return nil
}

func (cd *containerd) Name() string {
	return ContainerdDaemonName
}
