import { Link, useForm } from '@inertiajs/react'

import { ProjectForm } from './ProjectForm'
import {
  projectDateInputValue,
  toProjectMultipartData,
  type ProjectFormData,
  type ProjectItem,
  type ProjectValidationRules,
} from './types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type EditProps = {
  item: ProjectItem
  validationRules: ProjectValidationRules
}

export default function Edit({ item, validationRules }: EditProps) {
  const form = useForm<ProjectFormData>({
    published: Boolean(item.Published),
    title: String(item.Title ?? ''),
    slug: String(item.Slug ?? ''),
    startedAt: projectDateInputValue(item.StartedAt),
    status: String(item.Status ?? ''),
    description: String(item.Description ?? ''),
    metaTitle: String(item.MetaTitle ?? ''),
    metaDescription: String(item.MetaDescription ?? ''),
    cover: null,
    removeCover: false,
    metaImage: null,
    removeMetaImage: false,
    content: String(item.Content ?? ''),
    projectUrl: String(item.ProjectUrl ?? ''),
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="min-w-0">
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Edit project
          </p>
          <h1 className="mt-2 truncate font-heading text-3xl font-semibold tracking-tight">
            {item.Title || 'Untitled project'}
          </h1>
          <p className="mt-2 truncate text-sm text-muted-foreground">/{item.Slug}</p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminProjectShow(item.ID)} />}>
          View project
        </Button>
      </header>

      <ProjectForm
        form={form}
        validationRules={validationRules}
        cancelHref={routes.adminProjectShow(item.ID)}
        submitLabel="Save changes"
        processingLabel="Saving..."
        existingImageURL={String(item.ImageLink ?? '')}
        existingMetaImageURL={String(item.MetaImageLink ?? '')}
        onSubmit={() => {
          form.transform(toProjectMultipartData)
          form.put(routes.adminProjectUpdate(item.ID), { forceFormData: true })
        }}
      />
    </main>
  )
}
