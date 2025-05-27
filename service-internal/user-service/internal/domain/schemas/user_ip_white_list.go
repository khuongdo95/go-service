package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	entOrm "github.com/khuongdo95/go-pkg/orm/ent"
)

type UserIpWhiteList struct {
	ent.Schema
}

func (UserIpWhiteList) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("ip_address"),
	}
}

func (UserIpWhiteList) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entOrm.TimeMixin{},
	}
}

func (UserIpWhiteList) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("ip_white_list").
			Field("user_id").
			Unique().
			Required(),
	}
}
