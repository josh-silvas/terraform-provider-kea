// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

var (
	// Ensure provider defined types fully satisfy framework interfaces.
	_ datasource.DataSource              = &remoteSubnet4DataSource{}
	_ datasource.DataSourceWithConfigure = &remoteSubnet4DataSource{}
)

// NewRemoteSubnet4DataSource : Creates a new empty data source client.
func NewRemoteSubnet4DataSource() datasource.DataSource {
	return &remoteSubnet4DataSource{}
}

type (
	// remoteSubnet4DataSource defines the data source client.
	remoteSubnet4DataSource struct {
		client *kea.Client
	}

	// remoteSubnet4DataSourceSchema describes the data source data model.
	// Maps to the source schema data.
	remoteSubnet4DataSourceSchema struct {
		Prefix                    types.String                         `tfsdk:"prefix"`
		SubnetID                  types.Int64                          `tfsdk:"subnet_id"`
		Hostname                  types.String                         `tfsdk:"hostname"`
		ID                        types.Int64                          `tfsdk:"id"`
		OptionData                []remoteSubnet4DataSourceOptionModel `tfsdk:"option_data"`
		Pools                     types.List                           `tfsdk:"pools"`
		Relay                     types.List                           `tfsdk:"relay"`
		Subnet                    types.String                         `tfsdk:"subnet"`
		UserContext               types.Map                            `tfsdk:"user_context"`
		DdnsSendUpdates           types.Bool                           `tfsdk:"ddns_send_updates"`
		DdnsOverrideNoUpdate      types.Bool                           `tfsdk:"ddns_override_no_update"`
		DdnsOverrideClientUpdate  types.Bool                           `tfsdk:"ddns_override_client_update"`
		DdnsGeneratedPrefix       types.String                         `tfsdk:"ddns_generated_prefix"`
		DdnsQualifyingSuffix      types.String                         `tfsdk:"ddns_qualifying_suffix"`
		DdnsRevDNSName            types.String                         `tfsdk:"ddns_rev_dns_name"`
		DdnsUseConflictResolution types.Bool                           `tfsdk:"ddns_use_conflict_resolution"`
	}

	// optionDataModel : Represents a single option-data entry in Kea.
	remoteSubnet4DataSourceOptionModel struct {
		Code       types.Int64  `tfsdk:"code"`
		Data       types.String `tfsdk:"data"`
		Name       types.String `tfsdk:"name"`
		AlwaysSend types.Bool   `tfsdk:"always_send"`
	}
)

// Metadata : Defines the data source metadata.
func (d *remoteSubnet4DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_subnet4_data_source"
}

// Schema : Defines the data source schema.
func (d *remoteSubnet4DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Remote subnet4 data source",
		Attributes: map[string]schema.Attribute{
			"prefix": schema.StringAttribute{
				MarkdownDescription: "Prefix to fetch from Kea configuration-backend. e.g. 192.168.230.0/24`",
				Optional:            true,
			},
			"subnet_id": schema.Int64Attribute{
				MarkdownDescription: "Subnet4 ID to fetch from Kea configuration-backend. e.g. 1921682300`",
				Optional:            true,
			},
			"hostname": schema.StringAttribute{
				MarkdownDescription: "Hostname of the kea server to connect to. e.g. `kea.example.com`",
				Required:            true,
			},

			"id":     schema.Int64Attribute{Computed: true},
			"subnet": schema.StringAttribute{Computed: true},
			"pools":  schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"relay":  schema.ListAttribute{Computed: true, ElementType: types.StringType},
			"option_data": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code":        schema.Int64Attribute{Computed: true},
						"name":        schema.StringAttribute{Computed: true},
						"data":        schema.StringAttribute{Computed: true},
						"always_send": schema.BoolAttribute{Computed: true},
					},
				},
			},
			"user_context": schema.MapAttribute{
				MarkdownDescription: "Arbitrary string data to tie to the subnet. e.g. `{site = \"AUS\", name = \"Austin, Tx\"}`",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"ddns_send_updates": schema.BoolAttribute{
				MarkdownDescription: "When true, Kea sends DDNS updates for leases in this subnet.",
				Computed:            true,
			},
			"ddns_override_no_update": schema.BoolAttribute{
				MarkdownDescription: "When true, Kea sends DDNS updates even if the client sets the N flag.",
				Computed:            true,
			},
			"ddns_override_client_update": schema.BoolAttribute{
				MarkdownDescription: "When true, Kea sends DDNS updates even if the client requests to do it itself.",
				Computed:            true,
			},
			"ddns_generated_prefix": schema.StringAttribute{
				MarkdownDescription: "Prefix used when Kea generates forward DDNS names.",
				Computed:            true,
			},
			"ddns_qualifying_suffix": schema.StringAttribute{
				MarkdownDescription: "Suffix appended to generated forward DDNS names.",
				Computed:            true,
			},
			"ddns_rev_dns_name": schema.StringAttribute{
				MarkdownDescription: "Reverse DNS zone for PTR records.",
				Computed:            true,
			},
			"ddns_use_conflict_resolution": schema.BoolAttribute{
				MarkdownDescription: "When true, Kea uses DDNS conflict-resolution behavior for this subnet.",
				Computed:            true,
			},
		},
	}
}

