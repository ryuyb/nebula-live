package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserRole holds the schema definition for the UserRole entity.
type UserRole struct {
	ent.Schema
}

// Fields of the UserRole.
func (UserRole) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable(),
		field.Uint("user_id").
			Comment("用户ID"),
		field.Uint("role_id").
			Comment("角色ID"),
		field.Uint("assigned_by").
			Optional().
			Comment("分配者的用户ID"),
		field.Time("assigned_at").
			Default(time.Now).
			Comment("分配时间"),
	}
}

// Edges of the UserRole.
func (UserRole) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
		edge.To("role", Role.Type).
			Required().
			Unique().
			Field("role_id"),
		edge.To("assigner", User.Type).
			Field("assigned_by").
			Unique(),
	}
}

// Indexes of the UserRole.
func (UserRole) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
		index.Fields("role_id"),
		index.Fields("user_id", "role_id").Unique(), // 确保用户角色组合唯一
		index.Fields("assigned_by"),
		index.Fields("assigned_at"),
	}
}
