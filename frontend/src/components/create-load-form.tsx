import { Check, ChevronDown, Loader2 } from 'lucide-react'
import { useEffect, useRef, useState } from 'react'
import { Button } from '#/components/ui/button'
import { Input } from '#/components/ui/input'
import { Label } from '#/components/ui/label'
import { Textarea } from '#/components/ui/textarea'
import type { CreateLoadResponse, Load } from '#/lib/types'

type Primitive = string | number | boolean
type SectionKey = 'customer' | 'billTo' | 'pickup' | 'consignee' | 'carrier' | 'rateData' | 'specifications'
type FieldType = 'text' | 'number' | 'textarea' | 'checkbox'

type FieldConfig<T> = {
  key: Extract<keyof T, string>
  label: string
  type?: FieldType
  placeholder?: string
  step?: string
  mono?: boolean
  required?: boolean
}

// ---------- field configs (unchanged) ----------

const customerFields: FieldConfig<Load['customer']>[] = [
  { key: 'externalTMSId', label: 'External TMS ID', required: true },
  { key: 'name', label: 'Name' },
  { key: 'addressLine1', label: 'Address line 1' },
  { key: 'addressLine2', label: 'Address line 2' },
  { key: 'city', label: 'City' },
  { key: 'state', label: 'State' },
  { key: 'zipcode', label: 'Zipcode' },
  { key: 'country', label: 'Country' },
  { key: 'contact', label: 'Contact' },
  { key: 'phone', label: 'Phone' },
  { key: 'email', label: 'Email' },
  { key: 'refNumber', label: 'Reference number' },
]

const billToFields: FieldConfig<Load['billTo']>[] = [
  { key: 'externalTMSId', label: 'External TMS ID' },
  { key: 'name', label: 'Name' },
  { key: 'addressLine1', label: 'Address line 1' },
  { key: 'addressLine2', label: 'Address line 2' },
  { key: 'city', label: 'City' },
  { key: 'state', label: 'State' },
  { key: 'zipcode', label: 'Zipcode' },
  { key: 'country', label: 'Country' },
  { key: 'contact', label: 'Contact' },
  { key: 'phone', label: 'Phone' },
  { key: 'email', label: 'Email' },
]

const pickupFields: FieldConfig<Load['pickup']>[] = [
  { key: 'externalTMSId', label: 'External TMS ID', required: true },
  { key: 'warehouseId', label: 'Warehouse ID', required: true },
  { key: 'name', label: 'Name' },
  { key: 'contact', label: 'Contact' },
  { key: 'addressLine1', label: 'Address line 1' },
  { key: 'addressLine2', label: 'Address line 2' },
  { key: 'city', label: 'City' },
  { key: 'state', label: 'State' },
  { key: 'zipcode', label: 'Zipcode' },
  { key: 'country', label: 'Country' },
  { key: 'phone', label: 'Phone' },
  { key: 'email', label: 'Email' },
  { key: 'businessHours', label: 'Business hours' },
  { key: 'refNumber', label: 'Reference number' },
  { key: 'readyTime', label: 'Ready time', placeholder: 'RFC3339 timestamp', mono: true, required: true },
  { key: 'apptTime', label: 'Appointment time', placeholder: 'RFC3339 timestamp', mono: true, required: true },
  { key: 'apptNote', label: 'Appointment note', type: 'textarea' },
  { key: 'timezone', label: 'Timezone' },
]

const consigneeFields: FieldConfig<Load['consignee']>[] = [
  { key: 'externalTMSId', label: 'External TMS ID', required: true },
  { key: 'warehouseId', label: 'Warehouse ID', required: true },
  { key: 'name', label: 'Name' },
  { key: 'contact', label: 'Contact' },
  { key: 'addressLine1', label: 'Address line 1' },
  { key: 'addressLine2', label: 'Address line 2' },
  { key: 'city', label: 'City' },
  { key: 'state', label: 'State' },
  { key: 'zipcode', label: 'Zipcode' },
  { key: 'country', label: 'Country' },
  { key: 'phone', label: 'Phone' },
  { key: 'email', label: 'Email' },
  { key: 'businessHours', label: 'Business hours' },
  { key: 'refNumber', label: 'Reference number' },
  { key: 'mustDeliver', label: 'Must deliver', required: true },
  { key: 'apptTime', label: 'Appointment time', placeholder: 'RFC3339 timestamp', mono: true, required: true },
  { key: 'apptNote', label: 'Appointment note', type: 'textarea' },
  { key: 'timezone', label: 'Timezone' },
]

