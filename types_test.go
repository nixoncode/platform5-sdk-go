package platform5

import (
	"encoding/json"
	"testing"
)

func TestEnvelopeDecode(t *testing.T) {
	raw := json.RawMessage(`{
		"success": true,
		"message": "OK",
		"data": {"message_id": "m1", "to": "+2547", "sender_name": "B", "parts": 1, "cost": 1.0, "currency": "KES", "status": "queued"}
	}`)

	resp, err := decode[*SendSMSResponse](raw)
	if err != nil {
		t.Fatal(err)
	}
	if resp.MessageID != "m1" {
		t.Fatalf("expected m1, got %s", resp.MessageID)
	}
}

func TestJSONTags(t *testing.T) {
	raw := `{"to":"+2547","message":"Hello","from":"MyBrand"}`
	var req SendSMSRequest
	if err := json.Unmarshal([]byte(raw), &req); err != nil {
		t.Fatal(err)
	}
	if req.To != "+2547" {
		t.Fatalf("expected +2547, got %s", req.To)
	}
}
