import { Link, useForm } from '@inertiajs/react'

import { SubscriberForm } from './SubscriberForm'
import type { SubscriberFormData, SubscriberValidationRules } from './types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type CreateProps = {
  validationRules: SubscriberValidationRules
}

export default function Create({ validationRules }: CreateProps) {
  const form = useForm<SubscriberFormData>({
    email: '',
    subscribedAt: new Date().toISOString().slice(0, 10),
    referer: '',
    isVerified: false,
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Subscribers
          </p>
          <h1 className="mt-2 font-heading text-3xl font-semibold tracking-tight">
            New subscriber
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Add a reader to the newsletter audience.
          </p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminSubscriberIndex()} />}>
          Back to subscribers
        </Button>
      </header>

      <SubscriberForm
        form={form}
        validationRules={validationRules}
        cancelHref={routes.adminSubscriberIndex()}
        submitLabel="Create subscriber"
        processingLabel="Creating..."
        onSubmit={() => form.post(routes.adminSubscriberCreate())}
      />
    </main>
  )
}
