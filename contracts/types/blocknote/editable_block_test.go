package blocknote

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestArborizedEditableBlockRoundTripsExtendedBlockContent(t *testing.T) {
	blockID := uuid.New()
	payload := []byte(`{
		"id":"` + blockID.String() + `",
		"type":"paragraph",
		"props":{},
		"content":[{"type":"math","content":"x^2"}],
		"children":[]
	}`)

	var block ArborizedEditableBlock
	if err := json.Unmarshal(payload, &block); err != nil {
		t.Fatalf("unmarshal math inline content: %v", err)
	}

	inlineContent, ok := block.Content.(InlineContentList)
	if !ok || len(inlineContent) != 1 {
		t.Fatalf("expected one inline content item, got %#v", block.Content)
	}
	mathContent, ok := inlineContent[0].InlineContentUnion.(*Math)
	if !ok || mathContent.Content != "x^2" {
		t.Fatalf("expected math content x^2, got %#v", inlineContent[0].InlineContentUnion)
	}

	block.Type = "diagram"
	block.Props, _ = ParseProps("diagram", []byte(`{}`))
	block.Content = PlainContent("graph TD")

	encoded, err := json.Marshal(block)
	if err != nil {
		t.Fatalf("marshal calendar plain content: %v", err)
	}

	var roundTripped ArborizedEditableBlock
	if err := json.Unmarshal(encoded, &roundTripped); err != nil {
		t.Fatalf("unmarshal calendar plain content: %v", err)
	}
	plainContent, ok := roundTripped.Content.(PlainContent)
	if !ok || plainContent != "graph TD" {
		t.Fatalf("unexpected plain content: %#v", roundTripped.Content)
	}

	calendarProps, err := ParseProps("calendar", []byte(`{"calendarId":"work","anchorDate":"2026-08-01","timezone":"Asia/Taipei","view":"month"}`))
	if err != nil {
		t.Fatalf("parse calendar props: %v", err)
	}
	parsedCalendarProps := calendarProps.(*CalendarBlockProps)
	if parsedCalendarProps.CalendarId != "work" || parsedCalendarProps.AnchorDate != "2026-08-01" {
		t.Fatalf("unexpected calendar props: %#v", parsedCalendarProps)
	}
}
