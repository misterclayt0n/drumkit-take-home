import { Check, Loader2 } from 'lucide-react'
import { useState } from 'react'
import { Button } from '#/components/ui/button'
import { Input } from '#/components/ui/input'
import { Label } from '#/components/ui/label'
import { Textarea } from '#/components/ui/textarea'
import type { CreateLoadResponse, Load } from '#/lib/types'

export function CreateLoadForm({
  onSubmit,
  isPending,
  success,
  error,
}: {
  onSubmit: (payload: Load) => Promise<void> | void
  isPending: boolean
  success: CreateLoadResponse | null
  error: string | null
}) {
  const [draft, setDraft] = useState<Load>(() => makeDefaultLoad())

  function updateTopLevel<K extends keyof Load>(key: K, value: Load[K]) {
    setDraft((c) => ({ ...c, [key]: value }))
  }

  function updateSection(
    section: 'customer' | 'pickup' | 'consignee',
    field: string,
    value: string,
  ) {
    setDraft((c) => ({
      ...c,
      [section]: { ...c[section], [field]: value },
    }))
  }

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    await onSubmit(draft)
  }

  return (
    <div className="overflow-hidden rounded-xl border border-[var(--dk-line)] bg-white lg:sticky lg:top-20">
      {/* Header */}
      <div className="border-b border-[var(--dk-line)] px-5 py-4">
        <h2 className="text-base font-semibold text-[var(--dk-ink)]">
          Create load
        </h2>
        <p className="mt-0.5 text-sm text-[var(--dk-ink-soft)]">
          Prefilled with sandbox defaults. Sends to Turvo via the Go backend.
        </p>
      </div>

      <div className="px-5 py-5">
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

        <form className="space-y-5" onSubmit={handleSubmit}>
          {/* Identity */}
          <FormSection title="Identity">
            <div className="grid gap-3 sm:grid-cols-2">
              <Field label="External TMS Load ID">
                <Input
                  value={draft.externalTMSLoadID}
                  onChange={(e) => updateTopLevel('externalTMSLoadID', e.target.value)}
                  className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                />
              </Field>
              <Field label="Freight Load ID">
                <Input
                  value={draft.freightLoadID}
                  onChange={(e) => updateTopLevel('freightLoadID', e.target.value)}
                  className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                />
              </Field>
            </div>
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
          </FormSection>

          {/* Customer */}
          <FormSection title="Customer">
            <div className="grid gap-3 sm:grid-cols-2">
              <Field label="Customer ID">
                <Input
                  value={draft.customer.externalTMSId}
                  onChange={(e) => updateSection('customer', 'externalTMSId', e.target.value)}
                  className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                />
              </Field>
              <Field label="Name">
                <Input
                  value={draft.customer.name}
                  onChange={(e) => updateSection('customer', 'name', e.target.value)}
                  className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                />
              </Field>
            </div>
            <Field label="Reference number">
              <Input
                value={draft.customer.refNumber}
                onChange={(e) => updateSection('customer', 'refNumber', e.target.value)}
                className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                placeholder="Optional"
              />
            </Field>
          </FormSection>

          {/* Pickup */}
          <FormSection title="Pickup">
            <LocationFields
              value={draft.pickup}
              onChange={(f, v) => updateSection('pickup', f, v)}
              noteLabel="Pickup note"
            />
          </FormSection>

          {/* Delivery */}
          <FormSection title="Delivery">
            <LocationFields
              value={draft.consignee}
              onChange={(f, v) => updateSection('consignee', f, v)}
              noteLabel="Delivery note"
            />
          </FormSection>

          {/* Operational */}
          <FormSection title="Operational">
            <div className="grid gap-3 sm:grid-cols-2">
              <Field label="PO numbers">
                <Input
                  value={draft.poNums}
                  onChange={(e) => updateTopLevel('poNums', e.target.value)}
                  className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                />
              </Field>
              <Field label="Operator">
                <Input
                  value={draft.operator}
                  onChange={(e) => updateTopLevel('operator', e.target.value)}
                  className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
                  placeholder="Optional"
                />
              </Field>
            </div>
          </FormSection>

          {/* Actions */}
          <div className="flex gap-2 border-t border-[var(--dk-line)] pt-4">
            <Button
              type="button"
              variant="outline"
              onClick={() => setDraft(makeDefaultLoad())}
              className="h-9 rounded-lg border-[var(--dk-line-strong)] px-3 text-sm text-[var(--dk-ink)]"
            >
              Reset
            </Button>
            <Button
              type="submit"
              disabled={isPending}
              className="h-9 flex-1 rounded-lg bg-[var(--dk-ink)] px-4 text-sm font-medium text-white hover:bg-[var(--dk-ink)]/90"
            >
              {isPending ? (
                <>
                  <Loader2 className="size-3.5 animate-spin" />
                  Creating...
                </>
              ) : (
                'Create load'
              )}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}

function FormSection({
  title,
  children,
}: {
  title: string
  children: React.ReactNode
}) {
  return (
    <div className="space-y-3">
      <p className="text-xs font-medium tracking-wide text-[var(--dk-ink-soft)] uppercase">
        {title}
      </p>
      {children}
    </div>
  )
}

function Field({
  label,
  children,
}: {
  label: string
  children: React.ReactNode
}) {
  return (
    <div className="space-y-1.5">
      <Label className="text-xs font-medium text-[var(--dk-ink)]">{label}</Label>
      {children}
    </div>
  )
}

