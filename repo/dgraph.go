package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	conn := os.Getenv("HLX_DGRAPH_SERVER")
	if conn == "" {
		conn = "localhost:9080"
	}

	// Dial a gRPC connection.
	d, err := grpc.Dial(conn, grpc.WithInsecure())
	if err != nil {
		log.Fatal("While trying to dial gRPC")
	}

	return dgo.NewDgraphClient(
		api.NewDgraphClient(d),
	)
}

func (r *Repo) SaveCategories(ctx context.Context, cs []models.Category) ([]models.Category, error) {
	c := r.NewClient()

	// drop all data
	err := c.Alter(context.Background(), &api.Operation{DropAll: true})
	if err != nil {
		return nil, err
	}

	// set schema.
	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
			category: string .
			id: int @index(int) .
			name: string @index(trigram) .
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

	// load new data.
	cb, err := json.Marshal(cs)
	if err != nil {
		return nil, err
	}

	mu := &api.Mutation{CommitNow: true}
	mu.SetJson = cb
	assigned, err := txn.Mutate(ctx, mu)
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

	q := fmt.Sprintf(`{
		categories(func: uid(%s)) @recurse(depth: 100){
			expand(_all_)
		}
	}`, uids)

	resp, err := c.NewReadOnlyTxn().Query(ctx, q)
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

func (r *Repo) GetCategoryLeaf(ctx context.Context, name string, subLayerName string) ([]models.Category, error) {
	c := r.NewClient()

	var q string
	if subLayerName == "" {
		q = fmt.Sprintf(`{
			categories(func: regexp(name, /%s/i)) @filter(NOT has(children)) {
				expand(_all_)
			}
		}`, name)
	} else {
		q = fmt.Sprintf(`{
			var(func: regexp(name, /%s/i)) @recurse(depth: 100) {
				A as children
			}
			
			categories(func: uid(A)) @filter(regexp(name, /%s/i) AND NOT has(children)) {
				expand(_all_)
			}
		}`, subLayerName, name)
	}

	resp, err := c.NewReadOnlyTxn().Query(ctx, q)
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

	// drop all data
	err := c.Alter(context.Background(), &api.Operation{DropAll: true})
	if err != nil {
		return nil, err
	}

	// set schema.
	err = c.Alter(context.Background(), &api.Operation{
		Schema: `
			values: uid .
			value: string @index(term) .
			id: string @index(exact) .
			addonId: string @index(term) .
			label: string @index(trigram) .
			ico: string .
			pic: string .
			sources: uid .
			description: string @index(fulltext) .
			psycho: string .
		`,
	})
	if err != nil {
		return nil, err
	}

	txn := c.NewTxn()
	defer txn.Discard(ctx)

	// load new data.
	pb, err := json.Marshal(ps)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(pb))

	mu := &api.Mutation{CommitNow: true}
	mu.SetJson = pb
	assigned, err := txn.Mutate(ctx, mu)
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

	q := fmt.Sprintf(`{
		psychos(func: uid(%s)) @recurse(depth: 100){
			expand(_all_)
		}
	}`, uids)

	resp, err := c.NewReadOnlyTxn().Query(ctx, q)
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

func (r *Repo) GetPsychos(ctx context.Context, label string, subLayerLabel string) ([]models.Psycho, error) {
	c := r.NewClient()

	var q string
	if subLayerLabel == "" {
		q = fmt.Sprintf(`{
			psychos(func: regexp(label, /%s/i)) @filter(NOT has(values)) {
				expand(_all_)
			}
		}`, label)
	} else {
		q = fmt.Sprintf(`{
			var(func: regexp(label, /%s/i)) @recurse(depth: 100) {
				A as values
			}
			
			psychos(func: uid(A)) @filter(regexp(label, /%s/i) AND NOT has(values)) {
				expand(_all_)
			}
		}`, subLayerLabel, label)
	}

	resp, err := c.NewReadOnlyTxn().Query(ctx, q)
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
