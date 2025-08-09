package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable(),
		field.String("name").
			Unique().
			NotEmpty().
			MaxLen(50).
			Comment("角色名称，如：admin, user"),
		field.String("display_name").
			NotEmpty().
			MaxLen(100).
			Comment("显示名称，如：管理员, 普通用户"),
		field.String("description").
			Optional().
			MaxLen(500).
			Comment("角色描述"),
		field.Bool("is_system").
			Default(false).
			Comment("是否为系统角色（系统角色不可删除）"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return []ent.Edge{
		// 角色的用户角色关联
		edge.From("user_roles", UserRole.Type).
			Ref("role"),
		// 角色的权限关联
		edge.From("role_permissions", RolePermission.Type).
			Ref("role"),
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("is_system"),
		index.Fields("created_at"),
	}
}