const carrierFields: FieldConfig<Load['carrier']>[] = [
  { key: 'externalTMSId', label: 'External TMS ID' },
  { key: 'name', label: 'Name' },
  { key: 'mcNumber', label: 'MC number' },
  { key: 'dotNumber', label: 'DOT number' },
  { key: 'phone', label: 'Phone' },
  { key: 'email', label: 'Email' },
  { key: 'dispatcher', label: 'Dispatcher' },
  { key: 'scac', label: 'SCAC' },
  { key: 'sealNumber', label: 'Seal number' },
  { key: 'firstDriverName', label: 'First driver name' },
  { key: 'firstDriverPhone', label: 'First driver phone' },
  { key: 'secondDriverName', label: 'Second driver name' },
  { key: 'secondDriverPhone', label: 'Second driver phone' },
  { key: 'dispatchCity', label: 'Dispatch city' },
  { key: 'dispatchState', label: 'Dispatch state' },
  { key: 'externalTMSTruckId', label: 'External TMS truck ID' },
  { key: 'externalTMSTrailerId', label: 'External TMS trailer ID' },
  { key: 'confirmationSentTime', label: 'Confirmation sent', placeholder: 'RFC3339', mono: true },
  { key: 'confirmationReceivedTime', label: 'Confirmation received', placeholder: 'RFC3339', mono: true },
  { key: 'dispatchedTime', label: 'Dispatched time', placeholder: 'RFC3339', mono: true },
  { key: 'expectedPickupTime', label: 'Expected pickup', placeholder: 'RFC3339', mono: true },
  { key: 'pickupStart', label: 'Pickup start', placeholder: 'RFC3339', mono: true },
  { key: 'pickupEnd', label: 'Pickup end', placeholder: 'RFC3339', mono: true },
  { key: 'expectedDeliveryTime', label: 'Expected delivery', placeholder: 'RFC3339', mono: true },
  { key: 'deliveryStart', label: 'Delivery start', placeholder: 'RFC3339', mono: true },
  { key: 'deliveryEnd', label: 'Delivery end', placeholder: 'RFC3339', mono: true },
  { key: 'signedBy', label: 'Signed by' },
]

const rateDataFields: FieldConfig<Load['rateData']>[] = [
  { key: 'customerRateType', label: 'Customer rate type' },
  { key: 'customerNumHours', label: 'Customer hours', type: 'number', step: '0.01' },
  { key: 'customerLhRateUsd', label: 'Customer LH rate', type: 'number', step: '0.01' },
  { key: 'fscPercent', label: 'FSC %', type: 'number', step: '0.01' },
  { key: 'fscPerMile', label: 'FSC / mile', type: 'number', step: '0.01' },
  { key: 'carrierRateType', label: 'Carrier rate type' },
  { key: 'carrierNumHours', label: 'Carrier hours', type: 'number', step: '0.01' },
  { key: 'carrierLhRateUsd', label: 'Carrier LH rate', type: 'number', step: '0.01' },
  { key: 'carrierMaxRate', label: 'Carrier max rate', type: 'number', step: '0.01' },
  { key: 'netProfitUsd', label: 'Net profit', type: 'number', step: '0.01' },
  { key: 'profitPercent', label: 'Profit %', type: 'number', step: '0.01' },
]

