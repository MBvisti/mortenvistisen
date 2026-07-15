import { useEffect, useRef, useState } from 'react'
import { TableKit } from '@tiptap/extension-table'
import { Markdown } from '@tiptap/markdown'
import { EditorContent, useEditor } from '@tiptap/react'
import StarterKit from '@tiptap/starter-kit'
import {
  Bold,
  Code,
  Code2,
  Heading2,
  Heading3,
  Heading4,
  Italic,
  Link2,
  List,
  ListOrdered,
  PanelLeftClose,
  PanelLeftOpen,
  PanelRightOpen,
  PanelTopClose,
  PanelTopOpen,
  PanelBottomOpen,
  Quote,
  Redo2,
  Table2,
  Trash2,
  Undo2,
  Unlink,
} from 'lucide-react'

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import { Button } from '@/components/ui/button'
import { Field, FieldError, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

type MarkdownEditorProps = {
  value: string
  onChange: (value: string) => void
  invalid?: boolean
  ariaLabel?: string
}

type ToolbarButtonProps = {
  label: string
  active?: boolean
  disabled?: boolean
  onClick: () => void
  children: React.ReactNode
}

function ToolbarButton({
  label,
  active = false,
  disabled = false,
  onClick,
  children,
}: ToolbarButtonProps) {
  return (
    <Button
      type="button"
      variant={active ? 'secondary' : 'ghost'}
      size="icon-sm"
      disabled={disabled}
      onClick={onClick}
      aria-label={label}
      title={label}
      aria-pressed={active}
    >
      {children}
    </Button>
  )
}

export function MarkdownEditor({
  value,
  onChange,
  invalid = false,
  ariaLabel = 'Markdown content',
}: MarkdownEditorProps) {
  const onChangeRef = useRef(onChange)
  const [linkDialogOpen, setLinkDialogOpen] = useState(false)
  const [linkHref, setLinkHref] = useState('')
  const [linkError, setLinkError] = useState('')
  onChangeRef.current = onChange

  const editor = useEditor({
    immediatelyRender: false,
    extensions: [
      StarterKit.configure({
        heading: {
          levels: [2, 3, 4],
        },
        codeBlock: {
          enableTabIndentation: true,
          tabSize: 2,
        },
      }),
      TableKit,
      Markdown,
    ],
    content: value,
    contentType: 'markdown',
    editorProps: {
      attributes: {
        class:
          'tiptap min-h-full px-5 py-4 text-sm leading-7 text-foreground outline-none',
        'aria-label': ariaLabel,
      },
    },
    onUpdate: ({ editor: currentEditor }) => {
      onChangeRef.current(currentEditor.getMarkdown())
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
    return <div className="h-192 animate-pulse rounded-none border border-border bg-muted/20" />
  }

  const inCodeBlock = editor.isActive('codeBlock')
  const codeLanguage = String(editor.getAttributes('codeBlock').language ?? '')
  const inTable = editor.isActive('table')
  const inLink = editor.isActive('link')

  function openLinkDialog() {
    setLinkHref(String(editor?.getAttributes('link').href ?? ''))
    setLinkError('')
    setLinkDialogOpen(true)
  }

  function saveLink() {
    const href = linkHref.trim()

    if (href === '') {
      setLinkError('Enter a link URL.')
      return
    }
    if (!editor?.chain().focus().extendMarkRange('link').setLink({ href }).run()) {
      setLinkError('Enter a valid link URL.')
      return
    }

    setLinkDialogOpen(false)
  }

  return (
    <div
      className={cn(
        'overflow-hidden rounded-none border border-input bg-transparent transition-colors focus-within:border-ring focus-within:ring-2 focus-within:ring-ring/30',
        invalid && 'border-destructive ring-2 ring-destructive/20',
      )}
    >
      <div className="flex flex-wrap items-center gap-1 border-b border-border bg-muted/30 p-2">
        <ToolbarButton
          label="Heading 2"
          active={editor.isActive('heading', { level: 2 })}
          onClick={() => editor.chain().focus().toggleHeading({ level: 2 }).run()}
        >
          <Heading2 />
        </ToolbarButton>
        <ToolbarButton
          label="Heading 3"
          active={editor.isActive('heading', { level: 3 })}
          onClick={() => editor.chain().focus().toggleHeading({ level: 3 }).run()}
        >
          <Heading3 />
        </ToolbarButton>
        <ToolbarButton
          label="Heading 4"
          active={editor.isActive('heading', { level: 4 })}
          onClick={() => editor.chain().focus().toggleHeading({ level: 4 }).run()}
        >
          <Heading4 />
        </ToolbarButton>
        <ToolbarButton
          label="Bold"
          active={editor.isActive('bold')}
          onClick={() => editor.chain().focus().toggleBold().run()}
        >
          <Bold />
        </ToolbarButton>
        <ToolbarButton
          label="Italic"
          active={editor.isActive('italic')}
          onClick={() => editor.chain().focus().toggleItalic().run()}
        >
          <Italic />
        </ToolbarButton>
        <ToolbarButton
          label="Inline code"
          active={editor.isActive('code')}
          onClick={() => editor.chain().focus().toggleCode().run()}
        >
          <Code />
        </ToolbarButton>
        <ToolbarButton
          label="Add or edit link"
          active={inLink}
          onClick={openLinkDialog}
        >
          <Link2 />
        </ToolbarButton>
        <ToolbarButton
          label="Remove link"
          disabled={!inLink}
          onClick={() => editor.chain().focus().extendMarkRange('link').unsetLink().run()}
        >
          <Unlink />
        </ToolbarButton>
        <span className="mx-1 h-6 w-px bg-border" aria-hidden="true" />
        <ToolbarButton
          label="Bullet list"
          active={editor.isActive('bulletList')}
          onClick={() => editor.chain().focus().toggleBulletList().run()}
        >
          <List />
        </ToolbarButton>
        <ToolbarButton
          label="Numbered list"
          active={editor.isActive('orderedList')}
          onClick={() => editor.chain().focus().toggleOrderedList().run()}
        >
          <ListOrdered />
        </ToolbarButton>
        <ToolbarButton
          label="Quote"
          active={editor.isActive('blockquote')}
          onClick={() => editor.chain().focus().toggleBlockquote().run()}
        >
          <Quote />
        </ToolbarButton>
        <ToolbarButton
          label="Code block"
          active={inCodeBlock}
          onClick={() => editor.chain().focus().toggleCodeBlock().run()}
        >
          <Code2 />
        </ToolbarButton>
        <input
          type="text"
          value={codeLanguage}
          disabled={!inCodeBlock}
          onChange={(event) =>
            editor.commands.updateAttributes('codeBlock', {
              language: event.target.value || null,
            })
          }
          className="h-8 w-28 border border-input bg-background px-2 font-mono text-xs outline-none placeholder:text-muted-foreground focus:border-ring disabled:cursor-not-allowed disabled:opacity-50"
          aria-label="Code block language"
          placeholder="language"
          title="Code block language, for example go or typescript"
        />
        <ToolbarButton
          label="Insert table"
          active={inTable}
          onClick={() => editor.chain().focus().insertTable().run()}
        >
          <Table2 />
        </ToolbarButton>
        <ToolbarButton
          label="Add column before"
          disabled={!inTable}
          onClick={() => editor.chain().focus().addColumnBefore().run()}
        >
          <PanelLeftOpen />
        </ToolbarButton>
        <ToolbarButton
          label="Add column after"
          disabled={!inTable}
          onClick={() => editor.chain().focus().addColumnAfter().run()}
        >
          <PanelRightOpen />
        </ToolbarButton>
        <ToolbarButton
          label="Delete column"
          disabled={!inTable}
          onClick={() => editor.chain().focus().deleteColumn().run()}
        >
          <PanelLeftClose />
        </ToolbarButton>
        <ToolbarButton
          label="Add row before"
          disabled={!inTable}
          onClick={() => editor.chain().focus().addRowBefore().run()}
        >
          <PanelTopOpen />
        </ToolbarButton>
        <ToolbarButton
          label="Add row after"
          disabled={!inTable}
          onClick={() => editor.chain().focus().addRowAfter().run()}
        >
          <PanelBottomOpen />
        </ToolbarButton>
        <ToolbarButton
          label="Delete row"
          disabled={!inTable}
          onClick={() => editor.chain().focus().deleteRow().run()}
        >
          <PanelTopClose />
        </ToolbarButton>
        <ToolbarButton
          label="Delete table"
          disabled={!inTable}
          onClick={() => editor.chain().focus().deleteTable().run()}
        >
          <Trash2 />
        </ToolbarButton>
        <span className="mx-1 h-6 w-px bg-border" aria-hidden="true" />
        <ToolbarButton
          label="Undo"
          disabled={!editor.can().chain().focus().undo().run()}
          onClick={() => editor.chain().focus().undo().run()}
        >
          <Undo2 />
        </ToolbarButton>
        <ToolbarButton
          label="Redo"
          disabled={!editor.can().chain().focus().redo().run()}
          onClick={() => editor.chain().focus().redo().run()}
        >
          <Redo2 />
        </ToolbarButton>
      </div>

      <EditorContent
        editor={editor}
        className="h-192 overflow-y-auto [&_.tableWrapper]:overflow-x-auto [&_.tiptap_a]:text-primary [&_.tiptap_a]:underline [&_.tiptap_a]:underline-offset-4 [&_.tiptap_blockquote]:my-4 [&_.tiptap_blockquote]:border-l-2 [&_.tiptap_blockquote]:border-primary/50 [&_.tiptap_blockquote]:pl-4 [&_.tiptap_code]:rounded-sm [&_.tiptap_code]:bg-muted [&_.tiptap_code]:px-1 [&_.tiptap_code]:py-0.5 [&_.tiptap_code]:font-mono [&_.tiptap_h2]:mb-4 [&_.tiptap_h2]:mt-6 [&_.tiptap_h2]:font-heading [&_.tiptap_h2]:text-2xl [&_.tiptap_h2]:font-semibold [&_.tiptap_h3]:mb-3 [&_.tiptap_h3]:mt-5 [&_.tiptap_h3]:font-heading [&_.tiptap_h3]:text-xl [&_.tiptap_h3]:font-semibold [&_.tiptap_h4]:mb-2 [&_.tiptap_h4]:mt-4 [&_.tiptap_h4]:font-heading [&_.tiptap_h4]:text-lg [&_.tiptap_h4]:font-semibold [&_.tiptap_li]:my-1 [&_.tiptap_ol]:my-4 [&_.tiptap_ol]:list-decimal [&_.tiptap_ol]:pl-6 [&_.tiptap_p]:my-3 [&_.tiptap_pre]:my-5 [&_.tiptap_pre]:overflow-x-auto [&_.tiptap_pre]:border [&_.tiptap_pre]:border-border [&_.tiptap_pre]:bg-slate-950 [&_.tiptap_pre]:p-4 [&_.tiptap_pre]:font-mono [&_.tiptap_pre]:text-sm [&_.tiptap_pre]:leading-6 [&_.tiptap_pre]:text-slate-100 [&_.tiptap_pre_code]:bg-transparent [&_.tiptap_pre_code]:p-0 [&_.tiptap_table]:w-full [&_.tiptap_table]:border-collapse [&_.tiptap_td]:border [&_.tiptap_td]:border-border [&_.tiptap_td]:p-2 [&_.tiptap_th]:border [&_.tiptap_th]:border-border [&_.tiptap_th]:bg-muted [&_.tiptap_th]:p-2 [&_.tiptap_th]:text-left [&_.tiptap_ul]:my-4 [&_.tiptap_ul]:list-disc [&_.tiptap_ul]:pl-6"
      />

      <AlertDialog open={linkDialogOpen} onOpenChange={setLinkDialogOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{inLink ? 'Edit link' : 'Add link'}</AlertDialogTitle>
            <AlertDialogDescription>
              Enter the destination URL for this link.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <div className="contents">
            <Field data-invalid={Boolean(linkError)}>
              <FieldLabel htmlFor="editorLinkHref">Link URL</FieldLabel>
              <Input
                id="editorLinkHref"
                type="text"
                value={linkHref}
                onChange={(event) => {
                  setLinkHref(event.currentTarget.value)
                  setLinkError('')
                }}
                onKeyDown={(event) => {
                  if (event.key === 'Enter') {
                    event.preventDefault()
                    saveLink()
                  }
                }}
                placeholder="https://example.com"
                aria-invalid={Boolean(linkError)}
                autoFocus
              />
              <FieldError>{linkError}</FieldError>
            </Field>
            <AlertDialogFooter>
              <AlertDialogCancel type="button">Cancel</AlertDialogCancel>
              <AlertDialogAction type="button" onClick={saveLink}>
                Save link
              </AlertDialogAction>
            </AlertDialogFooter>
          </div>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  )
}
