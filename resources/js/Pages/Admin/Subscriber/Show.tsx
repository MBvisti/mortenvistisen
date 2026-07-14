import { useState } from 'react'
import { Link, router } from '@inertiajs/react'

import { formatSubscriberDate, type SubscriberItem } from './types'
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
import { Field, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { routes } from '@/routes'

type ShowProps = {
  item: SubscriberItem
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
  const [sendingConfirmation, setSendingConfirmation] = useState(false)
  const canDelete = deleteConfirmation === 'DELETE'

  function handleDeleteDialogOpenChange(open: boolean) {
    if (!open) {
      setDeleteConfirmation('')
    }
  }

  function deleteSubscriber() {
    if (!canDelete) {
      return
    }

    router.delete(routes.adminSubscriberDestroy(item.ID))
  }

  function resendConfirmation() {
    router.post(
      routes.adminSubscriberConfirmationCreate(item.ID),
      {},
      {
        preserveScroll: true,
        onStart: () => setSendingConfirmation(true),
        onFinish: () => setSendingConfirmation(false),
      },
    )
  }

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-3">
            <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
              Subscriber
            </p>
            <Badge
              variant="secondary"
              className={item.IsVerified ? 'text-emerald-400' : undefined}
            >
              {item.IsVerified ? 'Verified' : 'Unverified'}
            </Badge>
          </div>
          <h1 className="mt-3 truncate font-heading text-3xl font-semibold tracking-tight sm:text-4xl">
            {item.Email || 'Subscriber'}
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Subscribed {formatSubscriberDate(item.SubscribedAt).toLowerCase()}.
          </p>
        </div>
        <div className="flex flex-wrap gap-3">
          <Button variant="outline" render={<Link href={routes.adminSubscriberIndex()} />}>
            Back to subscribers
          </Button>
          <Button render={<Link href={routes.adminSubscriberEdit(item.ID)} />}>
            Edit subscriber
          </Button>
        </div>
      </header>

      <div className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
        <Card>
          <CardHeader>
            <CardTitle>Subscriber details</CardTitle>
            <CardDescription>Delivery address and acquisition source.</CardDescription>
          </CardHeader>
          <CardContent>
            <dl>
              <Detail label="Email">
                <a
                  href={`mailto:${item.Email}`}
                  className="break-all text-primary underline-offset-4 hover:underline"
                >
                  {item.Email}
                </a>
              </Detail>
              <Detail label="Referrer">
                {item.Referer ? (
                  <a
                    href={item.Referer}
                    target="_blank"
                    rel="noreferrer"
                    className="break-all text-primary underline-offset-4 hover:underline"
                  >
                    {item.Referer}
                  </a>
                ) : (
                  'Direct or unknown'
                )}
              </Detail>
            </dl>
          </CardContent>
        </Card>

        <div className="flex flex-col gap-6">
          <Card>
            <CardHeader>
              <CardTitle>Subscription</CardTitle>
              {!item.IsVerified && (
                <CardDescription>
                  This subscriber must confirm their email before receiving notifications.
                </CardDescription>
              )}
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Status">{item.IsVerified ? 'Verified' : 'Unverified'}</Detail>
                <Detail label="Subscribed">{formatSubscriberDate(item.SubscribedAt)}</Detail>
              </dl>
              {!item.IsVerified && (
                <Button
                  variant="outline"
                  className="mt-4 w-full"
                  disabled={sendingConfirmation}
                  onClick={resendConfirmation}
                >
                  {sendingConfirmation ? 'Sending confirmation...' : 'Resend confirmation email'}
                </Button>
              )}
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Record details</CardTitle>
            </CardHeader>
            <CardContent>
              <dl>
                <Detail label="Created">{formatSubscriberDate(item.CreatedAt)}</Detail>
                <Detail label="Updated">{formatSubscriberDate(item.UpdatedAt)}</Detail>
              </dl>
            </CardContent>
          </Card>

          <AlertDialog onOpenChange={handleDeleteDialogOpenChange}>
            <AlertDialogTrigger render={<Button variant="destructive" className="w-full" />}>
              Delete subscriber
            </AlertDialogTrigger>
            <AlertDialogContent>
              <AlertDialogHeader>
                <AlertDialogTitle>Delete this subscriber?</AlertDialogTitle>
                <AlertDialogDescription>
                  This permanently deletes {item.Email || 'this subscriber'}. This action cannot be undone.
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
                  onClick={deleteSubscriber}
                >
                  Delete subscriber
                </AlertDialogAction>
              </AlertDialogFooter>
            </AlertDialogContent>
          </AlertDialog>
        </div>
      </div>
    </main>
  )
}
