import {
  type FormEvent,
  type KeyboardEvent,
  type ReactNode,
  useState,
} from 'react'
import { Link, router } from '@inertiajs/react'
import { Search, X } from 'lucide-react'

import {
  formatJobDate,
  jobStateClass,
  jobStateLabel,
  jobStates,
  type JobFilters,
  type JobPagination,
  type JobState,
  type JobStats,
  type JobSummary,
} from './types'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { PollingControl } from '@/Components/PollingControl'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { routes } from '@/routes'

const jobIndexPollingProps = ['items', 'pagination', 'stats']

type IndexProps = {
  items: JobSummary[]
  pagination: JobPagination
  filters: JobFilters
  stats: JobStats
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

function pageHref(page: number, filters: JobFilters) {
  const params = new URLSearchParams()
  params.set('page', String(page))
  if (filters.Search) params.set('search', filters.Search)
  if (filters.State) params.set('state', filters.State)
  return `${routes.adminJobIndex()}?${params.toString()}`
}

function lastActivity(job: JobSummary) {
  return job.FinalizedAt ?? job.AttemptedAt ?? job.CreatedAt
}

export default function Index({ items, pagination, filters, stats }: IndexProps) {
  const [search, setSearch] = useState(filters.Search)
  const [state, setState] = useState(filters.State)
  const activeCount =
    stats.ByState.available +
    stats.ByState.running +
    stats.ByState.scheduled +
    stats.ByState.retryable +
    stats.ByState.pending
  const attentionCount = stats.ByState.discarded + stats.ByState.cancelled

  const statCards = [
    { label: 'All jobs', value: stats.Total, description: 'Retained River records' },
    { label: 'Active', value: activeCount, description: 'Queued, scheduled, or running' },
    { label: 'Running', value: stats.ByState.running, description: 'Currently being worked' },
    { label: 'Needs attention', value: attentionCount, description: 'Discarded or cancelled' },
  ]

  function applyFilters(event: FormEvent) {
    event.preventDefault()
    router.get(
      routes.adminJobIndex(),
      { search: search.trim() || undefined, state: state || undefined },
      { preserveState: true, replace: true },
    )
  }

  function clearFilters() {
    setSearch('')
    setState('')
    router.get(routes.adminJobIndex(), {}, { preserveState: true, replace: true })
  }

  const hasFilters = Boolean(filters.Search || filters.State)
  const firstItem = pagination.TotalCount === 0 ? 0 : (pagination.Page - 1) * pagination.PageSize + 1
  const lastItem = Math.min(pagination.Page * pagination.PageSize, pagination.TotalCount)

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="font-heading text-3xl font-semibold tracking-tight">Queue jobs</h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Inspect River activity and intervene when a job needs attention.
          </p>
        </div>
        <PollingControl only={jobIndexPollingProps} />
      </header>

      <section aria-label="Queue totals" className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        {statCards.map((stat) => (
          <Card key={stat.label} size="sm">
            <CardHeader>
              <CardDescription>{stat.label}</CardDescription>
            </CardHeader>
            <CardContent>
              <p className="font-heading text-4xl font-semibold tabular-nums">{stat.value}</p>
              <p className="mt-2 text-xs text-muted-foreground">{stat.description}</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <Card>
        <CardHeader>
          <CardTitle>Jobs</CardTitle>
          <CardDescription>
            {pagination.TotalCount === 0
              ? hasFilters
                ? 'No jobs match these filters.'
                : 'No jobs have been enqueued yet.'
              : `Showing ${firstItem}-${lastItem} of ${pagination.TotalCount} matching jobs.`}
          </CardDescription>
        </CardHeader>
        <CardContent className="px-0">
          <form
            onSubmit={applyFilters}
            className="flex flex-col gap-3 border-b px-6 pb-6 sm:flex-row sm:items-center"
          >
            <div className="relative min-w-0 flex-1">
              <Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
              <Input
                value={search}
                onChange={(event) => setSearch(event.target.value)}
                placeholder="Search by ID, kind, or queue"
                aria-label="Search jobs"
                className="pl-9"
              />
            </div>
            <select
              value={state}
              onChange={(event) => setState(event.target.value)}
              aria-label="Filter by state"
              style={{ colorScheme: 'dark' }}
              className="h-9 border border-input bg-popover px-3 text-sm text-popover-foreground outline-none focus-visible:ring-2 focus-visible:ring-ring"
            >
              <option value="" className="bg-popover text-popover-foreground">All states</option>
              {jobStates.map((jobState) => (
                <option
                  key={jobState}
                  value={jobState}
                  className="bg-popover text-popover-foreground"
                >
                  {jobStateLabel(jobState)} ({stats.ByState[jobState]})
                </option>
              ))}
            </select>
            <Button type="submit">Apply</Button>
            {hasFilters ? (
              <Button type="button" variant="ghost" onClick={clearFilters} className="gap-2">
                <X className="size-4" /> Clear
              </Button>
            ) : null}
          </form>

          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="pl-8">Job</TableHead>
                <TableHead>State</TableHead>
                <TableHead>Queue</TableHead>
                <TableHead>Attempts</TableHead>
                <TableHead>Scheduled</TableHead>
                <TableHead className="pr-8 text-right">Last activity</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-32 text-center text-muted-foreground">
                    {hasFilters ? 'No jobs match these filters.' : 'The queue is empty.'}
                  </TableCell>
                </TableRow>
              ) : (
                items.map((job) => (
                  <ClickableRow
                    key={job.ID}
                    href={routes.adminJobShow(job.ID)}
                    label={`View job ${job.ID}`}
                  >
                    <TableCell className="max-w-80 pl-8">
                      <p className="truncate font-mono text-sm font-medium text-foreground">{job.Kind}</p>
                      <p className="mt-1 text-xs text-muted-foreground">#{job.ID}</p>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className={jobStateClass(job.State)}>
                        {jobStateLabel(job.State)}
                      </Badge>
                    </TableCell>
                    <TableCell className="font-mono text-xs text-muted-foreground">{job.Queue}</TableCell>
                    <TableCell className="tabular-nums text-muted-foreground">
                      {job.Attempt} / {job.MaxAttempts}
                    </TableCell>
                    <TableCell className="text-muted-foreground">{formatJobDate(job.ScheduledAt)}</TableCell>
                    <TableCell className="pr-8 text-right text-muted-foreground">
                      {formatJobDate(lastActivity(job))}
                    </TableCell>
                  </ClickableRow>
                ))
              )}
            </TableBody>
          </Table>

          {pagination.TotalPages > 1 ? (
            <nav
              className="flex flex-col gap-3 border-t px-8 py-5 sm:flex-row sm:items-center sm:justify-between"
              aria-label="Job pagination"
            >
              <p className="text-xs text-muted-foreground">
                Page {pagination.Page} of {pagination.TotalPages}
              </p>
              <div className="flex gap-3">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={pagination.Page <= 1}
                  render={
                    pagination.Page > 1 ? (
                      <Link href={pageHref(pagination.Page - 1, filters)} preserveScroll />
                    ) : undefined
                  }
                >
                  Previous
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={pagination.Page >= pagination.TotalPages}
                  render={
                    pagination.Page < pagination.TotalPages ? (
                      <Link href={pageHref(pagination.Page + 1, filters)} preserveScroll />
                    ) : undefined
                  }
                >
                  Next
                </Button>
              </div>
            </nav>
          ) : null}
        </CardContent>
      </Card>
    </main>
  )
}
