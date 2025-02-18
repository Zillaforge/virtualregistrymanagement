package execution

import (
	cnt "VirtualRegistryManagement/constants"
	"VirtualRegistryManagement/modules/openstack/cinder"
	"VirtualRegistryManagement/modules/openstack/glance"

	"go.uber.org/zap"
)

const (
	_glanceResource = "glance"
	_cinderResource = "cinder"
)

// Glance ...
func (p *Connection) Glance(projectID string) *glance.Glance {
	pid := p.Pid(projectID)
	pc, err := p.providerClient(&pid)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.providerClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &glance.Glance{}
	}
	sc, err := p.serviceClient(pc, _glanceResource)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.serviceClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &glance.Glance{}
	}

	return glance.New(p.namespace, pid, sc)
}

// Cinder ...
func (p *Connection) Cinder(projectID string) *cinder.Cinder {
	pid := p.Pid(projectID)
	pc, err := p.providerClient(&pid)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.providerClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &cinder.Cinder{}
	}
	sc, err := p.serviceClient(pc, _cinderResource)
	if err != nil {
		zap.L().With(
			zap.String(cnt.Module, "p.serviceClient(...)"),
			zap.Any("connection", p),
			zap.String("project-id", projectID),
			zap.String("opstk-project-id", pid),
		).Error(err.Error())
		return &cinder.Cinder{}
	}
	return cinder.New(p.namespace, p.Pid(projectID), sc)
}
