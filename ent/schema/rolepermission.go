package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// RolePermission holds the schema definition for the RolePermission entity.
type RolePermission struct {
	ent.Schema
}

// Fields of the RolePermission.
func (RolePermission) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable(),
		field.Uint("role_id").
			Comment("角色ID"),
		field.Uint("permission_id").
			Comment("权限ID"),
		field.Uint("assigned_by").
			Optional().
			Comment("分配者的用户ID"),
		field.Time("assigned_at").
			Default(time.Now).
			Comment("分配时间"),
	}
}

// Edges of the RolePermission.
func (RolePermission) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("role", Role.Type).
			Required().
			Unique().
			Field("role_id"),
		edge.To("permission", Permission.Type).
			Required().
			Unique().
			Field("permission_id"),
		edge.To("assigner", User.Type).
			Field("assigned_by").
			Unique(),
	}
}

// Indexes of the RolePermission.
func (RolePermission) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role_id"),
		index.Fields("permission_id"),
		index.Fields("role_id", "permission_id").Unique(), // 确保角色权限组合唯一
		index.Fields("assigned_by"),
		index.Fields("assigned_at"),
	}
}
