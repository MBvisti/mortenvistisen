import { type KeyboardEvent, type ReactNode } from 'react'
import { Link, router } from '@inertiajs/react'

import {
  formatNewsletterDate,
  type NewsletterItem,
  type NewsletterPagination,
} from './types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
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
  items: NewsletterItem[]
  pagination: NewsletterPagination
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

function pageHref(page: number) {
  return `${routes.adminNewsletterIndex()}?page=${page}`
}

export default function Index({ items, pagination }: IndexProps) {
  const firstItem = pagination.TotalCount === 0 ? 0 : (pagination.Page - 1) * pagination.PageSize + 1
  const lastItem = Math.min(
    pagination.Page * pagination.PageSize,
    pagination.TotalCount,
  )
  const hasPreviousPage = pagination.Page > 1
  const hasNextPage = pagination.Page < pagination.TotalPages

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-heading text-3xl font-semibold tracking-tight">Newsletters</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Create, review, and publish newsletter content.
          </p>
        </div>
        <Button render={<Link href={routes.adminNewsletterNew()} />}>New newsletter</Button>
      </header>

      <Card>
        <CardHeader>
          <CardTitle>All newsletters</CardTitle>
          <CardDescription>
            {pagination.TotalCount === 0
              ? 'No newsletters yet.'
              : `Showing ${firstItem}-${lastItem} of ${pagination.TotalCount} newsletters.`}
          </CardDescription>
        </CardHeader>
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="pl-8">Title</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Slug</TableHead>
                <TableHead>Released</TableHead>
                <TableHead className="pr-8 text-right">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                    No newsletters yet. Create your first newsletter to get started.
                  </TableCell>
                </TableRow>
              ) : (
                items.map((item) => (
                  <ClickableRow
                    key={item.ID}
                    href={routes.adminNewsletterShow(item.ID)}
                    label={`View newsletter ${item.Title || 'Untitled newsletter'}`}
                  >
                    <TableCell className="max-w-80 pl-8">
                      <p className="truncate font-medium text-foreground">
                        {item.Title || 'Untitled newsletter'}
                      </p>
                      {item.MetaDescription ? (
                        <p className="mt-1 max-w-xl truncate text-xs text-muted-foreground">
                          {item.MetaDescription}
                        </p>
                      ) : null}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className={item.IsPublished ? 'text-emerald-400' : undefined}
                      >
                        {item.IsPublished ? 'Published' : 'Draft'}
                      </Badge>
                    </TableCell>
                    <TableCell className="max-w-64 truncate text-muted-foreground">
                      /{item.Slug}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatNewsletterDate(item.ReleasedAt)}
                    </TableCell>
                    <TableCell className="pr-8 text-right text-muted-foreground">
                      {formatNewsletterDate(item.UpdatedAt)}
                    </TableCell>
                  </ClickableRow>
                ))
              )}
            </TableBody>
          </Table>

          {pagination.TotalPages > 1 ? (
            <nav
              className="flex flex-col gap-3 border-t px-8 py-5 sm:flex-row sm:items-center sm:justify-between"
              aria-label="Newsletter pagination"
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
