import { useEffect, useState, type FormEvent } from 'react'
import { Link, type InertiaFormProps } from '@inertiajs/react'

import type { ProjectFormData, ProjectValidationRules } from './types'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from '@/components/ui/collapsible'
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'
import { ImageUpload } from '@/Components/ImageUpload'
import { MarkdownEditor } from '@/Components/MarkdownEditor'

function characterLength(value: string) {
  return Array.from(value).length
}

function slugify(value: string) {
  return value
    .normalize('NFKD')
    .replace(/[\u0300-\u036f]/g, '')
    .toLowerCase()
    .trim()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

function findRule(rules: ProjectValidationRules, field: string, code: string) {
  return rules[field]?.find((rule) => rule.code === code)
}

function ruleNumber(
  rules: ProjectValidationRules,
  field: string,
  code: string,
  param: string,
) {
  const value = findRule(rules, field, code)?.params?.[param]
  if (typeof value !== 'number') {
    throw new Error(`Missing numeric validation rule ${field}.${code}.${param}`)
  }

  return value
}

function MaximumLengthIndicator({ value, max }: { value: string; max: number }) {
  const length = characterLength(value)
  const tooLong = length > max

  return (
    <div className="flex flex-wrap items-center justify-between gap-2">
      <FieldDescription>Maximum {max} characters.</FieldDescription>
      <span
        className={`text-xs font-medium tabular-nums ${tooLong ? 'text-destructive' : 'text-muted-foreground'}`}
        aria-live="polite"
      >
        {length}/{max}{tooLong ? ' · Too long' : ''}
      </span>
    </div>
  )
}

function LengthIndicator({
  value,
  recommendedMin,
  recommendedMax,
  max,
}: {
  value: string
  recommendedMin: number
  recommendedMax: number
  max: number
}) {
  const length = characterLength(value)
  const tooLong = length > max

  let status = 'Valid'
  let statusClass = 'text-amber-400'

  if (tooLong) {
    status = 'Too long'
    statusClass = 'text-destructive'
  } else if (length === 0) {
    status = 'Optional'
    statusClass = 'text-muted-foreground'
  } else if (length >= recommendedMin && length <= recommendedMax) {
    status = 'Ideal length'
    statusClass = 'text-emerald-400'
  }

  return (
    <div className="flex flex-wrap items-center justify-between gap-2">
      <FieldDescription>
        Aim for {recommendedMin}-{recommendedMax} characters.
      </FieldDescription>
      <span className={`text-xs font-medium tabular-nums ${statusClass}`} aria-live="polite">
        {length}/{max} · {status}
      </span>
    </div>
  )
}

type ProjectFormProps = {
  form: InertiaFormProps<ProjectFormData>
  validationRules: ProjectValidationRules
  cancelHref: string
  submitLabel: string
  processingLabel: string
  existingImageURL: string
  existingMetaImageURL: string
  onSubmit: () => void
}

export function ProjectForm({
  form,
  validationRules,
  cancelHref,
  submitLabel,
  processingLabel,
  existingImageURL,
  existingMetaImageURL,
  onSubmit,
}: ProjectFormProps) {
  const [detailsOpen, setDetailsOpen] = useState(false)
  const [metadataOpen, setMetadataOpen] = useState(false)
  const [slugIsAutomatic, setSlugIsAutomatic] = useState(
    () => form.data.slug === '' || form.data.slug === slugify(form.data.title),
  )
  const publishing = form.data.published
  const titleRequired = Boolean(findRule(validationRules, 'title', 'required'))
  const titleMaxLength = ruleNumber(validationRules, 'title', 'max', 'max')
  const statusMaxLength = ruleNumber(validationRules, 'status', 'max', 'max')
  const metaTitleMaxLength = ruleNumber(validationRules, 'metaTitle', 'max', 'max')
  const metaTitleRecommendedMinLength = ruleNumber(
    validationRules,
    'metaTitle',
    'recommended_length',
    'min',
  )
  const metaTitleRecommendedMaxLength = ruleNumber(
    validationRules,
    'metaTitle',
    'recommended_length',
    'max',
  )
  const metaDescriptionMaxLength = ruleNumber(
    validationRules,
    'metaDescription',
    'max',
    'max',
  )
  const metaDescriptionRecommendedMinLength = ruleNumber(
    validationRules,
    'metaDescription',
    'recommended_length',
    'min',
  )
  const metaDescriptionRecommendedMaxLength = ruleNumber(
    validationRules,
    'metaDescription',
    'recommended_length',
    'max',
  )
  const projectURLMaxLength = ruleNumber(validationRules, 'projectUrl', 'max', 'max')
  const titleInvalid =
    (titleRequired && form.data.title.trim() === '') ||
    characterLength(form.data.title) > titleMaxLength ||
    Boolean(form.errors.title)
  const statusInvalid =
    characterLength(form.data.status) > statusMaxLength || Boolean(form.errors.status)
  const metaTitleInvalid =
    characterLength(form.data.metaTitle) > metaTitleMaxLength ||
    Boolean(form.errors.metaTitle)
  const metaDescriptionInvalid =
    characterLength(form.data.metaDescription) > metaDescriptionMaxLength ||
    Boolean(form.errors.metaDescription)
  const projectURLInvalid =
    characterLength(form.data.projectUrl) > projectURLMaxLength ||
    Boolean(form.errors.projectUrl)
  const publishRequirements = [
    { label: 'Start date', complete: form.data.startedAt !== '' },
    { label: 'Status', complete: form.data.status.trim() !== '' },
    { label: 'Description', complete: form.data.description.trim() !== '' },
    { label: 'Meta title', complete: form.data.metaTitle.trim() !== '' },
    { label: 'Meta description', complete: form.data.metaDescription.trim() !== '' },
    {
      label: 'Featured image',
      complete: Boolean(form.data.cover || (existingImageURL && !form.data.removeCover)),
    },
    { label: 'Content', complete: form.data.content.trim() !== '' },
  ]

  useEffect(() => {
    if (form.errors.title || form.errors.slug || form.errors.description) {
      setDetailsOpen(true)
    }

    if (form.errors.metaTitle || form.errors.metaDescription) {
      setMetadataOpen(true)
    }
  }, [
    form.errors.description,
    form.errors.metaDescription,
    form.errors.metaTitle,
    form.errors.slug,
    form.errors.title,
  ])

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    onSubmit()
  }

  return (
    <form onSubmit={submit} className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
      <div className="flex min-w-0 flex-col gap-6">
        <Collapsible open={detailsOpen} onOpenChange={setDetailsOpen} render={<Card />}>
          <CardHeader>
            <CardTitle>Project details</CardTitle>
            <CardDescription>
              Set the title, URL slug, and short summary visitors see first.
            </CardDescription>
            <CardAction>
              <CollapsibleTrigger render={<Button variant="ghost" size="xs" />}>
                {detailsOpen ? 'Collapse' : 'Expand'}
              </CollapsibleTrigger>
            </CardAction>
          </CardHeader>
          <CollapsibleContent>
            <CardContent>
              <FieldGroup className="gap-6">
                <Field data-invalid={titleInvalid}>
                  <FieldLabel htmlFor="title">Title</FieldLabel>
                  <Input
                    id="title"
                    value={form.data.title}
                    onChange={(event) => {
                      const title = event.currentTarget.value
                      form.setData('title', title)
                      form.clearErrors('title')

                      if (slugIsAutomatic) {
                        form.setData('slug', slugify(title))
                        form.clearErrors('slug')
                      }
                    }}
                    placeholder="A clear, descriptive project name"
                    aria-invalid={titleInvalid}
                  />
                  <MaximumLengthIndicator value={form.data.title} max={titleMaxLength} />
                  <FieldError>{form.errors.title}</FieldError>
                </Field>

                <Field data-invalid={Boolean(form.errors.slug)}>
                  <FieldLabel htmlFor="slug">Slug</FieldLabel>
                  <Input
                    id="slug"
                    value={form.data.slug}
                    onChange={(event) => {
                      const slug = event.currentTarget.value
                      form.setData('slug', slug)
                      form.clearErrors('slug')
                      setSlugIsAutomatic(slug === slugify(form.data.title))
                    }}
                    placeholder="a-clear-project-name"
                    aria-invalid={Boolean(form.errors.slug)}
                  />
                  <FieldDescription>
                    Generated from the title until you customize it.
                  </FieldDescription>
                  <FieldError>{form.errors.slug}</FieldError>
                </Field>

                <Field data-invalid={Boolean(form.errors.description)}>
                  <FieldLabel htmlFor="description">Description</FieldLabel>
                  <Textarea
                    id="description"
                    value={form.data.description}
                    onChange={(event) => {
                      form.setData('description', event.currentTarget.value)
                      form.clearErrors('description')
                    }}
                    placeholder="A concise introduction to the project"
                    className="min-h-24 resize-y"
                    required={publishing}
                    aria-invalid={Boolean(form.errors.description)}
                  />
                  <FieldError>{form.errors.description}</FieldError>
                </Field>
              </FieldGroup>
            </CardContent>
          </CollapsibleContent>
        </Collapsible>

        <Collapsible open={metadataOpen} onOpenChange={setMetadataOpen} render={<Card />}>
          <CardHeader>
            <CardTitle>Metadata</CardTitle>
            <CardDescription>
              Metadata used when the project appears in search results.
            </CardDescription>
            <CardAction>
              <CollapsibleTrigger render={<Button variant="ghost" size="xs" />}>
                {metadataOpen ? 'Collapse' : 'Expand'}
              </CollapsibleTrigger>
            </CardAction>
          </CardHeader>
          <CollapsibleContent>
            <CardContent>
              <FieldGroup className="gap-6">
                <Field data-invalid={metaTitleInvalid}>
                  <FieldLabel htmlFor="metaTitle">Meta title</FieldLabel>
                  <Input
                    id="metaTitle"
                    value={form.data.metaTitle}
                    onChange={(event) => {
                      form.setData('metaTitle', event.currentTarget.value)
                      form.clearErrors('metaTitle')
                    }}
                    placeholder="Defaults to the project title"
                    aria-invalid={metaTitleInvalid}
                    required={publishing}
                  />
                  <LengthIndicator
                    value={form.data.metaTitle}
                    recommendedMin={metaTitleRecommendedMinLength}
                    recommendedMax={metaTitleRecommendedMaxLength}
                    max={metaTitleMaxLength}
                  />
                  <FieldError>{form.errors.metaTitle}</FieldError>
                </Field>

                <Field data-invalid={metaDescriptionInvalid}>
                  <FieldLabel htmlFor="metaDescription">Meta description</FieldLabel>
                  <Textarea
                    id="metaDescription"
                    value={form.data.metaDescription}
                    onChange={(event) => {
                      form.setData('metaDescription', event.currentTarget.value)
                      form.clearErrors('metaDescription')
                    }}
                    placeholder="A short description for search results"
                    className="min-h-24 resize-y"
                    aria-invalid={metaDescriptionInvalid}
                    required={publishing}
                  />
                  <LengthIndicator
                    value={form.data.metaDescription}
                    recommendedMin={metaDescriptionRecommendedMinLength}
                    recommendedMax={metaDescriptionRecommendedMaxLength}
                    max={metaDescriptionMaxLength}
                  />
                  <FieldError>{form.errors.metaDescription}</FieldError>
                </Field>
              </FieldGroup>
            </CardContent>
          </CollapsibleContent>
        </Collapsible>

        <Card>
          <CardHeader>
            <CardTitle>Content</CardTitle>
            <CardDescription>
              Write visually while storing portable Markdown, including fenced code blocks.
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Field data-invalid={Boolean(form.errors.content)}>
              <FieldLabel htmlFor="content" className="sr-only">
                Content
              </FieldLabel>
              <MarkdownEditor
                value={form.data.content}
                ariaLabel="Project content"
                onChange={(content) => {
                  form.setData('content', content)
                  form.clearErrors('content')
                }}
                invalid={Boolean(form.errors.content)}
              />
              <FieldError>{form.errors.content}</FieldError>
            </Field>
          </CardContent>
        </Card>
      </div>

      <div className="flex flex-col gap-6 lg:sticky lg:top-6 lg:self-start">
        <Card>
          <CardHeader>
            <CardTitle>Publishing</CardTitle>
            <CardDescription>Control whether this project is published.</CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup className="gap-6">
              <div className="flex items-center gap-3">
                <Switch
                  id="published"
                  checked={form.data.published}
                  onCheckedChange={(checked) => {
                    form.setData('published', checked)
                    if (!checked) {
                      form.clearErrors(
                        'startedAt',
                        'status',
                        'description',
                        'metaTitle',
                        'metaDescription',
                        'cover',
                        'metaImage',
                        'content',
                      )
                    }
                  }}
                />
                <FieldLabel htmlFor="published">
                  {publishing ? 'Unpublish' : 'Publish'}
                </FieldLabel>
              </div>

              <Field data-invalid={Boolean(form.errors.startedAt)}>
                <FieldLabel htmlFor="startedAt">Started at</FieldLabel>
                <Input
                  id="startedAt"
                  type="date"
                  value={form.data.startedAt}
                  onChange={(event) => {
                    form.setData('startedAt', event.currentTarget.value)
                    form.clearErrors('startedAt')
                  }}
                  required={publishing}
                  aria-invalid={Boolean(form.errors.startedAt)}
                />
                <FieldError>{form.errors.startedAt}</FieldError>
              </Field>

              <Field data-invalid={statusInvalid}>
                <FieldLabel htmlFor="status">Status</FieldLabel>
                <Input
                  id="status"
                  value={form.data.status}
                  onChange={(event) => {
                    form.setData('status', event.currentTarget.value)
                    form.clearErrors('status')
                  }}
                  placeholder="Planned, active, or completed"
                  required={publishing}
                  aria-invalid={statusInvalid}
                />
                <MaximumLengthIndicator value={form.data.status} max={statusMaxLength} />
                <FieldError>{form.errors.status}</FieldError>
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>

        <ImageUpload
          id="cover"
          title="Featured image"
          description="The cover displayed on the project page."
          previewAlt="Project cover preview"
          existingURL={existingImageURL}
          file={form.data.cover}
          remove={form.data.removeCover}
          required={publishing}
          error={form.errors.cover}
          onFileChange={(file) => form.setData('cover', file)}
          onRemoveChange={(remove) => form.setData('removeCover', remove)}
          onError={(error) => error ? form.setError('cover', error) : form.clearErrors('cover')}
        />

        <ImageUpload
          id="metaImage"
          title="Social image"
          description="Optional sharing image. The featured image is used when this is empty."
          previewAlt="Project social image preview"
          existingURL={existingMetaImageURL}
          file={form.data.metaImage}
          remove={form.data.removeMetaImage}
          error={form.errors.metaImage}
          onFileChange={(file) => form.setData('metaImage', file)}
          onRemoveChange={(remove) => form.setData('removeMetaImage', remove)}
          onError={(error) => error ? form.setError('metaImage', error) : form.clearErrors('metaImage')}
        />

        <Card>
          <CardHeader>
            <CardTitle>Project link</CardTitle>
            <CardDescription>Add an optional external URL for the project.</CardDescription>
          </CardHeader>
          <CardContent>
            <Field data-invalid={projectURLInvalid}>
              <FieldLabel htmlFor="projectUrl">Project URL</FieldLabel>
              <Input
                id="projectUrl"
                type="url"
                value={form.data.projectUrl}
                onChange={(event) => {
                  form.setData('projectUrl', event.currentTarget.value)
                  form.clearErrors('projectUrl')
                }}
                placeholder="https://example.com"
                aria-invalid={projectURLInvalid}
              />
              <FieldError>{form.errors.projectUrl}</FieldError>
            </Field>
          </CardContent>
        </Card>

        {publishing && (
          <Card>
            <CardHeader>
              <CardTitle>Ready to publish?</CardTitle>
              <CardDescription>
                Complete every requirement before saving this project as published.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ul className="grid grid-cols-2 gap-x-3 gap-y-3 text-xs">
                {publishRequirements.map((requirement) => (
                  <li
                    key={requirement.label}
                    className={requirement.complete ? 'text-emerald-400' : 'text-muted-foreground'}
                  >
                    {requirement.complete ? '✓' : '○'} {requirement.label}
                  </li>
                ))}
              </ul>
            </CardContent>
          </Card>
        )}

        <div className="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end lg:flex-col-reverse">
          <Button variant="outline" render={<Link href={cancelHref} />}>
            Cancel
          </Button>
          <Button type="submit" disabled={form.processing}>
            {form.processing ? processingLabel : submitLabel}
          </Button>
        </div>
      </div>
    </form>
  )
}
