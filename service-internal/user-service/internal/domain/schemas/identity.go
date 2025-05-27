package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	entOrm "github.com/khuongdo95/go-pkg/orm/ent"
)

type Identity struct {
	ent.Schema
}

func (Identity) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("username").Optional(),
		field.String("password").Optional(),
	}
}

func (Identity) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entOrm.TimeMixin{},
	}
}

func (Identity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("identity").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (Identity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
	}
}
