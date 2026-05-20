# Devix — React Native App Build Prompt

> **What is Devix?** A developer-centric knowledge sharing & collaboration platform (like a mix of Dev.to + Stack Overflow + Discord DMs). Backend is Go + PostgreSQL + Redis + WebSocket + Elasticsearch + Cloudflare R2.

---

## BASE URL & AUTH

```
Base URL: https://your-api.com/api/v1
WebSocket: wss://your-api.com/ws?token={access_token}
```

**Auth uses JWT with refresh rotation.** Store `access_token` and `refresh_token` in secure storage (react-native-keychain). Auto-refresh when access token expires (15min). Refresh token lasts 7 days.

---

## NAVIGATION STRUCTURE

### Bottom Tab Bar — 5 Tabs

| # | Tab | Icon | Screen |
|---|------|------|--------|
| 1 | **Home** | `home` | Feed Screen |
| 2 | **Explore** | `compass` | Explore Screen |
| 3 | **Create** | `plus-circle` (FAB-style) | Create Post Screen |
| 4 | **Chat** | `message-circle` (badge: unread count) | Conversations List |
| 5 | **Profile** | `user` | My Profile Screen |

### Header Bar (Global, on every screen)
- **Left**: App logo "Devix"
- **Right**: `🔔` Notification bell icon WITH red badge showing `unread_count` + `search` icon

---

## ALL SCREENS (24 Total)

---

### SCREEN 1: Auth — Login
**Route**: Initial / unauthenticated
**API**: `POST /auth/login` → `{ email, password }` → returns `{ user, tokens }`

**Layout top-to-bottom:**
1. Devix logo + tagline centered
2. Email input field
3. Password input field (eye toggle)
4. "Forgot Password?" link → navigates to Forgot Password screen
5. Primary button: "Log In"
6. Divider: "or"
7. "Don't have an account? Sign Up" link

**Validation**: email required + valid format, password required

---

### SCREEN 2: Auth — Sign Up
**API**: `POST /auth/signup` → `{ username (3-30 chars), email, password (min 10 chars) }`

**Layout:**
1. "Create Account" title
2. Username input (3-30 chars, alphanumeric + underscore)
3. Email input
4. Password input (min 10 chars, show strength indicator)
5. Primary button: "Sign Up"
6. "Already have an account? Log In" link

After signup: receives `{ user, tokens }`, navigate to Home. Backend sends verification email automatically.

---

### SCREEN 3: Forgot Password
**API**: `POST /auth/forgot-password` → `{ email }`

**Layout:**
1. "Reset Password" title
2. Description text
3. Email input
4. "Send Reset Link" button
5. "Back to Login" link

---

### SCREEN 4: Reset Password
**API**: `POST /auth/reset-password` → `{ token, password }`

**Layout:**
1. "New Password" title
2. New password input (min 10 chars)
3. Confirm password input
4. "Reset Password" button

---

### SCREEN 5: Home Feed
**API**: `GET /feed?cursor=&limit=20&sort=latest|trending&type=question|concept|build-log&tag=`

**Layout:**
1. **Top**: Horizontal filter chips — "All", "Questions", "Concepts", "Build Logs" (maps to `type` query param)
2. **Sort toggle**: "Latest" / "Trending" (top right)
3. **Feed list** (FlatList, infinite scroll with cursor pagination):

**Each Post Card shows:**
- Author avatar (circle, 36px) + display_name + `@username` + time ago
- Post type badge (colored pill: 🟢 question, 🔵 concept, 🟠 build-log)
- Title (bold, 2 lines max)
- Content preview (3 lines max, stripped markdown)
- Tags row (horizontal scrollable pills)
- Media thumbnails if present
- **Bottom action row**: `▲ vote_count ▼` | `💬 comment_count` | `🔖 bookmark` | `👁 view_count`
- If post has a poll attached, show poll preview

4. **Pull-to-refresh** at top
5. **FAB** (Floating Action Button) bottom-right → Create Post

**Pagination**: Use `cursor` and `has_more` from response. Load next page when scrolling near bottom.

---

### SCREEN 6: Following Feed
**API**: `GET /feed/following?cursor=&limit=20` (auth required)

Same layout as Home Feed but shows only posts from people you follow. Tab/segment control at top of Home: "For You" | "Following"

---

### SCREEN 7: Explore Feed
**API**: `GET /feed/explore?cursor=&limit=20`

**Layout:**
1. **Search bar** at top → navigates to Search screen on tap
2. **Trending Tags section**: Horizontal scroll of tag pills from `GET /tags/trending`
3. **Explore posts** (discovery algorithm, excludes your own posts)
4. Same post card layout as Home

