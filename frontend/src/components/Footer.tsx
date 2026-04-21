import { getApiBaseUrl } from '#/lib/api'

export default function Footer() {
  return (
    <footer className="border-t border-[var(--dk-line)] bg-white/50">
      <div className="page-wrap flex flex-col gap-3 py-6 text-xs text-[var(--dk-ink-soft)] sm:flex-row sm:items-center sm:justify-between">
        <p className="m-0">
          Drumkit take-home &middot; TanStack Start + Go backend
        </p>
      </div>
    </footer>
  )
}
