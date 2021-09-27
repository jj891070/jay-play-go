package octopusdeploy

import (
	"github.com/dghubble/sling"
)

type rootService struct {
	service
}

func newRootService(sling *sling.Sling, uriTemplate string) *rootService {
	return &rootService{
		service: newService(ServiceRootService, sling, uriTemplate),
	}
}

func (s rootService) Get() (*RootResource, error) {
	path, err := getPath(s)
	if err != nil {
		return nil, err
	}

	resp, err := apiGet(s.getClient(), new(RootResource), path)
	if err != nil {
		return nil, err
	}

	return resp.(*RootResource), nil
}

var _ IService = &rootService{}
