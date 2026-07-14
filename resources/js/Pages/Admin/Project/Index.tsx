import { type KeyboardEvent, type ReactNode } from 'react'
import { Link, router } from '@inertiajs/react'

import { formatProjectDate, type ProjectItem } from './types'
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
  items: ProjectItem[]
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

export default function Index({ items }: IndexProps) {
  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-heading text-3xl font-semibold tracking-tight">Projects</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Create, review, and publish portfolio projects.
          </p>
        </div>
        <Button render={<Link href={routes.adminProjectNew()} />}>New project</Button>
      </header>

      <Card>
        <CardHeader>
          <CardTitle>All projects</CardTitle>
          <CardDescription>
            {items.length === 1 ? '1 project' : `${items.length} projects`} on this page.
          </CardDescription>
        </CardHeader>
        <CardContent className="px-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="pl-8">Title</TableHead>
                <TableHead>Publication</TableHead>
                <TableHead>Status</TableHead>
                <TableHead>Started</TableHead>
                <TableHead className="pr-8 text-right">Updated</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className="h-32 text-center text-muted-foreground">
                    No projects yet. Create your first project to get started.
                  </TableCell>
                </TableRow>
              ) : (
                items.map((item) => (
                  <ClickableRow
                    key={item.ID}
                    href={routes.adminProjectShow(item.ID)}
                    label={`View project ${item.Title || 'Untitled project'}`}
                  >
                    <TableCell className="max-w-80 pl-8">
                      <p className="truncate font-medium text-foreground">
                        {item.Title || 'Untitled project'}
                      </p>
                      {item.Description ? (
                        <p className="mt-1 max-w-xl truncate text-xs text-muted-foreground">
                          {item.Description}
                        </p>
                      ) : null}
                    </TableCell>
                    <TableCell>
                      <Badge
                        variant="secondary"
                        className={item.Published ? 'text-emerald-400' : undefined}
                      >
                        {item.Published ? 'Published' : 'Draft'}
                      </Badge>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {item.Status || 'Not set'}
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {formatProjectDate(item.StartedAt)}
                    </TableCell>
                    <TableCell className="pr-8 text-right text-muted-foreground">
                      {formatProjectDate(item.UpdatedAt)}
                    </TableCell>
                  </ClickableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>
    </main>
  )
}
