package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/sullivtr/terraform-provider-graphql/gql-server/graph/generated"
	"github.com/sullivtr/terraform-provider-graphql/gql-server/graph/model"
)

func (r *mutationResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	fmt.Println("Running create query")

	todo := model.Todo{
		Text: input.Text,
		ID:   fmt.Sprintf("T%d", rand.Int()),
		User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
	}

	content, _ := json.MarshalIndent(todo, "", " ")
	err := ioutil.WriteFile("./test.json", content, 0755)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	fmt.Println("Opening test json")

	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	fmt.Println(string(byteValue))

	var todoResult model.Todo
	_ = json.Unmarshal([]byte(byteValue), &todoResult)
	return &todoResult, nil
}

func (r *mutationResolver) UpdateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	todo := &model.Todo{
		Text: input.Text,
		ID:   fmt.Sprintf("T%d", rand.Int()),
		User: &model.User{ID: input.UserID, Name: "user " + input.UserID},
	}
	r.todos = append(r.todos, todo)
	return todo, nil
}

func (r *mutationResolver) DeleteTodo(ctx context.Context, input string) (*model.Todo, error) {
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var todoRead model.Todo
	_ = json.Unmarshal([]byte(byteValue), &todoRead)

	if todoRead.ID == input {
		content, _ := json.MarshalIndent(model.Todo{}, "", " ")
		_ = ioutil.WriteFile("./test.json", content, 0755)
	}

	byteValueReadAfterClear, _ := ioutil.ReadAll(jsonFile)

	var todoResultAfterClear model.Todo
	_ = json.Unmarshal([]byte(byteValueReadAfterClear), &todoResultAfterClear)
	return &todoResultAfterClear, nil
}

func (r *queryResolver) Todo(ctx context.Context) (*model.Todo, error) {
	jsonFile, err := os.Open("./test.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var todoResult model.Todo
	_ = json.Unmarshal([]byte(byteValue), &todoResult)
	return &todoResult, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
