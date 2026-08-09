package renderers

import (
	emailconfig "github.com/HiIamJeff67/notezy-backend/internal/email/configs"
)

type MarkdownEmailRenderer struct {
	Renderer
}

func NewMarkdownEmailRenderer(config emailconfig.RendererConfig) RendererInterface {
	return &MarkdownEmailRenderer{Renderer: newRenderer(config, "md")}
}
