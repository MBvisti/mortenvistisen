import type { TagItem } from '../Tag/types'

export type ArticleItem = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  FirstPublishedAt: string
  Published: boolean
  Title: string
  Excerpt: string
  MetaTitle: string
  MetaDescription: string
  Slug: string
  ImageLink: string
  MetaImageLink: string
  ReadTime: number
  Content: string
  Tags: TagItem[]
}

export type ArticleFormData = {
  published: boolean
  title: string
  excerpt: string
  metaTitle: string
  metaDescription: string
  slug: string
  cover: File | null
  removeCover: boolean
  metaImage: File | null
  removeMetaImage: boolean
  readTime: number
  content: string
  tagIds: number[]
  newTags: string[]
}

export function toArticleMultipartData(data: ArticleFormData) {
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

export type ArticleValidationRules = Record<string, ValidationRule[]>

export type ArticlePagination = {
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

export function formatArticleDate(value: string, fallback = 'Not set') {
  if (!value) {
    return fallback
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return fallback
  }

  return dateFormatter.format(date)
}

export function articleDateInputValue(value: string) {
  if (!value) {
    return ''
  }

  const date = new Date(value)
  if (Number.isNaN(date.getTime()) || date.getUTCFullYear() <= 1) {
    return ''
  }

  return value.slice(0, 10)
}
