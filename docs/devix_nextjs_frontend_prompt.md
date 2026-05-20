# Devix Next.js Frontend Build Spec

This document is the full, page-by-page frontend spec for a Next.js web app that connects to the Devix backend in this repo. It is aligned with the backend routes in [internal/router/router.go](internal/router/router.go) and module route files, plus the product details in [devix_react_native_prompt.md](devix_react_native_prompt.md).

Use this as the single source of truth when building the web UI. Every page below lists the API endpoints, required UI sections, and state handling so the frontend connects cleanly to the backend.

---

## Base URLs and Auth

- REST base URL: `https://your-api.com/api/v1`
- WebSocket: `wss://your-api.com/ws?token={access_token}`
- Authorization header: `Authorization: Bearer <access_token>`

### Token flow
- `POST /auth/login` and `POST /auth/signup` return `{ user, tokens }`.
- Tokens use rotation; access token expires quickly (15 minutes), refresh token lasts 7 days.
- When any request returns 401, call `POST /auth/refresh` with `{ refresh_token }`, then retry.
- On logout call `POST /auth/logout` with `{ refresh_token }`, then clear local tokens.

Storage guidance (match backend contract):
- Store `access_token` and `refresh_token` in client storage.
- If you implement a BFF or cookie approach, ensure `refresh_token` is still sent in request body.

---

## Required App Structure (Next.js App Router)

Suggested folder layout:

```
app/
  (public)/
    login/
    signup/
    forgot-password/
    reset-password/
  (app)/
    layout.tsx
    page.tsx
    explore/
    following/
    post/
      new/
      [id]/
      [id]/edit/
      [id]/analytics/
    drafts/
    search/
    tag/
      [slug]/
    chat/
      page.tsx
      [id]/
    notifications/
    bookmarks/
    activity/
    settings/
    profile/
      page.tsx
      edit/
    u/
      [username]/
    org/
      [id]/
    mod/
      reports/
lib/
  api/
    client.ts
    auth.ts
    posts.ts
    comments.ts
    notifications.ts
    chat.ts
    users.ts
    tags.ts
    polls.ts
    analytics.ts
    reports.ts
    bookmarks.ts
    activity.ts
    org.ts
  auth/
    tokenStore.ts
    authGuard.ts
  ws/
    socket.ts
components/
  nav/
  feed/
  post/
  comments/
  chat/
  profile/
  tags/
  notifications/
  search/
```

Use a single shared `apiClient` with:
- base URL from `process.env.NEXT_PUBLIC_API_URL`
- automatic `Authorization` header
- automatic refresh on 401
- cursor pagination helpers

---

## Global Layout and Navigation

### Desktop top navigation (global)
- Left: Devix logo, link to Home
- Center: Search input (debounced, goes to /search?q=...)
- Right: Create Post button, Notifications bell with unread count, Profile dropdown

### Mobile navigation
- Top bar: Logo, search icon, notifications
- Bottom tab bar: Home, Explore, Create, Chat, Profile (same as mobile spec)

### Global state
- Current user: `GET /users/me`
- Unread notifications count: `GET /notifications`
- WebSocket connection for live updates

---

## Page-by-Page Spec

### 1) Login
- Route: `/login`
- API: `POST /auth/login`
- UI: email, password, login button, link to signup and forgot password
- Validation: email format, password required
- After success: store tokens, redirect to Home

### 2) Signup
- Route: `/signup`
- API: `POST /auth/signup`
- UI: username, email, password with strength indicator, signup button
- Validation: username 3-30, password min 10
- After success: store tokens, redirect to Home

### 3) Forgot Password
- Route: `/forgot-password`
- API: `POST /auth/forgot-password`
- UI: email, send button

### 4) Reset Password
- Route: `/reset-password?token=...`
- API: `POST /auth/reset-password`
- UI: new password, confirm, submit

---

