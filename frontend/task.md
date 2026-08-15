# Rest object update

In file `src\tsx\personal-goals\goalCalendarPanel.tsx` instead of using `fetch` from JS standard library,
use ApiClient defined in `src\typescript\apiClient.ts`. ApiClient should return `GoalPostHeaderWithMethods` directly, without the need to convert objects in goalCalendarPanel.
