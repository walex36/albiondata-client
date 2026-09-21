<script>
  import { onDestroy } from 'svelte';
  import { Browser, Events } from '@wailsio/runtime';
  import { DashboardService } from '../../bindings/github.com/ao-data/albiondata-client/internal/dashboard/index.js';
  import CountersPanel from './CountersPanel.svelte';

  const WIDTH_KEY = 'sidebar-width';
  const MIN_WIDTH = 168;
  const MAX_WIDTH = 360;
  const DEFAULT_WIDTH = 208;

  function loadWidth() {
    const stored = Number(localStorage.getItem(WIDTH_KEY));
    if (Number.isFinite(stored) && stored >= MIN_WIDTH && stored <= MAX_WIDTH) return stored;
    return DEFAULT_WIDTH;
  }

  let width = $state(loadWidth());
  let resizing = $state(false);

  function startResize(e) {
    e.preventDefault();
    resizing = true;
    window.addEventListener('pointermove', onResize);
    window.addEventListener('pointerup', stopResize);
  }

  function onResize(e) {
    width = Math.min(MAX_WIDTH, Math.max(MIN_WIDTH, e.clientX));
  }

  function stopResize() {
    resizing = false;
    localStorage.setItem(WIDTH_KEY, String(width));
    window.removeEventListener('pointermove', onResize);
    window.removeEventListener('pointerup', stopResize);
  }

  onDestroy(() => {
    window.removeEventListener('pointermove', onResize);
    window.removeEventListener('pointerup', stopResize);
  });

  let status = $state({
    Version: '',
    UpdateAvailable: '',
    CaptureRunning: false,
    CaptureError: false,
    ServerID: 0,
    IngestBaseURL: '',
    CustomPublicIngest: false,
    DriverWarning: '',
    DriverHelpURL: '',
    EncryptionStatus: '',
    CharacterName: '',
    AlbionMarketSyncEnabled: false,
    SyncToken: '',
    AlbionMarketAPIUrl: '',
  });

  DashboardService.GetStatus().then((s) => (status = s));

  const unlisten = Events.On('status:changed', (evt) => {
    status = evt.data;
  });

  onDestroy(unlisten);

  const serverNames = { 0: 'Unknown', 1: 'West', 2: 'East', 3: 'Europe' };
  let serverLabel = $derived(
    status.CustomPublicIngest ? 'Private' : (serverNames[status.ServerID] ?? status.ServerID)
  );
  let badgeLabel = $derived(
    status.CaptureError ? 'Error' : status.CaptureRunning ? 'Capturing' : 'Ready'
  );
  let encryptionLabel = $derived(
    status.EncryptionStatus === 'encrypted'
      ? 'Encrypted'
      : status.EncryptionStatus === 'clear'
        ? 'Not Encrypted'
        : 'Encrypted?'
  );

  function openDriverHelp(e) {
    e.preventDefault();
    Browser.OpenURL(status.DriverHelpURL);
  }

  let showSettings = $state(false);
  let tokenInput = $state('');
  let apiUrlInput = $state('');
  let saveFeedback = $state('');
  let saveFeedbackTimer;

  function openSettings() {
    tokenInput = status.SyncToken || '';
    apiUrlInput = status.AlbionMarketAPIUrl || 'http://localhost:3001';
    saveFeedback = '';
    showSettings = true;
  }

  function closeSettings() {
    showSettings = false;
    clearTimeout(saveFeedbackTimer);
  }

  function saveSettings() {
    Events.Emit({
      name: 'albionmarket:saveconfig',
      data: {
        token: tokenInput.trim(),
        apiUrl: apiUrlInput.trim(),
      },
    });
    saveFeedback = 'Saved!';
    clearTimeout(saveFeedbackTimer);
    saveFeedbackTimer = setTimeout(() => {
      saveFeedback = '';
      showSettings = false;
    }, 750);
  }
</script>