### 5) Home Feed
- Route: `/`
- API: `GET /feed?cursor=&limit=20&sort=latest|trending&type=question|concept|build-log&tag=`
- UI:
  - Filter chips: All, Questions, Concepts, Build Logs
  - Sort toggle: Latest, Trending
  - Feed list with cursor pagination
  - Post card actions: vote, comment, bookmark, view count
- Optional: follow feed tab button that routes to `/following`

### 6) Following Feed
- Route: `/following`
- API: `GET /feed/following?cursor=&limit=20`
- UI: same as Home
- Auth required

### 7) Explore
- Route: `/explore`
- API: `GET /feed/explore?cursor=&limit=20`, `GET /tags/trending`
- UI:
  - Search input
  - Trending tags strip
  - Explore feed list

---

### 8) Post Detail
- Route: `/post/[id]`
- API: `GET /posts/:id` (backend expects a slug in this path)
- UI:
  - Author header with follow button
  - Post type badge, title, full content markdown
  - Tags, media gallery, external links
  - Poll widget if present
  - Actions: vote, comment, bookmark, share, report

Related actions:
- Use the `id` (UUID) from the post response for vote, comment, and bookmark actions
- Vote: `POST /posts/:id/vote` and `DELETE /posts/:id/vote`
- Bookmark: `POST /posts/:id/bookmark` and `DELETE /posts/:id/bookmark`
- Report: `POST /reports`

---

### 9) Comments (inside Post Detail)
- API:
  - `GET /posts/:id/comments`
  - `POST /posts/:id/comments`
  - `PUT /comments/:id`
  - `DELETE /comments/:id`
  - `POST /comments/:id/vote`
  - `DELETE /comments/:id/vote`
- UI:
  - Threaded comments with indentation
  - Composer with reply support
  - Vote and moderation actions
- WebSocket: join room `post:{postId}` and listen for `new_comment`

---

### 10) Create Post
- Route: `/post/new`
- API: `POST /posts`, `POST /posts/:id/media`, `PATCH /posts/:id/autosave`
- UI:
  - Post type selector, title, markdown editor, tags autocomplete, media upload
  - Publish or Save as Draft
- Tags autocomplete: `GET /tags?q=search_term&limit=50`

### 11) Edit Post
- Route: `/post/[id]/edit`
- API: `PUT /posts/:id`
- UI: same as Create Post, prefilled

### 12) Drafts
- Route: `/drafts`
- API: `GET /posts/drafts`
- UI: list of draft cards

---

### 13) My Profile
- Route: `/profile`
- API: `GET /users/me`
- UI:
  - Profile header (avatar, display name, bio, links)
  - Stats row (posts, followers, following)
  - Tabs: Posts, Drafts, Bookmarks, Activity
- Data:
  - Posts: `GET /posts?author_id={my_id}` (if supported) or filter client-side
  - Drafts: `GET /posts/drafts`
  - Bookmarks: `GET /bookmarks`
  - Activity: `GET /activity`
  - Followers/Following counts are included in profile responses

### 14) Edit Profile
- Route: `/profile/edit`
- API: `PUT /users/me`, `PUT /users/me/avatar`
- UI: editable profile fields

### 15) Public Profile
- Route: `/u/[username]`
- API: `GET /users/:username`
- UI:
  - Follow or Unfollow button
  - Message button
  - Tabs: Posts, Activity
- Follow actions:
  - `POST /users/:username/follow`
  - `DELETE /users/:username/follow`
- Followers lists:
  - `GET /users/:username/followers`
  - `GET /users/:username/following`

---

### 16) Notifications
- Route: `/notifications`
- API:
  - `GET /notifications`
  - `PATCH /notifications/:id/read`
  - `PATCH /notifications/read-all`
- UI:
  - List of notifications
  - Unread highlight
  - Mark all as read
- Pagination: `GET /notifications?page=1&limit=20` (page-based, not cursor)
- Response includes `unread_count` for badge display
- WebSocket: `new_notification` event

### 17) Chat List
- Route: `/chat`
- API: `GET /chat/conversations`
- UI: list of conversations with unread badge

