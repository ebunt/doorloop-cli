package mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"doorloop-pp-cli/internal/cli"
	"doorloop-pp-cli/internal/client"
	"doorloop-pp-cli/internal/cliutil"
	"doorloop-pp-cli/internal/config"
	"doorloop-pp-cli/internal/mcp/cobratree"
	"doorloop-pp-cli/internal/store"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterTools registers all API operations as MCP tools.
func RegisterTools(s *server.MCPServer) {
	s.AddTool(
		mcplib.NewTool("accounts_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("Account ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/accounts/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("accounts_list",
			mcplib.WithDescription("Optional: filter_group, filter_text."),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_text", mcplib.Description("Filter by account name")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/accounts", false, []mcpParamBinding{{PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_text", WireName: "filter_text", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("communication_logs_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/communication-logs", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("communication_logs_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/communication-logs/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("communication_logs_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/communication-logs/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("communication_logs_list",
			mcplib.WithDescription("Optional: filter_tenant, filter_lease, filter_property."),
			mcplib.WithString("filter_tenant", mcplib.Description("Filter by tenant ID")),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/communication-logs", false, []mcpParamBinding{{PublicName: "filter_tenant", WireName: "filter_tenant", Location: "query"}, {PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("communication_logs_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/communication-logs/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("expenses_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/expenses", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("expenses_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/expenses/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("expenses_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/expenses/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("expenses_list",
			mcplib.WithDescription("Optional: filter_property, filter_group, filter_date_from (plus 2 more)."),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_date_from", mcplib.Description("Expense date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_date_to", mcplib.Description("Expense date to (YYYY-MM-DD)")),
			mcplib.WithString("filter_vendor", mcplib.Description("Filter by vendor ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/expenses", false, []mcpParamBinding{{PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_date_from", WireName: "filter_date_from", Location: "query"}, {PublicName: "filter_date_to", WireName: "filter_date_to", Location: "query"}, {PublicName: "filter_vendor", WireName: "filter_vendor", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("expenses_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/expenses/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("files_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/files", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("files_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/files/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("files_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/files/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("files_list",
			mcplib.WithDescription("Optional: filter_property, filter_lease."),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/files", false, []mcpParamBinding{{PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("files_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/files/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_charges_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/lease-charges", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_charges_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/lease-charges/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_charges_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-charges/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_charges_list",
			mcplib.WithDescription("Optional: filter_lease, filter_property, filter_date_from (plus 1 more)."),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_date_from", mcplib.Description("Charge date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_date_to", mcplib.Description("Charge date to (YYYY-MM-DD)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-charges", false, []mcpParamBinding{{PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_date_from", WireName: "filter_date_from", Location: "query"}, {PublicName: "filter_date_to", WireName: "filter_date_to", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_charges_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/lease-charges/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_credits_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/lease-credits", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_credits_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/lease-credits/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_credits_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-credits/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_credits_list",
			mcplib.WithDescription("Optional: filter_lease, filter_property."),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-credits", false, []mcpParamBinding{{PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_credits_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/lease-credits/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_payments_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/lease-payments", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_payments_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/lease-payments/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_payments_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-payments/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_payments_list",
			mcplib.WithDescription("Optional: filter_lease, filter_property, filter_group (plus 2 more)."),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_date_from", mcplib.Description("Payment date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_date_to", mcplib.Description("Payment date to (YYYY-MM-DD)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-payments", false, []mcpParamBinding{{PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_date_from", WireName: "filter_date_from", Location: "query"}, {PublicName: "filter_date_to", WireName: "filter_date_to", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_payments_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/lease-payments/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_returned_payments_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/lease-returned-payments", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_returned_payments_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/lease-returned-payments/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_returned_payments_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-returned-payments/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("lease_returned_payments_list",
			mcplib.WithDescription("Optional: filter_lease, filter_property."),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/lease-returned-payments", false, []mcpParamBinding{{PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("lease_returned_payments_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/lease-returned-payments/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("leases_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/leases/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("leases_list",
			mcplib.WithDescription("Optional: filter_group, filter_property, filter_owner (plus 11 more)."),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_owner", mcplib.Description("Filter by owner ID")),
			mcplib.WithString("filter_text", mcplib.Description("Filter by lease name")),
			mcplib.WithString("filter_start_date_from", mcplib.Description("Start date range from (YYYY-MM-DD)")),
			mcplib.WithString("filter_start_date_to", mcplib.Description("Start date range to (YYYY-MM-DD)")),
			mcplib.WithString("filter_end_date_from", mcplib.Description("End date range from (YYYY-MM-DD)")),
			mcplib.WithString("filter_end_date_to", mcplib.Description("End date range to (YYYY-MM-DD)")),
			mcplib.WithString("filter_propertyClass", mcplib.Description("Filter by property class (RESIDENTIAL or COMMERCIAL)")),
			mcplib.WithString("filter_unit", mcplib.Description("Filter by unit ID")),
			mcplib.WithString("filter_tenant", mcplib.Description("Filter by tenant ID")),
			mcplib.WithString("filter_outstandingBalanceGreaterThan", mcplib.Description("Filter by minimum outstanding balance")),
			mcplib.WithString("filter_status", mcplib.Description("Filter by status (ACTIVE or INACTIVE)")),
			mcplib.WithString("filter_term", mcplib.Description("Filter by term (Rollover or AtWill)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/leases", false, []mcpParamBinding{{PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_owner", WireName: "filter_owner", Location: "query"}, {PublicName: "filter_text", WireName: "filter_text", Location: "query"}, {PublicName: "filter_start_date_from", WireName: "filter_start_date_from", Location: "query"}, {PublicName: "filter_start_date_to", WireName: "filter_start_date_to", Location: "query"}, {PublicName: "filter_end_date_from", WireName: "filter_end_date_from", Location: "query"}, {PublicName: "filter_end_date_to", WireName: "filter_end_date_to", Location: "query"}, {PublicName: "filter_propertyClass", WireName: "filter_propertyClass", Location: "query"}, {PublicName: "filter_unit", WireName: "filter_unit", Location: "query"}, {PublicName: "filter_tenant", WireName: "filter_tenant", Location: "query"}, {PublicName: "filter_outstandingBalanceGreaterThan", WireName: "filter_outstandingBalanceGreaterThan", Location: "query"}, {PublicName: "filter_status", WireName: "filter_status", Location: "query"}, {PublicName: "filter_term", WireName: "filter_term", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("leases_move_in",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/leases/move-in", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("leases_move_out",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/leases/move-out", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("leases_tenants",
			mcplib.WithDescription("Optional: filter_group, filter_property, filter_lease (plus 6 more)."),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_status", mcplib.Description("Filter by status")),
			mcplib.WithString("filter_text", mcplib.Description("Filter by tenant name/email")),
			mcplib.WithString("filter_movedInAt_from", mcplib.Description("Move-in date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_movedInAt_to", mcplib.Description("Move-in date to (YYYY-MM-DD)")),
			mcplib.WithString("filter_movedOutAt_from", mcplib.Description("Move-out date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_movedOutAt_To", mcplib.Description("Move-out date to (YYYY-MM-DD)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/leases/tenants", false, []mcpParamBinding{{PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_status", WireName: "filter_status", Location: "query"}, {PublicName: "filter_text", WireName: "filter_text", Location: "query"}, {PublicName: "filter_movedInAt_from", WireName: "filter_movedInAt_from", Location: "query"}, {PublicName: "filter_movedInAt_to", WireName: "filter_movedInAt_to", Location: "query"}, {PublicName: "filter_movedOutAt_from", WireName: "filter_movedOutAt_from", Location: "query"}, {PublicName: "filter_movedOutAt_To", WireName: "filter_movedOutAt_To", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("notes_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/notes", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("notes_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/notes/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("notes_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/notes/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("notes_list",
			mcplib.WithDescription("Optional: filter_property, filter_lease, filter_tenant."),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_tenant", mcplib.Description("Filter by tenant ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/notes", false, []mcpParamBinding{{PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_tenant", WireName: "filter_tenant", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("notes_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/notes/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("owners_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/owners", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("owners_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/owners/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("owners_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/owners/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("owners_list",
			mcplib.WithDescription("Optional: filter_text, filter_group."),
			mcplib.WithString("filter_text", mcplib.Description("Filter by owner name or email")),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/owners", false, []mcpParamBinding{{PublicName: "filter_text", WireName: "filter_text", Location: "query"}, {PublicName: "filter_group", WireName: "filter_group", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("owners_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/owners/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("portfolios_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/portfolios", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("portfolios_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/portfolios/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("portfolios_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/portfolios/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("portfolios_list",
			mcplib.WithDescription("Optional: filter_text."),
			mcplib.WithString("filter_text", mcplib.Description("Filter by portfolio name")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/portfolios", false, []mcpParamBinding{{PublicName: "filter_text", WireName: "filter_text", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("portfolios_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/portfolios/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("properties_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/properties/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("properties_list",
			mcplib.WithDescription("Optional: filter_text, filter_group, filter_class (plus 1 more)."),
			mcplib.WithString("filter_text", mcplib.Description("Filter by property name")),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_class", mcplib.Description("Filter by class (RESIDENTIAL or COMMERCIAL)")),
			mcplib.WithString("filter_owner", mcplib.Description("Filter by owner ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/properties", false, []mcpParamBinding{{PublicName: "filter_text", WireName: "filter_text", Location: "query"}, {PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_class", WireName: "filter_class", Location: "query"}, {PublicName: "filter_owner", WireName: "filter_owner", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("reports_list",
			mcplib.WithDescription("Optional: filter_property, filter_group, filter_date_from (plus 1 more)."),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_date_from", mcplib.Description("Report period from (YYYY-MM-DD)")),
			mcplib.WithString("filter_date_to", mcplib.Description("Report period to (YYYY-MM-DD)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/reports", false, []mcpParamBinding{{PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_date_from", WireName: "filter_date_from", Location: "query"}, {PublicName: "filter_date_to", WireName: "filter_date_to", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("tasks_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/tasks", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("tasks_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/tasks/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("tasks_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/tasks/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("tasks_list",
			mcplib.WithDescription("Optional: filter_property, filter_group, filter_assignee (plus 4 more)."),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_assignee", mcplib.Description("Filter by assigned user ID")),
			mcplib.WithString("filter_status", mcplib.Description("Filter by status")),
			mcplib.WithString("filter_text", mcplib.Description("Filter by task title")),
			mcplib.WithString("filter_due_date_from", mcplib.Description("Due date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_due_date_to", mcplib.Description("Due date to (YYYY-MM-DD)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/tasks", false, []mcpParamBinding{{PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_assignee", WireName: "filter_assignee", Location: "query"}, {PublicName: "filter_status", WireName: "filter_status", Location: "query"}, {PublicName: "filter_text", WireName: "filter_text", Location: "query"}, {PublicName: "filter_due_date_from", WireName: "filter_due_date_from", Location: "query"}, {PublicName: "filter_due_date_to", WireName: "filter_due_date_to", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("tasks_post_update",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/tasks/update", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("tasks_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/tasks/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("tenants_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/tenants", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("tenants_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/tenants/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("tenants_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/tenants/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("tenants_list",
			mcplib.WithDescription("Optional: filter_group, filter_property, filter_lease (plus 3 more)."),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_lease", mcplib.Description("Filter by lease ID")),
			mcplib.WithString("filter_text", mcplib.Description("Filter by tenant name or email")),
			mcplib.WithString("filter_unit", mcplib.Description("Filter by unit ID")),
			mcplib.WithString("filter_type", mcplib.Description("Filter by type (LEASE_TENANT or PROSPECT_TENANT)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/tenants", false, []mcpParamBinding{{PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_lease", WireName: "filter_lease", Location: "query"}, {PublicName: "filter_text", WireName: "filter_text", Location: "query"}, {PublicName: "filter_unit", WireName: "filter_unit", Location: "query"}, {PublicName: "filter_type", WireName: "filter_type", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("tenants_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/tenants/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("units_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/units/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("units_list",
			mcplib.WithDescription("Optional: filter_group, filter_property, filter_owner (plus 1 more)."),
			mcplib.WithString("filter_group", mcplib.Description("Filter by portfolio ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_owner", mcplib.Description("Filter by owner ID")),
			mcplib.WithString("filter_text", mcplib.Description("Filter by unit name")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/units", false, []mcpParamBinding{{PublicName: "filter_group", WireName: "filter_group", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_owner", WireName: "filter_owner", Location: "query"}, {PublicName: "filter_text", WireName: "filter_text", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("users_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/users", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("users_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/users/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("users_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/users/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("users_list",
			mcplib.WithDescription("Optional: filter_text."),
			mcplib.WithString("filter_text", mcplib.Description("Filter by user name or email")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/users", false, []mcpParamBinding{{PublicName: "filter_text", WireName: "filter_text", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("users_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/users/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_bills_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/vendor-bills", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_bills_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/vendor-bills/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_bills_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/vendor-bills/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_bills_list",
			mcplib.WithDescription("Optional: filter_vendor, filter_property, filter_date_from (plus 1 more)."),
			mcplib.WithString("filter_vendor", mcplib.Description("Filter by vendor ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithString("filter_date_from", mcplib.Description("Bill date from (YYYY-MM-DD)")),
			mcplib.WithString("filter_date_to", mcplib.Description("Bill date to (YYYY-MM-DD)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/vendor-bills", false, []mcpParamBinding{{PublicName: "filter_vendor", WireName: "filter_vendor", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}, {PublicName: "filter_date_from", WireName: "filter_date_from", Location: "query"}, {PublicName: "filter_date_to", WireName: "filter_date_to", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_bills_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/vendor-bills/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_credits_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/vendor-credits", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_credits_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/vendor-credits/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_credits_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/vendor-credits/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_credits_list",
			mcplib.WithDescription("Optional: filter_vendor, filter_property."),
			mcplib.WithString("filter_vendor", mcplib.Description("Filter by vendor ID")),
			mcplib.WithString("filter_property", mcplib.Description("Filter by property ID")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/vendor-credits", false, []mcpParamBinding{{PublicName: "filter_vendor", WireName: "filter_vendor", Location: "query"}, {PublicName: "filter_property", WireName: "filter_property", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("vendor_credits_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/vendor-credits/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendors_create",
			mcplib.WithDescription(""),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("POST", "/vendors", false, []mcpParamBinding{}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("vendors_delete",
			mcplib.WithDescription("Required: id. Destructive."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithDestructiveHintAnnotation(true),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("DELETE", "/vendors/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendors_get",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/vendors/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	s.AddTool(
		mcplib.NewTool("vendors_list",
			mcplib.WithDescription("Optional: filter_text."),
			mcplib.WithString("filter_text", mcplib.Description("Filter by vendor name")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("GET", "/vendors", false, []mcpParamBinding{{PublicName: "filter_text", WireName: "filter_text", Location: "query"}}, []string{}),
	)
	s.AddTool(
		mcplib.NewTool("vendors_update",
			mcplib.WithDescription("Required: id."),
			mcplib.WithString("id", mcplib.Required(), mcplib.Description("")),
			mcplib.WithOpenWorldHintAnnotation(true),
		),
		makeAPIHandler("PUT", "/vendors/{id}", false, []mcpParamBinding{{PublicName: "id", WireName: "id", Location: "path"}}, []string{"id"}),
	)
	// Search tool — faster than iterating list endpoints for finding specific items
	s.AddTool(
		mcplib.NewTool("search",
			mcplib.WithDescription("Full-text search across all synced data. Faster than paginating list endpoints. Requires sync first."),
			mcplib.WithString("query", mcplib.Required(), mcplib.Description("Search query (supports FTS5 syntax: AND, OR, NOT, quotes for phrases)")),
			mcplib.WithNumber("limit", mcplib.Description("Max results (default 25)")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
		),
		handleSearch,
	)
	// SQL tool — ad-hoc analysis on synced data without API calls
	s.AddTool(
		mcplib.NewTool("sql",
			mcplib.WithDescription("Run read-only SQL against local database. Use for ad-hoc analysis, aggregations, and joins across synced resources. Requires sync first."),
			mcplib.WithString("query", mcplib.Required(), mcplib.Description("SQL query (SELECT or WITH...SELECT). Tables match resource names.")),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
		),
		handleSQL,
	)

	// Context tool — front-loaded domain knowledge for agents.
	// Call this first to understand the API taxonomy, query patterns, and capabilities.
	s.AddTool(
		mcplib.NewTool("context",
			mcplib.WithDescription("Get API domain context: resource taxonomy, auth requirements, query tips, and unique capabilities. Call this first."),
			mcplib.WithReadOnlyHintAnnotation(true),
			mcplib.WithDestructiveHintAnnotation(false),
		),
		handleContext,
	)

	// Runtime Cobra-tree mirror — exposes every user-facing command that is
	// not already covered by a typed endpoint or framework MCP tool.
	cobratree.RegisterAll(s, cli.RootCmd(), cobratree.SiblingCLIPath)
}

type mcpParamBinding struct {
	PublicName string
	WireName   string
	Location   string
}

// makeAPIHandler creates a generic MCP tool handler for an API endpoint.
func makeAPIHandler(method, pathTemplate string, binaryResponse bool, bindings []mcpParamBinding, positionalParams []string) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		c, err := newMCPClient()
		if err != nil {
			return mcplib.NewToolResultError(err.Error()), nil
		}

		// mcp-go v0.47+ made CallToolParams.Arguments an `any` to support
		// non-map payloads; GetArguments() returns the map[string]any shape
		// we rely on here (or an empty map when the payload is something else).
		args := req.GetArguments()

		// positionalParams mixes real URL path params with CLI positional
		// args that map to query params (e.g. `search <query>` -> ?query=);
		// the placeholder check below disambiguates them at runtime.
		path := pathTemplate
		knownArgs := make(map[string]bool, len(bindings))
		pathParams := make(map[string]bool, len(positionalParams))
		params := make(map[string]string)
		bodyArgs := make(map[string]any)
		var headers map[string]string
		if binaryResponse {
			headers = map[string]string{client.BinaryResponseHeader: "true"}
		}
		for _, binding := range bindings {
			knownArgs[binding.PublicName] = true
			v, ok := args[binding.PublicName]
			if !ok {
				continue
			}
			switch binding.Location {
			case "path":
				placeholder := "{" + binding.WireName + "}"
				pathParams[binding.PublicName] = true
				path = strings.Replace(path, placeholder, fmt.Sprintf("%v", v), 1)
			case "body":
				bodyArgs[binding.WireName] = v
			default:
				params[binding.WireName] = fmt.Sprintf("%v", v)
			}
		}
		for _, p := range positionalParams {
			placeholder := "{" + p + "}"
			if !strings.Contains(pathTemplate, placeholder) {
				continue
			}
			pathParams[p] = true
			if v, ok := args[p]; ok {
				path = strings.Replace(path, placeholder, fmt.Sprintf("%v", v), 1)
			}
		}

		for k, v := range args {
			if pathParams[k] || knownArgs[k] {
				continue
			}
			switch method {
			case "POST", "PUT", "PATCH":
				bodyArgs[k] = v
			default:
				params[k] = fmt.Sprintf("%v", v)
			}
		}

		var data json.RawMessage
		switch method {
		case "GET":
			if binaryResponse {
				data, err = c.GetWithHeaders(path, params, headers)
				break
			}
			data, err = c.Get(path, params)
		case "POST":
			if binaryResponse {
				data, _, err = c.PostWithParamsAndHeaders(path, params, bodyArgs, headers)
				break
			}
			data, _, err = c.PostWithParams(path, params, bodyArgs)
		case "PUT":
			if binaryResponse {
				data, _, err = c.PutWithParamsAndHeaders(path, params, bodyArgs, headers)
				break
			}
			data, _, err = c.PutWithParams(path, params, bodyArgs)
		case "PATCH":
			if binaryResponse {
				data, _, err = c.PatchWithParamsAndHeaders(path, params, bodyArgs, headers)
				break
			}
			data, _, err = c.PatchWithParams(path, params, bodyArgs)
		case "DELETE":
			if binaryResponse {
				data, _, err = c.DeleteWithParamsAndHeaders(path, params, headers)
				break
			}
			data, _, err = c.DeleteWithParams(path, params)
		default:
			return mcplib.NewToolResultError("unsupported method: " + method), nil
		}

		if err != nil {
			msg := err.Error()
			switch {
			case strings.Contains(msg, "HTTP 409"):
				return mcplib.NewToolResultText("already exists (no-op)"), nil
			case strings.Contains(msg, "HTTP 400") && cliutil.LooksLikeAuthError(msg):
				return mcplib.NewToolResultError("authentication error: " + cliutil.SanitizeErrorBody(msg) +
					"\nhint: the API rejected the request — this usually means auth is missing or invalid." +
					"\n      Set your API key: export DOORLOOP_TOKEN=<your-key>" +
					"\n      Run 'doorloop-pp-cli doctor' to check auth status."), nil
			case strings.Contains(msg, "HTTP 401"):
				return mcplib.NewToolResultError("authentication failed: " + cliutil.SanitizeErrorBody(msg) +
					"\nhint: check your token." +
					"\n      Set it with: export DOORLOOP_TOKEN=<your-key>" +
					"\n      Run 'doorloop-pp-cli doctor' to check auth status."), nil
			case strings.Contains(msg, "HTTP 403"):
				return mcplib.NewToolResultError("permission denied: " + cliutil.SanitizeErrorBody(msg) +
					"\nhint: your credentials are valid but lack access to this resource." +
					"\n      Set it with: export DOORLOOP_TOKEN=<your-key>" +
					"\n      Run 'doorloop-pp-cli doctor' to check auth status."), nil
			case strings.Contains(msg, "HTTP 404"):
				if method == "DELETE" {
					return mcplib.NewToolResultText("already deleted (no-op)"), nil
				}
				return mcplib.NewToolResultError("not found: " + msg), nil
			case strings.Contains(msg, "HTTP 429"):
				return mcplib.NewToolResultError("rate limited: " + msg), nil
			default:
				return mcplib.NewToolResultError(msg), nil
			}
		}

		// For GET responses, wrap bare arrays with count metadata
		if method == "GET" {
			trimmed := strings.TrimSpace(string(data))
			if len(trimmed) > 0 && trimmed[0] == '[' {
				var items []json.RawMessage
				if json.Unmarshal(data, &items) == nil {
					wrapped := map[string]any{
						"count": len(items),
						"items": items,
					}
					out, _ := json.Marshal(wrapped)
					return mcplib.NewToolResultText(string(out)), nil
				}
			}
		}
		if binaryResponse {
			out, _ := json.Marshal(map[string]any{
				"content_encoding": "base64",
				"data_base64":      base64.StdEncoding.EncodeToString(data),
				"byte_count":       len(data),
			})
			return mcplib.NewToolResultText(string(out)), nil
		}
		return mcplib.NewToolResultText(string(data)), nil
	}
}

func newMCPClient() (*client.Client, error) {
	home, _ := os.UserHomeDir()
	cfgPath := filepath.Join(home, ".config", "doorloop-pp-cli", "config.toml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	c := client.New(cfg, 30*time.Second, 0)
	// Agents calling through MCP need fresh data every call. The on-disk
	// response cache survives across MCP server invocations, so a
	// DELETE/PATCH followed by a GET would otherwise return the
	// pre-mutation snapshot for up to the cache TTL. The interactive CLI
	// constructs its own client and is unaffected.
	c.NoCache = true
	return c, nil
}

func dbPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "doorloop-pp-cli", "data.db")
}

// Note: MCP tools use their own dbPath() because they are in a separate package (main, not cli).
// The CLI's defaultDBPath() in the cli package uses the same canonical path.

func handleSearch(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	args := req.GetArguments()
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcplib.NewToolResultError("query is required"), nil
	}

	limit := 25
	if v, ok := args["limit"].(float64); ok && v > 0 {
		limit = int(v)
	}

	db, err := store.OpenReadOnly(dbPath())
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("opening database: %v", err)), nil
	}
	defer db.Close()

	results, err := db.Search(query, limit)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("search failed: %v", err)), nil
	}

	data, _ := json.MarshalIndent(results, "", "  ")
	return mcplib.NewToolResultText(string(data)), nil
}

// validateReadOnlyQuery gates the MCP sql tool. The agent contract advertised
// to the host is ReadOnlyHintAnnotation(true); a false annotation on a
// mutating tool lets MCP hosts auto-approve writes and is treated as a real
// bug per the project's agent-native security model.
//
// The gate is an allowlist (SELECT or WITH only) applied AFTER stripping the
// leading whitespace, line comments, block comments, and semicolons that
// SQLite itself ignores before parsing. A naive HasPrefix check on a
// keyword blocklist is bypassable by prefixing the dangerous statement with
// "/* x */" or "-- x\n" — TrimSpace strips outer whitespace but does not
// understand SQL comment syntax. Combined with the empirical fact that
// modernc.org/sqlite's mode=ro does NOT block VACUUM INTO (writes a snapshot
// to a new file) or ATTACH DATABASE (opens a separate writable handle),
// such a bypass produces silent exfiltration to an attacker-chosen path.
//
// SELECT and WITH are the only allowed leading keywords. WITH supports
// SELECT-form CTEs; CTE-wrapped writes ("WITH x AS (...) INSERT ...") are
// caught by OpenReadOnly's mode=ro one layer down. PRAGMA, ATTACH, VACUUM,
// and every other DDL/DML keyword fail at this gate before reaching SQLite.
func validateReadOnlyQuery(query string) error {
	upper := strings.ToUpper(stripLeadingSQLNoise(query))
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "WITH") {
		return fmt.Errorf("only SELECT queries are allowed")
	}
	return nil
}

// stripLeadingSQLNoise removes leading whitespace, SQL line comments
// (-- to end of line), block comments (/* ... */), and statement
// separators (;) from query. SQLite skips these before parsing the first
// keyword, so a security gate that does not strip them mismatches what the
// driver actually executes.
func stripLeadingSQLNoise(query string) string {
	for {
		query = strings.TrimLeft(query, " \t\r\n;")
		switch {
		case strings.HasPrefix(query, "--"):
			if idx := strings.IndexByte(query, '\n'); idx >= 0 {
				query = query[idx+1:]
				continue
			}
			return ""
		case strings.HasPrefix(query, "/*"):
			if idx := strings.Index(query[2:], "*/"); idx >= 0 {
				query = query[2+idx+2:]
				continue
			}
			return ""
		default:
			return query
		}
	}
}

func handleSQL(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	args := req.GetArguments()
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return mcplib.NewToolResultError("query is required"), nil
	}

	if err := validateReadOnlyQuery(query); err != nil {
		return mcplib.NewToolResultError(err.Error()), nil
	}

	db, err := store.OpenReadOnly(dbPath())
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("opening database: %v", err)), nil
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return mcplib.NewToolResultError(fmt.Sprintf("query failed: %v", err)), nil
	}
	defer rows.Close()

	cols, _ := rows.Columns()
	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range values {
			ptrs[i] = &values[i]
		}
		rows.Scan(ptrs...)
		row := make(map[string]any)
		for i, col := range cols {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	data, _ := json.MarshalIndent(results, "", "  ")
	return mcplib.NewToolResultText(string(data)), nil
}

func handleContext(_ context.Context, _ mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	ctx := map[string]any{
		"api":         "doorloop",
		"description": "Every DoorLoop resource at your fingertips, plus offline delinquency reports, vacancy maps, and cash flow snapshots...",
		"archetype":   "payments",
		"tool_count":  93,
		// tool_surface tells agents which surface a capability lives on.
		"tool_surface": "MCP exposes typed endpoint tools plus a runtime mirror of user-facing CLI commands. Endpoint tools keep typed schemas; command-mirror tools shell out to the companion doorloop-pp-cli binary.",
		"auth": map[string]any{
			"type": "bearer_token",
			"env_vars": []map[string]any{
				{
					"name":        "DOORLOOP_TOKEN",
					"kind":        "per_call",
					"required":    true,
					"sensitive":   true,
					"description": "Set to your API credential.",
				},
			},
		},
		"resources": []map[string]any{
			{
				"name":        "accounts",
				"description": "Chart-of-accounts entries",
				"endpoints":   []string{"get", "list"},
				"syncable":    true,
			},
			{
				"name":        "communication_logs",
				"description": "Logged communications with tenants and owners",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "expenses",
				"description": "Property operating expenses",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "files",
				"description": "Documents and file attachments",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "lease_charges",
				"description": "Charges applied to leases (rent, fees, etc.)",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "lease_credits",
				"description": "Credits applied to leases",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "lease_payments",
				"description": "Rent payments recorded against leases",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "lease_returned_payments",
				"description": "Returned/bounced payments on leases",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "leases",
				"description": "Lease agreements linking tenants to units",
				"endpoints":   []string{"get", "list", "move_in", "move_out", "tenants"},
				"syncable":    true,
			},
			{
				"name":        "notes",
				"description": "Notes attached to records",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "owners",
				"description": "Property owners",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "portfolios",
				"description": "Portfolio groupings of properties",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "properties",
				"description": "Physical properties — residential or commercial",
				"endpoints":   []string{"get", "list"},
				"syncable":    true,
			},
			{
				"name":        "reports",
				"description": "Generated financial and operational reports",
				"endpoints":   []string{"list"},
				"syncable":    true,
			},
			{
				"name":        "tasks",
				"description": "Maintenance and operational tasks",
				"endpoints":   []string{"create", "delete", "get", "list", "post_update", "update"},
				"syncable":    true,
			},
			{
				"name":        "tenants",
				"description": "Tenants — active lease tenants and prospects",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "units",
				"description": "Individual rental units within properties",
				"endpoints":   []string{"get", "list"},
				"syncable":    true,
			},
			{
				"name":        "users",
				"description": "DoorLoop users (staff)",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "vendor_bills",
				"description": "Bills from vendors against properties",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "vendor_credits",
				"description": "Credits from vendors",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
			{
				"name":        "vendors",
				"description": "Vendors and contractors",
				"endpoints":   []string{"create", "delete", "get", "list", "update"},
				"syncable":    true,
			},
		},
		"query_tips": []string{
			"Pagination uses cursor-based paging. Pass after parameter for subsequent pages.",
			"Control page size with the limit parameter (default 100).",
			"Use the sql tool for ad-hoc analysis on synced data. Run sync first to populate the local database.",
			"Use the search tool for full-text search across all synced resources. Faster than iterating list endpoints.",
			"Prefer sql/search over repeated API calls when the data is already synced.",
		},
		// Command-mirror capabilities are exposed through MCP by shelling out
		// to the companion CLI binary.
		"command_mirror_capabilities": []map[string]string{
			{"name": "Sync", "command": "sync", "description": "Pull all DoorLoop entities into a local SQLite database for instant offline queries.", "rationale": "The Ruby and Python wrappers are stateless API clients — every query hits the network. Local sync unlocks compound...", "via": "mcp-command-mirror"},
			{"name": "Delinquency Report", "command": "delinquency", "description": "See every tenant with an outstanding balance, with their contact info and days past due — one command instead of...", "rationale": "DoorLoop filters leases by outstanding balance but doesn't join tenant contact info in a single call. Requires a...", "via": "mcp-command-mirror"},
			{"name": "Cash Flow Snapshot", "command": "cashflow", "description": "Get income, expenses, and net cash flow per property for any date range — faster than navigating DoorLoop's PDF...", "rationale": "Requires summing lease_payments (income) against expenses + vendor_bills (outgo) in local SQLite grouped by...", "via": "mcp-command-mirror"},
			{"name": "Lease Ledger", "command": "ledger", "description": "View a chronological bank-statement-style ledger for any lease: every charge, payment, credit, and returned payment...", "rationale": "Joins four separate financial tables (lease_charges, lease_payments, lease_credits, lease_returned_payments) on...", "via": "mcp-command-mirror"},
			{"name": "Vacancy Map", "command": "vacancy", "description": "Find every vacant unit across your portfolio — units with no active lease — with days vacant and last lease end...", "rationale": "DoorLoop has no vacancy API endpoint. Requires a LEFT JOIN of synced units against active leases in local SQLite to...", "via": "mcp-command-mirror"},
			{"name": "Expiring Leases Alert", "command": "expiring", "description": "List all leases expiring within N days with tenant contact info and whether a prospect already exists for the unit.", "rationale": "Queries synced leases by end_date window, joins tenants for contact info, and cross-references prospect tenants per...", "via": "mcp-command-mirror"},
			{"name": "Task Agenda", "command": "tasks agenda", "description": "See today's tasks bucketed as OVERDUE, DUE TODAY, and DUE THIS WEEK — with property and assignee — ready to act...", "rationale": "The DoorLoop API returns tasks with due_date but provides no urgency bucketing. Local SQLite bucketing by comparing...", "via": "mcp-command-mirror"},
			{"name": "Portfolio Health", "command": "portfolio-health", "description": "One-row-per-portfolio summary: total units, occupancy rate, total outstanding balance, and leases expiring within 30...", "rationale": "Aggregates units + leases + payments across three synced tables per portfolio. No single API call returns this...", "via": "mcp-command-mirror"},
		},
		"playbook": []map[string]string{
			{"topic": "Sync", "insight": "The Ruby and Python wrappers are stateless API clients — every query hits the network. Local sync unlocks compound queries (delinquency+expiry+vacancy) that require joining across multiple resources."},
			{"topic": "Delinquency Report", "insight": "DoorLoop filters leases by outstanding balance but doesn't join tenant contact info in a single call. Requires a local SQLite join across leases and tenants that no API endpoint provides."},
			{"topic": "Cash Flow Snapshot", "insight": "Requires summing lease_payments (income) against expenses + vendor_bills (outgo) in local SQLite grouped by property. No API aggregation endpoint exists; the reports API is PDF-only."},
			{"topic": "Lease Ledger", "insight": "Joins four separate financial tables (lease_charges, lease_payments, lease_credits, lease_returned_payments) on lease_id in local SQLite. No single API endpoint provides this combined view."},
			{"topic": "Vacancy Map", "insight": "DoorLoop has no vacancy API endpoint. Requires a LEFT JOIN of synced units against active leases in local SQLite to find units with no active lease match."},
			{"topic": "Expiring Leases Alert", "insight": "Queries synced leases by end_date window, joins tenants for contact info, and cross-references prospect tenants per unit — a three-table join no single API call can express."},
			{"topic": "Task Agenda", "insight": "The DoorLoop API returns tasks with due_date but provides no urgency bucketing. Local SQLite bucketing by comparing due_date to today produces the prioritized view maintenance coordinators need."},
			{"topic": "Portfolio Health", "insight": "Aggregates units + leases + payments across three synced tables per portfolio. No single API call returns this multi-entity shape; requires local SQLite aggregation."},
			{"topic": "Financial data", "insight": "Always use read-only operations for financial queries. Never use create/update tools for payment data without explicit user confirmation."},
			{"topic": "Reconciliation", "insight": "For reconciliation tasks, sync first then use sql for cross-referencing. API pagination over financial records is slow and rate-limited."},
		},
	}
	data, _ := json.MarshalIndent(ctx, "", "  ")
	return mcplib.NewToolResultText(string(data)), nil
}

// RegisterNovelFeatureTools is kept as a compatibility no-op for older MCP
// mains. New generated mains call RegisterTools only; RegisterTools now
// includes the runtime Cobra-tree mirror.
func RegisterNovelFeatureTools(s *server.MCPServer) {
	_ = s
}
