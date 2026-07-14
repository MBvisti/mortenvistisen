import { useState } from 'react'
import { Link, router } from '@inertiajs/react'
import { ArrowLeft, Ban, Play, RotateCcw, Trash2 } from 'lucide-react'

import { formatJobDate, jobStateClass, jobStateLabel, type Job } from './types'
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
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { JsonCodeBlock } from '@/Components/JsonCodeBlock'
import { PollingControl } from '@/Components/PollingControl'
import { routes } from '@/routes'

const jobShowPollingProps = ['item']

type ShowProps = {
  item: Job
}

function Detail({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div className="grid gap-1 border-b py-4 last:border-b-0 sm:grid-cols-[10rem_minmax(0,1fr)] sm:gap-6">
      <dt className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">{label}</dt>
      <dd className="min-w-0 text-sm text-foreground">{children}</dd>
    </div>
  )
}

export default function Show({ item }: ShowProps) {
  const [activeAction, setActiveAction] = useState<string | null>(null)

  function postAction(action: string, url: string) {
    setActiveAction(action)
    router.post(url, {}, { preserveScroll: true, onFinish: () => setActiveAction(null) })
  }

  function deleteJob() {
    setActiveAction('delete')
    router.delete(routes.adminJobDestroy(item.ID), {
      onFinish: () => setActiveAction(null),
    })
  }

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-5">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <Button
            variant="ghost"
            size="sm"
            render={<Link href={routes.adminJobIndex()} />}
            className="-ml-3 gap-2 self-start"
          >
            <ArrowLeft className="size-4" />
            All jobs
          </Button>

          <div className="flex flex-wrap items-center gap-3 sm:justify-end">
            {item.State !== 'completed' ? (
              <PollingControl only={jobShowPollingProps} suspended={activeAction !== null} />
            ) : null}

            {item.CanRun ? (
              <Button
                onClick={() => postAction('run', routes.adminJobRun(item.ID))}
                disabled={activeAction !== null}
                className="gap-2"
              >
                <Play className="size-4" />
                {activeAction === 'run' ? 'Starting...' : 'Run now'}
              </Button>
            ) : null}

            {item.CanRestart ? (
              <AlertDialog>
                <AlertDialogTrigger
                  render={
                    <Button variant="outline" disabled={activeAction !== null} className="gap-2">
                      <RotateCcw className="size-4" /> Restart
                    </Button>
                  }
                />
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Restart this job?</AlertDialogTitle>
                    <AlertDialogDescription>
                      River will make the job available immediately. Its existing attempt and error history will be kept.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Keep current state</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={() => postAction('restart', routes.adminJobRetry(item.ID))}
                    >
                      Restart job
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            ) : null}

            {item.CanCancel ? (
              <AlertDialog>
                <AlertDialogTrigger
                  render={
                    <Button variant="outline" disabled={activeAction !== null} className="gap-2">
                      <Ban className="size-4" /> Cancel
                    </Button>
                  }
                />
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Cancel this job?</AlertDialogTitle>
                    <AlertDialogDescription>
                      Queued jobs stop immediately. Running jobs receive a cancellation signal and may take time to stop.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Keep job</AlertDialogCancel>
                    <AlertDialogAction
                      variant="destructive"
                      onClick={() => postAction('cancel', routes.adminJobCancel(item.ID))}
                    >
                      Cancel job
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            ) : null}

            {item.CanDelete ? (
              <AlertDialog>
                <AlertDialogTrigger
                  render={
                    <Button variant="destructive" disabled={activeAction !== null} className="gap-2">
                      <Trash2 className="size-4" /> Delete
                    </Button>
                  }
                />
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Permanently delete this job?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This removes the River record and all attempt history. This action cannot be undone.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Keep job</AlertDialogCancel>
                    <AlertDialogAction variant="destructive" onClick={deleteJob}>
                      Delete permanently
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            ) : null}
          </div>
        </div>

        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="truncate font-mono text-2xl font-semibold tracking-tight sm:text-3xl">
              {item.Kind}
            </h1>
            <Badge variant="outline" className={jobStateClass(item.State)}>
              {jobStateLabel(item.State)}
            </Badge>
          </div>
          <p className="mt-2 font-mono text-sm text-muted-foreground">Job #{item.ID}</p>
        </div>
      </header>

      <section className="grid gap-6 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Execution</CardTitle>
            <CardDescription>How River will process this job.</CardDescription>
          </CardHeader>
          <CardContent>
            <dl>
              <Detail label="Queue"><span className="font-mono">{item.Queue}</span></Detail>
              <Detail label="Priority">{item.Priority} <span className="text-muted-foreground">(1 is highest)</span></Detail>
              <Detail label="Attempts"><span className="tabular-nums">{item.Attempt} of {item.MaxAttempts}</span></Detail>
              <Detail label="Tags">
                {item.Tags.length > 0 ? (
                  <div className="flex flex-wrap gap-2">
                    {item.Tags.map((tag) => <Badge key={tag} variant="secondary">{tag}</Badge>)}
                  </div>
                ) : <span className="text-muted-foreground">None</span>}
              </Detail>
            </dl>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Timeline</CardTitle>
            <CardDescription>Creation, scheduling, and execution timestamps.</CardDescription>
          </CardHeader>
          <CardContent>
            <dl>
              <Detail label="Created">{formatJobDate(item.CreatedAt)}</Detail>
              <Detail label="Scheduled">{formatJobDate(item.ScheduledAt)}</Detail>
              <Detail label="Last attempted">{formatJobDate(item.AttemptedAt)}</Detail>
              <Detail label="Finalized">{formatJobDate(item.FinalizedAt)}</Detail>
            </dl>
          </CardContent>
        </Card>
      </section>

      <section className="grid gap-6 xl:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Arguments</CardTitle>
            <CardDescription>The JSON payload passed to the worker.</CardDescription>
          </CardHeader>
          <CardContent>
            <JsonCodeBlock label="args.json" value={item.Args} />
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Metadata</CardTitle>
            <CardDescription>River metadata, including recorded job output when present.</CardDescription>
          </CardHeader>
          <CardContent>
            <JsonCodeBlock label="metadata.json" value={item.Metadata} />
          </CardContent>
        </Card>

        <Card className="xl:col-span-2">
          <CardHeader>
            <CardTitle>Attempt errors</CardTitle>
            <CardDescription>
              {item.Errors.length === 0
                ? 'No failures have been recorded.'
                : `${item.Errors.length} failed attempt${item.Errors.length === 1 ? '' : 's'} recorded.`}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <JsonCodeBlock label="errors.json" value={item.Errors} />
          </CardContent>
        </Card>
      </section>

      <Card>
        <CardHeader>
          <CardTitle>River identifiers</CardTitle>
          <CardDescription>Low-level values used for execution and uniqueness.</CardDescription>
        </CardHeader>
        <CardContent>
          <dl>
            <Detail label="Attempted by">
              {item.AttemptedBy.length > 0
                ? <span className="break-all font-mono text-xs">{item.AttemptedBy.join(', ')}</span>
                : <span className="text-muted-foreground">No client yet</span>}
            </Detail>
            <Detail label="Unique key">
              {item.UniqueKey
                ? <span className="break-all font-mono text-xs">{item.UniqueKey}</span>
                : <span className="text-muted-foreground">Not unique</span>}
            </Detail>
            <Detail label="Unique states">
              {item.UniqueStates.length > 0
                ? item.UniqueStates.map(jobStateLabel).join(', ')
                : <span className="text-muted-foreground">Not configured</span>}
            </Detail>
          </dl>
        </CardContent>
      </Card>
    </main>
  )
}
