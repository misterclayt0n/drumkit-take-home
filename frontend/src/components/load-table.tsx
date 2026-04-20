import {
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
} from '@tanstack/react-table'
import { ArrowRight, Package } from 'lucide-react'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '#/components/ui/table'
import type { Load } from '#/lib/types'

const columns: ColumnDef<Load>[] = [
  {
    header: 'Load',
    cell: ({ row }) => (
      <div className="min-w-[10rem]">
        <p className="text-sm font-medium text-[var(--dk-ink)]">
          {row.original.freightLoadID || 'Unnumbered'}
        </p>
        <p className="mt-0.5 font-mono text-xs text-[var(--dk-ink-soft)]">
          {row.original.externalTMSLoadID || '--'}
        </p>
      </div>
    ),
  },
  {
    header: 'Customer',
    cell: ({ row }) => (
      <div className="min-w-[10rem]">
        <p className="text-sm font-medium text-[var(--dk-ink)]">
          {row.original.customer.name || '--'}
        </p>
        <p className="mt-0.5 font-mono text-xs text-[var(--dk-ink-soft)]">
          {row.original.customer.externalTMSId || 'No ID'}
        </p>
      </div>
    ),
  },
  {
    header: 'Lane',
    cell: ({ row }) => (
      <div className="min-w-[14rem]">
        <p className="flex items-center gap-1.5 text-sm font-medium text-[var(--dk-ink)]">
          <span>{compactLocation(row.original.pickup.city, row.original.pickup.state)}</span>
          <ArrowRight className="size-3 text-[var(--dk-ink-soft)]" />
          <span>{compactLocation(row.original.consignee.city, row.original.consignee.state)}</span>
        </p>
        <p className="mt-0.5 text-xs text-[var(--dk-ink-soft)]">
          {formatDate(row.original.pickup.apptTime)} &rarr; {formatDate(row.original.consignee.apptTime)}
        </p>
      </div>
    ),
  },
  {
    header: 'Status',
    cell: ({ row }) => <StatusPill status={row.original.status} />,
  },
  {
    header: 'PO / Miles',
    cell: ({ row }) => (
      <div className="min-w-[8rem] text-sm">
        <p className="font-medium text-[var(--dk-ink)]">{row.original.poNums || '--'}</p>
        <p className="mt-0.5 font-mono text-xs text-[var(--dk-ink-soft)]">
          {row.original.routeMiles ? `${Math.round(row.original.routeMiles)} mi` : '--'}
        </p>
      </div>
    ),
  },
]

export function LoadTable({
  data,
  isLoading,
  onSelectLoad,
}: {
  data: Load[]
  isLoading: boolean
  onSelectLoad?: (load: Load) => void
}) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
  })

  if (isLoading) {
    return <LoadTableSkeleton />
  }

  if (data.length === 0) {
    return (
      <div className="flex min-h-64 flex-col items-center justify-center rounded-xl border border-dashed border-[var(--dk-line-strong)] bg-white px-6 text-center">
        <div className="mb-3 flex size-10 items-center justify-center rounded-lg bg-[var(--dk-surface-raised)]">
          <Package className="size-5 text-[var(--dk-ink-soft)]" />
        </div>
        <h3 className="text-sm font-semibold text-[var(--dk-ink)]">No loads found</h3>
        <p className="mt-1 max-w-xs text-sm text-[var(--dk-ink-soft)]">
          Try changing the status filter, or create a new load.
        </p>
      </div>
    )
  }

  return (
    <div className="overflow-hidden rounded-xl border border-[var(--dk-line)] bg-white">
      <Table>
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id} className="border-[var(--dk-line)] hover:bg-transparent">
              {headerGroup.headers.map((header) => (
                <TableHead
                  key={header.id}
                  className="h-10 bg-[var(--dk-surface-raised)] px-4 text-xs font-medium text-[var(--dk-ink-soft)]"
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows.map((row) => (
            <TableRow
              key={row.id}
              className="cursor-pointer border-[var(--dk-line)] transition-colors hover:bg-[var(--dk-surface-raised)]"
              onClick={() => onSelectLoad?.(row.original)}
            >
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id} className="px-4 py-3 align-top">
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}

function LoadTableSkeleton() {
  return (
    <div className="overflow-hidden rounded-xl border border-[var(--dk-line)] bg-white">
      <div className="grid grid-cols-5 gap-4 bg-[var(--dk-surface-raised)] px-4 py-2.5 text-xs font-medium text-[var(--dk-ink-soft)]">
        <span>Load</span>
        <span>Customer</span>
        <span>Lane</span>
        <span>Status</span>
        <span>PO / Miles</span>
      </div>
      <div className="divide-y divide-[var(--dk-line)]">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="grid grid-cols-5 gap-4 px-4 py-3">
            {Array.from({ length: 5 }).map((__, j) => (
              <div key={j} className="space-y-1.5">
                <div className="skeleton-shimmer h-4 w-3/4 rounded-md" />
                <div className="skeleton-shimmer h-3 w-1/2 rounded-md" />
              </div>
            ))}
          </div>
        ))}
      </div>
    </div>
  )
}

function StatusPill({ status }: { status: string }) {
  const isActive = status.toLowerCase() === 'tendered'
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ${
        isActive
          ? 'bg-gradient-to-r from-[rgba(252,6,67,0.08)] to-[rgba(254,139,87,0.08)] text-[var(--dk-red)]'
          : 'bg-[var(--dk-surface-raised)] text-[var(--dk-ink-soft)]'
      }`}
    >
      <span
        className={`inline-block size-1.5 rounded-full ${
          isActive ? 'bg-[var(--dk-red)]' : 'bg-[var(--dk-ink-soft)]/40'
        }`}
      />
      {status || 'Unknown'}
    </span>
  )
}

function compactLocation(city: string, state: string) {
  if (!city && !state) return 'TBD'
  return [city, state].filter(Boolean).join(', ')
}

function formatDate(value: string) {
  if (!value) return 'TBD'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  }).format(date)
}