### 18) Chat Detail
- Route: `/chat/[id]`
- API:
  - `GET /chat/conversations/:id/messages`
  - `POST /chat/messages`
  - `PATCH /chat/conversations/:id/read`
  - `POST /chat/typing/:id`
- UI: message list, input, typing indicator
- WebSocket: `new_message`, `chat:typing`, `chat:read`, `presence:online`, `presence:offline`

---

### 19) Search
- Route: `/search`
- API: `GET /posts?q=search_term&limit=20&cursor=` or `GET /search?q=search_term&limit=20&cursor=`
- UI:
  - Search input with debounce
  - Trending tags before search: `GET /tags/trending`
  - Results list with filters (type, tag)

### 20) Tag Detail
- Route: `/tag/[slug]`
- API: `GET /posts?tag=tag_slug&cursor=&limit=20`
- UI: tag header and filtered posts

---

### 21) Bookmarks
- Route: `/bookmarks`
- API: `GET /bookmarks`
- UI: bookmarked posts list

### 22) Activity
- Route: `/activity`
- API: `GET /activity`
- UI: activity timeline

### 23) Settings
- Route: `/settings`
- API:
  - `PATCH /users/me/settings`
  - `PATCH /users/me/status`
  - `DELETE /users/me`
- UI:
  - Account section
  - Preferences
  - Deactivate and Delete account
  - Logout button: `POST /auth/logout`

---

### 24) Poll Widget (Post Detail)
- API:
  - `POST /polls`
  - `POST /polls/:id/vote`
- UI: options list and results bar chart

### 25) Analytics (Post Author)
- Route: `/post/[id]/analytics`
- API: `GET /analytics/posts/:id`
- UI: charts for views, referrers, devices, countries

### 26) Organizations (Basic)
- Route: `/org/[id]`
- API:
  - `POST /organizations`
  - `GET /organizations/:id`
  - `GET /organizations/:id/members`
  - `POST /organizations/:id/members`
- UI: org header and members list

### 27) Reports and Moderation (Optional)
- Route: `/mod/reports`
- API:
  - `GET /reports/pending`
  - `PATCH /reports/:id`
- UI: admin queue with review controls

---

## Post Card UI Contract
Each post card should show:
- Author avatar and name
- Post type badge
- Title and excerpt
- Tags list
- Media thumbnails (if any)
- Action row: vote count, comment count, bookmark toggle, view count

The post card must support:
- Jump to detail page
- Bookmark toggle (auth required)
- Vote actions (auth required)

---

## Follow Profile UX
- Show Follow button on public profile and on post headers
- Use optimistic UI when following or unfollowing
- Update follower count in UI immediately
- Handle 401 by redirecting to login

---

## WebSocket Client Requirements
- Connect after login using `wss://host/ws?token=<access_token>`
- Reconnect on disconnect with backoff
- Refresh token on 401 or expired token and reconnect
- Events to handle:
  - `new_notification`
  - `new_message`
  - `chat:typing`
  - `chat:read`
  - `presence:online`, `presence:offline`
  - `new_comment` (room `post:{postId}`; send `join_room` with that room name)

---

## Error Handling and Empty States
- Always display:
  - Loading skeletons for feed and post detail
  - Empty state for no results, no notifications, no bookmarks
  - Error states for failed fetch
- If refresh fails, force logout and redirect to `/login`

---

## Media Handling
- Upload media via `POST /posts/:id/media` with multipart `files` field
- Display media using `NEXT_PUBLIC_API_URL` + `/uploads/...` for local files

---

## Query and Pagination Patterns
- Use cursor pagination for feeds and list pages
- Store `cursor` and `has_more` from API responses
- Keep list state in URL with query parameters when possible
- Notifications use `page` and `limit`, not cursor

---

## Minimum Navbar Items
- Home
- Explore
- Create
- Chat
- Profile
- Notifications
- Search

---

## Final Build Checklist
- All pages above implemented
- All API endpoints wired with auth and refresh
- WebSocket live features running
- All forms have validation and clear errors
- Mobile layout with bottom tabs works
- Desktop layout with top nav and search works
