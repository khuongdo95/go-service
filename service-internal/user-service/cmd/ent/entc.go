//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	// ex, err := entproto.NewExtension(
	// 	entproto.WithProtoDir("../../go-common-library/proto/user"), // adjust relative path if needed
	// )
	// if err != nil {
	// 	log.Fatalf("running ent codegen: %v", err)
	// }

	opts := []entc.Option{
		// entc.Extensions(ex),
	}

	if err := entc.Generate("../../internal/domain/schemas", &gen.Config{
		Target:  "../../internal/generated/ent",
		Package: "github.com/khuongdo95/go-service/internal/generated/ent",
		Features: []gen.Feature{
			gen.FeatureIntercept,
			gen.FeatureUpsert,
			gen.FeatureVersionedMigration,
		},
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
