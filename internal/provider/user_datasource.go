package provider

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/srikanthbhandary-teach/my-client"
)

var (
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{
		client: &client.Client{},
	}
}

type UserDataSource struct {
	client *client.Client
}

type UserInfo struct {
	Filter map[string]types.String `tfsdk:"filter" json:"filter"`
	ID     types.Int64             `tfsdk:"id" json:"number"`
	Name   types.String            `tfsdk:"name" json:"name"`
	Age    types.Int64             `tfsdk:"age" json:"age"`
	Users  []MyInfo                `tfsdk:"users"`
}

type MyInfo struct {
	ID   types.Int64  `tfsdk:"id" json:"number"`
	Name types.String `tfsdk:"name" json:"name"`
	Age  types.Int64  `tfsdk:"age" json:"age"`
}

// Configure adds the provider configured client to the data source.
func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected .Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Schema defines the schema for the data source.
func (d *UserDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Optional: true,
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"age": schema.Int64Attribute{
				Optional: true,
			},
			"filter": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"users": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"age": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserInfo
	_ = req.Config.Get(ctx, &state)

	users, err := d.client.GetMyInfo(state.Filter["id"].ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read users data",
			err.Error())
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Users data",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, user := range users {
		id, _ := strconv.Atoi(user.ID)
		userState := MyInfo{
			ID:   types.Int64Value(int64(id)),
			Name: types.StringValue(user.Name),
			Age:  types.Int64Value(int64(user.Age)),
		}
		state.Users = append(state.Users, userState)
	}

	// Set state
	sort.Slice(state.Users[:], func(i, j int) bool {
		return state.Users[i].ID.ValueInt64() < state.Users[j].ID.ValueInt64()
	})
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