function LocationFields({
  value,
  onChange,
  noteLabel,
}: {
  value: Load['pickup'] | Load['consignee']
  onChange: (field: string, value: string) => void
  noteLabel: string
}) {
  return (
    <>
      <div className="grid gap-3 sm:grid-cols-2">
        <Field label="Location ID">
          <Input
            value={value.externalTMSId}
            onChange={(e) => onChange('externalTMSId', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />
        </Field>
        <Field label="Name">
          <Input
            value={value.name}
            onChange={(e) => onChange('name', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />
        </Field>
      </div>

      <div className="grid gap-3 grid-cols-3">
        <Field label="City">
          <Input
            value={value.city}
            onChange={(e) => onChange('city', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />
        </Field>
        <Field label="State">
          <Input
            value={value.state}
            onChange={(e) => onChange('state', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />
        </Field>
        <Field label="Country">
          <Input
            value={value.country}
            onChange={(e) => onChange('country', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />
        </Field>
      </div>

      <div className="grid gap-3 sm:grid-cols-2">
        <Field label="Appointment (RFC3339)">
          <Input
            value={value.apptTime}
            onChange={(e) => onChange('apptTime', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm font-mono text-xs"
          />
        </Field>
        <Field label="Timezone">
          <Input
            value={value.timezone}
            onChange={(e) => onChange('timezone', e.target.value)}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />
        </Field>
      </div>

      <Field label={noteLabel}>
        <Textarea
          value={value.apptNote}
          onChange={(e) => onChange('apptNote', e.target.value)}
          className="min-h-16 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          placeholder="Optional"
        />
      </Field>
    </>
  )
}

function makeDefaultLoad(): Load {
  const pickupTime = futureIsoString(1)
  const deliveryTime = futureIsoString(2)

  return {
    externalTMSLoadID: 'local-test-load-001',
    freightLoadID: 'frontend-test-001',
    status: 'Tendered',
    customer: {
      externalTMSId: '834045',
      name: '37th St Bakery',
      addressLine1: '',
      addressLine2: '',
      city: '',
      state: '',
      zipcode: '',
      country: '',
      contact: '',
      phone: '',
      email: '',
      refNumber: 'local-ref-001',
    },
    billTo: {
      externalTMSId: '',
      name: '',
      addressLine1: '',
      addressLine2: '',
      city: '',
      state: '',
      zipcode: '',
      country: '',
      contact: '',
      phone: '',
      email: '',
    },
    pickup: {
      externalTMSId: '525513',
      name: '1611 CGT- Rockingham (North Carolina)',
      addressLine1: '',
      addressLine2: '',
      city: 'ROCKINGHAM',
      state: 'NC',
      zipcode: '',
      country: 'US',
      contact: '',
      phone: '',
      email: '',
      businessHours: '',
      refNumber: '',
      readyTime: pickupTime,
      apptTime: pickupTime,
      apptNote: 'Pickup note from the frontend',
      timezone: 'America/New_York',
      warehouseId: '525513',
    },
    consignee: {
      externalTMSId: '525541',
      name: 'AAA TRANS WORLD EXPRESS, INC',
      addressLine1: '',
      addressLine2: '',
      city: 'JAMAICA',
      state: 'NY',
      zipcode: '',
      country: 'US',
      contact: '',
      phone: '',
      email: '',
      businessHours: '',
      refNumber: '',
      mustDeliver: deliveryTime,
      apptTime: deliveryTime,
      apptNote: 'Delivery note from the frontend',
      timezone: 'America/New_York',
      warehouseId: '525541',
    },
    carrier: emptyCarrier(),
    rateData: emptyRateData(),
    specifications: emptySpecifications(),
    inPalletCount: 0,
    outPalletCount: 0,
    numCommodities: 0,
    totalWeight: 0,
    billableWeight: 0,
    poNums: 'LOCAL-PO-001',
    operator: '',
    routeMiles: 0,
  }
}

function emptyCarrier(): Load['carrier'] {
  return {
    mcNumber: '', dotNumber: '', name: '', phone: '', dispatcher: '',
    sealNumber: '', scac: '', firstDriverName: '', firstDriverPhone: '',
    secondDriverName: '', secondDriverPhone: '', email: '', dispatchCity: '',
    dispatchState: '', externalTMSTruckId: '', externalTMSTrailerId: '',
    confirmationSentTime: '', confirmationReceivedTime: '', dispatchedTime: '',
    expectedPickupTime: '', pickupStart: '', pickupEnd: '',
    expectedDeliveryTime: '', deliveryStart: '', deliveryEnd: '',
    signedBy: '', externalTMSId: '',
  }
}

function emptyRateData(): Load['rateData'] {
  return {
    customerRateType: '', customerNumHours: 0, customerLhRateUsd: 0,
    fscPercent: 0, fscPerMile: 0, carrierRateType: '', carrierNumHours: 0,
    carrierLhRateUsd: 0, carrierMaxRate: 0, netProfitUsd: 0, profitPercent: 0,
  }
}

function emptySpecifications(): Load['specifications'] {
  return {
    minTempFahrenheit: 0, maxTempFahrenheit: 0, liftgatePickup: false,
    liftgateDelivery: false, insidePickup: false, insideDelivery: false,
    tarps: false, oversized: false, hazmat: false, straps: false,
    permits: false, escorts: false, seal: false, customBonded: false,
    labor: false,
  }
}

function futureIsoString(daysAhead: number) {
  const d = new Date()
  d.setUTCDate(d.getUTCDate() + daysAhead)
  d.setUTCHours(14, 0, 0, 0)
  return d.toISOString().replace('.000', '')
}

function formatTimestamp(value: string) {
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  return new Intl.DateTimeFormat('en-US', {
    dateStyle: 'medium',
    timeStyle: 'short',
  }).format(d)
}