// Configure : Configures the data source client.
func (d *remoteSubnet4DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Fetch the Kea DHCP client from the provider.
	client, ok := req.ProviderData.(*kea.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kea.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.client = client
}

// Read : Reads the data source data into the Terraform state.
func (d *remoteSubnet4DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Define an empty configuration.
	var config remoteSubnet4DataSourceSchema

	// Read Terraform configuration data into the model
	// Also append any diagnostics to the diagnostics list.
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	// Validate that only one of `prefix` or `subnet_id` is specified.
	if (!config.Prefix.IsNull() && !config.SubnetID.IsNull()) || (config.Prefix.IsNull() && config.SubnetID.IsNull()) {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"One and only one of `prefix` or `subnet_id` must be specified.",
		)
	}

	// Validate that a `hostname` is specified.
	if config.Hostname.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid Configuration",
			"A `hostname` must be specified. DNS name or IP address of the Kea DHCP server.",
		)
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	var respData kea.RemoteSubnet4
	var err error
	if !config.Prefix.IsNull() {
		// nolint: contextcheck
		respData, err = d.client.RemoteSubnet4GetByPrefix(config.Hostname.ValueString(), config.Prefix.ValueString())
		if err != nil {
			// Only return an error if the error is NOT subnet not found.
			if !strings.Contains(err.Error(), "not found") {
				resp.Diagnostics.AddError(
					"RemoteSubnet4GetByPrefix",
					fmt.Sprintf("Unable to read example, got error: %s", err),
				)
				return
			}
		}
	} else {
		// nolint: contextcheck
		respData, err = d.client.RemoteSubnet4GetByID(config.Hostname.ValueString(), int(config.SubnetID.ValueInt64()))
		if err != nil {
			// Only return an error if the error is NOT subnet not found.
			if !strings.Contains(err.Error(), "not found") {
				resp.Diagnostics.AddError(
					"RemoteSubnet4GetByID",
					fmt.Sprintf("Unable to read example, got error: %s", err),
				)
				return
			}
		}
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Marshalling the response data taken from Kea, and write
	// it into the TF Subnets model.

	config.ID = types.Int64Value(int64(respData.ID))
	config.OptionData = func() []remoteSubnet4DataSourceOptionModel {
		r := make([]remoteSubnet4DataSourceOptionModel, 0)
		for _, v := range respData.OptionData {
			code := 0
			if v.Code != nil {
				code = *v.Code
			}
			r = append(r, remoteSubnet4DataSourceOptionModel{
				Code:       types.Int64Value(int64(code)),
				Data:       types.StringValue(v.Data),
				Name:       types.StringValue(v.Name),
				AlwaysSend: types.BoolValue(v.AlwaysSend),
			})
		}
		return r
	}()
	config.Pools = func() types.List {
		r := make([]attr.Value, 0)
		for _, v := range respData.Pools {
			r = append(r, types.StringValue(v.Pool))
		}
		retVal, diags := types.ListValue(types.StringType, r)
		resp.Diagnostics.Append(diags...)
		return retVal
	}()
	config.Relay = func() types.List {
		r := make([]attr.Value, 0)
		for _, v := range respData.Relay.IPAddresses {
			r = append(r, types.StringValue(v))
		}
		retVal, diags := types.ListValue(types.StringType, r)
		resp.Diagnostics.Append(diags...)
		return retVal
	}()
	config.Subnet = types.StringValue(respData.Subnet)
	if respData.UserContext != nil {
		config.UserContext = func() types.Map {
			fr := make(map[string]attr.Value)
			for k, v := range respData.UserContext {
				fr[k] = types.StringValue(fmt.Sprintf("%v", v))
			}
			mv, diags := types.MapValue(types.StringType, fr)
			resp.Diagnostics.Append(diags...)
			return mv
		}()
	}
	if respData.DdnsSendUpdates != nil {
		config.DdnsSendUpdates = types.BoolValue(*respData.DdnsSendUpdates)
	}
	if respData.DdnsOverrideNoUpdate != nil {
		config.DdnsOverrideNoUpdate = types.BoolValue(*respData.DdnsOverrideNoUpdate)
	}
	if respData.DdnsOverrideClientUpdate != nil {
		config.DdnsOverrideClientUpdate = types.BoolValue(*respData.DdnsOverrideClientUpdate)
	}
	if respData.DdnsGeneratedPrefix != "" {
		config.DdnsGeneratedPrefix = types.StringValue(respData.DdnsGeneratedPrefix)
	}
	if respData.DdnsQualifyingSuffix != "" {
		config.DdnsQualifyingSuffix = types.StringValue(respData.DdnsQualifyingSuffix)
	}
	if respData.DdnsRevDNSName != "" {
		config.DdnsRevDNSName = types.StringValue(respData.DdnsRevDNSName)
	}
	if respData.DdnsUseConflictResolution != nil {
		config.DdnsUseConflictResolution = types.BoolValue(*respData.DdnsUseConflictResolution)
	}

	// If there are any diagnostics errors, stop here.
	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
