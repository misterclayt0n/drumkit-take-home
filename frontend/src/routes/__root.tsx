import { HeadContent, Scripts, createRootRoute } from '@tanstack/react-router'
import Footer from '../components/Footer'
import Header from '../components/Header'
import { QueryProvider } from '../components/query-provider'

import appCss from '../styles.css?url'

export const Route = createRootRoute({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: 'Drumkit × Turvo Loads',
      },
      {
        name: 'description',
        content:
          'A Drumkit-style frontend for listing and creating Turvo loads through a Go backend.',
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
      {
        rel: 'icon',
        type: 'image/png',
        href: '/drumkit-favicon.png',
      },
    ],
  }),
  shellComponent: RootDocument,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body className="font-sans antialiased selection:bg-[rgba(252,6,67,0.12)] selection:text-[var(--dk-ink)]">
        <QueryProvider>
          <Header />
          {children}
          <Footer />
        </QueryProvider>
        <Scripts />
      </body>
    </html>
  )
}
