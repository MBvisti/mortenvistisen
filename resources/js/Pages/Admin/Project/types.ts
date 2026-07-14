export type ProjectItem = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  Published: boolean
  Title: string
  Slug: string
  StartedAt: string
  Status: string
  Description: string
  MetaTitle: string
  MetaDescription: string
  ImageLink: string
  MetaImageLink: string
  Content: string
  ProjectUrl: string
}

export type ProjectFormData = {
  published: boolean
  title: string
  slug: string
  startedAt: string
  status: string
  description: string
  metaTitle: string
  metaDescription: string
  cover: File | null
  removeCover: boolean
  metaImage: File | null
  removeMetaImage: boolean
  content: string
  projectUrl: string
}

export function toProjectMultipartData(data: ProjectFormData) {
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

export type ProjectValidationRules = Record<string, ValidationRule[]>

const dateFormatter = new Intl.DateTimeFormat('en', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
})

export function formatProjectDate(value: string, fallback = 'Not set') {
  if (!value) {
    return fallback
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return fallback
  }

  return dateFormatter.format(date)
}

export function projectDateInputValue(value: string) {
  if (!value) {
    return ''
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return ''
  }

  return value.slice(0, 10)
}
