import { useForm } from '@inertiajs/react'
import type { FormEvent } from 'react'

import Layout from '../../Layouts/Layout'
import { routes } from '../../routes'

type ConfirmEmailProps = {
  errors?: Record<string, string>
}

export default function ConfirmEmail({ errors = {} }: ConfirmEmailProps) {
  const form = useForm({ code: '' })

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    form.post(routes.subscriberConfirmationCreate())
  }

  return (
    <Layout>
      <section className="w-full max-w-md border border-[#2f3a37] bg-[#101414]/90 shadow-sm shadow-black/40">
        <div className="p-6 pb-0">
          <h1 className="text-xl font-semibold text-[#f2ead8]">Confirm Your Email</h1>
          <p className="mt-1 text-sm text-[#8f8a7d]">
            Enter the 6-digit confirmation code from your subscription email.
          </p>
        </div>
        <div className="p-6">
          <form onSubmit={submit} className="space-y-5">
            <div className="space-y-1">
              <label className="text-sm font-medium text-[#c7c0ad]" htmlFor="code">
                Confirmation code
              </label>
              <input
                id="code"
                value={form.data.code}
                onChange={(event) => form.setData('code', event.target.value.replace(/\D/g, ''))}
                className="flex h-11 w-full border border-[#2f3a37] bg-[#090c0d] px-3 py-1 text-center font-mono text-lg tracking-[0.35em] text-[#e4dfd2] shadow-inner shadow-black/35 focus:border-[#8df7a4] focus:outline-none focus:ring-2 focus:ring-[#8df7a4]/20"
                autoComplete="one-time-code"
                inputMode="numeric"
                maxLength={6}
                pattern="[0-9]{6}"
                required
                autoFocus
              />
              {errors.code && (
                <p className="text-sm font-medium text-[#ff875f]">{errors.code}</p>
              )}
            </div>
            <button
              type="submit"
              disabled={form.processing}
              className="inline-flex w-full items-center justify-center bg-[#ff6b1a] px-4 py-2 text-sm font-medium text-[#130f0b] shadow-sm shadow-black/40 hover:bg-[#ff8748] disabled:opacity-60"
            >
              {form.processing ? 'Confirming...' : 'Confirm email'}
            </button>
          </form>
        </div>
      </section>
    </Layout>
  )
}
