import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import './App.css'

type ShortenResponse = {
  short_url: string
  code: string
}

type URLItem = {
  long_url: string
  short_code: string
  clicks: number
  created_at: string
}

type URLListResponse = {
  urls: URLItem[]
}

type SortMode = 'newest' | 'oldest' | 'most-clicked'

const API_BASE = (import.meta.env.VITE_API_BASE_URL as string | undefined)?.replace(/\/$/, '') ?? ''

function apiPath(path: string): string {
  return API_BASE ? `${API_BASE}${path}` : path
}

function shortBaseUrl(): string {
  if (API_BASE) return API_BASE
  if (typeof window !== 'undefined') return window.location.origin
  return 'http://localhost:8080'
}

function normalizeUrl(value: string): string {
  const trimmed = value.trim()
  if (!trimmed) return ''
  if (/^https?:\/\//i.test(trimmed)) return trimmed
  return `https://${trimmed}`
}

function formatDate(value?: string): string {
  if (!value) return 'Unknown'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Unknown'
  return date.toLocaleString(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

function sortLinks(items: URLItem[], mode: SortMode): URLItem[] {
  const list = [...items]
  switch (mode) {
    case 'oldest':
      list.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
      return list
    case 'most-clicked':
      list.sort((a, b) => b.clicks - a.clicks)
      return list
    case 'newest':
    default:
      list.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime())
      return list
  }
}

function parseApiError(data: unknown, fallback: string): string {
  if (typeof data === 'object' && data !== null && 'error' in data) {
    const maybeError = (data as { error?: unknown }).error
    if (typeof maybeError === 'string' && maybeError.trim()) {
      return maybeError
    }
  }
  return fallback
}

async function copyText(value: string): Promise<void> {
  if (navigator.clipboard?.writeText) {
    await navigator.clipboard.writeText(value)
    return
  }

  const helper = document.createElement('textarea')
  helper.value = value
  helper.setAttribute('readonly', '')
  helper.style.position = 'fixed'
  helper.style.opacity = '0'
  document.body.appendChild(helper)
  helper.select()
  const copied = document.execCommand('copy')
  document.body.removeChild(helper)

  if (!copied) {
    throw new Error('Clipboard copy failed')
  }
}

function App() {
  const [urlInput, setUrlInput] = useState('')
  const [searchTerm, setSearchTerm] = useState('')
  const [sortMode, setSortMode] = useState<SortMode>('newest')
  const [shortened, setShortened] = useState<ShortenResponse | null>(null)
  const [recent, setRecent] = useState<URLItem[]>([])
  const [error, setError] = useState('')
  const [isShortening, setIsShortening] = useState(false)
  const [isLoadingRecent, setIsLoadingRecent] = useState(false)
  const [copied, setCopied] = useState(false)
  const [lastUpdated, setLastUpdated] = useState('')

  const shortUrlBase = useMemo(() => shortBaseUrl(), [])

  const visibleLinks = useMemo(() => {
    const lowered = searchTerm.trim().toLowerCase()
    const filtered = lowered
      ? recent.filter(
          (item) =>
            item.short_code.toLowerCase().includes(lowered) ||
            item.long_url.toLowerCase().includes(lowered),
        )
      : recent

    return sortLinks(filtered, sortMode)
  }, [recent, searchTerm, sortMode])

  const totalClicks = useMemo(
    () => recent.reduce((total, item) => total + item.clicks, 0),
    [recent],
  )

  async function loadRecentLinks() {
    setIsLoadingRecent(true)
    setError('')

    try {
      const response = await fetch(apiPath('/api/urls'))
      const data = (await response.json()) as URLListResponse | { error?: string }
      if (!response.ok) {
        throw new Error(parseApiError(data, 'Unable to fetch recent links.'))
      }

      setRecent((data as URLListResponse).urls)
      setLastUpdated(formatDate(new Date().toISOString()))
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Something went wrong while loading links.'
      setError(message)
    } finally {
      setIsLoadingRecent(false)
    }
  }

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const normalized = normalizeUrl(urlInput)

    if (!normalized) {
      setError('Please enter a URL.')
      return
    }

    setIsShortening(true)
    setError('')
    setCopied(false)

    try {
      const response = await fetch(apiPath('/api/shorten'), {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ url: normalized }),
      })

      const data = (await response.json()) as ShortenResponse | { error?: string }
      if (!response.ok) {
        throw new Error(parseApiError(data, 'Could not shorten this URL.'))
      }

      setShortened(data as ShortenResponse)
      setUrlInput('')
      await loadRecentLinks()
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Something went wrong while shortening the URL.'
      setError(message)
    } finally {
      setIsShortening(false)
    }
  }

  async function handleCopy() {
    if (!shortened) return

    try {
      await copyText(shortened.short_url)
      setCopied(true)
      setTimeout(() => setCopied(false), 1600)
    } catch {
      setError('Clipboard access is blocked in this browser context.')
    }
  }

  useEffect(() => {
    void loadRecentLinks()
  }, [])

  return (
    <div className="page-shell">
      <header className="hero">
        <p className="eyebrow">Fast links. Cleaner sharing.</p>
        <h1>URL Shortener</h1>
        <p className="hero-subtitle">
          Drop any long link and get a compact, trackable URL in seconds.
        </p>
      </header>

      <main className="content-grid">
        <section className="panel shorten-panel" aria-labelledby="shorten-title">
          <div className="panel-head">
            <h2 id="shorten-title">Create a short link</h2>
            <button
              type="button"
              className="ghost-btn"
              onClick={loadRecentLinks}
              disabled={isLoadingRecent}
            >
              {isLoadingRecent ? 'Refreshing...' : 'Refresh list'}
            </button>
          </div>

          <form onSubmit={handleSubmit} className="shorten-form">
            <label htmlFor="url-input">Long URL</label>
            <div className="input-row">
              <input
                id="url-input"
                type="url"
                value={urlInput}
                onChange={(event) => {
                  setUrlInput(event.target.value)
                  if (error) setError('')
                }}
                placeholder="https://example.com/very/long/path"
                autoComplete="off"
                required
              />
              <button type="submit" disabled={isShortening}>
                {isShortening ? 'Shortening...' : 'Shorten'}
              </button>
            </div>
          </form>

          {error ? <p className="status error">{error}</p> : null}

          {shortened ? (
            <div className="result-card" role="status" aria-live="polite">
              <p className="label">Latest short URL</p>
              <a href={shortened.short_url} target="_blank" rel="noreferrer">
                {shortened.short_url}
              </a>
              <div className="result-actions">
                <button type="button" onClick={handleCopy} className="ghost-btn">
                  {copied ? 'Copied' : 'Copy URL'}
                </button>
                <span>Code: {shortened.code}</span>
              </div>
            </div>
          ) : null}
        </section>

        <section className="panel stats-panel" aria-labelledby="stats-title">
          <h2 id="stats-title">Live stats</h2>
          <div className="stats-grid">
            <article>
              <p>Total links</p>
              <strong>{recent.length}</strong>
            </article>
            <article>
              <p>Total clicks</p>
              <strong>{totalClicks}</strong>
            </article>
            <article>
              <p>Last updated</p>
              <strong>{lastUpdated || 'Not yet loaded'}</strong>
            </article>
          </div>
        </section>

        <section className="panel list-panel" aria-labelledby="recent-title">
          <div className="panel-head">
            <h2 id="recent-title">Recent links</h2>
            <span>{visibleLinks.length} visible</span>
          </div>

          <div className="list-tools">
            <input
              type="search"
              value={searchTerm}
              onChange={(event) => setSearchTerm(event.target.value)}
              placeholder="Filter by code or URL"
              aria-label="Filter links"
            />
            <select
              value={sortMode}
              onChange={(event) => setSortMode(event.target.value as SortMode)}
              aria-label="Sort links"
            >
              <option value="newest">Newest</option>
              <option value="oldest">Oldest</option>
              <option value="most-clicked">Most clicked</option>
            </select>
          </div>

          {visibleLinks.length === 0 ? (
            <p className="status">No links match your filters yet.</p>
          ) : (
            <ul className="links-list">
              {visibleLinks.map((item) => (
                <li key={item.short_code}>
                  <a href={`${shortUrlBase}/${item.short_code}`} target="_blank" rel="noreferrer">
                    {shortUrlBase}/{item.short_code}
                  </a>
                  <p title={item.long_url}>{item.long_url}</p>
                  <div className="link-meta">
                    <span>{item.clicks} clicks</span>
                    <span>Created {formatDate(item.created_at)}</span>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </section>
      </main>
    </div>
  )
}

export default App