const specificationsFields: FieldConfig<Load['specifications']>[] = [
  { key: 'minTempFahrenheit', label: 'Min temp (F)', type: 'number', step: '0.01' },
  { key: 'maxTempFahrenheit', label: 'Max temp (F)', type: 'number', step: '0.01' },
  { key: 'liftgatePickup', label: 'Liftgate pickup', type: 'checkbox' },
  { key: 'liftgateDelivery', label: 'Liftgate delivery', type: 'checkbox' },
  { key: 'insidePickup', label: 'Inside pickup', type: 'checkbox' },
  { key: 'insideDelivery', label: 'Inside delivery', type: 'checkbox' },
  { key: 'tarps', label: 'Tarps', type: 'checkbox' },
  { key: 'oversized', label: 'Oversized', type: 'checkbox' },
  { key: 'hazmat', label: 'Hazmat', type: 'checkbox' },
  { key: 'straps', label: 'Straps', type: 'checkbox' },
  { key: 'permits', label: 'Permits', type: 'checkbox' },
  { key: 'escorts', label: 'Escorts', type: 'checkbox' },
  { key: 'seal', label: 'Seal', type: 'checkbox' },
  { key: 'customBonded', label: 'Custom bonded', type: 'checkbox' },
  { key: 'labor', label: 'Labor', type: 'checkbox' },
]

// ---------- component ----------

