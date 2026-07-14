import { useState } from 'react'
import { Link, router } from '@inertiajs/react'

import { formatProjectDate, type ProjectItem } from './types'
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Field, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { MarkdownContent } from '@/Components/MarkdownContent'
import { routes } from '@/routes'

type ShowProps = {
  item: ProjectItem
}

function Detail({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="grid gap-1 border-b py-4 last:border-b-0 sm:grid-cols-[10rem_minmax(0,1fr)] sm:gap-6">
      <dt className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
        {label}
      </dt>
      <dd className="min-w-0 text-sm text-foreground">{children}</dd>
    </div>
  )
}

export default function Show({ item }: ShowProps) {
  const [deleteConfirmation, setDeleteConfirmation] = useState('')
  const canDelete = deleteConfirmation === 'DELETE'

  function handleDeleteDialogOpenChange(open: boolean) {
    if (!open) {
      setDeleteConfirmation('')
    }
  }

  function deleteProject() {
    if (!canDelete) {
      return
    }

    router.delete(routes.adminProjectDestroy(item.ID))
  }

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Project
            </p>
            <Badge
              variant="secondary"
              className={item.Published ? 'text-emerald-400' : undefined}
            >
              {item.Published ? 'Published' : 'Draft'}
            </Badge>
            {item.Status ? <Badge variant="outline">{item.Status}</Badge> : null}
          </div>
          <h1 className="mt-3 font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
            {item.Title || 'Untitled project'}
          </h1>
          <p className="mt-2 truncate text-sm text-muted-foreground">/{item.Slug}</p>
        </div>
        <div className="flex flex-wrap gap-3">
          <Button variant="outline" render={<Link href={routes.adminProjectIndex()} />}>
            Back to projects
          </Button>
          <Button render={<Link href={routes.adminProjectEdit(item.ID)} />}>
            Edit project
          </Button>
        </div>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <div className="flex min-w-0 flex-col gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Project overview</CardTitle>
              <CardDescription>The short summary and public project link.</CardDescription>
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Description">{item.Description || 'Not set'}</Detail>
                <Detail label="Project URL">
                  {item.ProjectUrl ? (
                    <a
                      href={item.ProjectUrl}
                      target="_blank"
                      rel="noreferrer"
                      className="break-all text-primary underline-offset-4 hover:underline"
                    >
                      {item.ProjectUrl}
                    </a>
                  ) : (
                    'Not set'
                  )}
                </Detail>
              </dl>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Search metadata</CardTitle>
              <CardDescription>How this project is described to search engines.</CardDescription>
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Meta title">{item.MetaTitle || 'Not set'}</Detail>
                <Detail label="Description">{item.MetaDescription || 'Not set'}</Detail>
              </dl>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Project content</CardTitle>
              <CardDescription>A preview of the published project content.</CardDescription>
            </CardHeader>
            <CardContent>
              {item.Content ? (
                <MarkdownContent value={item.Content} ariaLabel="Project content preview" />
              ) : (
                <p className="text-sm text-muted-foreground">No project content yet.</p>
              )}
            </CardContent>
          </Card>
        </div>

        <div className="flex flex-col gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Publishing</CardTitle>
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Publication">{item.Published ? 'Published' : 'Draft'}</Detail>
                <Detail label="Status">{item.Status || 'Not set'}</Detail>
                <Detail label="Started">{formatProjectDate(item.StartedAt)}</Detail>
              </dl>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Record details</CardTitle>
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Created">{formatProjectDate(item.CreatedAt)}</Detail>
                <Detail label="Updated">{formatProjectDate(item.UpdatedAt)}</Detail>
              </dl>
            </CardContent>
          </Card>

          <AlertDialog onOpenChange={handleDeleteDialogOpenChange}>
            <AlertDialogTrigger render={<Button variant="destructive" className="w-full" />}>
              Delete project
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Delete this project?</AlertDialogTitle>
                <AlertDialogDescription>
                  This permanently deletes {item.Title || 'this project'}. This action cannot be
                  undone.
                </AlertDialogDescription>
              </AlertDialogHeader>
              <Field>
                <FieldLabel htmlFor="deleteConfirmation">Type DELETE to confirm</FieldLabel>
                <Input
                  id="deleteConfirmation"
                  value={deleteConfirmation}
                  onChange={(event) => setDeleteConfirmation(event.currentTarget.value)}
                  autoComplete="off"
                  autoFocus
                />
              </Field>
              <AlertDialogFooter>
                <AlertDialogCancel>Cancel</AlertDialogCancel>
                <AlertDialogAction
                  variant="destructive"
                  disabled={!canDelete}
                  onClick={deleteProject}
                >
                  Delete project
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </main>
  )
}
