package schemas

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/khuongdo95/go-pkg/common/utils"
	entOrm "github.com/khuongdo95/go-pkg/orm/ent"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").DefaultFunc(func() string {
			return utils.GenerateUUIDV7()
		}),
		field.String("name").Optional().Nillable(),
		field.String("email").Optional().Nillable(),
		field.String("tenant_id"),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entOrm.TimeMixin{},
		entOrm.SoftDeleteMixin{},
		entOrm.ModifierMixin{},
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("identity", Identity.Type),
		edge.To("ip_white_list", UserIpWhiteList.Type),
	}
}
