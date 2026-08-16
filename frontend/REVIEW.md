# Code Review: `src/tsx/personal-goals/goalBrowser.tsx` and `goalBrowser.narrow.tsx`

Scope: both files above, plus the components they compose
(`goalCalendarPanel.tsx`, `goalPostPanel.tsx`) and the route setup, since some
issues only become visible across those boundaries.

---

## Bugs

### 1. `setInterval` with no delay argument (goalBrowser.tsx:82) — **high priority**

```tsx
const timer = setInterval(() => {
    const container = document.getElementById(articleContainerId);
    setArticleContainerWidth(container?.clientWidth || 0);
});
```

With no delay, browsers clamp the interval to ~4 ms, so this forces a layout
read (`clientWidth`) and calls `setState` **hundreds of times per second**. It
also keeps running in narrow mode, where the element doesn't even exist.

Fix with a `ResizeObserver` on a ref:

```tsx
const articleContainerRef = useRef<HTMLDivElement>(null);

useEffect(() => {
    const el = articleContainerRef.current;
    if (!el) return;
    const ro = new ResizeObserver(([e]) => setArticleContainerWidth(e.contentRect.width));
    ro.observe(el);
    return () => ro.disconnect();
}, [isFullMode()]);
```

This also lets you delete the random-ID hack
(`useState('id-' + Math.random().toString(16))` + `document.getElementById`) —
a `ref` (or `useId` if you keep a DOM lookup) is sufficient.

### 2. Stale calendar posts when navigating between goals (goalCalendarPanel.tsx:31) — **high priority**

```tsx
useEffect(() => {
    loadPosts();
}, [props.reload]);
```

The route is a single `/personal-goals/:id`, so navigating from goal A to goal B
does **not** remount `GoalCalendarPanel`. The parent's `loadGoal` effect
correctly re-runs on `goalId` change (so the title updates), but the calendar
keeps showing goal A's posts.

Minimal fix:

```tsx
useEffect(() => {
    loadPosts();
}, [props.id, props.reload]);
```

Related: the `reload: number` + `Math.random()` trigger pattern is what makes
this fragile. A `key={goalId}` on the panel from the parent would be cleaner
and also reset the panel's internal state.

### 3. Dead prop: `goalManagerMode`

`GoalPostPanel` accepts a `goalManagerMode` prop but never uses it — it reads
`context.goalManagerMode` instead (goalPostPanel.tsx:50). And
`context.goalManagerMode` already exists (context.tsx:8). So in
goalBrowser.tsx you can delete:

- the `Cookie` import and the `isGoalManagerMode()` helper,
- the `goalManagerMode={isGoalManagerMode()}` prop at the `GoalPostPanel` call site,
- the `goalManagerMode: boolean` field from `GoalPostPanel`'s props.

Note `isGoalManagerMode()` in goalBrowser is also used in the page-title
effect — switch that to `context.goalManagerMode` (already in scope).

---

## Minor issues

- **Radix inconsistency**: `parseInt(goalId, 10)` in `loadGoal` vs bare
  `parseInt(goalId)` / `parseInt(activePostDate)` elsewhere. Use `10`
  everywhere, or compute `const goalIdNum = Number(goalId)` once and reuse it
  (right now the id is re-parsed in several places).
- **`receivePosts` wipes the query string**: `setSearchParams({ activePostDate: ... })`
  replaces all params. Use the functional form to preserve others:
  ```tsx
  setSearchParams((p) => { p.set('activePostDate', '' + newActivePostDate); return p; }, { replace: true });
  ```
- **Uncleaned `setTimeout`** in the `activePostDate` effect (goalBrowser.tsx:57):
  no cleanup on unmount. The "arm the transition after first paint" intent is
  better served by `requestAnimationFrame` or a one-time mount effect. The same
  effect also handles "hide calendar when a post is selected" — two concerns
  coupled in one effect.
- **exhaustive-deps**: `loadGoal` is used in the effect but only `goalId` is
  listed. Inline the fetch or wrap in `useCallback`.
- **`calendarVisible` initialized from `isFullMode()`** only once — harmless
  because the wide layout ignores it, but the coupling is accidental. Consider
  defaulting to `true` and letting only narrow mode drive it.

---

## Refactoring suggestions

1. **Simplify `GoalBrowserNarrow` props**: render-prop functions
   (`getGoalCalendarPanel: () => ReactElement`) are passed, but the parent
   calls them immediately anyway — no laziness is gained. Pass plain elements
   instead: `calendar: ReactElement; postPanel: ReactElement;` (call
   `getGoalCalendarPanel()` at the parent call site).
2. **Duplicated "active post slot"**: both layouts render
   `activePostDate ? getGoalPostPanel() : undefined`. Extract a shared helper
   (e.g. `getActivePostPanel()` or a small `<PostSlot>` component) to remove
   the duplication.
3. **Compute `activePostDateMs` once**:
   ```tsx
   const activePostDateMs = activePostDate ? parseInt(activePostDate, 10) : 0;
   ```
   and pass it to both panels instead of `parseInt(...) || 0` in three places.
4. **Simplify `GoalCalendarPanel.isLoading` bookkeeping**
   (goalCalendarPanel.tsx:12-14): the `isLoadingRef.current + 1` counter with a
   ref is a stale-closure workaround. Since loads are sequential per effect
   trigger, plain `setIsLoading(true)` / `setIsLoading(false)` around the
   try/finally is equivalent and much simpler (and the `Math.random()` reload
   trigger can go away if you key by `goalId`).

---

## Priority

1. `setInterval` polling — real performance problem; fix first.
2. Stale posts on goal→goal navigation — real functional bug.
3. Dead `goalManagerMode` prop — free simplification, removes a js-cookie import.
