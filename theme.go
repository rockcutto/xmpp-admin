package main

import "net/http"

const material3CSS = `
:root {
  color-scheme: light;
  --m3-primary: #6750a4;
  --m3-on-primary: #ffffff;
  --m3-primary-container: #eaddff;
  --m3-on-primary-container: #21005d;
  --m3-secondary-container: #e8def8;
  --m3-on-secondary-container: #1d192b;
  --m3-error: #ba1a1a;
  --m3-error-container: #ffdad6;
  --m3-on-error-container: #410002;
  --m3-surface: #fffbfe;
  --m3-surface-container-lowest: #ffffff;
  --m3-surface-container-low: #f7f2fa;
  --m3-surface-container: #f3edf7;
  --m3-surface-container-high: #ece6f0;
  --m3-surface-container-highest: #e6e0e9;
  --m3-on-surface: #1d1b20;
  --m3-on-surface-variant: #49454f;
  --m3-outline: #79747e;
  --m3-outline-variant: #cac4d0;
  --m3-focus: rgba(103, 80, 164, .18);
  --m3-radius-sm: 12px;
  --m3-radius-md: 16px;
  --m3-radius-lg: 24px;
  --m3-radius-full: 999px;
}

html[data-theme="dark"] {
  color-scheme: dark;
  --m3-primary: #d0bcff;
  --m3-on-primary: #381e72;
  --m3-primary-container: #4f378b;
  --m3-on-primary-container: #eaddff;
  --m3-secondary-container: #4a4458;
  --m3-on-secondary-container: #e8def8;
  --m3-error: #ffb4ab;
  --m3-error-container: #93000a;
  --m3-on-error-container: #ffdad6;
  --m3-surface: #141218;
  --m3-surface-container-lowest: #0f0d13;
  --m3-surface-container-low: #1d1b20;
  --m3-surface-container: #211f26;
  --m3-surface-container-high: #2b2930;
  --m3-surface-container-highest: #36343b;
  --m3-on-surface: #e6e0e9;
  --m3-on-surface-variant: #cac4d0;
  --m3-outline: #938f99;
  --m3-outline-variant: #49454f;
  --m3-focus: rgba(208, 188, 255, .18);
}

@media (prefers-color-scheme: dark) {
  html:not([data-theme]) {
    color-scheme: dark;
    --m3-primary: #d0bcff;
    --m3-on-primary: #381e72;
    --m3-primary-container: #4f378b;
    --m3-on-primary-container: #eaddff;
    --m3-secondary-container: #4a4458;
    --m3-on-secondary-container: #e8def8;
    --m3-error: #ffb4ab;
    --m3-error-container: #93000a;
    --m3-on-error-container: #ffdad6;
    --m3-surface: #141218;
    --m3-surface-container-lowest: #0f0d13;
    --m3-surface-container-low: #1d1b20;
    --m3-surface-container: #211f26;
    --m3-surface-container-high: #2b2930;
    --m3-surface-container-highest: #36343b;
    --m3-on-surface: #e6e0e9;
    --m3-on-surface-variant: #cac4d0;
    --m3-outline: #938f99;
    --m3-outline-variant: #49454f;
    --m3-focus: rgba(208, 188, 255, .18);
  }
}

html, body {
  min-height: 100%;
  background: var(--m3-surface);
  color: var(--m3-on-surface);
}

body {
  margin: 0;
  font-family: Inter, Roboto, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  font-size: 14px;
  line-height: 1.45;
  letter-spacing: .01em;
}

a { color: var(--m3-primary); }
a:hover { color: var(--m3-primary); opacity: .86; }

.page {
  min-height: 100vh;
  background:
    radial-gradient(circle at 12% -8%, color-mix(in srgb, var(--m3-primary) 7%, transparent), transparent 30rem),
    var(--m3-surface);
}

.navbar {
  min-height: 64px;
  padding: 8px 0;
  background: color-mix(in srgb, var(--m3-surface-container) 94%, transparent) !important;
  border: 0 !important;
  border-bottom: 1px solid var(--m3-outline-variant) !important;
  box-shadow: none !important;
  backdrop-filter: blur(18px);
}

.navbar-brand {
  color: var(--m3-on-surface) !important;
  font-size: 20px;
  font-weight: 650;
  letter-spacing: -.01em;
  text-decoration: none;
}

.page-header { padding: 20px 0 12px; }
.page-title {
  margin: 0;
  color: var(--m3-on-surface);
  font-size: clamp(28px, 3vw, 36px);
  font-weight: 500;
  line-height: 1.15;
  letter-spacing: -.025em;
}
.page-header .text-secondary { margin-top: 6px; }
.page-subtitle { max-width: 78ch; line-height: 1.4; }
.page-body { padding-top: 0; padding-bottom: 48px; }

.text-secondary, .form-hint { color: var(--m3-on-surface-variant) !important; }
.form-hint { margin-top: 7px; font-size: 12px; }

.card {
  overflow: hidden;
  background: var(--m3-surface-container-low);
  color: var(--m3-on-surface);
  border: 1px solid var(--m3-outline-variant);
  border-radius: var(--m3-radius-lg);
  box-shadow: none;
}

.card-header {
  min-height: 60px;
  padding: 18px 22px;
  background: transparent;
  border-bottom: 1px solid var(--m3-outline-variant);
}

.card-title {
  color: var(--m3-on-surface);
  font-size: 16px;
  font-weight: 600;
}
.card-body { padding: 22px; }

.btn {
  min-height: 40px;
  text-decoration: none;
  padding: 9px 18px;
  border-radius: var(--m3-radius-full) !important;
  border-width: 1px;
  font-weight: 600;
  box-shadow: none !important;
  transition: background-color .16s ease, color .16s ease, border-color .16s ease, transform .08s ease;
}
.btn:active { transform: scale(.985); }

.btn-primary {
  background: var(--m3-primary) !important;
  border-color: var(--m3-primary) !important;
  color: var(--m3-on-primary) !important;
}

.btn-outline-secondary {
  background: transparent !important;
  border-color: var(--m3-outline) !important;
  color: var(--m3-primary) !important;
}
.btn-outline-secondary:hover {
  background: var(--m3-secondary-container) !important;
  border-color: transparent !important;
  color: var(--m3-on-secondary-container) !important;
}

.btn-outline-danger {
  background: transparent !important;
  border-color: var(--m3-error) !important;
  color: var(--m3-error) !important;
}
.btn-outline-danger:hover {
  background: var(--m3-error-container) !important;
  color: var(--m3-on-error-container) !important;
}

.btn-group {
  padding: 3px;
  border-radius: var(--m3-radius-full);
  background: var(--m3-surface-container-high);
  gap: 2px;
}
.btn-group .btn {
  min-width: 44px;
  min-height: 34px;
  padding: 6px 12px;
  border: 0 !important;
}

.form-label {
  margin-bottom: 8px;
  color: var(--m3-on-surface);
  font-size: 13px;
  font-weight: 600;
}

.form-control, .input-group-text {
  min-height: 48px;
  background: var(--m3-surface-container-lowest) !important;
  color: var(--m3-on-surface) !important;
  border-color: var(--m3-outline) !important;
}

.form-control {
  border-radius: var(--m3-radius-sm) !important;
  padding: 11px 14px;
}

.input-group .form-control {
  border-top-right-radius: 0 !important;
  border-bottom-right-radius: 0 !important;
}

.input-group-text {
  border-top-right-radius: var(--m3-radius-sm) !important;
  border-bottom-right-radius: var(--m3-radius-sm) !important;
  color: var(--m3-on-surface-variant) !important;
}

.form-control:focus {
  border-color: var(--m3-primary) !important;
  box-shadow: 0 0 0 3px var(--m3-focus) !important;
}

.table {
  --tblr-table-bg: transparent;
  --tblr-table-color: var(--m3-on-surface);
  margin-bottom: 0;
  color: var(--m3-on-surface);
}

.table thead th {
  padding: 13px 14px;
  background: var(--m3-surface-container);
  color: var(--m3-on-surface-variant);
  border-bottom: 1px solid var(--m3-outline-variant);
  font-size: 12px;
  font-weight: 650;
  letter-spacing: .035em;
}

.table tbody td {
  padding: 13px 14px;
  color: var(--m3-on-surface);
  border-bottom: 1px solid var(--m3-outline-variant);
  vertical-align: middle;
}

.table tbody tr:last-child td { border-bottom: 0; }
.table tbody tr:hover td { background: color-mix(in srgb, var(--m3-primary) 5%, transparent); }

.alert {
  margin-bottom: 14px;
  padding: 14px 18px;
  border: 0;
  border-radius: var(--m3-radius-md);
  box-shadow: none;
}
.alert-info {
  background: var(--m3-secondary-container);
  color: var(--m3-on-secondary-container);
}
.alert-success {
  background: color-mix(in srgb, #76d275 24%, var(--m3-surface-container));
  color: var(--m3-on-surface);
}
.alert-danger {
  background: var(--m3-error-container);
  color: var(--m3-on-error-container);
}

.context-note {
  margin-bottom: 14px;
  padding: 9px 12px;
  border-left: 3px solid var(--m3-outline-variant);
  color: var(--m3-on-surface-variant);
  font-size: 12px;
  line-height: 1.45;
}

.badge {
  padding: 6px 10px;
  border-radius: var(--m3-radius-full);
  font-size: 11px;
  font-weight: 650;
}
.bg-green-lt {
  background: color-mix(in srgb, #76d275 24%, var(--m3-surface-container-high)) !important;
  color: var(--m3-on-surface) !important;
}
.bg-red-lt {
  background: var(--m3-error-container) !important;
  color: var(--m3-on-error-container) !important;
}
.bg-secondary-lt {
  background: var(--m3-secondary-container) !important;
  color: var(--m3-on-secondary-container) !important;
}
.bg-warning-lt {
  background: color-mix(in srgb, #f9a825 26%, var(--m3-surface-container-high)) !important;
  color: var(--m3-on-surface) !important;
}

code {
  padding: 2px 6px;
  border-radius: 6px;
  background: var(--m3-surface-container-high);
  color: var(--m3-on-surface);
}

hr { border-color: var(--m3-outline-variant); opacity: 1; }
dl.row { row-gap: 10px; }
dt { color: var(--m3-on-surface-variant); font-weight: 550; }
dd { color: var(--m3-on-surface); }
dd code, .ops-break, .cell-secondary {
  overflow-wrap: anywhere;
  word-break: break-word;
}
.cell-primary {
  color: var(--m3-on-surface);
  font-weight: 650;
}
.cell-secondary {
  margin-top: 3px;
  color: var(--m3-on-surface-variant);
  font-size: 12px;
  line-height: 1.35;
}
.numeric-cell { font-variant-numeric: tabular-nums; }

.m3-toolbar-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.m3-theme-toggle {
  width: 42px;
  min-width: 42px;
  height: 42px;
  min-height: 42px;
  padding: 0;
  display: inline-grid;
  place-items: center;
  border: 0 !important;
  border-radius: 50% !important;
  background: var(--m3-surface-container-high) !important;
  color: var(--m3-on-surface) !important;
}
.m3-theme-toggle:hover {
  background: var(--m3-secondary-container) !important;
  color: var(--m3-on-secondary-container) !important;
}

.m3-theme-icon::before {
  content: "☾";
  font-size: 20px;
  line-height: 20px;
}
html[data-theme="dark"] .m3-theme-icon::before { content: "☀"; }

.table-responsive { border-radius: inherit; }

.invites-table { min-width: 860px; }
.invite-date time { font-variant-numeric: tabular-nums; }
.invite-jid { overflow-wrap: anywhere; }
.invite-type {
  display: inline-flex;
  max-width: none;
  white-space: nowrap;
}
.invite-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}
.invite-actions form { margin: 0; }

.ops-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}
.ops-table { min-width: 680px; }
.ops-break { overflow-wrap: anywhere; word-break: break-word; }

.metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}
.metric-card {
  min-width: 0;
  padding: 18px 20px;
  background: var(--m3-surface-container-low);
  border: 1px solid var(--m3-outline-variant);
  border-radius: var(--m3-radius-md);
}
.metric-label {
  color: var(--m3-on-surface-variant);
  font-size: 12px;
  font-weight: 600;
}
.metric-value {
  margin-top: 6px;
  color: var(--m3-on-surface);
  font-size: 28px;
  font-weight: 550;
  line-height: 1.15;
  font-variant-numeric: tabular-nums;
}
.metric-value-small { font-size: 20px; }

* { box-sizing: border-box; }

.container-xl {
  width: min(1320px, calc(100% - 32px));
  margin-left: auto;
  margin-right: auto;
}

.row {
  display: flex;
  flex-wrap: wrap;
  margin-left: -8px;
  margin-right: -8px;
}

.row > * {
  width: 100%;
  padding-left: 8px;
  padding-right: 8px;
}

.row-cards {
  row-gap: 16px;
}

.col-5 { width: 41.666667%; }
.col-7 { width: 58.333333%; }

@media (min-width: 768px) and (max-width: 1179.98px) {
  .col-lg-4 {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 16px;
  }
  .col-lg-4 > .card.mt-3 { margin-top: 0 !important; }
}

@media (min-width: 1180px) {
  .col-lg-4 { width: 33.333333%; }
  .col-lg-8 { width: 66.666667%; min-width: 0; }
}

.d-inline { display: inline !important; }
.d-flex { display: flex !important; }
.flex-row { flex-direction: row !important; }
.flex-nowrap { flex-wrap: nowrap !important; }
.align-items-center { align-items: center !important; }
.justify-content-end { justify-content: flex-end !important; }
.order-md-last { order: 6; }

.w-100 { width: 100% !important; }
.text-end { text-align: right !important; }
.text-center { text-align: center !important; }
.text-nowrap { white-space: nowrap !important; }
.small { font-size: .875em !important; }
.fw-bold { font-weight: 700 !important; }

.mt-1 { margin-top: .25rem !important; }
.mt-2 { margin-top: .5rem !important; }
.mt-3 { margin-top: 1rem !important; }
.mb-0 { margin-bottom: 0 !important; }
.mb-1 { margin-bottom: .25rem !important; }
.mb-3 { margin-bottom: 1rem !important; }
.py-5 { padding-top: 3rem !important; padding-bottom: 3rem !important; }

.navbar .container-xl {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.navbar-nav {
  display: flex;
  align-items: center;
}

.input-group {
  display: flex;
  width: 100%;
}

.input-group .form-control {
  min-width: 0;
  flex: 1 1 auto;
}

.input-group-text {
  display: flex;
  align-items: center;
  padding: 0 14px;
  border: 1px solid var(--m3-outline);
  border-left: 0;
}

.table-responsive {
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.table {
  width: 100%;
  border-collapse: collapse;
}

.btn-list {
  display: flex;
  gap: 8px;
  align-items: center;
}

.btn-sm {
  min-height: 34px;
  padding: 6px 12px;
  font-size: 12px;
}

.alert a {
  overflow-wrap: anywhere;
}

@media (max-width: 991.98px) {
  .metric-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .page-header { padding-top: 18px; }
  .card { border-radius: var(--m3-radius-md); }
  .card-body, .card-header { padding-left: 18px; padding-right: 18px; }
}

@media (max-width: 767.98px) {
  .container-xl {
    width: calc(100% - 20px);
  }

  .navbar {
    min-height: 58px;
    padding: 6px 0;
  }
  .navbar .container-xl { gap: 8px; }
  .navbar-brand { font-size: 18px; }
  .m3-toolbar-actions { gap: 6px; }
  .m3-theme-toggle {
    width: 38px;
    min-width: 38px;
    height: 38px;
    min-height: 38px;
  }
  .btn-group { padding: 2px; }
  .btn-group .btn {
    min-width: 38px;
    min-height: 32px;
    padding: 5px 9px;
  }

  .page-header { padding: 16px 0 10px; }
  .page-title { font-size: 27px; }
  .page-subtitle { font-size: 13px; }
  .page-body { padding-bottom: 28px; }

  .card { border-radius: 14px; }
  .card-header {
    min-height: 54px;
    padding: 15px 16px;
  }
  .card-body { padding: 16px; }

  .input-group-text {
    max-width: 55%;
    padding-left: 10px;
    padding-right: 10px;
    overflow-wrap: anywhere;
    white-space: normal;
    line-height: 1.2;
  }

  .card-body dl.row {
    display: grid;
    grid-template-columns: minmax(110px, 42%) minmax(0, 1fr);
    column-gap: 12px;
    row-gap: 10px;
    margin-left: 0;
    margin-right: 0;
  }
  .card-body dl.row > dt,
  .card-body dl.row > dd {
    width: auto;
    min-width: 0;
    padding: 0;
    margin: 0;
  }

  .metric-grid { grid-template-columns: 1fr; gap: 10px; }
  .metric-card { padding: 15px 16px; }
  .metric-value { font-size: 25px; }
  .metric-value-small { font-size: 19px; }

  .table-responsive { overflow-x: visible; }
  .invites-table,
  .ops-table { min-width: 0; }

  .responsive-data-table thead { display: none; }
  .responsive-data-table,
  .responsive-data-table tbody,
  .responsive-data-table tr,
  .responsive-data-table td {
    display: block;
    width: 100%;
  }
  .responsive-data-table tbody tr {
    padding: 9px 14px 10px;
    border-bottom: 1px solid var(--m3-outline-variant);
  }
  .responsive-data-table tbody tr:last-child { border-bottom: 0; }
  .responsive-data-table tbody td {
    display: grid;
    grid-template-columns: minmax(94px, 36%) minmax(0, 1fr);
    align-items: start;
    gap: 12px;
    padding: 7px 0;
    border: 0;
    text-align: left !important;
    white-space: normal !important;
    background: transparent !important;
  }
  .responsive-data-table tbody td::before {
    content: attr(data-label);
    color: var(--m3-on-surface-variant);
    font-size: 11px;
    font-weight: 650;
    line-height: 1.35;
  }
  .responsive-data-table tbody td > * {
    grid-column: 2;
  }
  .responsive-data-table .table-actions-cell {
    display: block;
    padding-top: 10px;
  }
  .responsive-data-table .table-actions-cell::before { display: none; }
  .responsive-data-table .empty-row {
    padding: 0;
    border: 0;
  }
  .responsive-data-table .empty-row td {
    display: block;
    padding: 32px 16px !important;
    text-align: center !important;
  }
  .responsive-data-table .empty-row td::before { display: none; }
  .responsive-data-table .cell-primary { font-size: 14px; }
  .responsive-data-table .cell-secondary { font-size: 11px; }

  .invite-actions {
    justify-content: flex-start;
    flex-wrap: wrap;
  }
  .invite-actions .btn { min-height: 38px; }
  .invite-type { white-space: normal; }

  .health-table tbody td:last-child {
    padding-bottom: 9px;
  }
}

/* Material 3 app shell: compact desktop navigation rail, bottom navigation on phones. */
.nav-rail {
  position: fixed;
  inset: 0 auto 0 0;
  z-index: 40;
  width: 72px;
  padding: 10px 8px 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 14px;
  background: var(--m3-surface-container-lowest);
  border-right: 1px solid var(--m3-outline-variant);
}

.nav-rail-brand {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  flex: 0 0 auto;
  border: 1px solid var(--m3-outline-variant);
  border-radius: 15px;
  background: var(--m3-surface-container-high);
  color: var(--m3-on-surface);
  text-decoration: none;
  font-size: 15px;
  font-weight: 750;
  letter-spacing: -.03em;
}

.nav-rail-brand:hover {
  color: var(--m3-on-surface);
  background: var(--m3-surface-container-highest);
  opacity: 1;
}

.nav-rail-items {
  width: 100%;
  display: flex;
  flex: 1 1 auto;
  flex-direction: column;
  align-items: center;
  gap: 7px;
}

.nav-rail-link {
  width: 56px;
  min-height: 48px;
  padding: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--m3-on-surface-variant);
  border-radius: 17px;
  text-decoration: none;
  transition: background-color .16s ease, color .16s ease, transform .08s ease;
}

.nav-rail-link:hover {
  color: var(--m3-on-surface);
  background: var(--m3-surface-container);
  opacity: 1;
}

.nav-rail-link:active {
  transform: scale(.97);
}

.nav-rail-link.active {
  color: var(--m3-on-surface);
  background: transparent;
}

.nav-rail-icon {
  width: 40px;
  height: 40px;
  display: grid;
  place-items: center;
  border-radius: 14px;
}

.nav-rail-link.active .nav-rail-icon {
  width: 48px;
  height: 36px;
  border-radius: var(--m3-radius-full);
  background: var(--m3-primary-container);
  color: var(--m3-on-primary-container);
}

.nav-rail-icon svg {
  width: 22px;
  height: 22px;
  fill: none;
  stroke: currentColor;
  stroke-width: 1.8;
  stroke-linecap: round;
  stroke-linejoin: round;
}

.nav-rail-label {
  display: none;
}

.app-main {
  min-width: 0;
  min-height: 100vh;
  margin-left: 72px;
}

.app-main > .navbar {
  position: sticky;
  top: 0;
  z-index: 30;
  background: color-mix(in srgb, var(--m3-surface) 88%, transparent) !important;
  backdrop-filter: blur(18px);
}

@media (max-width: 767.98px) {
  .app-main {
    margin-left: 0;
    padding-bottom: calc(74px + env(safe-area-inset-bottom));
  }

  .nav-rail {
    inset: auto 0 0 0;
    width: 100%;
    height: calc(68px + env(safe-area-inset-bottom));
    padding: 6px 6px calc(5px + env(safe-area-inset-bottom));
    flex-direction: row;
    background: color-mix(in srgb, var(--m3-surface-container) 96%, transparent);
    border-right: 0;
    border-top: 1px solid var(--m3-outline-variant);
    box-shadow: 0 -8px 24px color-mix(in srgb, #000 8%, transparent);
    backdrop-filter: blur(20px);
  }

  .nav-rail-brand {
    display: none;
  }

  .nav-rail-items {
    height: 100%;
    flex-direction: row;
    justify-content: stretch;
    gap: 2px;
  }

  .nav-rail-link {
    min-width: 0;
    width: auto;
    min-height: 56px;
    padding: 3px 2px 2px;
    flex: 1 1 20%;
    flex-direction: column;
    justify-content: center;
    gap: 2px;
    border-radius: 14px;
  }

  .nav-rail-link:hover {
    background: transparent;
  }

  .nav-rail-link.active {
    background: transparent;
  }

  .nav-rail-icon {
    width: 38px;
    height: 30px;
    border-radius: 16px;
  }

  .nav-rail-link.active .nav-rail-icon {
    width: 48px;
    background: var(--m3-primary-container);
    color: var(--m3-on-primary-container);
  }

  .nav-rail-icon svg {
    width: 20px;
    height: 20px;
  }

  .nav-rail-label {
    display: block;
    max-width: 100%;
    color: var(--m3-on-surface-variant);
    font-size: 9px;
    font-weight: 600;
    line-height: 1.1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .nav-rail-link.active .nav-rail-label {
    color: var(--m3-on-surface);
    font-weight: 700;
  }

  .app-main > .navbar {
    position: sticky;
    top: 0;
  }
}

/* Task-first administration surfaces. */
.sr-only {
  position: absolute !important;
  width: 1px !important;
  height: 1px !important;
  padding: 0 !important;
  margin: -1px !important;
  overflow: hidden !important;
  clip: rect(0, 0, 0, 0) !important;
  white-space: nowrap !important;
  border: 0 !important;
}

.page-header {
  padding-top: 26px;
  padding-bottom: 16px;
}

.section-title {
  margin: 0 0 10px;
  color: var(--m3-on-surface);
  font-size: 17px;
  font-weight: 650;
}

.invite-workspace {
  display: grid;
  grid-template-columns: minmax(300px, 360px) minmax(0, 1fr);
  gap: 18px;
  align-items: start;
}

.invite-sidebar {
  display: grid;
  gap: 14px;
  min-width: 0;
}

.invite-feed-section {
  min-width: 0;
}

.list-toolbar {
  min-height: 44px;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  gap: 12px;
  justify-content: space-between;
}

.segmented-control {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 3px;
  border-radius: var(--m3-radius-full);
  background: var(--m3-surface-container-high);
}

.segment {
  min-height: 34px;
  padding: 6px 13px;
  border: 0;
  border-radius: var(--m3-radius-full);
  background: transparent;
  color: var(--m3-on-surface-variant);
  font: inherit;
  font-size: 12px;
  font-weight: 650;
  cursor: pointer;
}

.segment:hover {
  color: var(--m3-on-surface);
}

.segment.active {
  background: var(--m3-primary-container);
  color: var(--m3-on-primary-container);
}

.invite-feed {
  display: grid;
  gap: 10px;
}

.invite-card {
  position: relative;
  overflow: hidden;
  border: 1px solid var(--m3-outline-variant);
  border-radius: 18px;
  background: var(--m3-surface-container-low);
}

.invite-card.just-created {
  border-color: color-mix(in srgb, var(--m3-primary) 68%, var(--m3-outline-variant));
  box-shadow: 0 0 0 2px color-mix(in srgb, var(--m3-primary) 12%, transparent);
}

.invite-card[data-invite-state="expired"] {
  background: color-mix(in srgb, var(--m3-surface-container-low) 72%, var(--m3-surface));
}

.invite-card[data-invite-state="expired"] .invite-card-title {
  color: var(--m3-on-surface-variant);
  font-weight: 600;
}

.invite-ready-label {
  padding: 8px 16px;
  background: var(--m3-primary-container);
  color: var(--m3-on-primary-container);
  font-size: 11px;
  font-weight: 750;
  letter-spacing: .01em;
}

.invite-card-main {
  padding: 15px 16px 12px;
}

.invite-card-heading {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
}

.invite-card-title {
  color: var(--m3-on-surface);
  font-size: 15px;
  font-weight: 700;
}

.invite-card-account {
  margin-top: 4px;
  color: var(--m3-on-surface-variant);
  font-size: 12px;
}

.invite-meta-grid {
  margin-top: 13px;
  display: grid;
  grid-template-columns: minmax(0, 1.3fr) minmax(130px, .85fr) minmax(130px, .85fr);
  gap: 10px 16px;
}

.invite-meta-item {
  min-width: 0;
  display: grid;
  gap: 3px;
}

.invite-meta-label,
.config-label {
  color: var(--m3-on-surface-variant);
  font-size: 10px;
  font-weight: 650;
  letter-spacing: .02em;
}

.invite-meta-value {
  min-width: 0;
  color: var(--m3-on-surface);
  font-size: 12px;
  line-height: 1.35;
}

.invite-card-actions {
  padding: 10px 12px 12px;
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  border-top: 1px solid var(--m3-outline-variant);
}

.invite-card-actions form {
  margin: 0;
}

.config-summary-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 9px;
}

.config-summary-item {
  min-width: 0;
  padding: 10px 11px;
  border-radius: 12px;
  background: var(--m3-surface-container);
}

.config-value {
  margin-top: 5px;
  min-width: 0;
  color: var(--m3-on-surface);
  font-size: 13px;
  font-weight: 600;
  overflow-wrap: anywhere;
}

.technical-details {
  margin-top: 13px;
  border-top: 1px solid var(--m3-outline-variant);
  padding-top: 10px;
}

.technical-details > summary,
.inline-technical > summary {
  color: var(--m3-on-surface-variant);
  font-size: 11px;
  font-weight: 650;
  cursor: pointer;
  user-select: none;
}

.technical-grid {
  margin: 10px 0 0;
  display: grid;
  grid-template-columns: minmax(105px, .8fr) minmax(0, 1.2fr);
  gap: 7px 10px;
  font-size: 11px;
}

.technical-grid dt,
.technical-grid dd {
  min-width: 0;
  margin: 0;
}

.technical-grid dd {
  overflow-wrap: anywhere;
}

.data-section {
  min-width: 0;
}

.data-toolbar {
  justify-content: flex-start;
}

.search-box {
  width: min(520px, 100%);
  min-height: 42px;
  padding: 0 13px;
  display: flex;
  align-items: center;
  gap: 9px;
  border: 1px solid var(--m3-outline-variant);
  border-radius: var(--m3-radius-full);
  background: var(--m3-surface-container-lowest);
}

.search-box:focus-within {
  border-color: var(--m3-primary);
  box-shadow: 0 0 0 3px var(--m3-focus);
}

.search-box svg {
  width: 18px;
  height: 18px;
  flex: 0 0 auto;
  fill: none;
  stroke: var(--m3-on-surface-variant);
  stroke-width: 1.8;
  stroke-linecap: round;
}

.search-box input {
  width: 100%;
  min-width: 0;
  border: 0;
  outline: 0;
  background: transparent;
  color: var(--m3-on-surface);
  font: inherit;
  font-size: 13px;
}

.search-box input::placeholder {
  color: var(--m3-on-surface-variant);
}

.count-chip {
  min-width: 30px;
  height: 30px;
  padding: 0 9px;
  display: inline-grid;
  place-items: center;
  border-radius: var(--m3-radius-full);
  background: var(--m3-surface-container-high);
  color: var(--m3-on-surface-variant);
  font-size: 11px;
  font-weight: 700;
}

.compact-data-card {
  border-radius: 18px;
}

.compact-data-card .table thead th {
  padding-top: 10px;
  padding-bottom: 10px;
}

.compact-data-card .table tbody td {
  padding-top: 10px;
  padding-bottom: 10px;
}

.empty-state {
  padding: 48px 20px;
  text-align: center;
  color: var(--m3-on-surface-variant);
}

.empty-state.compact {
  padding: 28px 18px;
}

.empty-state-title {
  color: var(--m3-on-surface);
  font-size: 14px;
  font-weight: 700;
}

.empty-state-body {
  margin-top: 5px;
  font-size: 12px;
  line-height: 1.45;
}

.filter-empty[hidden] {
  display: none !important;
}

.health-section {
  margin-top: 20px;
}

.health-list {
  overflow: hidden;
  border: 1px solid var(--m3-outline-variant);
  border-radius: 18px;
  background: var(--m3-surface-container-low);
}

.health-row {
  padding: 13px 15px;
  border-bottom: 1px solid var(--m3-outline-variant);
}

.health-row:last-child {
  border-bottom: 0;
}

.health-row.health-error {
  box-shadow: inset 3px 0 0 var(--m3-error);
}

.health-row.health-warn {
  box-shadow: inset 3px 0 0 color-mix(in srgb, #f9a825 78%, var(--m3-outline));
}

.health-row-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.health-name {
  color: var(--m3-on-surface);
  font-size: 13px;
  font-weight: 700;
}

.health-summary {
  margin-top: 4px;
  color: var(--m3-on-surface-variant);
  font-size: 12px;
  line-height: 1.4;
}

.inline-technical {
  margin-top: 7px;
}

.inline-technical code {
  display: block;
  margin-top: 6px;
  padding: 7px 9px;
  font-size: 10px;
  line-height: 1.4;
}

.task-error .inline-technical {
  margin-top: 8px;
}

@media (max-width: 1099.98px) {
  .invite-workspace {
    display: flex;
    flex-direction: column;
  }

  .invite-sidebar {
    display: contents;
  }

  .create-invite-card { order: 1; width: 100%; }
  .invite-feed-section { order: 2; width: 100%; }
  .invitation-settings-card { order: 3; width: 100%; }

  .config-summary-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 767.98px) {
  .page-header {
    padding-top: 18px;
    padding-bottom: 12px;
  }

  .list-toolbar {
    align-items: stretch;
    flex-direction: column;
  }

  .segmented-control {
    width: 100%;
  }

  .segment {
    flex: 1 1 0;
    padding-left: 8px;
    padding-right: 8px;
  }

  .invite-meta-grid {
    grid-template-columns: 1fr;
    gap: 8px;
  }

  .invite-card-actions .btn {
    min-height: 38px;
  }

  .config-summary-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .technical-grid {
    grid-template-columns: 1fr;
    gap: 3px;
  }

  .technical-grid dd {
    margin-bottom: 7px;
  }

  .data-toolbar {
    flex-direction: row;
    align-items: center;
  }

  .search-box {
    min-width: 0;
  }

  .health-row {
    padding: 12px 13px;
  }

  .nav-rail-link.active .nav-rail-icon {
    width: 48px;
    height: 30px;
  }
}

@media (max-width: 420px) {
  .config-summary-grid {
    grid-template-columns: 1fr;
  }

  .invite-card-heading {
    gap: 8px;
  }

  .invite-card-actions {
    align-items: stretch;
  }

  .invite-card-actions .btn,
  .invite-card-actions form,
  .invite-card-actions form .btn {
    width: 100%;
  }
}

@media print {
  .d-print-none { display: none !important; }
}
`

