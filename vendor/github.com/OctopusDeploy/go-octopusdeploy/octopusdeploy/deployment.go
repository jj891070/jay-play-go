package octopusdeploy

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Deployment struct {
	Changes                  []*ReleaseChanges `json:"Changes"`
	ChangesMarkdown          string            `json:"ChangesMarkdown,omitempty"`
	ChannelID                string            `json:"ChannelId,omitempty"`
	Comments                 string            `json:"Comments,omitempty"`
	Created                  *time.Time        `json:"Created,omitempty"`
	DeployedBy               string            `json:"DeployedBy,omitempty"`
	DeployedByID             string            `json:"DeployedById,omitempty"`
	DeployedToMachineIDs     []string          `json:"DeployedToMachineIds"`
	DeploymentProcessID      string            `json:"DeploymentProcessId,omitempty"`
	EnvironmentID            string            `json:"EnvironmentId" validate:"required"`
	ExcludedMachineIDs       []string          `json:"ExcludedMachineIds"`
	FailureEncountered       bool              `json:"FailureEncountered,omitempty"`
	ForcePackageDownload     bool              `json:"ForcePackageDownload,omitempty"`
	ForcePackageRedeployment bool              `json:"ForcePackageRedeployment,omitempty"`
	FormValues               map[string]string `json:"FormValues,omitempty"`
	ManifestVariableSetID    string            `json:"ManifestVariableSetId,omitempty"`
	Name                     string            `json:"Name,omitempty"`
	ProjectID                string            `json:"ProjectId,omitempty"`
	QueueTime                *time.Time        `json:"QueueTime,omitempty"`
	QueueTimeExpiry          *time.Time        `json:"QueueTimeExpiry,omitempty"`
	ReleaseID                string            `json:"ReleaseId" validate:"required"`
	SkipActions              []string          `json:"SkipActions"`
	SpaceID                  string            `json:"SpaceId,omitempty"`
	SpecificMachineIDs       []string          `json:"SpecificMachineIds"`
	TaskID                   string            `json:"TaskId,omitempty"`
	TenantID                 string            `json:"TenantId,omitempty"`
	TentacleRetentionPeriod  *RetentionPeriod  `json:"TentacleRetentionPeriod,omitempty"`
	UseGuidedFailure         bool              `json:"UseGuidedFailure,omitempty"`

	resource
}

// Deployments defines a collection of deployment instances with built-in
// support for paged results.
type Deployments struct {
	Items []*Deployment `json:"Items"`
	PagedResults
}

// NewDeployment initializes a deployment with a name, environment ID, and
// release ID.
func NewDeployment(environmentID string, releaseID string) *Deployment {
	return &Deployment{
		EnvironmentID: environmentID,
		ReleaseID:     releaseID,
		resource:      *newResource(),
	}
}

// Validate checks the state of the deployment and returns an error if invalid.
func (d *Deployment) Validate() error {
	return validator.New().Struct(d)
}
