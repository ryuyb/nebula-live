package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserPushSetting holds the schema definition for the UserPushSetting entity.
type UserPushSetting struct {
	ent.Schema
}

// Fields of the UserPushSetting.
func (UserPushSetting) Fields() []ent.Field {
	return []ent.Field{
		field.Uint("id").
			Unique().
			Immutable(),
		field.Uint("user_id").
			Comment("关联的用户ID"),
		field.Enum("provider").
			Values("bark").
			Comment("推送服务提供商"),
		field.Bool("enabled").
			Default(false).
			Comment("是否启用此推送设置"),
		field.String("device_id").
			NotEmpty().
			Comment("设备ID或推送标识符"),
		field.String("device_name").
			Optional().
			MaxLen(100).
			Comment("设备名称，用于用户识别"),
		field.JSON("settings", map[string]interface{}{}).
			Optional().
			Comment("提供商特定的设置，JSON格式存储"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the UserPushSetting.
func (UserPushSetting) Edges() []ent.Edge {
	return []ent.Edge{
		// 关联到用户
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id"),
	}
}

// Indexes of the UserPushSetting.
func (UserPushSetting) Indexes() []ent.Index {
	return []ent.Index{
		// 用户和提供商的组合索引，支持一个用户多个同类型设备
		index.Fields("user_id", "provider"),
		index.Fields("user_id"),
		index.Fields("provider"),
		index.Fields("enabled"),
		index.Fields("created_at"),
		// 设备ID的唯一性索引，防止重复添加同一设备
		index.Fields("provider", "device_id").Unique(),
	}
}