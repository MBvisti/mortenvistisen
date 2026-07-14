import { useEffect, useState, type KeyboardEvent, type MouseEvent, type ReactNode } from 'react'
import { Link, router } from '@inertiajs/react'

import {
  formatSubscriberDate,
  type SubscriberItem,
  type SubscriberPagination,
} from './types'
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
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Field, FieldLabel } from '@/components/ui/field'
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
  items: SubscriberItem[]
  pagination: SubscriberPagination
}

function ClickableRow({
  href,
  label,
  children,
}: {
  href: string
  label: string
  children: ReactNode
}) {
  function visit() {
    router.visit(href)
  }

  function handleKeyDown(event: KeyboardEvent<HTMLTableRowElement>) {
    if (event.key !== 'Enter' && event.key !== ' ') {
      return
    }

    event.preventDefault()
    visit()
  }

  return (
    <TableRow
      role="link"
      tabIndex={0}
      aria-label={label}
      onClick={visit}
      onKeyDown={handleKeyDown}
      className="cursor-pointer focus-visible:bg-muted/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-inset focus-visible:ring-ring"
    >
      {children}
    </TableRow>
  )
}

function stopRowNavigation(event: MouseEvent | KeyboardEvent) {
  event.stopPropagation()
}

function pageHref(page: number) {
  return `${routes.adminSubscriberIndex()}?page=${page}`
}

export default function Index({ items, pagination }: IndexProps) {
  const [selectedIDs, setSelectedIDs] = useState<Set<number>>(new Set())
  const [deleteConfirmation, setDeleteConfirmation] = useState('')
  const visibleIDs = items.map((item) => item.ID)
  const allSelected = items.length > 0 && visibleIDs.every((id) => selectedIDs.has(id))
  const someSelected = visibleIDs.some((id) => selectedIDs.has(id)) && !allSelected
  const selectedCount = selectedIDs.size
  const canDelete = selectedCount > 0 && deleteConfirmation === 'DELETE'
  const firstItem = pagination.TotalCount === 0 ? 0 : (pagination.Page - 1) * pagination.PageSize + 1
  const lastItem = Math.min(
    pagination.Page * pagination.PageSize,
    pagination.TotalCount,
  )
  const hasPreviousPage = pagination.Page > 1
  const hasNextPage = pagination.Page < pagination.TotalPages

  useEffect(() => {
    const currentIDs = new Set(visibleIDs)
    setSelectedIDs((current) => new Set([...current].filter((id) => currentIDs.has(id))))
  }, [items])

  function selectAll(checked: boolean) {
    setSelectedIDs(checked ? new Set(visibleIDs) : new Set())
  }

  function selectSubscriber(id: number, checked: boolean) {
    setSelectedIDs((current) => {
      const next = new Set(current)
      if (checked) {
        next.add(id)
      } else {
        next.delete(id)
      }
      return next
    })
  }

  function handleDeleteDialogOpenChange(open: boolean) {
    if (!open) {
      setDeleteConfirmation('')
    }
  }

  function deleteSelectedSubscribers() {
    if (!canDelete) {
      return
    }

    router.delete(routes.adminSubscriberBulkDestroy(), {
      data: { subscriberIds: [...selectedIDs] },
      preserveScroll: true,
      onSuccess: () => setSelectedIDs(new Set()),
    })
  }

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-heading text-3xl font-semibold tracking-tight">Subscribers</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Review and manage the audience receiving published article notifications.
          </p>
        </div>
        <Button render={<Link href={routes.adminSubscriberNew()} />}>New subscriber</Button>
      </header>

      <Card>
        <CardHeader>
          <CardTitle>All subscribers</CardTitle>
          <CardDescription>
            {pagination.TotalCount === 0
              ? 'No subscribers yet.'
              : `Showing ${firstItem}-${lastItem} of ${pagination.TotalCount} subscribers.`}
          </CardDescription>
          {selectedCount > 0 ? (
            <CardAction>
              <AlertDialog onOpenChange={handleDeleteDialogOpenChange}>
                <AlertDialogTrigger render={<Button variant="destructive" size="sm" />}>
                  Delete selected ({selectedCount})
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Delete selected subscribers?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This permanently deletes {selectedCount} selected {selectedCount === 1 ? 'subscriber' : 'subscribers'}. This action cannot be undone.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <Field>
                    <FieldLabel htmlFor="bulkDeleteConfirmation">Type DELETE to confirm</FieldLabel>
                    <Input
                      id="bulkDeleteConfirmation"
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
                      onClick={deleteSelectedSubscribers}
                    >
                      Delete {selectedCount} {selectedCount === 1 ? 'subscriber' : 'subscribers'}
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            </CardAction>
          ) : null}
        </CardHeader>
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-14 pl-8">
                  <Checkbox
                    checked={allSelected}
                    indeterminate={someSelected}
                    onCheckedChange={selectAll}
                    aria-label="Select all subscribers on this page"
                  />
                </TableHead>
                <TableHead>Email</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Subscribed</TableHead>
                <TableHead>Referrer</TableHead>
                <TableHead className="pr-8 text-right">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-32 text-center text-muted-foreground">
                    No subscribers yet. Add your first subscriber to get started.
                  </TableCell>
                </TableRow>
              ) : (
                items.map((item) => (
                  <ClickableRow
                    key={item.ID}
                    href={routes.adminSubscriberShow(item.ID)}
                    label={`View subscriber ${item.Email}`}
                  >
                    <TableCell
                      className="w-14 pl-8"
                      onClick={stopRowNavigation}
                      onKeyDown={stopRowNavigation}
                    >
                      <Checkbox
                        checked={selectedIDs.has(item.ID)}
                        onCheckedChange={(checked) => selectSubscriber(item.ID, checked)}
                        aria-label={`Select ${item.Email}`}
                      />
                    </TableCell>
                    <TableCell className="max-w-80">
                      <p className="truncate font-medium text-foreground">{item.Email}</p>
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className={item.IsVerified ? 'text-emerald-400' : undefined}
                      >
                        {item.IsVerified ? 'Verified' : 'Unverified'}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatSubscriberDate(item.SubscribedAt)}
                    </TableCell>
                    <TableCell className="max-w-64 truncate text-muted-foreground">
                      {item.Referer || 'Direct'}
                    </TableCell>
                    <TableCell className="pr-8 text-right text-muted-foreground">
                      {formatSubscriberDate(item.UpdatedAt)}
                    </TableCell>
                  </ClickableRow>
                ))
              )}
            </TableBody>
          </Table>

          {pagination.TotalPages > 1 ? (
            <nav
              className="flex flex-col gap-3 border-t px-8 py-5 sm:flex-row sm:items-center sm:justify-between"
              aria-label="Subscriber pagination"
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
                  <Button variant="outline" size="sm" disabled>
                    Previous
                  </Button>
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
                  <Button variant="outline" size="sm" disabled>
                    Next
                  </Button>
                )}
              </div>
            </nav>
          ) : null}
        </CardContent>
      </Card>
    </main>
  )
}
