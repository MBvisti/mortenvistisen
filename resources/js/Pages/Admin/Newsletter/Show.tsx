import { useState } from 'react'
import { Link, router } from '@inertiajs/react'

import { formatNewsletterDate, type NewsletterItem } from './types'
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
  item: NewsletterItem
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

  function deleteNewsletter() {
    if (!canDelete) {
      return
    }

    router.delete(routes.adminNewsletterDestroy(item.ID))
  }

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Newsletter
            </p>
            <Badge
              variant="secondary"
              className={item.IsPublished ? 'text-emerald-400' : undefined}
            >
              {item.IsPublished ? 'Published' : 'Draft'}
            </Badge>
          </div>
          <h1 className="mt-3 font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
            {item.Title || 'Untitled newsletter'}
          </h1>
          <p className="mt-2 truncate text-sm text-muted-foreground">/{item.Slug}</p>
        </div>
        <div className="flex flex-wrap gap-3">
          <Button variant="outline" render={<Link href={routes.adminNewsletterIndex()} />}>
            Back to newsletters
          </Button>
          <Button render={<Link href={routes.adminNewsletterEdit(item.ID)} />}>
            Edit newsletter
          </Button>
        </div>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <div className="flex min-w-0 flex-col gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Featured image</CardTitle>
              <CardDescription>The image displayed with this newsletter.</CardDescription>
            </CardHeader>
            <CardContent>
              {item.ImageLink ? (
                <a href={item.ImageLink} target="_blank" rel="noreferrer" className="block">
                  <img
                    src={item.ImageLink}
                    alt={`Featured image for ${item.Title || 'this newsletter'}`}
                    className="aspect-video w-full bg-muted object-cover"
                  />
                </a>
              ) : (
                <div className="flex aspect-video items-center justify-center bg-muted text-sm text-muted-foreground">
                  No featured image set.
                </div>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Search metadata</CardTitle>
              <CardDescription>
                How this newsletter is described to search engines.
              </CardDescription>
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
              <CardTitle>Newsletter</CardTitle>
              <CardDescription>A preview of the published newsletter content.</CardDescription>
            </CardHeader>
            <CardContent>
              {item.Content ? (
                <MarkdownContent value={item.Content} ariaLabel="Newsletter content preview" />
              ) : (
                <p className="text-sm text-muted-foreground">No newsletter content yet.</p>
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
                <Detail label="Status">{item.IsPublished ? 'Published' : 'Draft'}</Detail>
                <Detail label="First released">
                  {formatNewsletterDate(item.ReleasedAt)}
                </Detail>
              </dl>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Record details</CardTitle>
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Created">{formatNewsletterDate(item.CreatedAt)}</Detail>
                <Detail label="Updated">{formatNewsletterDate(item.UpdatedAt)}</Detail>
              </dl>
            </CardContent>
          </Card>

          <AlertDialog onOpenChange={handleDeleteDialogOpenChange}>
            <AlertDialogTrigger render={<Button variant="destructive" className="w-full" />}>
              Delete newsletter
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Delete this newsletter?</AlertDialogTitle>
                <AlertDialogDescription>
                  This permanently deletes {item.Title || 'this newsletter'}. This action cannot
                  be undone.
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
                  onClick={deleteNewsletter}
                >
                  Delete newsletter
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </main>
  )
}
