import { Link, useForm } from '@inertiajs/react'

import { NewsletterForm } from './NewsletterForm'
import type {
  NewsletterFormData,
  NewsletterItem,
  NewsletterValidationRules,
} from './types'
import { toNewsletterMultipartData } from './types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type EditProps = {
  item: NewsletterItem
  validationRules: NewsletterValidationRules
}

export default function Edit({ item, validationRules }: EditProps) {
  const form = useForm<NewsletterFormData>({
    title: String(item.Title ?? ''),
    slug: String(item.Slug ?? ''),
    metaTitle: String(item.MetaTitle ?? ''),
    metaDescription: String(item.MetaDescription ?? ''),
    isPublished: Boolean(item.IsPublished),
    cover: null,
    removeCover: false,
    metaImage: null,
    removeMetaImage: false,
    content: String(item.Content ?? ''),
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="min-w-0">
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Edit newsletter
          </p>
          <h1 className="mt-2 truncate font-heading text-3xl font-semibold tracking-tight">
            {item.Title || 'Untitled newsletter'}
          </h1>
          <p className="mt-2 truncate text-sm text-muted-foreground">/{item.Slug}</p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminNewsletterShow(item.ID)} />}>
          View newsletter
        </Button>
      </header>

      <NewsletterForm
        form={form}
        validationRules={validationRules}
        cancelHref={routes.adminNewsletterShow(item.ID)}
        submitLabel="Save changes"
        processingLabel="Saving..."
        existingImageURL={String(item.ImageLink ?? '')}
        existingMetaImageURL={String(item.MetaImageLink ?? '')}
        onSubmit={() => {
          form.transform(toNewsletterMultipartData)
          form.put(routes.adminNewsletterUpdate(item.ID), { forceFormData: true })
        }}
      />
    </main>
  )
}
