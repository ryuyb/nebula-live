package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Permission holds the schema definition for the Permission entity.
type Permission struct {
	ent.Schema
}

// Fields of the Permission.
func (Permission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable(),
		field.String("name").
			Unique().
			NotEmpty().
			MaxLen(100).
			Comment("权限名称，如：user:read, user:write"),
		field.String("display_name").
			NotEmpty().
			MaxLen(100).
			Comment("显示名称，如：查看用户, 修改用户"),
		field.String("description").
			Optional().
			MaxLen(500).
			Comment("权限描述"),
		field.String("resource").
			NotEmpty().
			MaxLen(50).
			Comment("资源名称，如：user, post, system"),
		field.String("action").
			NotEmpty().
			MaxLen(50).
			Comment("操作名称，如：read, write, delete, manage"),
		field.Bool("is_system").
			Default(false).
			Comment("是否为系统权限（系统权限不可删除）"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Permission.
func (Permission) Edges() []ent.Edge {
	return []ent.Edge{
		// 权限的角色权限关联
		edge.From("role_permissions", RolePermission.Type).
			Ref("permission"),
	}
}

// Indexes of the Permission.
func (Permission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").Unique(),
		index.Fields("resource"),
		index.Fields("action"),
		index.Fields("resource", "action").Unique(),
		index.Fields("is_system"),
		index.Fields("created_at"),
	}
}