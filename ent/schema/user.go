package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable(),
		field.String("username").
			Unique().
			NotEmpty().
			MaxLen(50),
		field.String("email").
			Unique().
			NotEmpty().
			MaxLen(100),
		field.String("password").
			NotEmpty().
			Sensitive(), // 敏感字段，不会在日志中显示
		field.String("nickname").
			Optional().
			MaxLen(100),
		field.String("avatar").
			Optional().
			MaxLen(500),
		field.Enum("status").
			Values("active", "inactive", "banned").
			Default("active"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// 用户的用户角色关联
		edge.From("user_roles", UserRole.Type).
			Ref("user"),
		// 作为分配者的用户角色关联
		edge.From("assigned_user_roles", UserRole.Type).
			Ref("assigner"),
		// 作为分配者的角色权限关联
		edge.From("assigned_role_permissions", RolePermission.Type).
			Ref("assigner"),
		// 用户的推送设置
		edge.From("push_settings", UserPushSetting.Type).
			Ref("user"),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username"),
		index.Fields("email"),
		index.Fields("status"),
		index.Fields("created_at"),
	}
}
