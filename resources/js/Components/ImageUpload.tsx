import { useEffect, useRef, useState } from 'react'
import { ImageIcon, Trash2 } from 'lucide-react'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Field, FieldDescription, FieldError, FieldLabel } from '@/components/ui/field'

const maxImageBytes = 10 * 1024 * 1024

type ImageUploadProps = {
  id: string
  title: string
  description: string
  previewAlt: string
  existingURL: string
  file: File | null
  remove: boolean
  required?: boolean
  error?: string
  onFileChange: (file: File | null) => void
  onRemoveChange: (remove: boolean) => void
  onError: (error?: string) => void
}

export function ImageUpload({
  id,
  title,
  description,
  previewAlt,
  existingURL,
  file,
  remove,
  required = false,
  error,
  onFileChange,
  onRemoveChange,
  onError,
}: ImageUploadProps) {
  const [previewURL, setPreviewURL] = useState(existingURL)
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    if (remove) {
      setPreviewURL('')
      return
    }
    if (!file) {
      setPreviewURL(existingURL)
      return
    }
    const objectURL = URL.createObjectURL(file)
    setPreviewURL(objectURL)
    return () => URL.revokeObjectURL(objectURL)
  }, [existingURL, file, remove])

  function removeDisplayedImage() {
    if (file) {
      onFileChange(null)
      onRemoveChange(false)
    } else if (existingURL) {
      onRemoveChange(true)
    }
    onError()
    if (inputRef.current) {
      inputRef.current.value = ''
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <input
          ref={inputRef}
          id={id}
          type="file"
          accept="image/jpeg,image/png,image/webp,image/avif,image/gif"
          className="peer sr-only"
          onChange={(event) => {
            const selected = event.currentTarget.files?.[0] ?? null
            if (selected && selected.size > maxImageBytes) {
              onFileChange(null)
              onError('The selected image must be 10 MB or smaller.')
              event.currentTarget.value = ''
              return
            }
            onFileChange(selected)
            if (selected) {
              onRemoveChange(false)
            }
            onError()
          }}
        />
        {previewURL ? (
          <div className="group relative overflow-hidden border border-border bg-muted">
            <img src={previewURL} alt={previewAlt} className="aspect-video w-full object-cover" />
            <Button
              type="button"
              variant="destructive"
              size="icon-xs"
              className="absolute right-2 top-2 opacity-0 shadow-md transition-opacity group-hover:opacity-100 group-focus-within:opacity-100"
              disabled={required && !file}
              aria-label={file ? `Discard selected ${title.toLowerCase()}` : `Remove ${title.toLowerCase()}`}
              title={required && !file ? `Unpublish before removing the ${title.toLowerCase()}.` : undefined}
              onClick={removeDisplayedImage}
            >
              <Trash2 aria-hidden="true" />
            </Button>
          </div>
        ) : (
          <label
            htmlFor={id}
            className="flex aspect-video cursor-pointer flex-col items-center justify-center gap-3 border border-dashed border-border bg-muted px-4 text-center text-muted-foreground transition-colors hover:border-primary/60 hover:text-foreground peer-focus-visible:border-primary peer-focus-visible:ring-2 peer-focus-visible:ring-ring/30"
          >
            <ImageIcon className="size-8" aria-hidden="true" />
            <span className="text-sm font-medium">Click to add an image</span>
            <span className="text-xs">JPEG, PNG, WebP, AVIF, or GIF. Maximum 10 MB.</span>
          </label>
        )}
        <Field data-invalid={Boolean(error)}>
          <FieldLabel htmlFor={id} className="sr-only">{title}</FieldLabel>
          {file && <FieldDescription className="truncate">Selected: {file.name}</FieldDescription>}
          {required && existingURL && !file && (
            <FieldDescription>Unpublish before removing this image.</FieldDescription>
          )}
          {remove && (
            <div className="flex items-center justify-between gap-3 border border-destructive/30 bg-destructive/10 p-3 text-xs text-destructive">
              <span>The image will be deleted when changes are saved.</span>
              <Button type="button" variant="ghost" size="xs" onClick={() => onRemoveChange(false)}>
                Undo
              </Button>
            </div>
          )}
          <FieldError>{error}</FieldError>
        </Field>
      </CardContent>
    </Card>
  )
}
