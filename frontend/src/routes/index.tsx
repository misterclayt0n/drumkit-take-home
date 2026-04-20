import { createFileRoute } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ChevronLeft, ChevronRight, Plus, RotateCw } from 'lucide-react'
import { useCallback, useMemo, useState } from 'react'
import { CreateLoadForm } from '#/components/create-load-form'
import { LoadDetailModal } from '#/components/load-detail-modal'
import { LoadTable } from '#/components/load-table'
import { Button } from '#/components/ui/button'
import { Input } from '#/components/ui/input'
import { createLoad, listLoads } from '#/lib/api'
import type { Load } from '#/lib/types'

export const Route = createFileRoute('/')({ component: App })

const DEFAULT_PAGE_SIZE = 20

function App() {
  const queryClient = useQueryClient()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(DEFAULT_PAGE_SIZE)
  const [status, setStatus] = useState('')
  const [customerId, setCustomerId] = useState('')
  const [pickupDateSearchFrom, setPickupDateSearchFrom] = useState('')
  const [pickupDateSearchTo, setPickupDateSearchTo] = useState('')
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [selectedLoad, setSelectedLoad] = useState<Load | null>(null)

  const loadsQuery = useQuery({
    queryKey: [
      'loads',
      { page, pageSize, status, customerId, pickupDateSearchFrom, pickupDateSearchTo },
    ],
    queryFn: () =>
      listLoads({
        page,
        limit: pageSize,
        status: status || undefined,
        customerId: customerId || undefined,
        pickupDateSearchFrom: pickupDateSearchFrom || undefined,
        pickupDateSearchTo: pickupDateSearchTo || undefined,
      }),
  })

  const createMutation = useMutation({
    mutationFn: (payload: Load) => createLoad(payload),
    onSuccess: async () => {
      setPage(1)
      await queryClient.invalidateQueries({ queryKey: ['loads'] })
    },
  })

  const loads = loadsQuery.data?.data ?? []
  const pagination = loadsQuery.data?.pagination

  const stats = useMemo(() => {
    const customers = new Set(loads.map((l) => l.customer.name).filter(Boolean)).size
    return { customers, total: pagination?.total ?? 0 }
  }, [loads, pagination])

  const closeDetailModal = useCallback(() => setSelectedLoad(null), [])
  const closeCreateModal = useCallback(() => setShowCreateModal(false), [])

  return (
    <main id="top" className="overflow-x-hidden w-full max-w-full">
      {/* Top bar */}
      <section className="page-wrap pt-8 pb-6">
        <div className="flex flex-col gap-6 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="text-2xl font-semibold tracking-tight text-[var(--dk-ink)] sm:text-3xl">
              Loads
            </h1>
            <p className="mt-1 text-sm text-[var(--dk-ink-soft)]">
              {stats.total > 0 ? (
                <>
                  <span className="font-mono font-medium text-[var(--dk-ink)]">{stats.total}</span> loads from Turvo
                  {stats.customers > 0 && (
                    <> &middot; <span className="font-mono font-medium text-[var(--dk-ink)]">{stats.customers}</span> customers on this page</>
                  )}
                </>
              ) : (
                'Loading from Turvo...'
              )}
            </p>
          </div>

          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              onClick={() => loadsQuery.refetch()}
              disabled={loadsQuery.isFetching}
              className="h-9 gap-1.5 rounded-lg border-[var(--dk-line-strong)] bg-white px-3 text-sm text-[var(--dk-ink)] hover:bg-[var(--dk-surface-raised)]"
            >
              <RotateCw className={`size-3.5 ${loadsQuery.isFetching ? 'animate-spin' : ''}`} />
              <span className="hidden sm:inline">Refresh</span>
            </Button>

            <Button
              onClick={() => setShowCreateModal(true)}
              className="h-9 gap-1.5 rounded-lg bg-[var(--dk-ink)] px-3.5 text-sm font-medium text-white hover:bg-[var(--dk-ink)]/90"
            >
              <Plus className="size-3.5" />
              New load
            </Button>
          </div>
        </div>

        {/* Filters */}
        <div className="mt-4 grid gap-3 rounded-xl border border-[var(--dk-line)] bg-white p-4 sm:grid-cols-2 lg:grid-cols-5">
          <div className="relative">
            <select
              value={status}
              onChange={(e) => { setStatus(e.target.value); setPage(1) }}
              className="h-9 w-full appearance-none rounded-lg border border-[var(--dk-line-strong)] bg-white pl-3 pr-8 text-sm font-medium text-[var(--dk-ink)] outline-none transition focus:border-[var(--dk-red)] focus:ring-2 focus:ring-[var(--ring)]"
            >
              <option value="">All statuses</option>
              <option value="Tendered">Tendered</option>
              <option value="Covered">Covered</option>
            </select>
            <ChevronRight className="pointer-events-none absolute right-2 top-1/2 size-3.5 -translate-y-1/2 rotate-90 text-[var(--dk-ink-soft)]" />
          </div>

          <Input
            value={customerId}
            onChange={(e) => { setCustomerId(e.target.value); setPage(1) }}
            placeholder="Customer ID"
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />

          <Input
            type="date"
            value={pickupDateSearchFrom}
            onChange={(e) => { setPickupDateSearchFrom(e.target.value); setPage(1) }}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />

          <Input
            type="date"
            value={pickupDateSearchTo}
            onChange={(e) => { setPickupDateSearchTo(e.target.value); setPage(1) }}
            className="h-9 rounded-lg border-[var(--dk-line-strong)] bg-white text-sm"
          />

          <div className="relative">
            <select
              value={pageSize}
              onChange={(e) => { setPageSize(Number(e.target.value)); setPage(1) }}
              className="h-9 w-full appearance-none rounded-lg border border-[var(--dk-line-strong)] bg-white pl-3 pr-8 text-sm font-medium text-[var(--dk-ink)] outline-none transition focus:border-[var(--dk-red)] focus:ring-2 focus:ring-[var(--ring)]"
            >
              <option value={20}>20 per page</option>
              <option value={50}>50 per page</option>
              <option value={100}>100 per page</option>
            </select>
            <ChevronRight className="pointer-events-none absolute right-2 top-1/2 size-3.5 -translate-y-1/2 rotate-90 text-[var(--dk-ink-soft)]" />
          </div>
        </div>
      </section>

      {/* Table */}
      <section id="loads" className="page-wrap pb-16">
        <div className="space-y-4">
          {loadsQuery.error && (
            <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm">
              <p className="font-medium text-red-900">Failed to load data</p>
              <p className="mt-1 text-red-700">
                {loadsQuery.error instanceof Error ? loadsQuery.error.message : 'Unknown error'}
              </p>
            </div>
          )}

          <LoadTable
            data={loads}
            isLoading={loadsQuery.isLoading || loadsQuery.isFetching}
            onSelectLoad={setSelectedLoad}
          />

          {pagination && pagination.pages > 1 && (
            <div className="flex items-center justify-between rounded-xl border border-[var(--dk-line)] bg-white px-4 py-3">
              <p className="text-sm text-[var(--dk-ink-soft)]">
                Page <span className="font-mono font-medium text-[var(--dk-ink)]">{pagination.page}</span> of{' '}
                <span className="font-mono font-medium text-[var(--dk-ink)]">{pagination.pages}</span>
              </p>
              <div className="flex gap-1.5">
                <Button
                  variant="outline" size="sm"
                  className="h-8 gap-1 rounded-lg border-[var(--dk-line-strong)] px-3 text-[var(--dk-ink)]"
                  disabled={page <= 1 || loadsQuery.isFetching}
                  onClick={() => setPage((c) => Math.max(1, c - 1))}
                >
                  <ChevronLeft className="size-3.5" /> Prev
                </Button>
                <Button
                  variant="outline" size="sm"
                  className="h-8 gap-1 rounded-lg border-[var(--dk-line-strong)] px-3 text-[var(--dk-ink)]"
                  disabled={page >= pagination.pages || loadsQuery.isFetching}
                  onClick={() => setPage((c) => c + 1)}
                >
                  Next <ChevronRight className="size-3.5" />
                </Button>
              </div>
            </div>
          )}
        </div>
      </section>

      {/* Modals */}
      <LoadDetailModal load={selectedLoad} onClose={closeDetailModal} />
      {showCreateModal && (
        <CreateLoadForm
          onSubmit={async (payload) => { await createMutation.mutateAsync(payload) }}
          isPending={createMutation.isPending}
          success={createMutation.data ?? null}
          error={createMutation.error instanceof Error ? createMutation.error.message : null}
          onClose={closeCreateModal}
        />
      )}
    </main>
  )
}
