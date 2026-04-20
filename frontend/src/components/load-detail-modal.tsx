import { useEffect, useRef } from 'react'
import { ArrowRight, X } from 'lucide-react'
import type { Load } from '#/lib/types'

export function LoadDetailModal({
  load,
  onClose,
}: {
  load: Load | null
  onClose: () => void
}) {
  const overlayRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!load) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKey)
    document.body.style.overflow = 'hidden'
    return () => {
      document.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }, [load, onClose])

  if (!load) return null

  return (
    <div
      ref={overlayRef}
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/40 px-4 py-12 backdrop-blur-sm"
      onClick={(e) => {
        if (e.target === overlayRef.current) onClose()
      }}
    >
      <div className="relative w-full max-w-2xl rounded-2xl border border-[var(--dk-line)] bg-white shadow-[0_24px_64px_rgba(0,0,0,0.16)]">
        {/* Header */}
        <div className="flex items-start justify-between border-b border-[var(--dk-line)] px-6 py-5">
          <div>
            <div className="flex items-center gap-2.5">
              <h2 className="text-lg font-semibold text-[var(--dk-ink)]">
                {load.freightLoadID || 'Unnumbered load'}
              </h2>
              <StatusPill status={load.status} />
            </div>
            <p className="mt-1 font-mono text-xs text-[var(--dk-ink-soft)]">
              {load.externalTMSLoadID || '--'}
            </p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="flex size-8 items-center justify-center rounded-lg text-[var(--dk-ink-soft)] transition-colors hover:bg-[var(--dk-surface-raised)] hover:text-[var(--dk-ink)]"
          >
            <X className="size-4" />
          </button>
        </div>

        {/* Body */}
        <div className="divide-y divide-[var(--dk-line)] px-6">
          {/* Lane */}
          <div className="py-5">
            <SectionLabel>Lane</SectionLabel>
            <div className="mt-3 flex items-center gap-3 text-sm font-medium text-[var(--dk-ink)]">
              <span>{formatLocation(load.pickup.city, load.pickup.state)}</span>
              <ArrowRight className="size-3.5 text-[var(--dk-ink-soft)]" />
              <span>{formatLocation(load.consignee.city, load.consignee.state)}</span>
            </div>
            <div className="mt-2 grid gap-3 sm:grid-cols-2">
              <Detail label="Pickup" value={load.pickup.name || '--'} sub={formatDateTime(load.pickup.apptTime)} />
              <Detail label="Delivery" value={load.consignee.name || '--'} sub={formatDateTime(load.consignee.apptTime)} />
            </div>
          </div>

          {/* Customer */}
          <div className="py-5">
            <SectionLabel>Customer</SectionLabel>
            <div className="mt-3 grid gap-3 sm:grid-cols-2">
              <Detail label="Name" value={load.customer.name || '--'} sub={`ID: ${load.customer.externalTMSId || '--'}`} />
              <Detail label="Reference" value={load.customer.refNumber || '--'} />
              <Detail label="Bill to" value={load.billTo.name || '--'} sub={formatLocation(load.billTo.city, load.billTo.state)} />
              <Detail label="Contact" value={load.customer.contact || '--'} sub={load.customer.email || load.customer.phone || '--'} />
            </div>
          </div>

          {/* Operational */}
          <div className="py-5">
            <SectionLabel>Operational</SectionLabel>
            <div className="mt-3 grid gap-3 sm:grid-cols-3">
              <Detail label="PO numbers" value={load.poNums || '--'} />
              <Detail label="Route miles" value={load.routeMiles ? `${Math.round(load.routeMiles)} mi` : '--'} />
              <Detail label="Operator" value={load.operator || '--'} />
            </div>
          </div>

          {/* Freight */}
          <div className="py-5">
            <SectionLabel>Freight</SectionLabel>
            <div className="mt-3 grid grid-cols-3 gap-3 sm:grid-cols-6">
              <Metric label="Weight" value={formatWeight(load.totalWeight)} />
              <Metric label="Billable" value={formatWeight(load.billableWeight)} />
              <Metric label="In pallets" value={String(load.inPalletCount || 0)} />
              <Metric label="Out pallets" value={String(load.outPalletCount || 0)} />
              <Metric label="Commodities" value={String(load.numCommodities || 0)} />
              <Metric label="Profit" value={formatCurrency(load.rateData.netProfitUsd)} />
            </div>
          </div>

          {/* Specs */}
          <div className="py-5">
            <SectionLabel>Specifications</SectionLabel>
            <div className="mt-3 flex flex-wrap gap-1.5">
              {formatSpecFlags(load).map((flag) => (
                <span
                  key={flag}
                  className="rounded-md bg-[var(--dk-surface-raised)] px-2 py-1 text-xs font-medium text-[var(--dk-ink-soft)]"
                >
                  {flag}
                </span>
              ))}
              <span className="rounded-md bg-[var(--dk-surface-raised)] px-2 py-1 text-xs font-medium text-[var(--dk-ink-soft)]">
                {formatTemp(load.specifications.minTempFahrenheit, load.specifications.maxTempFahrenheit)}
              </span>
            </div>
          </div>

          {/* Carrier */}
          {load.carrier.name && (
            <div className="py-5">
              <SectionLabel>Carrier</SectionLabel>
              <div className="mt-3 grid gap-3 sm:grid-cols-2">
                <Detail label="Name" value={load.carrier.name} sub={`MC: ${load.carrier.mcNumber || '--'} / DOT: ${load.carrier.dotNumber || '--'}`} />
                <Detail label="Driver" value={load.carrier.firstDriverName || '--'} sub={load.carrier.firstDriverPhone || '--'} />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

function SectionLabel({ children }: { children: React.ReactNode }) {
  return (
    <p className="text-xs font-medium tracking-wide text-[var(--dk-ink-soft)] uppercase">
      {children}
    </p>
  )
}

function Detail({ label, value, sub }: { label: string; value: string; sub?: string }) {
  return (
    <div>
      <p className="text-xs text-[var(--dk-ink-soft)]">{label}</p>
      <p className="mt-0.5 text-sm font-medium text-[var(--dk-ink)]">{value}</p>
      {sub && <p className="mt-0.5 text-xs text-[var(--dk-ink-soft)]">{sub}</p>}
    </div>
  )
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-lg bg-[var(--dk-surface-raised)] px-2.5 py-2 text-center">
      <p className="text-[10px] text-[var(--dk-ink-soft)]">{label}</p>
      <p className="mt-0.5 font-mono text-sm font-semibold text-[var(--dk-ink)]">{value}</p>
    </div>
  )
}

function StatusPill({ status }: { status: string }) {
  const isActive = status.toLowerCase() === 'tendered'
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ${
        isActive
          ? 'bg-gradient-to-r from-[rgba(252,6,67,0.08)] to-[rgba(254,139,87,0.08)] text-[var(--dk-red)]'
          : 'bg-[var(--dk-surface-raised)] text-[var(--dk-ink-soft)]'
      }`}
    >
      <span className={`inline-block size-1.5 rounded-full ${isActive ? 'bg-[var(--dk-red)]' : 'bg-[var(--dk-ink-soft)]/40'}`} />
      {status || 'Unknown'}
    </span>
  )
}

function formatLocation(city?: string, state?: string) {
  const parts = [city, state].filter(Boolean)
  return parts.length > 0 ? parts.join(', ') : '--'
}

function formatDateTime(value: string) {
  if (!value) return 'TBD'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return new Intl.DateTimeFormat('en-US', { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' }).format(d)
}

function formatWeight(v: number) {
  return v ? `${Math.round(v).toLocaleString()} lb` : '--'
}

function formatCurrency(v: number) {
  if (!v) return '$0'
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD', maximumFractionDigits: 0 }).format(v)
}

function formatTemp(min: number, max: number) {
  if (!min && !max) return 'No temp spec'
  if (min === max) return `${min}°F`
  return `${min}°F – ${max}°F`
}

function formatSpecFlags(load: Load) {
  const flags = [
    load.specifications.hazmat ? 'Hazmat' : '',
    load.specifications.liftgatePickup || load.specifications.liftgateDelivery ? 'Liftgate' : '',
    load.specifications.insidePickup || load.specifications.insideDelivery ? 'Inside service' : '',
    load.specifications.oversized ? 'Oversized' : '',
    load.specifications.tarps ? 'Tarps' : '',
    load.specifications.straps ? 'Straps' : '',
    load.specifications.seal ? 'Seal' : '',
    load.specifications.customBonded ? 'Custom bonded' : '',
  ].filter(Boolean)
  return flags.length > 0 ? flags : ['No special handling']
}