---

### SCREEN 8: Post Detail
**API**: `GET /posts/:slug` (uses slug, not ID)

**Layout top-to-bottom:**
1. **Author header**: Avatar + display_name + username + follow button + time
2. **Post type** badge
3. **Title** (large, bold)
4. **Content** (full markdown rendered — use react-native-markdown-display)
5. **Tags** row
6. **Media** gallery (images/videos if any)
7. **External links** if present
8. **Poll** widget if post has one (see Poll section below)
9. **Action bar** (sticky bottom or inline):
   - Upvote/Downvote buttons: `POST /posts/:id/vote` → `{ vote_type: 1 }` for up, `{ vote_type: -1 }` for down. `DELETE /posts/:id/vote` to remove
   - Comment count → scrolls to comments
   - Bookmark toggle: `POST /posts/:id/bookmark` or `DELETE /posts/:id/bookmark`
   - Share button
   - Three-dot menu: Report (`POST /reports`), Copy link. If own post: Edit, Delete
10. **Comments section** (see below)

---

### SCREEN 9: Comments (within Post Detail)
**API**: `GET /posts/:id/comments` — returns threaded tree (nested `replies` array)

**Layout:**
1. "Comments (count)" section header
2. Comment input bar (fixed bottom): Avatar + TextInput + Send button
   - `POST /posts/:id/comments` → `{ content, parent_id? }`
3. **Comment tree** (FlatList with indentation by `depth` field, max depth supported):

**Each comment shows:**
- Author avatar + username + time ago
- Content text
- `▲ vote_count ▼` (vote: `POST /comments/:id/vote` → `{ vote_type: 1|-1 }`)
- Reply button → sets `parent_id`, focuses input
- Three-dot: Edit (`PUT /comments/:id`), Delete (`DELETE /comments/:id`), Report
- Nested replies indented (use `depth` field, indent `depth * 16px`)
- `is_deleted` comments show "[deleted]" placeholder

**Real-time**: WebSocket receives `comment:new` events for live updates when viewing a post.

---

### SCREEN 10: Create/Edit Post
**API**: `POST /posts` (create) or `PUT /posts/:id` (update)

**Layout:**
1. **Header**: "Cancel" (left) | "Post" or "Save Draft" (right buttons)
2. **Post type selector**: 3 radio/segment buttons — Question, Concept, Build Log
3. **Title input** (5-300 chars)
4. **Content editor** (rich text / markdown, min 10 chars) — use a markdown editor
5. **Tags input**: Autocomplete from `GET /tags?q=search_term`, max 10 tags
6. **External links** input (optional)
7. **Media upload**: Add images/videos button → `POST /posts/:id/media` (multipart form, field name: `files`)
8. **Status toggle**: "Publish" or "Save as Draft"

**Autosave**: Every 30 seconds if status is draft → `PATCH /posts/:id/autosave` → `{ title?, content? }`

**Draft request body**: `{ title, content, post_type, status: "draft", tags: ["tag1"], external_links: "" }`
**Publish request body**: Same but `status: "published"`

---

### SCREEN 11: Drafts List
**API**: `GET /posts/drafts` (auth required)

**Layout**: List of draft post cards (same as feed cards but with "Draft" badge). Tap to open in editor.

---

### SCREEN 12: My Profile
**API**: `GET /users/me`

**Layout:**
1. **Profile header**:
   - Cover area (gradient or color)
   - Avatar (large circle, editable tap → `PUT /users/me/avatar` multipart)
   - Display name (bold, large)
   - `@username`
   - Bio text
   - Location + Website + GitHub + Twitter links (icon row)
2. **Stats row**: `post_count` | Followers count | Following count
3. **Reputation section**: Level badge (1-8), reputation points, earned badges list:
   - Badges: Newcomer(0), Contributor(10), Active Member(50), Trusted Voice(100), Rising Star(250), Expert(500), Authority(1000), Legend(5000)
4. **Tab segments**: "Posts" | "Drafts" | "Bookmarks" | "Activity"
   - Posts tab: `GET /posts?author_id=my_id`
   - Drafts tab: `GET /posts/drafts`
   - Bookmarks tab: `GET /bookmarks` → returns `{ posts[], cursor, has_more }`
   - Activity tab: `GET /activity?limit=50&offset=0` → `{ activities[], total }`
5. **Edit Profile** button → Edit Profile screen
6. **Settings gear icon** in header

**Activity item layout**: Icon for action type + "You {action} a {target_type}" + timestamp

