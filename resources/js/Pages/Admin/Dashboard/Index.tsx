import { type KeyboardEvent, type ReactNode } from 'react'
import { Link, router } from '@inertiajs/react'

import { Badge } from '@/components/ui/badge'
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { routes } from '@/routes'

type DashboardStats = {
  Subscribers: number
  Articles: number
  Newsletters: number
}

type Subscriber = {
  ID: number
  UpdatedAt: string
  Email: string
  IsVerified: boolean
}

type Article = {
  ID: number
  UpdatedAt: string
  Published: boolean
  Title: string
}

type Newsletter = {
  ID: number
  UpdatedAt: string
  Title: string
  IsPublished: boolean
}

type IndexProps = {
  stats: DashboardStats
  subscribers: Subscriber[]
  articles: Article[]
  newsletters: Newsletter[]
}

const dateFormatter = new Intl.DateTimeFormat('en', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
})

function formatDate(value: string) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : dateFormatter.format(date)
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

function EmptyRow({ columns, label }: { columns: number; label: string }) {
  return (
    <TableRow>
      <TableCell colSpan={columns} className="h-24 text-center text-muted-foreground">
        {label}
      </TableCell>
    </TableRow>
  )
}

export default function Index({ stats, subscribers, articles, newsletters }: IndexProps) {
  const statCards = [
    { label: 'Subscribers', value: stats.Subscribers },
    { label: 'Articles', value: stats.Articles },
    { label: 'Newsletters', value: stats.Newsletters },
  ]

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <div>
        <h1 className="font-heading text-3xl font-semibold tracking-tight">Dashboard</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          A quick overview of your audience and published content.
        </p>
      </div>

      <section aria-label="Resource totals" className="grid gap-4 sm:grid-cols-3">
        {statCards.map((stat) => (
          <Card key={stat.label} size="sm">
            <CardHeader>
              <CardDescription>{stat.label}</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="font-heading text-4xl font-semibold tabular-nums">{stat.value}</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <section className="grid gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Recent subscribers</CardTitle>
            <CardDescription>The five most recently updated subscribers.</CardDescription>
            <CardAction>
              <Link
                href={routes.adminSubscriberIndex()}
                className="text-xs font-semibold uppercase tracking-widest text-muted-foreground hover:text-foreground"
              >
                View all
              </Link>
            </CardAction>
          </CardHeader>
          <CardContent className="px-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="pl-8">Email</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="pr-8 text-right">Updated</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {subscribers.length === 0 ? (
                  <EmptyRow columns={3} label="No subscribers yet." />
                ) : (
                  subscribers.map((subscriber) => (
                    <ClickableRow
                      key={subscriber.ID}
                      href={routes.adminSubscriberShow(subscriber.ID)}
                      label={`View subscriber ${subscriber.Email}`}
                    >
                      <TableCell className="pl-8 font-medium text-foreground">
                        {subscriber.Email}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className={subscriber.IsVerified ? 'text-emerald-400' : undefined}
                        >
                          {subscriber.IsVerified ? 'Verified' : 'Unverified'}
                        </Badge>
                      </TableCell>
                      <TableCell className="pr-8 text-right text-muted-foreground">
                        {formatDate(subscriber.UpdatedAt)}
                      </TableCell>
                    </ClickableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Recent articles</CardTitle>
            <CardDescription>The five most recently updated articles.</CardDescription>
            <CardAction>
              <Link
                href={routes.adminArticleIndex()}
                className="text-xs font-semibold uppercase tracking-widest text-muted-foreground hover:text-foreground"
              >
                View all
              </Link>
            </CardAction>
          </CardHeader>
          <CardContent className="px-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="pl-8">Title</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="pr-8 text-right">Updated</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {articles.length === 0 ? (
                  <EmptyRow columns={3} label="No articles yet." />
                ) : (
                  articles.map((article) => (
                    <ClickableRow
                      key={article.ID}
                      href={routes.adminArticleShow(article.ID)}
                      label={`View article ${article.Title}`}
                    >
                      <TableCell className="max-w-72 truncate pl-8 font-medium text-foreground">
                        {article.Title}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className={article.Published ? 'text-emerald-400' : undefined}
                        >
                          {article.Published ? 'Published' : 'Draft'}
                        </Badge>
                      </TableCell>
                      <TableCell className="pr-8 text-right text-muted-foreground">
                        {formatDate(article.UpdatedAt)}
                      </TableCell>
                    </ClickableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>

        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle>Recent newsletters</CardTitle>
            <CardDescription>The five most recently updated newsletters.</CardDescription>
            <CardAction>
              <Link
                href={routes.adminNewsletterIndex()}
                className="text-xs font-semibold uppercase tracking-widest text-muted-foreground hover:text-foreground"
              >
                View all
              </Link>
            </CardAction>
          </CardHeader>
          <CardContent className="px-0">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="pl-8">Title</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="pr-8 text-right">Updated</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {newsletters.length === 0 ? (
                  <EmptyRow columns={3} label="No newsletters yet." />
                ) : (
                  newsletters.map((newsletter) => (
                    <ClickableRow
                      key={newsletter.ID}
                      href={routes.adminNewsletterShow(newsletter.ID)}
                      label={`View newsletter ${newsletter.Title}`}
                    >
                      <TableCell className="max-w-72 truncate pl-8 font-medium text-foreground">
                        {newsletter.Title}
                      </TableCell>
                      <TableCell>
                        <Badge
                          variant="secondary"
                          className={newsletter.IsPublished ? 'text-emerald-400' : undefined}
                        >
                          {newsletter.IsPublished ? 'Published' : 'Draft'}
                        </Badge>
                      </TableCell>
                      <TableCell className="pr-8 text-right text-muted-foreground">
                        {formatDate(newsletter.UpdatedAt)}
                      </TableCell>
                    </ClickableRow>
                  ))
                )}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </section>
    </main>
  )
}
