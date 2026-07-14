import { Link, useForm } from '@inertiajs/react'

import { SubscriberForm } from './SubscriberForm'
import {
  subscriberDateInputValue,
  type SubscriberFormData,
  type SubscriberItem,
  type SubscriberValidationRules,
} from './types'
import { Button } from '@/components/ui/button'
import { routes } from '@/routes'

type EditProps = {
  item: SubscriberItem
  validationRules: SubscriberValidationRules
}

export default function Edit({ item, validationRules }: EditProps) {
  const form = useForm<SubscriberFormData>({
    email: String(item.Email ?? ''),
    subscribedAt: subscriberDateInputValue(item.SubscribedAt),
    referer: String(item.Referer ?? ''),
    isVerified: Boolean(item.IsVerified),
  })

  return (
    <main className="mx-auto flex max-w-7xl flex-col gap-8 px-4 py-8 sm:px-6 lg:px-8">
      <header className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div className="min-w-0">
          <p className="text-xs font-semibold uppercase tracking-widest text-muted-foreground">
            Edit subscriber
          </p>
          <h1 className="mt-2 truncate font-heading text-3xl font-semibold tracking-tight">
            {item.Email || 'Subscriber'}
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            Update delivery details and verification status.
          </p>
        </div>
        <Button variant="outline" render={<Link href={routes.adminSubscriberShow(item.ID)} />}>
          View subscriber
        </Button>
      </header>

      <SubscriberForm
        form={form}
        validationRules={validationRules}
        cancelHref={routes.adminSubscriberShow(item.ID)}
        submitLabel="Save changes"
        processingLabel="Saving..."
        onSubmit={() => form.put(routes.adminSubscriberUpdate(item.ID))}
      />
    </main>
  )
}