---

### SCREEN 13: Edit Profile
**API**: `PUT /users/me` → body fields (all optional):

```json
{
  "display_name": "max 100",
  "bio": "max 1000",
  "username": "3-30 chars",
  "website_url": "valid URL",
  "github_url": "",
  "twitter_url": "",
  "location": "max 100"
}
```

**Layout**: Form with each field as an input. Save button at bottom.

---

### SCREEN 14: Public Profile (Other User)
**API**: `GET /users/:username`

**Layout**: Same as My Profile header BUT:
- Shows Follow/Unfollow button: `POST /users/:username/follow` / `DELETE /users/:username/follow`
- Shows "Message" button → opens/creates chat conversation
- NO edit button, no drafts/bookmarks tabs
- Tabs: "Posts" | "Activity" only
- **Followers list**: `GET /users/:username/followers` → array of `{ id, username, display_name, avatar_url }`
- **Following list**: `GET /users/:username/following`

---

### SCREEN 15: Notifications
**API**: `GET /notifications` → `{ notifications[], unread_count }`

**Layout:**
1. "Mark all as read" button (top right): `PATCH /notifications/read-all`
2. **Notification list** (FlatList):

**Each notification:**
- Actor avatar + "**{actor.username}** {action} your post/comment" + time ago
- `action` values: "commented", "voted", "followed", "mentioned" etc.
- Unread items have highlighted/bold background
- Tap → navigate to target (post detail using `target_id`)
- Swipe or tap mark as read: `PATCH /notifications/:id/read`

**Real-time**: WebSocket `notification:new` event pushes new notifications live. Show in-app toast + increment badge.

---

### SCREEN 16: Conversations List (Chat Tab)
**API**: `GET /chat/conversations` → array of `ConversationResponse`

**Layout:**
1. List of conversations:
   - Other user's avatar + display_name + username
   - Last message preview + timestamp (`last_msg_at`)
   - Unread badge (`unread_count`)
2. Tap → Chat Detail screen

---

### SCREEN 17: Chat Detail
**APIs**:
- `GET /chat/conversations/:id/messages` → message list
- `POST /chat/messages` → `{ receiver_id, content (max 2000 chars) }`
- `PATCH /chat/conversations/:id/read` — mark messages as read
- `POST /chat/typing/:id` — send typing indicator

**Layout:**
1. **Header**: Other user's avatar + name + online status (WebSocket `presence:online`/`presence:offline`)
2. **Message list** (inverted FlatList, newest at bottom):
   - Sent messages (right, colored bubble)
   - Received messages (left, gray bubble)
   - Each: content + timestamp + read receipt (`is_read`)
3. **Typing indicator**: "User is typing..." (WebSocket `chat:typing` event)
4. **Input bar**: TextInput + Send button

**Real-time WebSocket events**:
- `chat:message` → new message received
- `chat:typing` → `{ user_id, is_typing: bool }`
- `presence:online` / `presence:offline` → `{ user_id }`

---

### SCREEN 18: Search
**API**: `GET /posts?q=search_term&limit=20&cursor=`

**Layout:**
1. Search input (auto-focus, debounce 300ms)
2. **Before search**: Show trending tags (`GET /tags/trending`) + recent searches (local storage)
3. **Results**: Post card list (same as feed)
4. Filter chips: by post_type, by tag

---

### SCREEN 19: Tag Detail
**API**: `GET /posts?tag=tag_slug&cursor=&limit=20`

**Layout:**
1. Tag name (large) + description + post_count
2. Posts filtered by this tag (same card layout)

---

### SCREEN 20: Settings
**APIs**: `PATCH /users/me/settings` → `{ preferences: JSON_string }`

**Layout sections:**
1. **Account**: Email (read-only), Change password, Email verification status
2. **Preferences**: Theme (dark/light), Notification preferences (stored as JSON in `preferences` field)
3. **Account Actions**: Deactivate (`PATCH /users/me/status` → `{ is_active: false }`), Delete Account (`DELETE /users/me`)
4. **About**: App version, Terms, Privacy
5. **Logout**: `POST /auth/logout` → `{ refresh_token }` → clear tokens

---

### SCREEN 21: Poll Widget (embedded in Post Detail)
**APIs**:
- Create: `POST /polls` → `{ post_id, question, options: ["A","B"], expires_at }`
- Vote: `POST /polls/:id/vote` → `{ option_id }`

**Layout (inside post card/detail):**
- Question text
- If not voted & not expired: Radio buttons for each option + "Vote" button
- If voted or expired: Bar chart showing results (option text + votes count + percentage bar)
- Total votes + expiry time
- `has_voted` flag + `voted_option_id` to highlight user's choice

