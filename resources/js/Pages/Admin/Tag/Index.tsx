import { useMemo, useState, type FormEvent } from 'react'
import { Link, router, useForm } from '@inertiajs/react'
import { Check, Pencil, Search, Trash2, X } from 'lucide-react'

import type { TagItem, TagPagination, TagValidationRules } from './types'
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
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Field, FieldDescription, FieldError, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { routes } from '@/routes'

type IndexProps = {
  items: TagItem[]
  pagination: TagPagination
  validationRules: TagValidationRules
}

function maximumTitleLength(rules: TagValidationRules) {
  const value = rules.title?.find((rule) => rule.code === 'max')?.params?.max
  return typeof value === 'number' ? value : 255
}

function pageHref(page: number) {
  return `${routes.adminTagIndex()}?page=${page}`
}

function DeleteTag({ item, page }: { item: TagItem; page: number }) {
  const [confirmation, setConfirmation] = useState('')
  const canDelete = confirmation === 'DELETE'

  return (
    <AlertDialog onOpenChange={(open) => !open && setConfirmation('')}>
      <AlertDialogTrigger
        render={(
          <Button
            variant="ghost"
            size="icon-xs"
            aria-label={`Delete ${item.Title}`}
          />
        )}
      >
        <Trash2 />
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Delete “{item.Title}”?</AlertDialogTitle>
          <AlertDialogDescription>
            {item.ArticleCount > 0
              ? `This removes the tag from ${item.ArticleCount} ${item.ArticleCount === 1 ? 'article' : 'articles'} and permanently deletes it.`
              : 'This permanently deletes the tag.'}{' '}
            This action cannot be undone.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <Field>
          <FieldLabel htmlFor={`deleteTag-${item.ID}`}>Type DELETE to confirm</FieldLabel>
          <Input
            id={`deleteTag-${item.ID}`}
            value={confirmation}
            onChange={(event) => setConfirmation(event.currentTarget.value)}
            autoComplete="off"
          />
        </Field>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction
            variant="destructive"
            disabled={!canDelete}
            onClick={() => router.delete(
              `${routes.adminTagDestroy(item.ID)}?page=${page}`,
              { preserveScroll: true },
            )}
          >
            Delete tag
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  )
}

