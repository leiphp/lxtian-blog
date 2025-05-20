package mysql

type TxyRolePermissions struct {
	RoleId uint64 `db:"role_id"`
	PermId uint64 `db:"perm_id"`
}