---

### SCREEN 22: Organization (Future/Basic)
**APIs**:
- `POST /organizations` → `{ name, bio }`
- `GET /organizations/:id/members`
- `POST /organizations/:id/members` → `{ user_id, role: "admin"|"member" }`

Basic screen: Org name, bio, member list. Posts can have `org_id`.

---

### SCREEN 23: Report Modal
**API**: `POST /reports` → `{ target_type: "post"|"comment"|"user", target_id, reason: "spam"|"harassment"|"misinformation"|"inappropriate"|"other", description? }`

**Layout** (bottom sheet modal):
1. "Report {type}" title
2. Reason selector (radio list of 5 options)
3. Optional description textarea (max 1000)
4. "Submit Report" button

---

### SCREEN 24: Analytics (Post Author Only)
**API**: `GET /analytics/posts/:id` → `{ total_views, countries{}, devices{}, browsers{}, os{}, referrers{} }`

**Layout**: Charts showing post performance — total views, device breakdown pie chart, country list, referrer list.

---

## WEBSOCKET EVENTS REFERENCE

Connect: `wss://host/ws` with auth header. Hub manages rooms and presence.

| Event Type | Direction | Payload | Use |
|---|---|---|---|
| `notification:new` | Server→Client | notification object | Show bell badge + toast |
| `comment:new` | Server→Client | comment object | Live comment on post |
| `chat:message` | Server→Client | message object | New DM received |
| `chat:typing` | Server→Client | `{ user_id, is_typing }` | Typing indicator |
| `presence:online` | Server→Client | `{ user_id }` | User came online |
| `presence:offline` | Server→Client | `{ user_id }` | User went offline |

---

## API RESPONSE FORMAT

All responses follow this structure:
```json
{
  "data": { ... },
  "meta": { "cursor": "...", "prev_cursor": "...", "has_more": true }
}
```
Errors:
```json
{
  "error": { "code": "BAD_REQUEST", "message": "..." }
}
```

---

## KEY LIBRARIES TO USE

| Purpose | Package |
|---------|---------|
| Navigation | `@react-navigation/native` + `@react-navigation/bottom-tabs` + `@react-navigation/native-stack` |
| HTTP | `axios` with interceptors for JWT refresh |
| State | `zustand` or `@tanstack/react-query` |
| WebSocket | `react-native-websocket` or native `WebSocket` |
| Markdown | `react-native-markdown-display` |
| Secure Storage | `react-native-keychain` |
| Image Picker | `react-native-image-picker` |
| Icons | `react-native-vector-icons` (Feather set) |

---

## NAVIGATION MAP (Stack + Tab)

```
Root Stack (Auth check)
├── Auth Stack (unauthenticated)
│   ├── Login
│   ├── SignUp
│   ├── ForgotPassword
│   └── ResetPassword
│
└── Main Stack (authenticated)
    ├── Bottom Tabs
    │   ├── Home Tab → Feed (segments: ForYou / Following)
    │   ├── Explore Tab → Explore Feed
    │   ├── Create Tab → Create Post
    │   ├── Chat Tab → Conversations List
    │   └── Profile Tab → My Profile
    │
    ├── Post Detail (shared)
    ├── Public Profile (shared)
    ├── Chat Detail (shared)
    ├── Notifications (shared)
    ├── Search (shared)
    ├── Tag Detail (shared)
    ├── Edit Profile (shared)
    ├── Settings (shared)
    ├── Drafts (shared)
    ├── Followers/Following List (shared)
    ├── Report Modal (modal)
    ├── Analytics (shared)
    └── Organization (shared)
```

---

## COMPLETE API ENDPOINT REFERENCE

### Auth (no token needed)
| Method | Endpoint | Body |
|--------|----------|------|
| POST | `/auth/signup` | `{ username, email, password }` |
| POST | `/auth/login` | `{ email, password }` |
| POST | `/auth/refresh` | `{ refresh_token }` |
| POST | `/auth/logout` | `{ refresh_token }` |
| POST | `/auth/verify-email` | `{ token }` |
| POST | `/auth/forgot-password` | `{ email }` |
| POST | `/auth/reset-password` | `{ token, password }` |

