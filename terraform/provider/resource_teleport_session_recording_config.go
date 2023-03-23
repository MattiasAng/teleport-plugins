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

	apitypes "github.com/gravitational/teleport/api/types"
	tfschema "github.com/gravitational/teleport-plugins/terraform/tfschema"
	"github.com/gravitational/teleport-plugins/lib/backoff"
	"github.com/gravitational/trace"
	"github.com/jonboulle/clockwork"
)

// resourceTeleportSessionRecordingConfigType is the resource metadata type
type resourceTeleportSessionRecordingConfigType struct{}

// resourceTeleportSessionRecordingConfig is the resource
type resourceTeleportSessionRecordingConfig struct {
	p Provider
}

// GetSchema returns the resource schema
func (r resourceTeleportSessionRecordingConfigType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfschema.GenSchemaSessionRecordingConfigV2(ctx)
}

// NewResource creates the empty resource
func (r resourceTeleportSessionRecordingConfigType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceTeleportSessionRecordingConfig{
		p: *(p.(*Provider)),
	}, nil
}

// Create creates the provision token
func (r resourceTeleportSessionRecordingConfig) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sessionRecordingConfig := &apitypes.SessionRecordingConfigV2{}
	diags = tfschema.CopySessionRecordingConfigV2FromTerraform(ctx, plan, sessionRecordingConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sessionRecordingConfig.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error setting SessionRecordingConfig defaults", trace.Wrap(err), "session_recording_config"))
		return
	}

	

	sessionRecordingConfigBefore, err := r.p.Client.GetSessionRecordingConfig(ctx)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
		return
	}

	err = r.p.Client.SetSessionRecordingConfig(ctx, sessionRecordingConfig)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error creating SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
		return
	}

	var sessionRecordingConfigI apitypes.SessionRecordingConfig

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		sessionRecordingConfigI, err = r.p.Client.GetSessionRecordingConfig(ctx)
		if err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))	
			return
		}
		if sessionRecordingConfigBefore.GetMetadata().ID != sessionRecordingConfigI.GetMetadata().ID || false {
			break
		}
		if bErr := backoff.Do(ctx); bErr != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
			return
		}
		if tries >= r.p.RetryConfig.MaxTries {
			diagMessage := fmt.Sprintf("Error reading SessionRecordingConfig (tried %d times) - state outdated, please import resource", tries)
			resp.Diagnostics.Append(diagFromWrappedErr(diagMessage, trace.Wrap(err), "session_recording_config"))
			return
		}
	}
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))	
		return
	}

	sessionRecordingConfig, ok := sessionRecordingConfigI.(*apitypes.SessionRecordingConfigV2)
	if !ok {
		resp.Diagnostics.Append(
			diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Errorf("Can not convert %T to SessionRecordingConfigV2", sessionRecordingConfigI), "session_recording_config"),
		)
		return
	}

	diags = tfschema.CopySessionRecordingConfigV2ToTerraform(ctx, *sessionRecordingConfig, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Attrs["id"] = types.String{Value: "session_recording_config"}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read reads teleport SessionRecordingConfig
func (r resourceTeleportSessionRecordingConfig) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sessionRecordingConfigI, err := r.p.Client.GetSessionRecordingConfig(ctx)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))	
		return
	}

	sessionRecordingConfig := sessionRecordingConfigI.(*apitypes.SessionRecordingConfigV2)
	diags = tfschema.CopySessionRecordingConfigV2ToTerraform(ctx, *sessionRecordingConfig, &state)
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

// Update updates teleport SessionRecordingConfig
func (r resourceTeleportSessionRecordingConfig) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.IsConfigured(resp.Diagnostics) {
		return
	}

	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sessionRecordingConfig := &apitypes.SessionRecordingConfigV2{}
	diags = tfschema.CopySessionRecordingConfigV2FromTerraform(ctx, plan, sessionRecordingConfig)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := sessionRecordingConfig.CheckAndSetDefaults()
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
		return
	}

	sessionRecordingConfigBefore, err := r.p.Client.GetSessionRecordingConfig(ctx)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
		return
	}

	err = r.p.Client.SetSessionRecordingConfig(ctx, sessionRecordingConfig)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
		return
	}

	var sessionRecordingConfigI apitypes.SessionRecordingConfig

	tries := 0
	backoff := backoff.NewDecorr(r.p.RetryConfig.Base, r.p.RetryConfig.Cap, clockwork.NewRealClock())
	for {
		tries = tries + 1
		sessionRecordingConfigI, err = r.p.Client.GetSessionRecordingConfig(ctx)
		if err != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
			return
		}
		if sessionRecordingConfigBefore.GetMetadata().ID != sessionRecordingConfigI.GetMetadata().ID || false {
			break
		}
		if bErr := backoff.Do(ctx); bErr != nil {
			resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))
			return
		}
		if tries >= r.p.RetryConfig.MaxTries {
			diagMessage := fmt.Sprintf("Error reading SessionRecordingConfig (tried %d times) - state outdated, please import resource", tries)
			resp.Diagnostics.AddError(diagMessage, "session_recording_config")
			return
		}
	}
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error reading SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))	
		return
	}

	sessionRecordingConfig = sessionRecordingConfigI.(*apitypes.SessionRecordingConfigV2)
	diags = tfschema.CopySessionRecordingConfigV2ToTerraform(ctx, *sessionRecordingConfig, &plan)
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

// Delete deletes Teleport SessionRecordingConfig
func (r resourceTeleportSessionRecordingConfig) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	err := r.p.Client.ResetSessionRecordingConfig(ctx)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error deleting SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))	
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports SessionRecordingConfig state
func (r resourceTeleportSessionRecordingConfig) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	sessionRecordingConfigI, err := r.p.Client.GetSessionRecordingConfig(ctx)
	if err != nil {
		resp.Diagnostics.Append(diagFromWrappedErr("Error updating SessionRecordingConfig", trace.Wrap(err), "session_recording_config"))	
		return
	}

	sessionRecordingConfig := sessionRecordingConfigI.(*apitypes.SessionRecordingConfigV2)

	var state types.Object

	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = tfschema.CopySessionRecordingConfigV2ToTerraform(ctx, *sessionRecordingConfig, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Attrs["id"] = types.String{Value: sessionRecordingConfig.Metadata.Name}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
