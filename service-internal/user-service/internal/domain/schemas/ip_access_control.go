package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	entOrm "github.com/khuongdo95/go-pkg/orm/ent"
)

type IPAccessControl struct {
	ent.Schema
}

func (IPAccessControl) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Unique(),
		field.String("ip_address"),
		field.Bool("tcp_stage").Default(false),
		field.Bool("tcp_live").Default(false),
		field.Bool("rabbitmq_stage").Default(false),
		field.Bool("rabbitmq_live").Default(false),
		field.Bool("translations_rmq").Default(false),
		field.Bool("translations_tcp").Default(false),
		field.Bool("active").Default(true),
	}
}

func (IPAccessControl) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entOrm.TimeMixin{},
		entOrm.ModifierMixin{},
	}
}
