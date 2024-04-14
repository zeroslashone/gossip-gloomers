package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	messages := []int{}
	node := maelstrom.NewNode()

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages = append(messages, int(body["message"].(float64)))

		body["type"] = "broadcast_ok"
		body["in_reply_to"] = body["msg_id"]
		delete(body, "message")

		return node.Reply(msg, body)
	})

	node.Handle("read", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["messages"] = messages
		body["type"] = "read_ok"
		body["in_reply_to"] = body["msg_id"]

		return node.Reply(msg, body)
	})

	node.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		body["type"] = "topology_ok"
		body["in_reply_to"] = body["msg_id"]
		delete(body, "topology")

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