export default function Index({ items, pagination, validationRules }: IndexProps) {
  const createForm = useForm({ title: '' })
  const editForm = useForm({ title: '' })
  const [query, setQuery] = useState('')
  const [editingID, setEditingID] = useState<number | null>(null)
  const titleMaxLength = maximumTitleLength(validationRules)
  const firstItem = pagination.TotalCount === 0
    ? 0
    : (pagination.Page - 1) * pagination.PageSize + 1
  const lastItem = Math.min(
    pagination.Page * pagination.PageSize,
    pagination.TotalCount,
  )
  const hasPreviousPage = pagination.Page > 1
  const hasNextPage = pagination.Page < pagination.TotalPages
  const filteredItems = useMemo(() => {
    const normalized = query.trim().toLocaleLowerCase()
    if (normalized === '') {
      return items
    }
    return items.filter((item) => item.Title.toLocaleLowerCase().includes(normalized))
  }, [items, query])

  function createTag(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    createForm.post(routes.adminTagCreate(), {
      preserveScroll: true,
      onSuccess: () => createForm.reset(),
    })
  }

  function beginEditing(item: TagItem) {
    setEditingID(item.ID)
    editForm.setData('title', item.Title)
    editForm.clearErrors()
  }

  function cancelEditing() {
    setEditingID(null)
    editForm.reset()
    editForm.clearErrors()
  }

  function updateTag(event: FormEvent<HTMLFormElement>, id: number) {
    event.preventDefault()
    editForm.put(`${routes.adminTagUpdate(id)}?page=${pagination.Page}`, {
      preserveScroll: true,
      onSuccess: cancelEditing,
    })
  }

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header>
        <h1 className="font-heading text-3xl font-semibold tracking-tight">Tags</h1>
        <p className="mt-2 max-w-2xl text-sm text-muted-foreground">
          Create a shared vocabulary for articles, rename tags, and remove tags that are no longer useful.
        </p>
      </header>

      <div className="grid items-start gap-6 lg:grid-cols-[20rem_minmax(0,1fr)]">
        <Card className="lg:sticky lg:top-6">
          <CardHeader>
            <CardTitle>Create tag</CardTitle>
            <CardDescription>New tags become available in every article form.</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={createTag} className="space-y-4">
              <Field data-invalid={Boolean(createForm.errors.title)}>
                <FieldLabel htmlFor="newTagTitle">Title</FieldLabel>
                <Input
                  id="newTagTitle"
                  value={createForm.data.title}
                  onChange={(event) => {
                    createForm.setData('title', event.currentTarget.value)
                    createForm.clearErrors('title')
                  }}
                  maxLength={titleMaxLength}
                  placeholder="Engineering"
                  aria-invalid={Boolean(createForm.errors.title)}
                />
                <FieldDescription>
                  Titles must be unique, regardless of capitalization.
                </FieldDescription>
                <FieldError>{createForm.errors.title}</FieldError>
              </Field>
              <Button
                type="submit"
                className="w-full"
                disabled={createForm.processing || createForm.data.title.trim() === ''}
              >
                {createForm.processing ? 'Creating...' : 'Create tag'}
              </Button>
            </form>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>All tags</CardTitle>
            <CardDescription>
              {pagination.TotalCount === 0
                ? 'No tags yet.'
                : `Showing ${firstItem}-${lastItem} of ${pagination.TotalCount} tags.`}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-5 px-0">
            <div className="relative mx-6">
              <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                value={query}
                onChange={(event) => setQuery(event.currentTarget.value)}
                placeholder="Search this page"
                className="pl-9"
                aria-label="Search tags on this page"
              />
            </div>

            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="pl-6">Title</TableHead>
                  <TableHead>Articles</TableHead>
                  <TableHead className="pr-6 text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredItems.length === 0 ? (
                  <TableRow>
                    <TableCell colSpan={3} className="h-32 text-center text-muted-foreground">
                      {pagination.TotalCount === 0
                        ? 'Create your first tag to get started.'
                        : 'No tags match your search.'}
                    </TableCell>
                  </TableRow>
                ) : filteredItems.map((item) => (
                  <TableRow key={item.ID}>
                    <TableCell className="pl-6 font-medium">
                      {editingID === item.ID ? (
                        <form
                          onSubmit={(event) => updateTag(event, item.ID)}
                          className="flex max-w-md items-start gap-2"
                        >
                          <div className="flex-1">
                            <Input
                              value={editForm.data.title}
                              onChange={(event) => {
                                editForm.setData('title', event.currentTarget.value)
                                editForm.clearErrors('title')
                              }}
                              maxLength={titleMaxLength}
                              aria-label={`Rename ${item.Title}`}
                              aria-invalid={Boolean(editForm.errors.title)}
                              autoFocus
                            />
                            <FieldError>{editForm.errors.title}</FieldError>
                          </div>
                          <Button
                            type="submit"
                            size="icon-sm"
                            disabled={editForm.processing || editForm.data.title.trim() === ''}
                            aria-label="Save tag title"
                          >
                            <Check />
                          </Button>
                          <Button
                            type="button"
                            variant="outline"
                            size="icon-sm"
                            onClick={cancelEditing}
                            aria-label="Cancel editing"
                          >
                            <X />
                          </Button>
                        </form>
                      ) : item.Title}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {item.ArticleCount}
                    </TableCell>
                    <TableCell className="pr-6">
                      {editingID !== item.ID ? (
                        <div className="flex justify-end gap-1">
                          <Button
                            type="button"
                            variant="ghost"
                            size="icon-xs"
                            onClick={() => beginEditing(item)}
                            aria-label={`Rename ${item.Title}`}
                          >
                            <Pencil />
                          </Button>
                          <DeleteTag item={item} page={pagination.Page} />
                        </div>
                      ) : null}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>

            {pagination.TotalPages > 1 ? (
              <nav
                className="flex flex-col gap-3 border-t px-6 pt-5 sm:flex-row sm:items-center sm:justify-between"
                aria-label="Tag pagination"
              >
                <p className="text-xs text-muted-foreground">
                  Page {pagination.Page} of {pagination.TotalPages}
                </p>
                <div className="flex gap-3">
                  {hasPreviousPage ? (
                    <Button
                      variant="outline"
                      size="sm"
                      render={<Link href={pageHref(pagination.Page - 1)} preserveScroll />}
                    >
                      Previous
                    </Button>
                  ) : (
                    <Button variant="outline" size="sm" disabled>Previous</Button>
                  )}
                  {hasNextPage ? (
                    <Button
                      variant="outline"
                      size="sm"
                      render={<Link href={pageHref(pagination.Page + 1)} preserveScroll />}
                    >
                      Next
                    </Button>
                  ) : (
                    <Button variant="outline" size="sm" disabled>Next</Button>
                  )}
                </div>
              </nav>
            ) : null}
          </CardContent>
        </Card>
      </div>
    </main>
  )
}