export function CreateLoadForm({
  onSubmit,
  isPending,
  success,
  error,
  onClose,
}: {
  onSubmit: (payload: Load) => Promise<void> | void
  isPending: boolean
  success: CreateLoadResponse | null
  error: string | null
  onClose: () => void
}) {
  const [draft, setDraft] = useState<Load>(() => makeDefaultLoad())
  const overlayRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    document.addEventListener('keydown', onKey)
    document.body.style.overflow = 'hidden'
    return () => {
      document.removeEventListener('keydown', onKey)
      document.body.style.overflow = ''
    }
  }, [onClose])

  function updateTopLevel<K extends keyof Load>(key: K, value: Load[K]) {
    setDraft((c) => ({ ...c, [key]: value }))
  }

  function updateSection(section: SectionKey, field: string, value: Primitive) {
    setDraft((c) => ({ ...c, [section]: { ...(c[section] as Record<string, Primitive>), [field]: value } }))
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    await onSubmit(draft)
  }

  return (
    <div
      ref={overlayRef}
      className="fixed inset-0 z-50 flex items-start justify-center overflow-y-auto bg-black/40 px-4 py-8 backdrop-blur-sm"
      onClick={(e) => { if (e.target === overlayRef.current) onClose() }}
    >
      <div className="relative w-full max-w-3xl rounded-2xl border border-[var(--dk-line)] bg-white shadow-[0_24px_64px_rgba(0,0,0,0.16)]">
        {/* Fixed header */}
        <div className="sticky top-0 z-10 flex items-center justify-between rounded-t-2xl border-b border-[var(--dk-line)] bg-white px-6 py-4">
          <div>
            <h2 className="text-lg font-semibold text-[var(--dk-ink)]">Create load</h2>
            <p className="mt-0.5 text-sm text-[var(--dk-ink-soft)]">
              Full Drumkit load schema. Expand sections as needed.
            </p>
            <p className="mt-1 text-xs text-[var(--dk-ink-soft)]">
              <span className="font-semibold text-[var(--dk-red)]">*</span> required for the current Turvo adapter. For pickup and consignee, provide one of External TMS ID or Warehouse ID, and one of the marked time fields in each section.
            </p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="flex size-8 items-center justify-center rounded-lg text-[var(--dk-ink-soft)] transition-colors hover:bg-[var(--dk-surface-raised)] hover:text-[var(--dk-ink)]"
          >
            <span className="text-lg leading-none">&times;</span>
          </button>
        </div>

        <div className="px-6 py-5">
          {/* Alerts */}
          {success && (
            <div className="mb-5 flex items-start gap-2.5 rounded-lg border border-emerald-200 bg-emerald-50 px-3 py-2.5 text-sm">
              <Check className="mt-0.5 size-4 shrink-0 text-emerald-600" />
              <div>
                <p className="font-medium text-emerald-900">Load created</p>
                <p className="mt-0.5 font-mono text-xs text-emerald-700">
                  Turvo ID: {success.id} &middot; {formatTimestamp(success.createdAt)}
                </p>
              </div>
            </div>
          )}
          {error && (
            <div className="mb-5 rounded-lg border border-red-200 bg-red-50 px-3 py-2.5 text-sm">
              <p className="font-medium text-red-900">Creation failed</p>
              <p className="mt-0.5 text-red-700">{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit}>
            {/* Top-level fields - always visible */}
            <AccordionSection title="Identity & freight" defaultOpen count={11}>
              <div className="grid gap-3 sm:grid-cols-3">
                <PrimitiveField label="External TMS Load ID" type="text" value={draft.externalTMSLoadID} onChange={(v) => updateTopLevel('externalTMSLoadID', v as string)} />
                <PrimitiveField label="Freight Load ID" type="text" value={draft.freightLoadID} onChange={(v) => updateTopLevel('freightLoadID', v as string)} />
                <Field label="Status">
                  <select
                    value={draft.status}
                    onChange={(e) => updateTopLevel('status', e.target.value)}
                    className="flex h-9 w-full rounded-lg border border-[var(--dk-line-strong)] bg-white px-3 text-sm text-[var(--dk-ink)] outline-none transition focus:border-[var(--dk-red)] focus:ring-2 focus:ring-[var(--ring)]"
                  >
                    <option value="Tendered">Tendered</option>
                    <option value="Covered">Covered</option>
                  </select>
                </Field>
                <PrimitiveField label="PO numbers" type="text" value={draft.poNums} onChange={(v) => updateTopLevel('poNums', v as string)} />
                <PrimitiveField label="Operator" type="text" value={draft.operator} onChange={(v) => updateTopLevel('operator', v as string)} />
                <PrimitiveField label="Route miles" type="number" step="0.01" value={draft.routeMiles} onChange={(v) => updateTopLevel('routeMiles', v as number)} />
                <PrimitiveField label="Total weight" type="number" step="0.01" value={draft.totalWeight} onChange={(v) => updateTopLevel('totalWeight', v as number)} />
                <PrimitiveField label="Billable weight" type="number" step="0.01" value={draft.billableWeight} onChange={(v) => updateTopLevel('billableWeight', v as number)} />
                <PrimitiveField label="In pallets" type="number" value={draft.inPalletCount} onChange={(v) => updateTopLevel('inPalletCount', v as number)} />
                <PrimitiveField label="Out pallets" type="number" value={draft.outPalletCount} onChange={(v) => updateTopLevel('outPalletCount', v as number)} />
                <PrimitiveField label="Commodities" type="number" value={draft.numCommodities} onChange={(v) => updateTopLevel('numCommodities', v as number)} />
              </div>
            </AccordionSection>

            <ObjectAccordion title="Customer" defaultOpen count={customerFields.length} value={draft.customer} fields={customerFields} onChange={(f, v) => updateSection('customer', f, v)} />
            <ObjectAccordion title="Pickup" defaultOpen count={pickupFields.length} helperText="One of External TMS ID or Warehouse ID is required. One of Ready time or Appointment time is required." value={draft.pickup} fields={pickupFields} onChange={(f, v) => updateSection('pickup', f, v)} />
            <ObjectAccordion title="Consignee (delivery)" defaultOpen count={consigneeFields.length} helperText="One of External TMS ID or Warehouse ID is required. One of Must deliver or Appointment time is required." value={draft.consignee} fields={consigneeFields} onChange={(f, v) => updateSection('consignee', f, v)} />
            <ObjectAccordion title="Bill to" count={billToFields.length} value={draft.billTo} fields={billToFields} onChange={(f, v) => updateSection('billTo', f, v)} />
            <ObjectAccordion title="Carrier" count={carrierFields.length} value={draft.carrier} fields={carrierFields} onChange={(f, v) => updateSection('carrier', f, v)} />
            <ObjectAccordion title="Rate data" count={rateDataFields.length} value={draft.rateData} fields={rateDataFields} onChange={(f, v) => updateSection('rateData', f, v)} />
            <ObjectAccordion title="Specifications" count={specificationsFields.length} value={draft.specifications} fields={specificationsFields} onChange={(f, v) => updateSection('specifications', f, v)} />

            {/* Sticky footer */}
            <div className="sticky bottom-0 -mx-6 mt-6 flex gap-2 border-t border-[var(--dk-line)] bg-white px-6 py-4">
              <Button
                type="button"
                variant="outline"
                onClick={() => setDraft(makeDefaultLoad())}
                className="h-9 rounded-lg border-[var(--dk-line-strong)] px-3 text-sm text-[var(--dk-ink)]"
              >
                Reset defaults
              </Button>
              <Button
                type="submit"
                disabled={isPending}
                className="h-9 flex-1 rounded-lg bg-[var(--dk-ink)] px-4 text-sm font-medium text-white hover:bg-[var(--dk-ink)]/90"
              >
                {isPending ? (
                  <><Loader2 className="size-3.5 animate-spin" /> Creating...</>
                ) : (
                  'Create load'
                )}
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  )
}

// ---------- accordion ----------

function AccordionSection({
  title,
  defaultOpen = false,
  count,
  required = false,
  children,
}: {
  title: string
  defaultOpen?: boolean
  count: number
  required?: boolean
  children: React.ReactNode
}) {
  const [open, setOpen] = useState(defaultOpen)

  return (
    <div className="border-b border-[var(--dk-line)]">
      <button
        type="button"
        onClick={() => setOpen(!open)}
        className="flex w-full items-center justify-between py-3 text-left"
      >
        <div className="flex items-center gap-2">
          <span className="text-sm font-semibold text-[var(--dk-ink)]">
            {title}
            {required && <span className="ml-1 text-[var(--dk-red)]">*</span>}
          </span>
          <span className="rounded-md bg-[var(--dk-surface-raised)] px-1.5 py-0.5 text-[10px] font-medium text-[var(--dk-ink-soft)]">
            {count}
          </span>
        </div>
        <ChevronDown className={`size-4 text-[var(--dk-ink-soft)] transition-transform ${open ? 'rotate-180' : ''}`} />
      </button>
      {open && <div className="pb-4">{children}</div>}
    </div>
  )
}

function ObjectAccordion({
  title,
  defaultOpen = false,
  count,
  helperText,
  value,
  fields,
  onChange,
}: {
  title: string
  defaultOpen?: boolean
  count: number
  helperText?: string
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  fields: FieldConfig<any>[]
  onChange: (field: string, value: Primitive) => void
}) {
  return (
    <AccordionSection title={title} defaultOpen={defaultOpen} count={count} required={fields.some((field) => field.required)}>
      {helperText && <p className="mb-3 text-xs text-[var(--dk-ink-soft)]">{helperText}</p>}
      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
        {fields.map((field) => (
          <PrimitiveField
            key={field.key}
            label={field.label}
            type={field.type ?? 'text'}
            placeholder={field.placeholder}
            step={field.step}
            mono={field.mono}
            value={value[field.key] as Primitive}
            required={field.required}
            onChange={(v) => onChange(field.key, v)}
          />
        ))}
      </div>
    </AccordionSection>
  )
}

// ---------- field primitives ----------

function Field({ label, required, children }: { label: string; required?: boolean; children: React.ReactNode }) {
  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-medium text-[var(--dk-ink)]">
        {label}
        {required && <span className="ml-1 text-[var(--dk-red)]">*</span>}
      </Label>
      {children}
    </div>
  )
}

