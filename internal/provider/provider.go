package provider

import (
	"context"

	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/srikanthbhandary-teach/my-client"
)

var (
	_ provider.Provider = &myprovider{}
)

// hashicupsProviderModel maps provider schema data to a Go type.
type hashicupsProviderModel struct {
	ApiKey  types.String `tfsdk:"apikey"`
	BaseUrl types.String `tfsdk:"baseurl"`
}

type myprovider struct {
	version string
}

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &myprovider{
			version: version,
		}
	}
}

func (p *myprovider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "myprovider"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *myprovider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Optional: true,
			},
			"baseurl": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *myprovider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config hashicupsProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown apikey",
			"The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	apiKey := os.Getenv("APIKEY")

	if !config.ApiKey.IsNull() {
		apiKey = config.ApiKey.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Missing apikey",
			"The provider cannot create the client because of the missing apikey. "+
				"Either target apply the source of the value first,set the value statically in the configuration, or use the APIKEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	baseUrl := config.BaseUrl.ValueString()

	if baseUrl == "" {
		baseUrl = "http://localhost:8080"
	}

	client := client.NewClient(baseUrl, apiKey)
	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *myprovider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *myprovider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
