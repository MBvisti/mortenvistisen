import { Link, useForm } from '@inertiajs/react'

import { ProjectForm } from './ProjectForm'
import { toProjectMultipartData, type ProjectFormData, type ProjectValidationRules } from './types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type CreateProps = {
  validationRules: ProjectValidationRules
}

export default function Create({ validationRules }: CreateProps) {
  const form = useForm<ProjectFormData>({
    published: false,
    title: '',
    slug: '',
    startedAt: '',
    status: '',
    description: '',
    metaTitle: '',
    metaDescription: '',
    cover: null,
    removeCover: false,
    metaImage: null,
    removeMetaImage: false,
    content: '',
    projectUrl: '',
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Projects
          </p>
          <h1 className="mt-2 font-heading text-3xl font-semibold tracking-tight">New project</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Create a new draft and prepare it for publication.
          </p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminProjectIndex()} />}>
          Back to projects
        </Button>
      </header>

      <ProjectForm
        form={form}
        validationRules={validationRules}
        cancelHref={routes.adminProjectIndex()}
        submitLabel="Create project"
        processingLabel="Creating..."
        existingImageURL=""
        existingMetaImageURL=""
        onSubmit={() => {
          form.transform(toProjectMultipartData)
          form.post(routes.adminProjectCreate(), { forceFormData: true })
        }}
      />
    </main>
  )
}
