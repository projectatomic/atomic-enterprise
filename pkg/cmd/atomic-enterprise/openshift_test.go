package openshift

import (
	"strings"
	"testing"
)

func TestCommandFor(t *testing.T) {
	cmd := CommandFor("atomic-enterprise-router")
	if !strings.HasPrefix(cmd.Use, "atomic-enterprise-router ") {
		t.Errorf("expected command to start with prefix: %#v", cmd)
	}

	cmd = CommandFor("unknown")
	if cmd.Use != "atomic-enterprise" {
		t.Errorf("expected command to be atomic-enterprise: %#v", cmd)
	}
}
