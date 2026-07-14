import { useCallback, useEffect, useRef, useState } from 'react'
import { router } from '@inertiajs/react'
import { Circle, LoaderCircle } from 'lucide-react'

import { Button } from '@/components/ui/button'

type PollingControlProps = {
  only: string[]
  intervalMs?: number
  suspended?: boolean
}

export function PollingControl({
  only,
  intervalMs = 5_000,
  suspended = false,
}: PollingControlProps) {
  const [enabled, setEnabled] = useState(false)
  const refreshInFlight = useRef(false)

  const poll = useCallback(() => {
    if (refreshInFlight.current) {
      return
    }

    refreshInFlight.current = true
    router.reload({
      only,
      onFinish: () => {
        refreshInFlight.current = false
      },
    })
  }, [only])

  useEffect(() => {
    if (!enabled || suspended) {
      return
    }

    poll()
    const intervalID = window.setInterval(poll, intervalMs)
    return () => window.clearInterval(intervalID)
  }, [enabled, intervalMs, poll, suspended])

  const isPolling = enabled && !suspended

  return (
    <div className="flex items-center gap-3" aria-live="polite">
      <div className="flex items-center gap-2 text-xs font-medium uppercase tracking-[0.16em] text-muted-foreground">
        {isPolling ? (
          <LoaderCircle className="size-3.5 animate-spin" />
        ) : (
          <Circle className="size-3.5" />
        )}
        <span>{isPolling ? 'Polling' : 'Paused'}</span>
      </div>
      <Button
        type="button"
        variant="secondary"
        size="sm"
        onClick={() => setEnabled((current) => !current)}
        disabled={suspended}
      >
        {enabled ? 'Pause polling' : 'Resume polling'}
      </Button>
    </div>
  )
}
