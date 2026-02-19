package views

import "mortenvistisen/services"

func markdownToHTML(content string) string {
	return services.MarkdownToHTML(content)
}
