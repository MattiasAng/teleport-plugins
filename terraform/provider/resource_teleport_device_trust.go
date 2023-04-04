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
	tfschema "github.com/gravitational/teleport-plugins/terraform/tfschema/devicetrust/v1"
	"github.com/gravitational/teleport-plugins/lib/backoff"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// resourceTeleportDeviceV1Type is the resource metadata type
type resourceTeleportDeviceV1Type struct{}

// resourceTeleportDeviceV1 is the resource
type resourceTeleportDeviceV1 struct {
	p Provider
}

// GetSchema returns the resource schema
func (r resourceTeleportDeviceV1Type) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfschema.GenSchemaDeviceV1(ctx)
}

// NewResource creates the empty resource
func (r resourceTeleportDeviceV1Type) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTeleportDeviceV1{
		p: *(p.(*Provider)),
	}, nil
}

// Create creates the DeviceV1
func (r resourceTeleportDeviceV1) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustedDevice := &apitypes.DeviceV1{}
	diags = tfschema.CopyDeviceV1FromTerraform(ctx, plan, trustedDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	

	_, err := r.p.Client.GetDeviceResource(ctx, trustedDevice.Metadata.Name)
	if !trace.IsNotFound(err) {
		if err == nil {
			n := trustedDevice.Metadata.Name
			existErr := fmt.Sprintf("DeviceV1 exists in Teleport. Either remove it (tctl rm device/%v)"+
				" or import it to the existing state (terraform import teleport_app.%v %v)", n, n, n)

			resp.Diagnostics.Append(diagFromErr("DeviceV1 exists in Teleport", trace.Errorf(existErr)))
			return
		}

		resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", trace.Wrap(err), "device"))
		return
	}

	

	dev, err := r.p.Client.UpsertDeviceResource(ctx, trustedDevice)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error creating DeviceV1", trace.Wrap(err), "device"))
		return
	}

	id := dev.Metadata.Name
	// Not really an inferface, just using the same name for easier templating.
	var trustedDeviceI *apitypes.DeviceV1
	

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		trustedDeviceI, err = r.p.Client.GetDeviceResource(ctx, id)
		if trace.IsNotFound(err) {
			if bErr := backoff.Do(ctx); bErr != nil {
				resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", trace.Wrap(err), "device"))
				return
			}
			if tries >= r.p.RetryConfig.MaxTries {
				diagMessage := fmt.Sprintf("Error reading DeviceV1 (tried %d times) - state outdated, please import resource", tries)
				resp.Diagnostics.Append(diagFromWrappedErr(diagMessage, trace.Wrap(err), "device"))
				return
			}
			continue
		}
		break
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", trace.Wrap(err), "device"))
		return
	}

	trustedDevice = trustedDeviceI
	

	diags = tfschema.CopyDeviceV1ToTerraform(ctx, *trustedDevice, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Attrs["id"] = types.String{Value: trustedDevice.Metadata.Name}
    plan.Attrs["Metadata.name"] = types.String{Value: trustedDevice.Metadata.Name}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads teleport DeviceV1
func (r resourceTeleportDeviceV1) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
    p :=  path.Root("id")
	diags = req.State.GetAttribute(ctx, p, &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustedDeviceI, err := r.p.Client.GetDeviceResource(ctx, id.Value)
	if trace.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", trace.Wrap(err), "device"))
		return
	}

	trustedDevice := trustedDeviceI
	diags = tfschema.CopyDeviceV1ToTerraform(ctx, *trustedDevice, &state)
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

// Update updates teleport DeviceV1
func (r resourceTeleportDeviceV1) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustedDevice := &apitypes.DeviceV1{}
	diags = tfschema.CopyDeviceV1FromTerraform(ctx, plan, trustedDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


    var id types.String
    p := path.Root("id")
    diags = req.State.GetAttribute(ctx, p , &id)
    resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
    trustedDevice.Metadata.Name = id.Value
    

    name := trustedDevice.Metadata.Name

	

	trustedDeviceBefore, err := r.p.Client.GetDeviceResource(ctx, name)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", err, "device"))
		return
	}

	_, err = r.p.Client.UpsertDeviceResource(ctx, trustedDevice)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating DeviceV1", err, "device"))
		return
	}

	

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		trustedDevice, err = r.p.Client.GetDeviceResource(ctx, name)
		if err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", err, "device"))
			return
		}

        if trustedDeviceBefore.Metadata.Name != trustedDevice.Metadata.Name || true {
			break
		}

        

		if err := backoff.Do(ctx); err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", trace.Wrap(err), "device"))
			return
		}
		if tries >= r.p.RetryConfig.MaxTries {
			diagMessage := fmt.Sprintf("Error reading DeviceV1 (tried %d times) - state outdated, please import resource", tries)
			resp.Diagnostics.AddError(diagMessage, "device")
			return
		}
	}

	diags = tfschema.CopyDeviceV1ToTerraform(ctx, *trustedDevice, &plan)
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

// Delete deletes Teleport DeviceV1
func (r resourceTeleportDeviceV1) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var id types.String

    p :=  path.Root("id")
	diags := req.State.GetAttribute(ctx, p , &id)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.p.Client.DeleteDeviceResource(ctx, id.Value)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error deleting DeviceV1", trace.Wrap(err), "device"))
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports DeviceV1 state
func (r resourceTeleportDeviceV1) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	trustedDevice, err := r.p.Client.GetDeviceResource(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading DeviceV1", trace.Wrap(err), "device"))
		return
	}

	

	var state types.Object

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = tfschema.CopyDeviceV1ToTerraform(ctx, *trustedDevice, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Attrs["id"] = types.String{Value: trustedDevice.Metadata.Name}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
