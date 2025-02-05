package acceptance

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// RequiresImportErrorStep returns a Test Step which expects a Requires Import
// error to be returned when running this step
func (td TestData) RequiresImportErrorStep(configBuilder func(data TestData) string) resource.TestStep {
	config := configBuilder(td)
	return resource.TestStep{
		Config:      config,
		ExpectError: RequiresImportError(td.ResourceType),
	}
}

func RequiresImportError(resourceName string) *regexp.Regexp {
	message := "to be managed via Terraform this resource needs to be imported into the State. Please see the resource documentation for %q for more information."
	return regexp.MustCompile(fmt.Sprintf(message, resourceName))
}

// ImportStep returns a Test Step which Imports the Resource, optionally
// ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStep(importStateIdFunc resource.ImportStateIdFunc, importStateCheckFunc resource.ImportStateCheckFunc, ignore ...string) resource.TestStep {
	return td.ImportStepFor(td.ResourceName, importStateIdFunc, importStateCheckFunc, ignore...)
}

// ImportStepFor returns a Test Step which Imports a given resource by name,
// optionally ignoring any fields which may not be imported (for example, as they're
// not returned from the API)
func (td TestData) ImportStepFor(resourceName string, importStateIdFunc resource.ImportStateIdFunc, importStateCheckFunc resource.ImportStateCheckFunc, ignore ...string) resource.TestStep {
	if strings.HasPrefix(resourceName, "data.") {
		return resource.TestStep{
			ResourceName: resourceName,
			SkipFunc: func() (bool, error) {
				return false, fmt.Errorf("Data Sources (%q) do not support import - remove the ImportStep / ImportStepFor`", resourceName)
			},
		}
	}

	step := resource.TestStep{
		ResourceName:      resourceName,
		ImportState:       true,
		ImportStateVerify: false,
		ImportStateIdFunc: importStateIdFunc,
		ImportStateCheck:  importStateCheckFunc,
	}

	if len(ignore) > 0 {
		step.ImportStateVerifyIgnore = ignore
	}

	return step
}
