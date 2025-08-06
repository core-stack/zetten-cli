package mirror

import (
	"github.com/kardianos/service"
)

type ZettenService struct{}

func (p *ZettenService) Start(s service.Service) error {
	go StartMirrorService("mirror.yaml")
	return nil
}

func (p *ZettenService) Stop(s service.Service) error {
	StopMirrorService()
	return nil
}
