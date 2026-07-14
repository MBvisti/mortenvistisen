export type NewsletterItem = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  Title: string
  Slug: string
  MetaTitle: string
  MetaDescription: string
  IsPublished: boolean
  ReleasedAt: string
  ImageLink: string
  MetaImageLink: string
  Content: string
}

export type NewsletterFormData = {
  title: string
  slug: string
  metaTitle: string
  metaDescription: string
  isPublished: boolean
  cover: File | null
  removeCover: boolean
  metaImage: File | null
  removeMetaImage: boolean
  content: string
}

export function toNewsletterMultipartData(data: NewsletterFormData) {
  const { cover, metaImage, ...payload } = data

  return {
    payload: JSON.stringify(payload),
    cover,
    metaImage,
  }
}

export type ValidationRule = {
  code: string
  message: string
  params?: Record<string, unknown>
}

export type NewsletterValidationRules = Record<string, ValidationRule[]>

export type NewsletterPagination = {
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

export function formatNewsletterDate(value: string, fallback = 'Not set') {
  if (!value) {
    return fallback
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return fallback
  }

  return dateFormatter.format(date)
}
