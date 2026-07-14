import { Link, useForm } from '@inertiajs/react'

import { NewsletterForm } from './NewsletterForm'
import {
  toNewsletterMultipartData,
  type NewsletterFormData,
  type NewsletterValidationRules,
} from './types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type CreateProps = {
  validationRules: NewsletterValidationRules
}

export default function Create({ validationRules }: CreateProps) {
  const form = useForm<NewsletterFormData>({
    title: '',
    slug: '',
    metaTitle: '',
    metaDescription: '',
    isPublished: false,
    cover: null,
    removeCover: false,
    metaImage: null,
    removeMetaImage: false,
    content: '',
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Newsletters
          </p>
          <h1 className="mt-2 font-heading text-3xl font-semibold tracking-tight">
            New newsletter
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Create a new draft and prepare it for publication.
          </p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminNewsletterIndex()} />}>
          Back to newsletters
        </Button>
      </header>

      <NewsletterForm
        form={form}
        validationRules={validationRules}
        cancelHref={routes.adminNewsletterIndex()}
        submitLabel="Create newsletter"
        processingLabel="Creating..."
        existingImageURL=""
        existingMetaImageURL=""
        onSubmit={() => {
          form.transform(toNewsletterMultipartData)
          form.post(routes.adminNewsletterCreate(), { forceFormData: true })
        }}
      />
    </main>
  )
}
