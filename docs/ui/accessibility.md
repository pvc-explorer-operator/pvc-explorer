# UI Accessibility

This document covers the accessibility decisions and patterns implemented across the Vue 3 UI. It is intended for contributors who need to understand or extend these patterns.

## Landmark structure

Every view has a `<main>` element and an `<h1>`. Where the visual design already communicates the page title (e.g. the sidebar breadcrumb), the `<h1>` uses `class="visually-hidden"` so it is present for screen readers but not rendered visually.

`AppLayout.vue` includes a skip-to-main link as its first child — it becomes visible on keyboard focus:

```html
<a href="#main-content" class="skip-link">Skip to content</a>
```

The `.visually-hidden` utility class (and its focusable variant) lives in `ui/src/style.css`.

## Focus management

`:focus-visible` outlines are applied globally in `ui/src/layout/layout.css`:

```css
:focus-visible {
  outline: 3px solid var(--p-primary-500);
  outline-offset: 2px;
}
```

This uses `:focus-visible` rather than `:focus` so mouse users are not affected.

When the mobile sidebar overlay is open, the main content area receives `:inert` to prevent background keyboard focus:

```html
<div class="layout-main-container" :inert="layoutState.mobileMenuActive">
```

## Font sizing

`html { font-size: 87.5% }` is set in `layout.css` instead of `font-size: 14px`. This preserves the user's browser font size preference (WCAG 1.4.4). All spacing tokens use `rem` so the entire UI scales proportionally.

## Contrast

`prefers-contrast: more` is handled in `ui/src/theme/variables.css`:

- Light mode: `--surface-border` → `var(--p-surface-400)`, `--text-color-secondary` → `var(--p-surface-700)`
- Dark mode: `--surface-border` → `var(--p-surface-500)`, `--text-color-secondary` → `var(--p-surface-300)`

PrimeVue preset overrides in `ui/src/main.ts` set `text.color` to `{surface.900}` and `text.mutedColor` to `{surface.600}` for light mode to meet contrast ratios.

## Interactive elements

**Icon-only buttons** — All buttons with no visible text label carry `aria-label` matching their `title`. See `AppTopbar.vue`, `HomeView.vue`, `ScopeListView.vue`.

**Clickable cards** — `AppCard.vue` and `ScopeCard.vue` use `<a>` wrappers with `href` and `@click.prevent`. The status dot carries `aria-hidden="true"`.

**Active nav item** — `AppMenuItem.vue` sets `:aria-current="route.path === item.to ? 'page' : undefined"` on the active `<router-link>`.

## Form controls

All `<label>` elements are linked to their inputs via `for`/`id` pairs. `CreateScopeView.vue` uses a `<fieldset>` + `<legend>` for the label-selector chip group.

`LabelAutocomplete.vue` implements the full ARIA combobox pattern: `role="combobox"`, `aria-autocomplete="list"`, `:aria-expanded`, `:aria-controls`, `:aria-activedescendant`, `<ul role="listbox">`, `<li role="option">`.

## Modals and dialogs

`FileExplorer.vue` modals, `KeyboardShortcutsModal.vue`, and `SearchDialog.vue` use native `<dialog>` with `showModal()`. This gives built-in focus trapping, Escape key handling, top-layer rendering, and `::backdrop` support without JavaScript.

`WakeUpDialog.vue` uses PrimeVue `<Dialog>` (ARIA-compliant, built-in focus trap).

`prefers-reduced-motion: reduce` guards are applied to all dialog CSS transitions.

## Motion

All `<TransitionGroup>` card animations and scroll-driven section reveals are wrapped in:

```css
@media (prefers-reduced-motion: no-preference) { … }
```

Scroll-driven reveals (`animation-timeline: view()`) additionally require an `@supports` guard for Firefox, which renders all sections statically as an acceptable degradation.

## Page titles

`router.afterEach` in `src/router/index.ts` sets `document.title` from `route.meta.title` so each view announces itself to screen readers and browser history correctly.
