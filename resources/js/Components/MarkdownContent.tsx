import { useEffect } from 'react'
import { TableKit } from '@tiptap/extension-table'
import { Markdown } from '@tiptap/markdown'
import { EditorContent, useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'

type MarkdownContentProps = {
  value: string
  ariaLabel?: string
}

export function MarkdownContent({
  value,
  ariaLabel = 'Rendered Markdown content',
}: MarkdownContentProps) {
  const editor = useEditor({
    immediatelyRender: false,
    editable: false,
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [2, 3, 4],
        },
      }),
      TableKit,
      Markdown,
    ],
    content: value,
    contentType: 'markdown',
    editorProps: {
      attributes: {
        class: 'tiptap text-sm leading-7 text-foreground outline-none',
        'aria-label': ariaLabel,
      },
    },
  })

  useEffect(() => {
    if (!editor || editor.getMarkdown() === value) {
      return
    }

    editor.commands.setContent(value, {
      contentType: 'markdown',
      emitUpdate: false,
    })
  }, [editor, value])

  if (!editor) {
    return <div className="h-24 animate-pulse bg-muted/20" aria-hidden="true" />
  }

  return (
    <EditorContent
      editor={editor}
      className="[&_.tableWrapper]:overflow-x-auto [&_.tiptap_a]:text-primary [&_.tiptap_a]:underline [&_.tiptap_a]:underline-offset-4 [&_.tiptap_blockquote]:my-4 [&_.tiptap_blockquote]:border-l-2 [&_.tiptap_blockquote]:border-primary/50 [&_.tiptap_blockquote]:pl-4 [&_.tiptap_code]:rounded-sm [&_.tiptap_code]:bg-muted [&_.tiptap_code]:px-1 [&_.tiptap_code]:py-0.5 [&_.tiptap_code]:font-mono [&_.tiptap_h2]:mb-4 [&_.tiptap_h2]:mt-6 [&_.tiptap_h2]:font-heading [&_.tiptap_h2]:text-2xl [&_.tiptap_h2]:font-semibold [&_.tiptap_h3]:mb-3 [&_.tiptap_h3]:mt-5 [&_.tiptap_h3]:font-heading [&_.tiptap_h3]:text-xl [&_.tiptap_h3]:font-semibold [&_.tiptap_h4]:mb-2 [&_.tiptap_h4]:mt-4 [&_.tiptap_h4]:font-heading [&_.tiptap_h4]:text-lg [&_.tiptap_h4]:font-semibold [&_.tiptap_li]:my-1 [&_.tiptap_ol]:my-4 [&_.tiptap_ol]:list-decimal [&_.tiptap_ol]:pl-6 [&_.tiptap_p]:my-3 [&_.tiptap_pre]:my-5 [&_.tiptap_pre]:overflow-x-auto [&_.tiptap_pre]:border [&_.tiptap_pre]:border-border [&_.tiptap_pre]:bg-slate-950 [&_.tiptap_pre]:p-4 [&_.tiptap_pre]:font-mono [&_.tiptap_pre]:text-sm [&_.tiptap_pre]:leading-6 [&_.tiptap_pre]:text-slate-100 [&_.tiptap_pre_code]:bg-transparent [&_.tiptap_pre_code]:p-0 [&_.tiptap_table]:w-full [&_.tiptap_table]:border-collapse [&_.tiptap_td]:border [&_.tiptap_td]:border-border [&_.tiptap_td]:p-2 [&_.tiptap_th]:border [&_.tiptap_th]:border-border [&_.tiptap_th]:bg-muted [&_.tiptap_th]:p-2 [&_.tiptap_th]:text-left [&_.tiptap_ul]:my-4 [&_.tiptap_ul]:list-disc [&_.tiptap_ul]:pl-6"
    />
  )
}
