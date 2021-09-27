package octopusdeploy

import (
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
)

// account is the embedded struct used for all accounts.
type account struct {
	AccountType            AccountType `validate:"required,oneof=None UsernamePassword SshKeyPair AzureSubscription AzureServicePrincipal AmazonWebServicesAccount AmazonWebServicesRoleAccount GoogleCloudAccount Token"`
	Description            string
	EnvironmentIDs         []string
	Name                   string `validate:"required,notblank,notall"`
	SpaceID                string `validate:"omitempty,notblank"`
	TenantedDeploymentMode TenantedDeploymentMode
	TenantIDs              []string
	TenantTags             []string

	resource
}

// newAccount creates and initializes an account.
func newAccount(name string, accountType AccountType) *account {
	return &account{
		AccountType:            accountType,
		EnvironmentIDs:         []string{},
		Name:                   name,
		TenantedDeploymentMode: TenantedDeploymentMode("Untenanted"),
		TenantIDs:              []string{},
		TenantTags:             []string{},
		resource:               *newResource(),
	}
}

// GetAccountType returns the type of this account.
func (a *account) GetAccountType() AccountType {
	return a.AccountType
}

// GetDescription returns the description of the account.
func (a *account) GetDescription() string {
	return a.Description
}

func (a *account) GetEnvironmentIDs() []string {
	return a.EnvironmentIDs
}

// GetName returns the name of the account.
func (a *account) GetName() string {
	return a.Name
}

// GetSpaceID returns the space ID of this account.
func (a *account) GetSpaceID() string {
	return a.SpaceID
}

func (a *account) GetTenantedDeploymentMode() TenantedDeploymentMode {
	return a.TenantedDeploymentMode
}

func (a *account) GetTenantIDs() []string {
	return a.TenantIDs
}

func (a *account) GetTenantTags() []string {
	return a.TenantTags
}

// SetDescription sets the description of the account.
func (a *account) SetDescription(description string) {
	a.Description = description
}

// SetName sets the name of the account.
func (a *account) SetName(name string) {
	a.Name = name
}

// SetSpaceID sets the space ID of this account.
func (a *account) SetSpaceID(spaceID string) {
	a.SpaceID = spaceID
}

// Validate checks the state of the account and returns an error if
// invalid.
func (a *account) Validate() error {
	v := validator.New()
	err := v.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return err
	}
	err = v.RegisterValidation("notall", NotAll)
	if err != nil {
		return err
	}
	return v.Struct(a)
}

var _ IAccount = &account{}
