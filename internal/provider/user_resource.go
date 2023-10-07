package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/srikanthbhandary-teach/my-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &UserResource{}
	_ resource.ResourceWithConfigure = &UserResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &UserResource{
		client: &client.Client{},
	}
}

// UserResource is the resource implementation.
type UserResource struct {
	client *client.Client
}

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Unable to setup the client %v", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the resource.
func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Optional: true,
			},
			"age": schema.Int64Attribute{
				Optional: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan MyInfo
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	age := plan.Age.ValueInt64()
	err := r.client.CreateMyInfo(plan.ID.String(), plan.Name.ValueString(), int(age))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating User",
			"Could not create User, unexpected error: "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MyInfo
	_ = req.State.Get(ctx, &state)

	users, err := r.client.GetMyInfo(state.ID.String())
	if err != nil {
		tflog.Info(ctx, "Read HashiCups Coffees", map[string]any{"success": err.Error()})

		resp.Diagnostics.AddError(
			"UnableSRi to Read HashiCups Coffees",
			state.ID.String())
		return
	}
	tflog.SetField(ctx, "Buddy", users)
	tflog.Info(ctx, "Read HashiCups Coffees", map[string]any{"success": users})

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read HashiCups Coffees",
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
		state.ID = userState.ID
		state.Name = userState.Name
		state.Age = userState.Age
		break
	}

	var f interface{}
	tflog.Info(ctx, "ReadGEt", map[string]any{"Data1": f})

	// Set state

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan MyInfo
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Update User
	err := r.client.UpdateMyInfo(plan.ID.String(), plan.Name.ValueString(), int(plan.Age.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User",
			"Could not update user ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan MyInfo
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Update User
	err := r.client.DeleteMyInfo(plan.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating User",
			"Could not update user ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
}