<aside class="sidebar" style="width: {width}px">
  <div
    class="resize-handle"
    class:active={resizing}
    onpointerdown={startResize}
    role="separator"
    aria-orientation="vertical"
    aria-label="Resize sidebar"
  ></div>
  <div class="brand">
    <span class="wordmark">Albion Data Client</span>
  </div>

  <div class="group">
    <div class="pill-row">
      <span
        class="status-pill"
        class:running={status.CaptureRunning && !status.CaptureError}
        class:error={status.CaptureError}
      >
        <span class="lantern"><span class="glow"></span></span>
        {badgeLabel}
      </span>

      <span
        class="status-pill"
        class:encrypted={status.EncryptionStatus === 'encrypted'}
        class:clear={status.EncryptionStatus === 'clear'}
      >
        <span class="lantern"><span class="glow"></span></span>
        {encryptionLabel}
      </span>
    </div>

    <div class="field">
      <span class="label">Server</span>
      <span class="value">{serverLabel}</span>
    </div>

    {#if status.CharacterName}
      <div class="field">
        <span class="label">Character</span>
        <span class="value">{status.CharacterName}</span>
      </div>
    {/if}

    <div class="field">
      <span class="label">Albion Market</span>
      <div class="sync-row">
        <span
          class="sync-pill"
          class:synced={status.AlbionMarketSyncEnabled}
        >
          {status.AlbionMarketSyncEnabled ? 'Token Linked' : 'No Token'}
        </span>
        <button class="config-btn" type="button" onclick={openSettings} title="Configure Albion Market Token">
          ⚙️
        </button>
      </div>
    </div>
  </div>

  {#if status.DriverWarning}
    <div class="group driver-warning-group">
      <div class="driver-warning">
        <span class="driver-warning-text">{status.DriverWarning}</span>
        {#if status.DriverHelpURL}
          <a class="driver-warning-link" href={status.DriverHelpURL} onclick={openDriverHelp}>
            Get Npcap
          </a>
        {/if}
      </div>
    </div>
  {/if}

  <CountersPanel />

  {#if status.UpdateAvailable}
    <div class="group update-group">
      <span class="update-badge">
        <span class="label">Update available</span>
        <span class="value">{status.UpdateAvailable}</span>
      </span>
    </div>
  {/if}

  {#if showSettings}
    <div class="modal-backdrop" onclick={closeSettings} role="presentation">
      <div
        class="modal-card"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.key === 'Escape' && closeSettings()}
        role="dialog"
        aria-modal="true"
        tabindex="-1"
      >
        <div class="modal-header">
          <h3 class="modal-title">Albion Market Settings</h3>
          <button class="close-btn" type="button" onclick={closeSettings} aria-label="Close">✕</button>
        </div>
        <div class="modal-body">
          <label class="input-label" for="sync-token-input">
            Sync Token
            <span class="input-hint">Used to authenticate your skills and Destiny Board synchronization.</span>
          </label>
          <input
            id="sync-token-input"
            class="text-input"
            type="password"
            placeholder="am_sync_..."
            bind:value={tokenInput}
          />

          <label class="input-label" for="api-url-input">
            API URL
            <span class="input-hint">Albion Market backend URL (default: http://localhost:3001).</span>
          </label>
          <input
            id="api-url-input"
            class="text-input"
            type="text"
            placeholder="http://localhost:3001"
            bind:value={apiUrlInput}
          />
        </div>
        <div class="modal-footer">
          {#if saveFeedback}
            <span class="save-feedback">{saveFeedback}</span>
          {/if}
          <button class="btn btn-secondary" type="button" onclick={closeSettings}>Cancel</button>
          <button class="btn btn-primary" type="button" onclick={saveSettings}>Save</button>
        </div>
      </div>
    </div>
  {/if}
</aside>

<style>
  .sidebar {
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 1.75rem;
    flex-shrink: 0;
    padding: 1.5rem 1.4rem 1rem;
    background: var(--bg-raised);
    border-right: 1px solid var(--border);
    overflow-y: auto;
    overflow-x: hidden;
  }
  .resize-handle {
    position: absolute;
    top: 0;
    right: -3px;
    width: 6px;
    height: 100%;
    cursor: col-resize;
    touch-action: none;
    z-index: 1;
  }
  .resize-handle::after {
    content: '';
    position: absolute;
    top: 0;
    left: 2px;
    width: 2px;
    height: 100%;
    background: transparent;
    transition: background-color 0.15s ease;
  }
  .resize-handle:hover::after,
  .resize-handle.active::after {
    background: var(--blue);
  }
  .brand {
    display: flex;
    min-width: 0;
  }
  .wordmark {
    font-family: var(--font-display);
    font-size: 1.05rem;
    font-weight: 600;
    letter-spacing: 0.01em;
    color: var(--text);
    line-height: 1.2;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .group {
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
    padding-top: 1.15rem;
    border-top: 1px solid var(--border);
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }
  .label {
    font-size: 0.66rem;
    font-weight: 600;
    letter-spacing: 0.06em;
    text-transform: uppercase;
    color: var(--text-faint);
  }
  .value {
    font-family: var(--font-mono);
    font-size: 0.85rem;
    color: var(--text);
  }
  .pill-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.5rem;
  }
  .sync-row {
    display: flex;
    align-items: center;
    margin-top: 0.15rem;
  }
  .sync-pill {
    display: inline-flex;
    align-items: center;
    padding: 0.2rem 0.55rem;
    border-radius: var(--radius-sm);
    font-size: 0.72rem;
    font-weight: 600;
    letter-spacing: 0.03em;
    background: rgba(255, 255, 255, 0.05);
    color: var(--text-faint);
    border: 1px solid var(--border);
  }
  .sync-pill.synced {
    background: var(--blue-soft);
    color: var(--blue-bright);
    border-color: rgba(47, 140, 255, 0.4);
  }
  .status-pill {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    align-self: flex-start;
    padding: 0.3rem 0.7rem 0.3rem 0.55rem;
    border-radius: 999px;
    font-size: 0.78rem;
    font-weight: 500;
    background: var(--ready-soft);
    color: var(--ready);
  }
  .lantern {
    position: relative;
    display: inline-flex;
    width: 7px;
    height: 7px;
  }
  .lantern::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: 50%;
    background: currentColor;
  }
  .lantern .glow {
    position: absolute;
    inset: -4px;
    border-radius: 50%;
    background: currentColor;
    opacity: 0;
  }
  .status-pill.running {
    background: var(--blue-soft);
    color: var(--blue-bright);
  }
  .status-pill.running .glow {
    animation: breathe 2.4s ease-in-out infinite;
  }
  .status-pill.error,
  .status-pill.encrypted {
    background: var(--ember-soft);
    color: var(--ember);
  }
  .status-pill.clear {
    background: var(--moss-soft);
    color: var(--moss);
  }
  @keyframes breathe {
    0%, 100% { opacity: 0; transform: scale(0.6); }
    50% { opacity: 0.45; transform: scale(1.4); }
  }
  @media (prefers-reduced-motion: reduce) {
    .status-pill.running .glow {
      animation: none;
    }
  }
  .driver-warning-group {
    border-top: 1px solid var(--border);
  }
  .driver-warning {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    padding: 0.6rem 0.7rem;
    border-radius: var(--radius-sm);
    background: var(--blue-soft);
    border: 1px solid rgba(47, 140, 255, 0.3);
  }
  .driver-warning-text {
    font-size: 0.78rem;
    line-height: 1.5;
    color: var(--blue);
  }
  .driver-warning-link {
    align-self: flex-start;
    font-size: 0.72rem;
    font-weight: 600;
    letter-spacing: 0.02em;
    text-transform: uppercase;
    color: var(--blue);
    text-decoration: none;
    border-bottom: 1px solid rgba(47, 140, 255, 0.5);
  }
  .driver-warning-link:hover {
    border-bottom-color: var(--blue);
  }
  .update-group {
    border-top: 1px solid var(--border);
  }
  .update-badge {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
    padding: 0.6rem 0.7rem;
    border-radius: var(--radius-sm);
    background: var(--blue-soft);
    border: 1px solid rgba(47, 140, 255, 0.25);
  }
  .update-badge .label {
    color: var(--blue);
    opacity: 0.9;
  }
  .update-badge .value {
    font-family: var(--font-body);
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text);
  }

  .config-btn {
    background: transparent;
    border: none;
    font-size: 0.85rem;
    cursor: pointer;
    margin-left: 0.4rem;
    opacity: 0.65;
    padding: 0.1rem 0.25rem;
    border-radius: var(--radius-sm);
    transition: opacity 0.15s, background 0.15s;
    line-height: 1;
  }
  .config-btn:hover {
    opacity: 1;
    background: rgba(255, 255, 255, 0.08);
  }

  .modal-backdrop {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.65);
    backdrop-filter: blur(3px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 9999;
  }
  .modal-card {
    background: var(--bg-raised);
    border: 1px solid var(--border-strong);
    border-radius: var(--radius);
    width: 360px;
    max-width: 90vw;
    box-shadow: 0 12px 32px rgba(0, 0, 0, 0.5);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
  .modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0.85rem 1.15rem;
    border-bottom: 1px solid var(--border);
  }
  .modal-title {
    margin: 0;
    font-size: 0.95rem;
    font-weight: 600;
    color: var(--text);
  }
  .close-btn {
    background: transparent;
    border: none;
    color: var(--text-faint);
    font-size: 0.95rem;
    cursor: pointer;
    padding: 0.2rem 0.4rem;
    border-radius: var(--radius-sm);
  }
  .close-btn:hover {
    color: var(--text);
    background: rgba(255, 255, 255, 0.06);
  }
  .modal-body {
    padding: 1.15rem;
    display: flex;
    flex-direction: column;
    gap: 0.9rem;
  }
  .input-label {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    font-size: 0.76rem;
    font-weight: 600;
    color: var(--text);
    letter-spacing: 0.02em;
  }
  .input-hint {
    font-size: 0.68rem;
    font-weight: 400;
    color: var(--text-muted);
    line-height: 1.35;
  }
  .text-input {
    background: var(--bg-sunken);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text);
    font-family: var(--font-mono);
    font-size: 0.82rem;
    padding: 0.5rem 0.65rem;
    outline: none;
    transition: border-color 0.15s;
  }
  .text-input:focus {
    border-color: var(--blue);
  }
  .modal-footer {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    gap: 0.6rem;
    padding: 0.75rem 1.15rem;
    background: rgba(0, 0, 0, 0.18);
    border-top: 1px solid var(--border);
  }
  .save-feedback {
    font-size: 0.75rem;
    font-weight: 600;
    color: var(--moss);
    margin-right: auto;
  }
  .btn {
    font-size: 0.78rem;
    font-weight: 600;
    padding: 0.4rem 0.85rem;
    border-radius: var(--radius-sm);
    cursor: pointer;
    border: 1px solid transparent;
    transition: background 0.15s, border-color 0.15s;
  }
  .btn-secondary {
    background: transparent;
    color: var(--text-muted);
    border-color: var(--border);
  }
  .btn-secondary:hover {
    color: var(--text);
    background: rgba(255, 255, 255, 0.05);
  }
  .btn-primary {
    background: var(--blue);
    color: #fff;
  }
  .btn-primary:hover {
    background: var(--blue-bright);
  }
</style>
