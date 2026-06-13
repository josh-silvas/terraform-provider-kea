package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/josh-silvas/terraform-provider-kea/tools/kea"
)

func TestApplyDdnsConfig(t *testing.T) {
	trueVal := true
	config := remoteSubnet4ResourceSchema{
		DdnsSendUpdates:           types.BoolValue(true),
		DdnsOverrideNoUpdate:      types.BoolValue(true),
		DdnsOverrideClientUpdate:  types.BoolValue(true),
		DdnsGeneratedPrefix:       types.StringValue("host"),
		DdnsQualifyingSuffix:      types.StringValue("example.com."),
		DdnsRevDNSName:            types.StringValue("225.168.192.in-addr.arpa."),
		DdnsUseConflictResolution: types.BoolValue(true),
	}

	subnet := kea.NewRemoteSubnet4{ID: 10002260, Subnet: "192.168.255.0/24"}
	applyDdnsConfig(config, &subnet)

	if subnet.DdnsSendUpdates == nil || *subnet.DdnsSendUpdates != trueVal {
		t.Fatalf("expected ddns-send-updates true, got %#v", subnet.DdnsSendUpdates)
	}
	if subnet.DdnsOverrideNoUpdate == nil || *subnet.DdnsOverrideNoUpdate != trueVal {
		t.Fatalf("expected ddns-override-no-update true, got %#v", subnet.DdnsOverrideNoUpdate)
	}
	if subnet.DdnsOverrideClientUpdate == nil || *subnet.DdnsOverrideClientUpdate != trueVal {
		t.Fatalf("expected ddns-override-client-update true, got %#v", subnet.DdnsOverrideClientUpdate)
	}
	if subnet.DdnsGeneratedPrefix != "host" {
		t.Fatalf("expected ddns-generated-prefix, got %q", subnet.DdnsGeneratedPrefix)
	}
	if subnet.DdnsQualifyingSuffix != "example.com." {
		t.Fatalf("expected ddns-qualifying-suffix, got %q", subnet.DdnsQualifyingSuffix)
	}
	if subnet.DdnsRevDNSName != "225.168.192.in-addr.arpa." {
		t.Fatalf("expected ddns-rev-dns-name, got %q", subnet.DdnsRevDNSName)
	}
	if subnet.DdnsUseConflictResolution == nil || *subnet.DdnsUseConflictResolution != trueVal {
		t.Fatalf("expected ddns-use-conflict-resolution true, got %#v", subnet.DdnsUseConflictResolution)
	}
}

func TestSetDdnsStateFromKea(t *testing.T) {
	trueVal := true
	respData := kea.RemoteSubnet4{
		DdnsSendUpdates:           &trueVal,
		DdnsOverrideNoUpdate:      &trueVal,
		DdnsOverrideClientUpdate:  &trueVal,
		DdnsGeneratedPrefix:       "host",
		DdnsQualifyingSuffix:      "example.com.",
		DdnsRevDNSName:            "225.168.192.in-addr.arpa.",
		DdnsUseConflictResolution: &trueVal,
	}

	config := remoteSubnet4ResourceSchema{}
	setDdnsStateFromKea(respData, &config)

	if config.DdnsSendUpdates.ValueBool() != true {
		t.Fatalf("expected ddns_send_updates true, got %v", config.DdnsSendUpdates)
	}
	if config.DdnsGeneratedPrefix.ValueString() != "host" {
		t.Fatalf("expected ddns_generated_prefix, got %q", config.DdnsGeneratedPrefix.ValueString())
	}
	if config.DdnsQualifyingSuffix.ValueString() != "example.com." {
		t.Fatalf("expected ddns_qualifying_suffix, got %q", config.DdnsQualifyingSuffix.ValueString())
	}
	if config.DdnsRevDNSName.ValueString() != "225.168.192.in-addr.arpa." {
		t.Fatalf("expected ddns_rev_dns_name, got %q", config.DdnsRevDNSName.ValueString())
	}
}
