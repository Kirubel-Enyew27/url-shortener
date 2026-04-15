import { useState, useEffect, type FormEvent } from 'react';
import axios from 'axios';
import './App.css';
import type { URL, ShortenResponse } from './types';

const API_BASE = 'http://localhost:8080';

function App() {
  const [longUrl, setLongUrl] = useState('');
  const [shortUrl, setShortUrl] = useState('');
  const [urls, setUrls] = useState<URL[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    fetchUrls();
  }, []);

  const fetchUrls = async () => {
    try {
      const res = await axios.get<{ urls: URL[] }>(`${API_BASE}/api/urls`);
      setUrls(res.data.urls);
    } catch {
      // Backend not running yet
    }
  };

  const shortenUrl = async (e: FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setShortUrl('');
    setCopied(false);

    try {
      const res = await axios.post<ShortenResponse>(`${API_BASE}/api/shorten`, {
        url: longUrl,
      });
      setShortUrl(res.data.short_url);
      setLongUrl('');
      fetchUrls();
} catch (err: unknown) {
      if (err && typeof err === 'object' && 'response' in err) {
        const axiosErr = err as { response?: { data?: { error?: string } } };
        setError(axiosErr.response?.data?.error || 'Failed to shorten URL');
      } else {
        setError('Failed to shorten URL');
      }
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(shortUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const getShortLink = (code: string) => `${API_BASE}/${code}`;

  return (
    <main className="container">
      <header className="header">
        <h1>URL Shortener</h1>
        <p>Convert long links into short, shareable URLs</p>
      </header>

      <section className="shorten-section">
        <form onSubmit={shortenUrl} className="shorten-form">
          <div className="input-group">
            <input
              type="url"
              value={longUrl}
              onChange={(e) => setLongUrl(e.target.value)}
              placeholder="Paste your long URL here..."
              required
            />
            <button type="submit" disabled={loading}>
              {loading ? 'Shortening...' : 'Shorten'}
            </button>
          </div>
        </form>

        {error && <div className="error-message">{error}</div>}

        {shortUrl && (
          <div className="result">
            <span className="result-label">Your short URL:</span>
            <div className="result-link">
              <a href={shortUrl} target="_blank" rel="noopener noreferrer">
                {shortUrl}
              </a>
              <button onClick={copyToClipboard} className="copy-btn">
                {copied ? 'Copied!' : 'Copy'}
              </button>
            </div>
          </div>
        )}
      </section>

      <section className="dashboard">
        <h2>Your URLs</h2>
        {urls.length === 0 ? (
          <p className="empty-state">No URLs created yet. Shorten your first link above!</p>
        ) : (
          <div className="urls-table">
            <div className="table-header">
              <span>Short URL</span>
              <span>Original URL</span>
              <span>Clicks</span>
            </div>
            <div className="table-body">
              {urls.map((url) => (
                <div key={url.short_code} className="table-row">
                  <a
                    href={getShortLink(url.short_code)}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="short-link"
                  >
                    {getShortLink(url.short_code)}
                  </a>
                  <span className="long-link" title={url.long_url}>
                    {url.long_url}
                  </span>
                  <span className="clicks">{url.clicks}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </section>

      <footer className="footer">
        <p>Redirects tracked automatically</p>
      </footer>
    </main>
  );
}

export default App;