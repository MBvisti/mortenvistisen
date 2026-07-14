import type { FormEvent } from 'react'
import { Link, type InertiaFormProps } from '@inertiajs/react'

import type { SubscriberFormData, SubscriberValidationRules } from './types'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Field,
  FieldDescription,
  FieldError,
  FieldGroup,
  FieldLabel,
} from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'

function characterLength(value: string) {
  return Array.from(value).length
}

function findRule(rules: SubscriberValidationRules, field: string, code: string) {
  return rules[field]?.find((rule) => rule.code === code)
}

function ruleNumber(
  rules: SubscriberValidationRules,
  field: string,
  code: string,
  param: string,
) {
  const value = findRule(rules, field, code)?.params?.[param]
  if (typeof value !== 'number') {
    throw new Error(`Missing numeric validation rule ${field}.${code}.${param}`)
  }

  return value
}

function MaximumLengthIndicator({ value, max }: { value: string; max: number }) {
  const length = characterLength(value)
  const tooLong = length > max

  return (
    <div className="flex flex-wrap items-center justify-between gap-2">
      <FieldDescription>Maximum {max} characters.</FieldDescription>
      <span
        className={`text-xs font-medium tabular-nums ${tooLong ? 'text-destructive' : 'text-muted-foreground'}`}
        aria-live="polite"
      >
        {length}/{max}{tooLong ? ' · Too long' : ''}
      </span>
    </div>
  )
}

type SubscriberFormProps = {
  form: InertiaFormProps<SubscriberFormData>
  validationRules: SubscriberValidationRules
  cancelHref: string
  submitLabel: string
  processingLabel: string
  onSubmit: () => void
}

export function SubscriberForm({
  form,
  validationRules,
  cancelHref,
  submitLabel,
  processingLabel,
  onSubmit,
}: SubscriberFormProps) {
  const emailRequired = Boolean(findRule(validationRules, 'email', 'required'))
  const subscribedAtRequired = Boolean(findRule(validationRules, 'subscribedAt', 'required'))
  const emailMaxLength = ruleNumber(validationRules, 'email', 'max', 'max')
  const refererMaxLength = ruleNumber(validationRules, 'referer', 'max', 'max')
  const emailInvalid =
    (emailRequired && form.data.email.trim() === '') ||
    characterLength(form.data.email) > emailMaxLength ||
    Boolean(form.errors.email)
  const refererInvalid =
    characterLength(form.data.referer) > refererMaxLength || Boolean(form.errors.referer)

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    onSubmit()
  }

  return (
    <form onSubmit={submit} className="grid gap-6 lg:grid-cols-[minmax(0,1fr)_20rem]">
      <Card>
        <CardHeader>
          <CardTitle>Subscriber details</CardTitle>
          <CardDescription>
            Set the address used for newsletter delivery and its optional referrer.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <FieldGroup className="gap-6">
            <Field data-invalid={emailInvalid}>
              <FieldLabel htmlFor="email">Email address</FieldLabel>
              <Input
                id="email"
                type="email"
                value={form.data.email}
                onChange={(event) => {
                  form.setData('email', event.currentTarget.value)
                  form.clearErrors('email')
                }}
                placeholder="reader@example.com"
                required={emailRequired}
                maxLength={emailMaxLength}
                aria-invalid={emailInvalid}
                autoComplete="email"
                autoFocus
              />
              <MaximumLengthIndicator value={form.data.email} max={emailMaxLength} />
              <FieldError>{form.errors.email}</FieldError>
            </Field>

            <Field data-invalid={refererInvalid}>
              <FieldLabel htmlFor="referer">Referrer</FieldLabel>
              <Input
                id="referer"
                type="text"
                value={form.data.referer}
                onChange={(event) => {
                  form.setData('referer', event.currentTarget.value)
                  form.clearErrors('referer')
                }}
                placeholder="/newsletter or campaign name"
                maxLength={refererMaxLength}
                aria-invalid={refererInvalid}
              />
              <FieldDescription>
                Optional source where this subscriber joined the list.
              </FieldDescription>
              <FieldError>{form.errors.referer}</FieldError>
            </Field>
          </FieldGroup>
        </CardContent>
      </Card>

      <div className="flex flex-col gap-6 lg:sticky lg:top-6 lg:self-start">
        <Card>
          <CardHeader>
            <CardTitle>Subscription</CardTitle>
            <CardDescription>Record when the subscription began and its status.</CardDescription>
          </CardHeader>
          <CardContent>
            <FieldGroup className="gap-6">
              <Field data-invalid={Boolean(form.errors.subscribedAt)}>
                <FieldLabel htmlFor="subscribedAt">Subscribed at</FieldLabel>
                <Input
                  id="subscribedAt"
                  type="date"
                  value={form.data.subscribedAt}
                  onChange={(event) => {
                    form.setData('subscribedAt', event.currentTarget.value)
                    form.clearErrors('subscribedAt')
                  }}
                  required={subscribedAtRequired}
                  aria-invalid={Boolean(form.errors.subscribedAt)}
                />
                <FieldError>{form.errors.subscribedAt}</FieldError>
              </Field>

              <div className="flex items-center gap-3">
                <Switch
                  id="isVerified"
                  checked={form.data.isVerified}
                  onCheckedChange={(checked) => form.setData('isVerified', checked)}
                />
                <FieldLabel htmlFor="isVerified">
                  {form.data.isVerified ? 'Verified' : 'Unverified'}
                </FieldLabel>
              </div>
              <FieldDescription>
                Only verified subscribers receive published article notifications.
              </FieldDescription>
            </FieldGroup>
          </CardContent>
        </Card>

        <div className="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end lg:flex-col-reverse">
          <Button variant="outline" render={<Link href={cancelHref} />}>
            Cancel
          </Button>
          <Button type="submit" disabled={form.processing}>
            {form.processing ? processingLabel : submitLabel}
          </Button>
        </div>
      </div>
    </form>
  )
}
