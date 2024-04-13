package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	node := maelstrom.NewNode()

	node.Handle("generate", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		body["type"] = "generate_ok"
		uuid, _ := uuid.NewV7()
		sha := sha256.Sum256([]byte(msg.Src + uuid.String()))
		base64Encoded := make([]byte, 48)
		base64.StdEncoding.Encode(base64Encoded, sha[:])

		body["id"] = string(base64Encoded)
		body["in_reply_to"] = body["msg_id"]

		return node.Reply(msg, body)
	})

	if err := node.Run(); err != nil {
		log.Fatal(err)
	}
}
