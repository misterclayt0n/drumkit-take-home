export default function Header() {
  return (
    <header className="sticky top-0 z-40 border-b border-[var(--dk-line)] bg-white/80 backdrop-blur-xl">
      <div className="page-wrap flex h-14 items-center justify-between">
        <a href="#top" className="flex items-center gap-2 no-underline">
          <img src="/drumkit_logo.svg" alt="Drumkit" className="h-6" />
          <span className="rounded-md bg-[var(--dk-surface-raised)] px-1.5 py-0.5 text-[11px] font-medium text-[var(--dk-ink-soft)]">
            Take-home
          </span>
        </a>

        <nav className="flex items-center gap-1">
          <a
            href="#loads"
            className="rounded-lg px-3 py-1.5 text-sm text-[var(--dk-ink-soft)] transition-colors hover:bg-[var(--dk-surface-raised)] hover:text-[var(--dk-ink)]"
          >
            Loads
          </a>
          <a
            href="#create"
            className="rounded-lg px-3 py-1.5 text-sm text-[var(--dk-ink-soft)] transition-colors hover:bg-[var(--dk-surface-raised)] hover:text-[var(--dk-ink)]"
          >
            Create
          </a>
        </nav>
      </div>
    </header>
  )
}
