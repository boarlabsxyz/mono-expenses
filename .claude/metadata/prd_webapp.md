# Web Application Product Requirements Document

## Functional Requirements Checklist

### 1. Purpose, Users & Roles (RBAC)
- **Application Purpose**: What is the core value proposition and primary purpose of this web application?
- **Target Users**: Who are the primary user types (internal staff, external customers, partners, admins)? Any guest/anonymous access?
- **Problems Solved**: What specific business problems does it solve for each user type?
- **Core Workflows**: What are the main user journeys and use cases for each persona?
- **Integration Context**: How does it fit into users' existing business processes and other systems?
- **Role-Based Access**: What can each role view/create/edit/approve/delete/export/admin?
- **Action Scope**: Scope of actions (own data, team resources, organization-wide, tenant-specific)
- **Security Controls**: Dual-approval requirements, MFA needs, audit trail requirements
- **User Lifecycle**: User provisioning, role changes, suspension, data retention on deletion (invite → active → suspended → transfer → delete)

### 2. Information Architecture & Navigation
- Primary sections and hierarchy (IA sitemap)
- Global nav items, breadcrumbs, search entry points
- Deep links and shareable URLs (preserve filters/sort)

### 3. Screens, Workflows & States
- Critical workflows (rank 1–3) with happy path and edge paths
- Page templates needed (dashboards, list/detail, wizard, settings)
- Required empty/loading/error/no-permission states

### 4. Forms, Inputs & Validation
- Field rules (required, formats, masks, uniqueness)
- Inline vs toast errors; autosave/drafts; multi-step wizards; undo
- Calculations/derived fields; time zone/rounding rules

### 5. Domain Model & Relationships
- Core entities, fields, relationships (1-N/M-N), cascades on delete/merge
- State machines per entity (e.g., Draft → Pending → Approved). Who can transition?

### 6. Search, Filter, Sort & Saved Views
- Searchable fields (exact vs fuzzy), operators, default sort
- Filters, saved views (per-user vs shared), column chooser, pagination/virtualization
- Exports: formats, columns, row limits; scheduled exports?

### 7. Real-time & Collaboration
- Live updates (polling, WebSockets, SSE)
- Presence/typing indicators, record locking vs optimistic concurrency, comments/mentions

### 8. Files & Uploads
- Allowed types, max size, chunked uploads/resume, virus scanning
- Thumbnails/previews; metadata extraction; storage locations

### 9. Notifications
- Triggers → channels (in-app/toast, banner, email, push, webhook)
- Throttling/digests; per-user preferences; templates/localization

### 10. Auth, Sessions & Identity
- Sign-in methods (passwordless, email+password, SSO: OIDC/SAML)
- MFA, device/session management, session timeout/remember-me
- Multi-account/tenant switching; domain capture/custom domains

### 11. Permissions, Privacy & Audit
- Field-level visibility/masking (PII), row-level scoping
- Admin/back-office access rules; audit log scope & retention

### 12. URLs, Routing & Sharing
- Route patterns, 404/403/401 handling; guard protected routes
- Linkable UI state (filters/tabs) and expiring share links

### 13. UI/UX Standards
- Component library, spacing/typography rules, modals/sheets, toasts, tables
- Keyboard shortcuts, focus management, skip links; copy tone & microcopy

### 14. Responsive, Accessibility & Internationalization
- Breakpoints and layout changes; mobile gestures
- WCAG 2.1 AA targets (focus order, labels, contrast, reduced motion)
- Languages, RTL, number/date/currency formats, time zones

### 15. Performance & Resilience (Functional NFRs)
- Performance budgets (LCP/INP targets), lazy-load, code-split, caching
- Offline/PWA? What works vs not; retry/backoff, idempotency keys
- Max list sizes; concurrency rules/versioning

### 16. Integrations & APIs
- External systems; data contracts (endpoints, payloads, auth, rate limits)
- Webhooks: events, retries, signatures
- Imports/sync: mapping, dedupe, conflict resolution

### 17. Security Behaviors (User-visible)
- CSRF protection, XSS sanitization (rich text), CSP/CORS constraints
- Download headers for exports; redaction in logs; session fixation prevention
- Destructive actions: confirm + type-to-confirm + require role

### 18. Feature Flags, Rollout & Migration
- Flag gating by user/role/tenant; seed data; backfill/migration steps; rollback plan
- Versioned experiments and kill switches

### 19. Analytics & Consent
- Events/properties, funnels, dashboards; PII guardrails
- Consent banner/cookies by region; Do-Not-Track behavior

### 20. Error Handling & Support
- Error taxonomy and messages (user-fixable vs system)
- Maintenance mode, graceful degradation, support links

### 21. Browser Support & Packaging
- Supported browsers/versions; progressive enhancement rules
- Asset/CDN strategy; self-host vs cloud controls

### 22. Testing & CI (Functional Expectations)
- Golden screenshots/visual diff, a11y tests, contract tests for APIs
- Example user journeys for smoke tests