function PrimitiveField({
  label, type = 'text', value, onChange, placeholder, step, mono, required,
}: {
  label: string; type?: FieldType; value: Primitive; onChange: (v: Primitive) => void
  placeholder?: string; step?: string; mono?: boolean; required?: boolean
}) {
  if (type === 'checkbox') {
    return (
      <label className="flex h-9 items-center gap-2 rounded-lg border border-[var(--dk-line-strong)] bg-white px-3 text-sm text-[var(--dk-ink)]">
        <input type="checkbox" checked={Boolean(value)} onChange={(e) => onChange(e.target.checked)} className="size-4 rounded border-[var(--dk-line-strong)]" />
        <span>
          {label}
          {required && <span className="ml-1 text-[var(--dk-red)]">*</span>}
        </span>
      </label>
    )
  }
  if (type === 'textarea') {
    return (
      <Field label={label} required={required}>
        <Textarea value={String(value)} onChange={(e) => onChange(e.target.value)} className="min-h-16 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm" placeholder={placeholder} />
      </Field>
    )
  }
  if (type === 'number') {
    return (
      <Field label={label} required={required}>
        <Input type="number" step={step} value={String(value)} onChange={(e) => onChange(parseNumber(e.target.value))} className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm" placeholder={placeholder} />
      </Field>
    )
  }
  return (
    <Field label={label} required={required}>
      <Input value={String(value)} onChange={(e) => onChange(e.target.value)} className={`h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm ${mono ? 'font-mono text-xs' : ''}`} placeholder={placeholder} />
    </Field>
  )
}

