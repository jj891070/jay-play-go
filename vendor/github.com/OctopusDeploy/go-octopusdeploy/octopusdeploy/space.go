package octopusdeploy

import "github.com/go-playground/validator/v10"

type Spaces struct {
	Items []*Space `json:"Items"`
	PagedResults
}

type Space struct {
	Description              string   `json:"Description,omitempty"`
	IsDefault                bool     `json:"IsDefault,omitempty"`
	Name                     string   `json:"Name" validate:"required,max=20"`
	SpaceManagersTeamMembers []string `json:"SpaceManagersTeamMembers,omitempty"`
	SpaceManagersTeams       []string `json:"SpaceManagersTeams,omitempty"`
	TaskQueueStopped         bool     `json:"TaskQueueStopped,omitempty"`

	resource
}

// NewSpace initializes a Space with a name.
func NewSpace(name string) *Space {
	return &Space{
		Name:     name,
		resource: *newResource(),
	}
}

// Validate checks the state of the space and returns an error if
// invalid.
func (s *Space) Validate() error {
	return validator.New().Struct(s)
}
