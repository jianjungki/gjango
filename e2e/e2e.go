package e2e

import (
	"tiktok_tools/manager"
	"tiktok_tools/model"
)

// SetupDatabase creates the schema, populates it with data and returns with superadmin user
func SetupDatabase(m *manager.Manager) (*model.User, error) {
	models := manager.GetModels()
	m.CreateSchema(models...)
	m.CreateRoles()
	return m.CreateSuperAdmin("superuser@example.org", "testpassword")
}