function parseNumber(v: string) {
  const t = v.trim()
  if (!t) return 0
  const n = Number(t)
  return Number.isFinite(n) ? n : 0
}

// ---------- defaults ----------

function makeDefaultLoad(): Load {
  const pickupTime = futureIsoString(1)
  const deliveryTime = futureIsoString(2)
  return {
    externalTMSLoadID: 'local-test-load-001',
    freightLoadID: 'frontend-test-001',
    status: 'Tendered',
    customer: { externalTMSId: '834045', name: '37th St Bakery', addressLine1: '', addressLine2: '', city: '', state: '', zipcode: '', country: '', contact: '', phone: '', email: '', refNumber: 'local-ref-001' },
    billTo: { externalTMSId: '', name: '', addressLine1: '', addressLine2: '', city: '', state: '', zipcode: '', country: '', contact: '', phone: '', email: '' },
    pickup: { externalTMSId: '525513', name: '1611 CGT- Rockingham (North Carolina)', addressLine1: '', addressLine2: '', city: 'ROCKINGHAM', state: 'NC', zipcode: '', country: 'US', contact: '', phone: '', email: '', businessHours: '', refNumber: '', readyTime: pickupTime, apptTime: pickupTime, apptNote: 'Pickup note from the frontend', timezone: 'America/New_York', warehouseId: '525513' },
    consignee: { externalTMSId: '525541', name: 'AAA TRANS WORLD EXPRESS, INC', addressLine1: '', addressLine2: '', city: 'JAMAICA', state: 'NY', zipcode: '', country: 'US', contact: '', phone: '', email: '', businessHours: '', refNumber: '', mustDeliver: deliveryTime, apptTime: deliveryTime, apptNote: 'Delivery note from the frontend', timezone: 'America/New_York', warehouseId: '525541' },
    carrier: { mcNumber: '', dotNumber: '', name: '', phone: '', dispatcher: '', sealNumber: '', scac: '', firstDriverName: '', firstDriverPhone: '', secondDriverName: '', secondDriverPhone: '', email: '', dispatchCity: '', dispatchState: '', externalTMSTruckId: '', externalTMSTrailerId: '', confirmationSentTime: '', confirmationReceivedTime: '', dispatchedTime: '', expectedPickupTime: '', pickupStart: '', pickupEnd: '', expectedDeliveryTime: '', deliveryStart: '', deliveryEnd: '', signedBy: '', externalTMSId: '' },
    rateData: { customerRateType: '', customerNumHours: 0, customerLhRateUsd: 0, fscPercent: 0, fscPerMile: 0, carrierRateType: '', carrierNumHours: 0, carrierLhRateUsd: 0, carrierMaxRate: 0, netProfitUsd: 0, profitPercent: 0 },
    specifications: { minTempFahrenheit: 0, maxTempFahrenheit: 0, liftgatePickup: false, liftgateDelivery: false, insidePickup: false, insideDelivery: false, tarps: false, oversized: false, hazmat: false, straps: false, permits: false, escorts: false, seal: false, customBonded: false, labor: false },
    inPalletCount: 0, outPalletCount: 0, numCommodities: 0, totalWeight: 0, billableWeight: 0, poNums: 'LOCAL-PO-001', operator: '', routeMiles: 0,
  }
}

function futureIsoString(days: number) {
  const d = new Date()
  d.setUTCDate(d.getUTCDate() + days)
  d.setUTCHours(14, 0, 0, 0)
  return d.toISOString().replace('.000', '')
}

function formatTimestamp(v: string) {
  const d = new Date(v)
  if (Number.isNaN(d.getTime())) return v
  return new Intl.DateTimeFormat('en-US', { dateStyle: 'medium', timeStyle: 'short' }).format(d)
}
