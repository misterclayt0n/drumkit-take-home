import { createFileRoute } from '@tanstack/react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { ChevronLeft, ChevronRight, Plus, RotateCw, X } from 'lucide-react'
import { useMemo, useState } from 'react'
import { CreateLoadForm } from '#/components/create-load-form'
import { LoadTable } from '#/components/load-table'
import { Button } from '#/components/ui/button'
import { createLoad, listLoads } from '#/lib/api'
import type { Load } from '#/lib/types'

export const Route = createFileRoute('/')({ component: App })

const PAGE_SIZE = 20

function App() {
  const queryClient = useQueryClient()
  const [page, setPage] = useState(1)
  const [status, setStatus] = useState('')
  const [showCreatePanel, setShowCreatePanel] = useState(false)

  const loadsQuery = useQuery({
    queryKey: ['loads', { page, status, limit: PAGE_SIZE }],
    queryFn: () =>
      listLoads({
        page,
        limit: PAGE_SIZE,
        status: status || undefined,
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

  return (
    <main id="top" className="overflow-x-hidden w-full max-w-full">
      {/* Top bar with stats and actions */}
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
            {/* Status filter */}
            <div className="relative">
              <select
                value={status}
                onChange={(e) => {
                  setStatus(e.target.value)
                  setPage(1)
                }}
                className="h-9 appearance-none rounded-lg border border-[var(--dk-line-strong)] bg-white pl-3 pr-8 text-sm font-medium text-[var(--dk-ink)] outline-none transition focus:border-[var(--dk-red)] focus:ring-2 focus:ring-[var(--ring)]"
              >
                <option value="">All statuses</option>
                <option value="Tendered">Tendered</option>
                <option value="Covered">Covered</option>
              </select>
              <ChevronRight className="pointer-events-none absolute right-2 top-1/2 size-3.5 -translate-y-1/2 rotate-90 text-[var(--dk-ink-soft)]" />
            </div>

            {/* Refresh */}
            <Button
              variant="outline"
              onClick={() => loadsQuery.refetch()}
              disabled={loadsQuery.isFetching}
              className="h-9 gap-1.5 rounded-lg border-[var(--dk-line-strong)] bg-white px-3 text-sm text-[var(--dk-ink)] hover:bg-[var(--dk-surface-raised)]"
            >
              <RotateCw className={`size-3.5 ${loadsQuery.isFetching ? 'animate-spin' : ''}`} />
              <span className="hidden sm:inline">Refresh</span>
            </Button>

            {/* Create button */}
            <Button
              onClick={() => setShowCreatePanel(!showCreatePanel)}
              className="h-9 gap-1.5 rounded-lg bg-[var(--dk-ink)] px-3.5 text-sm font-medium text-white hover:bg-[var(--dk-ink)]/90"
            >
              {showCreatePanel ? (
                <X className="size-3.5" />
              ) : (
                <Plus className="size-3.5" />
              )}
              {showCreatePanel ? 'Close' : 'New load'}
            </Button>
          </div>
        </div>
      </section>

      {/* Main content */}
      <section id="loads" className="page-wrap pb-16">
        <div className={`grid gap-6 ${showCreatePanel ? 'lg:grid-cols-[minmax(0,1fr)_420px]' : ''} lg:items-start`}>
          {/* Table area */}
          <div className="space-y-4">
            {loadsQuery.error && (
              <div className="rounded-xl border border-red-200 bg-red-50 px-4 py-3 text-sm">
                <p className="font-medium text-red-900">Failed to load data</p>
                <p className="mt-1 text-red-700">
                  {loadsQuery.error instanceof Error ? loadsQuery.error.message : 'Unknown error'}
                </p>
              </div>
            )}

            <LoadTable data={loads} isLoading={loadsQuery.isLoading || loadsQuery.isFetching} />

            {/* Pagination */}
            {pagination && pagination.pages > 1 && (
              <div className="flex items-center justify-between rounded-xl border border-[var(--dk-line)] bg-white px-4 py-3">
                <p className="text-sm text-[var(--dk-ink-soft)]">
                  Page <span className="font-mono font-medium text-[var(--dk-ink)]">{pagination.page}</span> of{' '}
                  <span className="font-mono font-medium text-[var(--dk-ink)]">{pagination.pages}</span>
                </p>
                <div className="flex gap-1.5">
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-8 gap-1 rounded-lg border-[var(--dk-line-strong)] px-3 text-[var(--dk-ink)]"
                    disabled={page <= 1 || loadsQuery.isFetching}
                    onClick={() => setPage((c) => Math.max(1, c - 1))}
                  >
                    <ChevronLeft className="size-3.5" />
                    Prev
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    className="h-8 gap-1 rounded-lg border-[var(--dk-line-strong)] px-3 text-[var(--dk-ink)]"
                    disabled={page >= pagination.pages || loadsQuery.isFetching}
                    onClick={() => setPage((c) => c + 1)}
                  >
                    Next
                    <ChevronRight className="size-3.5" />
                  </Button>
                </div>
              </div>
            )}
          </div>

          {/* Create panel - slides in */}
          {showCreatePanel && (
            <div id="create">
              <CreateLoadForm
                onSubmit={async (payload) => {
                  await createMutation.mutateAsync(payload)
                }}
                isPending={createMutation.isPending}
                success={createMutation.data ?? null}
                error={createMutation.error instanceof Error ? createMutation.error.message : null}
              />
            </div>
          )}
        </div>
      </section>
    </main>
  )
}
