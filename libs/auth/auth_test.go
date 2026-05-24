package auth

import (
	"context"
	"testing"
)

func TestGetTenantID(t *testing.T) {
	ctx := context.WithValue(context.Background(), TenantIDKey, "test-tenant")
	if GetTenantID(ctx) != "test-tenant" {
		t.Errorf("expected test-tenant, got %s", GetTenantID(ctx))
	}

	if GetTenantID(context.Background()) != "" {
		t.Errorf("expected empty string, got %s", GetTenantID(context.Background()))
	}
}
