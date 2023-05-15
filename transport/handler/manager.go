package handler

import "test/service"

type Manager struct {
	service *service.Manager
}

func NewManagerHandler(service *service.Manager) *Manager {
	return &Manager{
		service: service,
	}
}
