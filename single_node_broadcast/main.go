package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type TopologyBody struct {
	MsgId    int                 `json:"msg_id"`
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func main() {
	messages := []int{}
	node := maelstrom.NewNode()
	neighbors := []string{}

	node.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		message := int(body["message"].(float64))
		present := false
		for _, m := range messages {
			if m == message {
				present = true
				break
			}
		}
		if !present {
			messages = append(messages, message)

			for _, nid := range neighbors {
				node.Send(nid, msg.Body)
			}
		}

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
		var body TopologyBody
		respBody := map[string]any{}

		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}
		neighbors = body.Topology[node.ID()]
		log.Println("Neighbors are: ", neighbors)

		respBody["type"] = "topology_ok"
		respBody["in_reply_to"] = body.MsgId
		respBody["msg_id"] = body.MsgId

		return node.Reply(msg, respBody)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
