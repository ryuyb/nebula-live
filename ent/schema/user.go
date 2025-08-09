package schema

import (
	"time"

	"entgo.io/ent"
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
	return nil
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