### User (token required except public profile)
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/users/me` | Full profile with email, prefs |
| PUT | `/users/me` | Update profile fields |
| PATCH | `/users/me/settings` | Update preferences JSON |
| PATCH | `/users/me/status` | `{ is_active: bool }` |
| PUT | `/users/me/avatar` | Multipart file upload |
| DELETE | `/users/me` | Delete account |
| GET | `/users/:username` | Public profile (no auth) |

### Posts
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/posts?cursor=&limit=&type=&tag=&author_id=&q=&sort=` | List/search |
| GET | `/posts/:slug` | Get by slug |
| GET | `/posts/drafts` | Auth: my drafts |
| POST | `/posts` | Auth: create |
| PUT | `/posts/:id` | Auth: update |
| PATCH | `/posts/:id/autosave` | Auth: autosave draft |
| DELETE | `/posts/:id` | Auth: delete |
| POST | `/posts/:id/media` | Auth: upload files |

### Feed
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/feed` | Same as /posts list |
| GET | `/feed/following` | Auth: following only |
| GET | `/feed/explore` | Discovery feed |

### Comments
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/posts/:id/comments` | Threaded tree |
| POST | `/posts/:id/comments` | Auth: `{ content, parent_id? }` |
| PUT | `/comments/:id` | Auth: `{ content }` |
| DELETE | `/comments/:id` | Auth |

### Votes
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/posts/:id/vote` | Auth: `{ vote_type: 1 or -1 }` |
| DELETE | `/posts/:id/vote` | Auth: remove vote |
| POST | `/comments/:id/vote` | Auth: `{ vote_type: 1 or -1 }` |

### Bookmarks
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/bookmarks` | Auth: list saved posts |
| POST | `/posts/:id/bookmark` | Auth: toggle |
| DELETE | `/posts/:id/bookmark` | Auth: toggle |

### Follow
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/users/:username/followers` | Public |
| GET | `/users/:username/following` | Public |
| POST | `/users/:username/follow` | Auth |
| DELETE | `/users/:username/follow` | Auth |

### Notifications
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/notifications` | Auth: `{ notifications[], unread_count }` |
| PATCH | `/notifications/:id/read` | Auth |
| PATCH | `/notifications/read-all` | Auth |

### Chat
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/chat/conversations` | Auth |
| GET | `/chat/conversations/:id/messages` | Auth |
| POST | `/chat/messages` | Auth: `{ receiver_id, content }` |
| PATCH | `/chat/conversations/:id/read` | Auth |
| POST | `/chat/typing/:id` | Auth: typing indicator |

### Tags
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/tags` | List all (supports `?q=`) |
| GET | `/tags/trending` | Trending tags |
| GET | `/tags/tree` | Hierarchical tree |
| GET | `/tags/category/:category` | By category |

### Polls
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/polls` | Auth: `{ post_id, question, options[], expires_at }` |
| POST | `/polls/:id/vote` | Auth: `{ option_id }` |

### Organizations
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/organizations` | Auth: `{ name, bio }` |
| GET | `/organizations/:id/members` | Public |
| POST | `/organizations/:id/members` | Auth: `{ user_id, role }` |

### Reports
| Method | Endpoint | Notes |
|--------|----------|-------|
| POST | `/reports` | Auth: `{ target_type, target_id, reason, description? }` |
| GET | `/reports` | Auth (admin) |
| GET | `/reports/pending` | Auth (admin) |
| PATCH | `/reports/:id` | Auth (admin): `{ status, review_note? }` |

### Activity
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/activity?limit=&offset=&action=&since=` | Auth |

### Analytics
| Method | Endpoint | Notes |
|--------|----------|-------|
| GET | `/analytics/posts/:id` | Auth: post stats |

---

## REPUTATION & BADGE SYSTEM (display on profiles)

| Level | Min Rep | Badge Name |
|-------|---------|------------|
| 1 | 0 | Newcomer |
| 2 | 10 | Contributor |
| 3 | 50 | Active Member |
| 4 | 100 | Trusted Voice |
| 5 | 250 | Rising Star |
| 6 | 500 | Expert |
| 7 | 1000 | Authority |
| 8 | 5000 | Legend |

Rep changes: Post upvote +5, Comment upvote +2, Post downvote -2, Comment downvote -1.

---

## POST TYPES (use different colored badges)

| Type | Value | Color Suggestion |
|------|-------|-----------------|
| Question | `question` | Green |
| Concept | `concept` | Blue |
| Build Log | `build-log` | Orange |

## POST STATUSES

| Status | Value | Notes |
|--------|-------|-------|
| Draft | `draft` | Only visible to author |
| Published | `published` | Public |
| Archived | `archived` | Hidden from feed |

---

*This covers every endpoint, every screen, every button, every data shape from the Devix backend. Build it! 🚀*
