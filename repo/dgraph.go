package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"
	"github.com/kind84/hlx/models"
	"google.golang.org/grpc"
)

type Repo struct {
	ConnStr string
}

func (r *Repo) NewClient() *dgo.Dgraph {
	// Dial a gRPC connection.
	d, err := grpc.Dial("localhost:9080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func (r *Repo) SaveCategories(ctx context.Context, cs []models.Category) ([]models.Category, error) {
	c := r.NewClient()

	// set schema.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			category: string .
			id: int @index(int) .
			name: string @index(term) .
			level: int @index(int) .
			children: uid .
			pic: string .
			type: string @index(term) .
		`,
	})
	if err != nil {
		return nil, err
	}

	txn := c.NewTxn()
	defer txn.Discard(ctx)

	// drop existing nodes.
	q := `{
		UIDS(func: has(category)) {
			uid
		}
	}`

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	type decode struct {
		UIDS []struct {
			UID string `json:"uid"`
		} `json:"uids"`
	}

	var d decode
	err = json.Unmarshal(resp.Json, &d)
	if err != nil {
		return nil, err
	}

	if len(d.UIDS) > 0 {
		db, err := json.Marshal(d.UIDS)
		if err != nil {
			return nil, err
		}

		del := &api.Mutation{DeleteJson: db}
		_, err = txn.Mutate(ctx, del)
		if err != nil {
			return nil, err
		}
	}

	// load new data.
	cb, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{}
	mu.SetJson = cb
	assigned, err := txn.Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return nil, err
	}

	// send loaded data in response.
	var sb strings.Builder
	for _, u := range assigned.Uids {
		sb.WriteString(u + ", ")
	}
	uids := sb.String()
	uids = uids[:sb.Len()-2]

	q = fmt.Sprintf(`{
		categories(func: uid(%s)) @recurse(depth: 100){
			uid
			id
			name
			level
			children
			pic
			type
			category
		}
	}`, uids)

	resp, err = c.NewReadOnlyTxn().Query(ctx, q)
	if err != nil {
		return nil, err
	}

	type Root struct {
		Categories []models.Category `json:"categories"`
	}

	var rt Root
	err = json.Unmarshal(resp.Json, &rt)
	if err != nil {
		return nil, err
	}

	return rt.Categories, nil
}

func (r *Repo) GetCategories(ctx context.Context, name string, level *int) ([]models.Category, error) {
	c := r.NewClient()

	vars := map[string]string{"$name": name}
	var q string
	if level == nil {
		q = `query Categories($name: string, $level: string){
			categories(func: allofterms(name, $name)) @recurse(depth: 100) {
				uid
				id
				name
				level 
				children
				pic
				type
			}
		}`
	} else {
		vars["$level"] = strconv.Itoa(*level)
		q = `query Categories($name: string, $level: string){
			categories(func: allofterms(name, $name)) @filter(eq(level, $level)) @recurse(depth: 100) {
				uid
				id
				name
				level 
				children
				pic
				type
			}
		}`
	}

	resp, err := c.NewReadOnlyTxn().QueryWithVars(ctx, q, vars)
	if err != nil {
		return nil, err
	}

	type Root struct {
		Categories []models.Category `json:"categories"`
	}

	var rt Root
	err = json.Unmarshal(resp.Json, &rt)
	if err != nil {
		return nil, err
	}

	return rt.Categories, nil
}

func (r *Repo) SavePsychos(ctx context.Context, ps []models.Psycho) ([]models.Psycho, error) {
	c := r.NewClient()

	// set schema.
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
			values: uid .
			value: string @index(term) .
			id: string @index(exact) .
			addonId: string @index(term) .
			label: string @index(term) .
			ico: string .
			pic: string .
			sources: uid .
			psycho: string .
		`,
	})
	if err != nil {
		return nil, err
	}

	txn := c.NewTxn()
	defer txn.Discard(ctx)

	// drop existing nodes.
	q := `{
		UIDS(func: has(psycho)) {
			uid
		}
	}`

	resp, err := txn.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	type decode struct {
		UIDS []struct {
			UID string `json:"uid"`
		} `json:"uids"`
	}

	var d decode
	err = json.Unmarshal(resp.Json, &d)
	if err != nil {
		return nil, err
	}

	if len(d.UIDS) > 0 {
		db, err := json.Marshal(d.UIDS)
		if err != nil {
			return nil, err
		}

		del := &api.Mutation{DeleteJson: db}
		_, err = txn.Mutate(ctx, del)
		if err != nil {
			return nil, err
		}
	}

	// load new data.
	pb, err := json.Marshal(ps)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{}
	mu.SetJson = pb
	assigned, err := txn.Mutate(ctx, mu)
	if err != nil {
		return nil, err
	}

	err = txn.Commit(ctx)
	if err != nil {
		return nil, err
	}

	// send loaded data in response.
	var sb strings.Builder
	for _, u := range assigned.Uids {
		sb.WriteString(u + ", ")
	}
	uids := sb.String()
	uids = uids[:sb.Len()-2]

	q = fmt.Sprintf(`{
		psychos(func: uid(%s)) @recurse(depth: 100){
			values
			value
			id
			addonId
			label
			ico
			pic
			sources {
				id
				description
			}
			psychos
		}
	}`, uids)

	resp, err = c.NewReadOnlyTxn().Query(ctx, q)
	if err != nil {
		return nil, err
	}

	type Root struct {
		Psychos []models.Psycho `json:"psychos"`
	}

	var rt Root
	err = json.Unmarshal(resp.Json, &rt)
	if err != nil {
		return nil, err
	}

	return rt.Psychos, nil
}
