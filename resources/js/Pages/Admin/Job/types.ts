export const jobStates = [
  'available',
  'running',
  'scheduled',
  'retryable',
  'pending',
  'completed',
  'discarded',
  'cancelled',
] as const

export type JobState = (typeof jobStates)[number]

export type JobSummary = {
  ID: number
  State: JobState
  Attempt: number
  MaxAttempts: number
  AttemptedAt: string | null
  CreatedAt: string
  FinalizedAt: string | null
  ScheduledAt: string
  Priority: number
  Kind: string
  Queue: string
}

export type Job = {
  ID: number
  State: JobState
  Attempt: number
  MaxAttempts: number
  AttemptedAt: string | null
  AttemptedBy: string[]
  CreatedAt: string
  FinalizedAt: string | null
  ScheduledAt: string
  Priority: number
  Args: unknown
  Errors: Array<{
    at: string
    attempt: number
    error: string
    trace: string
  }>
  Kind: string
  Metadata: unknown
  Queue: string
  Tags: string[]
  UniqueKey: string
  UniqueStates: JobState[]
  CanRun: boolean
  CanRestart: boolean
  CanCancel: boolean
  CanDelete: boolean
}

export type JobPagination = {
  Page: number
  PageSize: number
  TotalCount: number
  TotalPages: number
}

export type JobFilters = {
  State: string
  Search: string
}

export type JobStats = {
  Total: number
  ByState: Record<JobState, number>
}

const dateFormatter = new Intl.DateTimeFormat('en', {
  day: 'numeric',
  month: 'short',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
})

export function formatJobDate(value: string | null) {
  if (!value) {
    return 'Not yet'
  }

  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? value : dateFormatter.format(date)
}

export function jobStateClass(state: JobState) {
  switch (state) {
    case 'completed':
      return 'border-emerald-500/30 bg-emerald-500/10 text-emerald-400'
    case 'running':
      return 'border-sky-500/30 bg-sky-500/10 text-sky-400'
    case 'available':
      return 'border-cyan-500/30 bg-cyan-500/10 text-cyan-400'
    case 'scheduled':
    case 'pending':
      return 'border-amber-500/30 bg-amber-500/10 text-amber-400'
    case 'retryable':
      return 'border-orange-500/30 bg-orange-500/10 text-orange-400'
    case 'discarded':
    case 'cancelled':
      return 'border-red-500/30 bg-red-500/10 text-red-400'
  }
}

export function jobStateLabel(state: JobState) {
  return state.charAt(0).toUpperCase() + state.slice(1)
}
