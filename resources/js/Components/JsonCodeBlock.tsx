import { useMemo, useState } from 'react'
import { Check, Copy } from 'lucide-react'

import { Button } from '@/components/ui/button'

type JsonCodeBlockProps = {
  value: unknown
  label: string
}

const tokenPattern = /("(?:\\u[\da-fA-F]{4}|\\[^u]|[^\\"])*"(?:\s*:)?|\b(?:true|false|null)\b|-?\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)/g

function tokenClass(token: string) {
  if (token.startsWith('"')) {
    return token.trimEnd().endsWith(':') ? 'text-sky-300' : 'text-emerald-300'
  }
  if (token === 'true' || token === 'false') {
    return 'text-violet-300'
  }
  if (token === 'null') {
    return 'text-muted-foreground'
  }
  if (/^-?\d/.test(token)) {
    return 'text-amber-300'
  }
  return undefined
}

export function JsonCodeBlock({ value, label }: JsonCodeBlockProps) {
  const [copied, setCopied] = useState(false)
  const formatted = useMemo(() => JSON.stringify(value, null, 2) ?? String(value), [value])
  const tokens = useMemo(() => formatted.split(tokenPattern), [formatted])

  async function copyJson() {
    try {
      await navigator.clipboard.writeText(formatted)
      setCopied(true)
      window.setTimeout(() => setCopied(false), 1600)
    } catch {
      setCopied(false)
    }
  }

  return (
    <div className="overflow-hidden border bg-[#0b0f12]">
      <div className="flex items-center justify-between border-b border-white/10 px-4 py-2">
        <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
          {label}
        </span>
        <Button variant="ghost" size="sm" onClick={copyJson} className="h-7 gap-2 px-2 text-xs">
          {copied ? <Check className="size-3.5" /> : <Copy className="size-3.5" />}
          {copied ? 'Copied' : 'Copy'}
        </Button>
      </div>
      <pre className="max-h-[32rem] overflow-auto p-4 font-mono text-xs leading-6 text-foreground sm:text-sm">
        <code>
          {tokens.map((token, index) => (
            <span key={`${index}-${token.slice(0, 12)}`} className={tokenClass(token)}>
              {token}
            </span>
          ))}
        </code>
      </pre>
    </div>
  )
}
