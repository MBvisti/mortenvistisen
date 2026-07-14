import { Link, useForm } from '@inertiajs/react'

import { ArticleForm } from './ArticleForm'
import {
  type ArticleFormData,
  type ArticleItem,
  type ArticleValidationRules,
  toArticleMultipartData,
} from './types'
import type { TagItem } from '../Tag/types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type EditProps = {
  item: ArticleItem
  validationRules: ArticleValidationRules
  tagOptions: TagItem[]
}

export default function Edit({ item, validationRules, tagOptions }: EditProps) {
  const form = useForm<ArticleFormData>({
    published: Boolean(item.Published),
    title: String(item.Title ?? ''),
    excerpt: String(item.Excerpt ?? ''),
    metaTitle: String(item.MetaTitle ?? ''),
    metaDescription: String(item.MetaDescription ?? ''),
    slug: String(item.Slug ?? ''),
    cover: null,
    removeCover: false,
    metaImage: null,
    removeMetaImage: false,
    readTime: Number(item.ReadTime ?? 0),
    content: String(item.Content ?? ''),
    tagIds: (item.Tags ?? []).map((tag) => tag.ID),
    newTags: [],
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="min-w-0">
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Edit article
          </p>
          <h1 className="mt-2 truncate font-heading text-3xl font-semibold tracking-tight">
            {item.Title || 'Untitled article'}
          </h1>
          <p className="mt-2 truncate text-sm text-muted-foreground">/{item.Slug}</p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminArticleShow(item.ID)} />}>
          View article
        </Button>
      </header>

      <ArticleForm
        form={form}
        validationRules={validationRules}
        tagOptions={tagOptions}
        cancelHref={routes.adminArticleShow(item.ID)}
        submitLabel="Save changes"
        processingLabel="Saving..."
        existingImageURL={String(item.ImageLink ?? '')}
        existingMetaImageURL={String(item.MetaImageLink ?? '')}
        onSubmit={() => {
          form.transform(toArticleMultipartData)
          form.put(routes.adminArticleUpdate(item.ID), { forceFormData: true })
        }}
      />
    </main>
  )
}
