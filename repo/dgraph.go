package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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
	// Dial a gRPC connection. The address to dial to can be configured when
	// setting up the dgraph cluster.
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
	txn := c.NewTxn()
	defer txn.Discard(ctx)

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
