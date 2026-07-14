export type TagItem = {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  Title: string
  ArticleCount: number
}

export type TagValidationRule = {
  code: string
  message: string
  params?: Record<string, unknown>
}

export type TagValidationRules = Record<string, TagValidationRule[]>

export type TagPagination = {
  Page: number
  PageSize: number
  TotalCount: number
  TotalPages: number
}
