import { Link, useForm } from '@inertiajs/react'

import { ArticleForm } from './ArticleForm'
import {
  toArticleMultipartData,
  type ArticleFormData,
  type ArticleValidationRules,
} from './types'
import type { TagItem } from '../Tag/types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type CreateProps = {
  validationRules: ArticleValidationRules
  tagOptions: TagItem[]
}

export default function Create({ validationRules, tagOptions }: CreateProps) {
  const form = useForm<ArticleFormData>({
    published: false,
    title: '',
    excerpt: '',
    metaTitle: '',
    metaDescription: '',
    slug: '',
    cover: null,
    removeCover: false,
    metaImage: null,
    removeMetaImage: false,
    readTime: 0,
    content: '',
    tagIds: [],
    newTags: [],
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Articles
          </p>
          <h1 className="mt-2 font-heading text-3xl font-semibold tracking-tight">New article</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Create a new draft and prepare it for publication.
          </p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminArticleIndex()} />}>
          Back to articles
        </Button>
      </header>

      <ArticleForm
        form={form}
        validationRules={validationRules}
        tagOptions={tagOptions}
        cancelHref={routes.adminArticleIndex()}
        submitLabel="Create article"
        processingLabel="Creating..."
        existingImageURL=""
        existingMetaImageURL=""
        onSubmit={() => {
          form.transform(toArticleMultipartData)
          form.post(routes.adminArticleCreate(), { forceFormData: true })
        }}
      />
    </main>
  )
}
