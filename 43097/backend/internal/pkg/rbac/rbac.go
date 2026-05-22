package rbac

const (
	RoleAdmin     = "admin"
	RoleFrontDesk = "front_desk"
	RoleUser      = "user"
)

type Permission struct {
	Resource string
	Action   string
}

var rolePermissions = map[string][]Permission{
	RoleAdmin: {
		{Resource: "*", Action: "*"},
	},
	RoleFrontDesk: {
		{Resource: "booking", Action: "create"},
		{Resource: "booking", Action: "read"},
		{Resource: "booking", Action: "update"},
		{Resource: "booking", Action: "delete"},
		{Resource: "checkin", Action: "create"},
		{Resource: "checkin", Action: "read"},
		{Resource: "checkin", Action: "update"},
		{Resource: "checkin", Action: "delete"},
		{Resource: "room", Action: "read"},
		{Resource: "member", Action: "read"},
		{Resource: "member", Action: "update"},
	},
	RoleUser: {
		{Resource: "booking", Action: "create"},
		{Resource: "booking", Action: "read"},
		{Resource: "room", Action: "read"},
		{Resource: "user", Action: "read"},
		{Resource: "user", Action: "update"},
	},
}

func CheckPermission(userRole string, resource, action string) bool {
	permissions, exists := rolePermissions[userRole]
	if !exists {
		return false
	}

	for _, perm := range permissions {
		if (perm.Resource == "*" || perm.Resource == resource) &&
			(perm.Action == "*" || perm.Action == action) {
			return true
		}
	}

	return false
}

func GetRolePermissions(role string) []Permission {
	return rolePermissions[role]
}

func IsValidRole(role string) bool {
	_, exists := rolePermissions[role]
	return exists
}

func HasAdminRole(role string) bool {
	return role == RoleAdmin
}

func HasFrontDeskRole(role string) bool {
	return role == RoleAdmin || role == RoleFrontDesk
}
