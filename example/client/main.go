package main

import (
	"context"
	"fmt"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

func main() {
	// just for this demo
	ctx := context.Background()
	json := `{ "message": "hello" }`
	data := []byte(json)
	store := "statestore"

	// create the client
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// publish a message to the topic messagebus
	if err := client.PublishEvent(ctx, "messagebus", "demo", data); err != nil {
		panic(err)
	}
	fmt.Println("data published")

	// save state with the key key1
	fmt.Printf("saving data: %s\n", string(data))
	if err := client.SaveState(ctx, store, "key1", data); err != nil {
		panic(err)
	}
	fmt.Println("data saved")

	// get state for key key1
	item, err := client.GetState(ctx, store, "key1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("data retrieved [key:%s etag:%s]: %s\n", item.Key, item.Etag, string(item.Value))

	// save state with options
	item2 := &dapr.SetStateItem{
		Etag: "2",
		Key:  item.Key,
		Metadata: map[string]string{
			"created-on": time.Now().UTC().String(),
		},
		Value: item.Value,
		Options: &dapr.StateOptions{
			Concurrency: dapr.StateConcurrencyLastWrite,
			Consistency: dapr.StateConsistencyStrong,
		},
	}
	if err := client.SaveBulkState(ctx, store, item2); err != nil {
		panic(err)
	}
	fmt.Println("data item saved")

	// delete state for key key1
	if err := client.DeleteState(ctx, store, "key1"); err != nil {
		panic(err)
	}
	fmt.Println("data deleted")

	// invoke a method called EchoMethod on another dapr enabled service
	content := &dapr.DataContent{
		ContentType: "text/plain",
		Data:        []byte("hellow"),
	}
	resp, err := client.InvokeMethodWithContent(ctx, "serving", "echo", "post", content)
	if err != nil {
		panic(err)
	}
	fmt.Printf("service method invoked, response: %s", string(resp))

	in := &dapr.InvokeBindingRequest{
		Name:      "example-http-binding",
		Operation: "create",
	}
	if err := client.InvokeOutputBinding(ctx, in); err != nil {
		panic(err)
	}
	fmt.Println("output binding invoked")
	fmt.Println("DONE (CTRL+C to Exit)")
}
