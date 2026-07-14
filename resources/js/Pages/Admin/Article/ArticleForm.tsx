import { useEffect, useMemo, useState, type FormEvent } from 'react'
import { Link, type InertiaFormProps } from '@inertiajs/react'
import { Plus, Search, X } from 'lucide-react'

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

import type { ArticleFormData, ArticleValidationRules } from './types'
import type { TagItem } from '../Tag/types'

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

function findRule(rules: ArticleValidationRules, field: string, code: string) {
  return rules[field]?.find((rule) => rule.code === code)
}

function ruleNumber(
  rules: ArticleValidationRules,
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

function LengthIndicator({
  value,
  recommendedMin,
  max,
  required = false,
}: {
  value: string
  recommendedMin: number
  max: number
  required?: boolean
}) {
  const length = characterLength(value)
  const missing = required && value.trim() === ''
  const tooLong = length > max

  let status = 'Valid'
  let statusClass = 'text-amber-400'

  if (missing) {
    status = 'Required'
    statusClass = 'text-destructive'
  } else if (tooLong) {
    status = 'Too long'
    statusClass = 'text-destructive'
  } else if (length === 0) {
    status = 'Optional'
    statusClass = 'text-muted-foreground'
  } else if (length >= recommendedMin) {
    status = 'Ideal length'
    statusClass = 'text-emerald-400'
  }

  return (
    <div className="flex flex-wrap items-center justify-between gap-2">
      <FieldDescription>
        Aim for {recommendedMin}-{max} characters.
      </FieldDescription>
      <span className={`text-xs font-medium tabular-nums ${statusClass}`} aria-live="polite">
        {length}/{max} · {status}
      </span>
    </div>
  )
}

type ArticleFormProps = {
  form: InertiaFormProps<ArticleFormData>
  validationRules: ArticleValidationRules
  tagOptions: TagItem[]
  cancelHref: string
  submitLabel: string
  processingLabel: string
  existingImageURL: string
  existingMetaImageURL: string
  onSubmit: () => void
}

function TagSelector({
  form,
  tagOptions,
}: {
  form: InertiaFormProps<ArticleFormData>
  tagOptions: TagItem[]
}) {
  const [query, setQuery] = useState('')
  const normalizedQuery = query.trim().toLocaleLowerCase()
  const selectedIDs = new Set(form.data.tagIds)
  const selectedTags = tagOptions.filter((tag) => selectedIDs.has(tag.ID))
  const matchingTags = useMemo(
    () => tagOptions.filter((tag) => (
      !selectedIDs.has(tag.ID) &&
      (normalizedQuery === '' || tag.Title.toLocaleLowerCase().includes(normalizedQuery))
    )).slice(0, 8),
    [form.data.tagIds, normalizedQuery, tagOptions],
  )
  const titleAlreadyExists = tagOptions.some(
    (tag) => tag.Title.trim().toLocaleLowerCase() === normalizedQuery,
  ) || form.data.newTags.some(
    (title) => title.trim().toLocaleLowerCase() === normalizedQuery,
  )
  const canCreate = normalizedQuery !== '' && !titleAlreadyExists

  function selectTag(id: number) {
    form.setData('tagIds', [...form.data.tagIds, id])
    form.clearErrors('tagIds')
    setQuery('')
  }

  function createTag() {
    const title = query.trim()
    if (!canCreate || title === '') {
      return
    }

    form.setData('newTags', [...form.data.newTags, title])
    form.clearErrors('newTags')
    setQuery('')
  }

  function removeTag(id: number) {
    form.setData('tagIds', form.data.tagIds.filter((tagID) => tagID !== id))
  }

  function removeNewTag(title: string) {
    form.setData('newTags', form.data.newTags.filter((tagTitle) => tagTitle !== title))
  }

  return (
    <Field data-invalid={Boolean(form.errors.tagIds || form.errors.newTags)}>
      <FieldLabel htmlFor="tagSearch">Tags</FieldLabel>
      <div className="flex flex-wrap gap-2">
        {selectedTags.map((tag) => (
          <span
            key={tag.ID}
            className="inline-flex items-center gap-2 border border-border bg-muted px-2 py-1 text-xs"
          >
            {tag.Title}
            <button
              type="button"
              onClick={() => removeTag(tag.ID)}
              className="text-muted-foreground hover:text-foreground"
              aria-label={`Remove ${tag.Title}`}
            >
              <X className="size-3" />
            </button>
          </span>
        ))}
        {form.data.newTags.map((title) => (
          <span
            key={title}
            className="inline-flex items-center gap-2 border border-primary/40 bg-primary/10 px-2 py-1 text-xs"
          >
            {title}
            <span className="text-[0.625rem] uppercase tracking-wider text-muted-foreground">
              New
            </span>
            <button
              type="button"
              onClick={() => removeNewTag(title)}
              className="text-muted-foreground hover:text-foreground"
              aria-label={`Remove new tag ${title}`}
            >
              <X className="size-3" />
            </button>
          </span>
        ))}
      </div>
      <div className="relative">
        <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          id="tagSearch"
          value={query}
          onChange={(event) => setQuery(event.currentTarget.value)}
          onKeyDown={(event) => {
            if (event.key === 'Enter' && canCreate) {
              event.preventDefault()
              createTag()
            }
          }}
          placeholder="Search or create a tag"
          className="pl-9"
          autoComplete="off"
        />
      </div>
      {(matchingTags.length > 0 || canCreate) && (
        <div className="max-h-52 overflow-y-auto border border-border">
          {matchingTags.map((tag) => (
            <button
              key={tag.ID}
              type="button"
              onClick={() => selectTag(tag.ID)}
              className="flex w-full items-center justify-between border-b border-border px-3 py-2 text-left text-sm last:border-b-0 hover:bg-muted"
            >
              <span>{tag.Title}</span>
              <Plus className="size-3.5 text-muted-foreground" />
            </button>
          ))}
          {canCreate && (
            <button
              type="button"
              onClick={createTag}
              className="flex w-full items-center gap-2 px-3 py-2 text-left text-sm text-primary hover:bg-muted"
            >
              <Plus className="size-3.5" />
              Create “{query.trim()}”
            </button>
          )}
        </div>
      )}
      <FieldDescription>
        Select existing tags or type a new title and press Enter.
      </FieldDescription>
      <FieldError>{form.errors.tagIds || form.errors.newTags}</FieldError>
    </Field>
  )
}

export function ArticleForm({
  form,
  validationRules,
  tagOptions,
  cancelHref,
  submitLabel,
  processingLabel,
  existingImageURL,
  existingMetaImageURL,
  onSubmit,
}: ArticleFormProps) {
  const [detailsOpen, setDetailsOpen] = useState(false)
  const [metadataOpen, setMetadataOpen] = useState(false)
  const [tagsOpen, setTagsOpen] = useState(false)
  const [slugIsAutomatic, setSlugIsAutomatic] = useState(
    () => form.data.slug === '' || form.data.slug === slugify(form.data.title),
  )
  const publishing = form.data.published
  const titleRequired = Boolean(findRule(validationRules, 'title', 'required'))
  const titleMaxLength = ruleNumber(validationRules, 'title', 'max', 'max')
  const titleRecommendedMinLength = ruleNumber(
    validationRules,
    'title',
    'recommended_length',
    'min',
  )
  const metaTitleMaxLength = ruleNumber(validationRules, 'metaTitle', 'max', 'max')
  const metaTitleRecommendedMinLength = ruleNumber(
    validationRules,
    'metaTitle',
    'recommended_length',
    'min',
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
  const readTimeMin = ruleNumber(validationRules, 'readTime', 'min', 'min')
  const titleInvalid =
    (titleRequired && form.data.title.trim() === '') ||
    characterLength(form.data.title) > titleMaxLength ||
    Boolean(form.errors.title)
  const metaTitleInvalid =
    characterLength(form.data.metaTitle) > metaTitleMaxLength ||
    Boolean(form.errors.metaTitle)
  const metaDescriptionInvalid =
    characterLength(form.data.metaDescription) > metaDescriptionMaxLength ||
    Boolean(form.errors.metaDescription)
  const publishRequirements = [
    { label: 'Excerpt', complete: form.data.excerpt.trim() !== '' },
    { label: 'Meta title', complete: form.data.metaTitle.trim() !== '' },
    { label: 'Meta description', complete: form.data.metaDescription.trim() !== '' },
    {
      label: 'Featured image',
      complete: Boolean(form.data.cover || (existingImageURL && !form.data.removeCover)),
    },
    { label: 'Read time', complete: form.data.readTime >= readTimeMin },
    { label: 'Content', complete: form.data.content.trim() !== '' },
  ]

  useEffect(() => {
    if (form.errors.title || form.errors.slug || form.errors.excerpt) {
      setDetailsOpen(true)
    }

    if (form.errors.metaTitle || form.errors.metaDescription) {
      setMetadataOpen(true)
    }

    if (form.errors.tagIds || form.errors.newTags) {
      setTagsOpen(true)
    }
  }, [
    form.errors.excerpt,
    form.errors.metaDescription,
    form.errors.metaTitle,
    form.errors.newTags,
    form.errors.slug,
    form.errors.tagIds,
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
            <CardTitle>Article details</CardTitle>
            <CardDescription>
              Set the title, URL slug, and short summary readers see first.
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
                    placeholder="A clear, descriptive title"
                    aria-invalid={titleInvalid}
                  />
                  <LengthIndicator
                    value={form.data.title}
                    recommendedMin={titleRecommendedMinLength}
                    max={titleMaxLength}
                    required={titleRequired}
                  />
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
                    placeholder="a-clear-descriptive-title"
                    aria-invalid={Boolean(form.errors.slug)}
                  />
                  <FieldDescription>
                    Generated from the title until you customize it.
                  </FieldDescription>
                  <FieldError>{form.errors.slug}</FieldError>
                </Field>

                <Field data-invalid={Boolean(form.errors.excerpt)}>
                  <FieldLabel htmlFor="excerpt">Excerpt</FieldLabel>
                  <Textarea
                    id="excerpt"
                    value={form.data.excerpt}
                    onChange={(event) => {
                      form.setData('excerpt', event.currentTarget.value)
                      form.clearErrors('excerpt')
                    }}
                    placeholder="A concise introduction to the article"
                    className="min-h-24 resize-y"
                    required={publishing}
                    aria-invalid={Boolean(form.errors.excerpt)}
                  />
                  <FieldError>{form.errors.excerpt}</FieldError>
                </Field>
              </FieldGroup>
            </CardContent>
          </CollapsibleContent>
        </Collapsible>

        <Collapsible open={metadataOpen} onOpenChange={setMetadataOpen} render={<Card />}>
          <CardHeader>
            <CardTitle>Metadata</CardTitle>
            <CardDescription>
              Metadata used when the article appears in search results.
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
                    placeholder="Defaults to the article title"
                    aria-invalid={metaTitleInvalid}
                    required={publishing}
                  />
                  <LengthIndicator
                    value={form.data.metaTitle}
                    recommendedMin={metaTitleRecommendedMinLength}
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
                    max={metaDescriptionMaxLength}
                  />
                  <FieldError>{form.errors.metaDescription}</FieldError>
                </Field>
              </FieldGroup>
            </CardContent>
          </CollapsibleContent>
        </Collapsible>

        <Collapsible open={tagsOpen} onOpenChange={setTagsOpen} render={<Card />}>
          <CardHeader>
            <CardTitle>Tags</CardTitle>
            <CardDescription>Organize this article for browsing and discovery.</CardDescription>
            <CardAction>
              <CollapsibleTrigger render={<Button variant="ghost" size="xs" />}>
                {tagsOpen ? 'Collapse' : 'Expand'}
              </CollapsibleTrigger>
            </CardAction>
          </CardHeader>
          <CollapsibleContent>
            <CardContent>
              <TagSelector form={form} tagOptions={tagOptions} />
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
              <FieldLabel htmlFor="content" className="sr-only">Content</FieldLabel>
              <MarkdownEditor
                value={form.data.content}
                ariaLabel="Article content"
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
            <CardDescription>Control when and how this article is published.</CardDescription>
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
                        'excerpt',
                        'metaTitle',
                        'metaDescription',
                        'cover',
                        'metaImage',
                        'readTime',
                        'content',
                      )
                    }
                  }}
                />
                <FieldLabel htmlFor="published">
                  {publishing ? 'Unpublish' : 'Publish'}
                </FieldLabel>
              </div>

              <Field data-invalid={Boolean(form.errors.readTime)}>
                <FieldLabel htmlFor="readTime">Read time</FieldLabel>
                <Input
                  id="readTime"
                  type="number"
                  min={publishing ? readTimeMin : 0}
                  value={form.data.readTime}
                  onChange={(event) => {
                    form.setData('readTime', Number(event.currentTarget.value))
                    form.clearErrors('readTime')
                  }}
                  aria-invalid={Boolean(form.errors.readTime)}
                />
                <FieldDescription>Estimated reading time in minutes.</FieldDescription>
                <FieldError>{form.errors.readTime}</FieldError>
              </Field>
            </FieldGroup>
          </CardContent>
        </Card>

        <ImageUpload
          id="cover"
          title="Featured image"
          description="The cover displayed on the article page."
          previewAlt="Article cover preview"
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
          previewAlt="Article social image preview"
          existingURL={existingMetaImageURL}
          file={form.data.metaImage}
          remove={form.data.removeMetaImage}
          error={form.errors.metaImage}
          onFileChange={(file) => form.setData('metaImage', file)}
          onRemoveChange={(remove) => form.setData('removeMetaImage', remove)}
          onError={(error) => error ? form.setError('metaImage', error) : form.clearErrors('metaImage')}
        />

        {publishing && (
          <Card>
            <CardHeader>
              <CardTitle>Ready to publish?</CardTitle>
              <CardDescription>
                Complete every requirement before saving this article as published.
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