const themeJS = `
(() => {
  const key = "xmpp-admin-theme";
  const root = document.documentElement;
  const system = window.matchMedia("(prefers-color-scheme: dark)");

  function storedTheme() {
    try {
      const value = localStorage.getItem(key);
      return value === "light" || value === "dark" ? value : null;
    } catch (_) {
      return null;
    }
  }

  function effectiveTheme() {
    return storedTheme() || (system.matches ? "dark" : "light");
  }

  function apply(theme) {
    root.dataset.theme = theme;
    const button = document.querySelector("[data-theme-toggle]");
    if (button) {
      const next = theme === "dark" ? "light" : "dark";
      const label = next === "dark" ? button.dataset.labelDark : button.dataset.labelLight;
      button.setAttribute("aria-label", label);
      button.setAttribute("title", label);
    }
  }

  function formatHumanTimes() {
    const lang = (root.lang || "en").toLowerCase();
    const locale = lang.startsWith("ru") ? "ru-RU" : "en-GB";
    const now = new Date();
    const timeFormat = new Intl.DateTimeFormat(locale, { hour: "2-digit", minute: "2-digit" });
    const dateFormat = new Intl.DateTimeFormat(locale, { day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" });
    const relative = new Intl.RelativeTimeFormat(locale, { numeric: "auto" });

    const dayNumber = (date) => Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()) / 86400000;

    document.querySelectorAll("time[data-human-time]").forEach((el) => {
      const date = new Date(el.getAttribute("datetime"));
      if (Number.isNaN(date.getTime())) return;
      const days = Math.round(dayNumber(date) - dayNumber(now));
      const kind = el.dataset.humanTime;

      if (kind === "created") {
        if (days === 0) {
          el.textContent = (lang.startsWith("ru") ? "Сегодня, " : "Today, ") + timeFormat.format(date);
        } else if (days === -1) {
          el.textContent = (lang.startsWith("ru") ? "Вчера, " : "Yesterday, ") + timeFormat.format(date);
        } else {
          el.textContent = dateFormat.format(date);
        }
        return;
      }

      if (kind === "expiry" && Math.abs(days) <= 14) {
        el.textContent = relative.format(days, "day") + " · " + timeFormat.format(date);
      } else {
        el.textContent = dateFormat.format(date);
      }
    });
  }

  async function copyText(value) {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(value);
      return;
    }
    const area = document.createElement("textarea");
    area.value = value;
    area.setAttribute("readonly", "");
    area.style.position = "fixed";
    area.style.opacity = "0";
    document.body.appendChild(area);
    area.select();
    document.execCommand("copy");
    area.remove();
  }

  function wireCopyButtons() {
    document.querySelectorAll("[data-copy-value]").forEach((button) => {
      button.addEventListener("click", async () => {
        const value = button.dataset.copyValue || "";
        if (!value) return;
        try {
          await copyText(value);
          const original = button.dataset.copyLabel || button.textContent;
          button.textContent = button.dataset.copiedLabel || original;
          window.setTimeout(() => { button.textContent = original; }, 1400);
        } catch (_) {
          // Keep the original label if the browser blocks clipboard access.
        }
      });
    });
  }

  function wireInviteFilters() {
    const feed = document.querySelector("[data-invite-feed]");
    if (!feed) return;
    const buttons = document.querySelectorAll("[data-invite-filter]");
    const empty = feed.querySelector("[data-invite-filter-empty]");
    const cards = feed.querySelectorAll(".invite-card[data-invite-state]");
    if (cards.length === 0) return;

    const applyFilter = (filter) => {
      let visible = 0;
      cards.forEach((card) => {
        const show = filter === "all" || card.dataset.inviteState === filter;
        card.hidden = !show;
        if (show) visible += 1;
      });
      if (empty) empty.hidden = visible !== 0;
      buttons.forEach((button) => button.classList.toggle("active", button.dataset.inviteFilter === filter));
    };

    buttons.forEach((button) => {
      button.addEventListener("click", () => applyFilter(button.dataset.inviteFilter || "all"));
    });
  }

  function wireListFilters() {
    document.querySelectorAll("[data-list-filter]").forEach((input) => {
      const id = input.dataset.listFilter;
      const list = document.getElementById(id);
      const empty = document.querySelector('[data-filter-empty="' + id + '"]');
      if (!list) return;
      const rows = list.querySelectorAll("[data-filter-row]");
      if (rows.length === 0) return;

      const run = () => {
        const query = input.value.trim().toLocaleLowerCase();
        let visible = 0;
        rows.forEach((row) => {
          const haystack = (row.dataset.filterText || row.textContent || "").toLocaleLowerCase();
          const show = query === "" || haystack.includes(query);
          row.hidden = !show;
          if (show) visible += 1;
        });
        if (empty) empty.hidden = query === "" || visible !== 0;
      };

      input.addEventListener("input", run);
    });
  }

  apply(effectiveTheme());

  document.addEventListener("DOMContentLoaded", () => {
    apply(effectiveTheme());
    formatHumanTimes();
    wireCopyButtons();
    wireInviteFilters();
    wireListFilters();

    const created = document.querySelector(".invite-card.just-created");
    if (created) {
      const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
      window.setTimeout(() => {
        created.scrollIntoView({ block: "center", behavior: reducedMotion ? "auto" : "smooth" });
      }, 60);
      try {
        const current = new URL(window.location.href);
        if (current.searchParams.get("created") === "1") {
          current.searchParams.delete("created");
          history.replaceState(null, "", current.pathname + current.search + current.hash);
        }
      } catch (_) {
        // URL cleanup is cosmetic only.
      }
    }

    const button = document.querySelector("[data-theme-toggle]");
    if (button) {
      button.addEventListener("click", () => {
        const next = effectiveTheme() === "dark" ? "light" : "dark";
        try {
          localStorage.setItem(key, next);
        } catch (_) {
          // Storage can be disabled by browser policy; switching still works for this page.
        }
        apply(next);
      });
    }
  });

  if (system.addEventListener) {
    system.addEventListener("change", () => {
      if (!storedTheme()) apply(effectiveTheme());
    });
  }
})();
`

func (a *App) material3CSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write([]byte(material3CSS))
}

func (a *App) themeJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write([]byte(themeJS))
}
