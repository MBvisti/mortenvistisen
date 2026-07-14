export type SubscriberItem = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  Email: string
  SubscribedAt: string
  Referer: string
  IsVerified: boolean
}

export type SubscriberFormData = {
  email: string
  subscribedAt: string
  referer: string
  isVerified: boolean
}

export type ValidationRule = {
  code: string
  message: string
  params?: Record<string, unknown>
}

export type SubscriberValidationRules = Record<string, ValidationRule[]>

export type SubscriberPagination = {
  Page: number
  PageSize: number
  TotalCount: number
  TotalPages: number
}

const dateFormatter = new Intl.DateTimeFormat('en', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
})

export function formatSubscriberDate(value: string, fallback = 'Not set') {
  if (!value) {
    return fallback
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return fallback
  }

  return dateFormatter.format(date)
}

export function subscriberDateInputValue(value: string) {
  if (!value) {
    return ''
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return ''
  }

  return value.slice(0, 10)
}
