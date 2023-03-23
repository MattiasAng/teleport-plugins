// Code generated by _gen/main.go DO NOT EDIT
/*
Copyright 2015-2022 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/path"

	apitypes "github.com/gravitational/teleport/api/types"
	tfschema "github.com/gravitational/teleport-plugins/terraform/tfschema"
	"github.com/gravitational/teleport-plugins/lib/backoff"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// resourceTeleportRoleType is the resource metadata type
type resourceTeleportRoleType struct{}

// resourceTeleportRole is the resource
type resourceTeleportRole struct {
	p Provider
}

// GetSchema returns the resource schema
func (r resourceTeleportRoleType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfschema.GenSchemaRoleV6(ctx)
}

// NewResource creates the empty resource
func (r resourceTeleportRoleType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTeleportRole{
		p: *(p.(*Provider)),
	}, nil
}

// Create creates the Role
func (r resourceTeleportRole) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := &apitypes.RoleV6{}
	diags = tfschema.CopyRoleV6FromTerraform(ctx, plan, role)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	

	_, err := r.p.Client.GetRole(ctx, role.Metadata.Name)
	if !trace.IsNotFound(err) {
		if err == nil {
			n := role.Metadata.Name
			existErr := fmt.Sprintf("Role exists in Teleport. Either remove it (tctl rm role/%v)"+
				" or import it to the existing state (terraform import teleport_app.%v %v)", n, n, n)

			resp.Diagnostics.Append(diagFromErr("Role exists in Teleport", trace.Errorf(existErr)))
			return
		}

		resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Wrap(err), "role"))
		return
	}

	err = role.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error setting Role defaults", trace.Wrap(err), "role"))
		return
	}

	err = r.p.Client.UpsertRole(ctx, role)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error creating Role", trace.Wrap(err), "role"))
		return
	}

	id := role.Metadata.Name
	var roleI apitypes.Role

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		roleI, err = r.p.Client.GetRole(ctx, id)
		if trace.IsNotFound(err) {
			if bErr := backoff.Do(ctx); bErr != nil {
				resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Wrap(err), "role"))
				return
			}
			if tries >= r.p.RetryConfig.MaxTries {
				diagMessage := fmt.Sprintf("Error reading Role (tried %d times) - state outdated, please import resource", tries)
				resp.Diagnostics.Append(diagFromWrappedErr(diagMessage, trace.Wrap(err), "role"))
				return
			}
			continue
		}
		break
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Wrap(err), "role"))
		return
	}

	role, ok := roleI.(*apitypes.RoleV6)
	if !ok {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Errorf("Can not convert %T to RoleV6", roleI), "role"))
		return
	}

	diags = tfschema.CopyRoleV6ToTerraform(ctx, *role, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Attrs["id"] = types.String{Value: role.Metadata.Name}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads teleport Role
func (r resourceTeleportRole) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	diags = req.State.GetAttribute(ctx, path.Root("metadata").AtName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	roleI, err := r.p.Client.GetRole(ctx, id.Value)
	if trace.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Wrap(err), "role"))
		return
	}

	role := roleI.(*apitypes.RoleV6)
	diags = tfschema.CopyRoleV6ToTerraform(ctx, *role, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates teleport Role
func (r resourceTeleportRole) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	role := &apitypes.RoleV6{}
	diags = tfschema.CopyRoleV6FromTerraform(ctx, plan, role)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := role.Metadata.Name

	err := role.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating Role", err, "role"))
		return
	}

	roleBefore, err := r.p.Client.GetRole(ctx, name)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", err, "role"))
		return
	}

	err = r.p.Client.UpsertRole(ctx, role)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating Role", err, "role"))
		return
	}

	var roleI apitypes.Role

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		roleI, err = r.p.Client.GetRole(ctx, name)
		if err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", err, "role"))
			return
		}
		if roleBefore.GetMetadata().ID != roleI.GetMetadata().ID || false {
			break
		}

		if err := backoff.Do(ctx); err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Wrap(err), "role"))
			return
		}
		if tries >= r.p.RetryConfig.MaxTries {
			diagMessage := fmt.Sprintf("Error reading Role (tried %d times) - state outdated, please import resource", tries)
			resp.Diagnostics.AddError(diagMessage, "role")
			return
		}
	}

	role = roleI.(*apitypes.RoleV6)
	diags = tfschema.CopyRoleV6ToTerraform(ctx, *role, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes Teleport Role
func (r resourceTeleportRole) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var id types.String
	diags := req.State.GetAttribute(ctx, path.Root("metadata").AtName("name"), &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.Client.DeleteRole(ctx, id.Value)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error deleting RoleV6", trace.Wrap(err), "role"))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports Role state
func (r resourceTeleportRole) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	roleI, err := r.p.Client.GetRole(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading Role", trace.Wrap(err), "role"))
		return
	}

	role := roleI.(*apitypes.RoleV6)

	var state types.Object

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = tfschema.CopyRoleV6ToTerraform(ctx, *role, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Attrs["id"] = types.String{Value: role.Metadata.Name}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
