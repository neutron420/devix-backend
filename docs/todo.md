# Backend TODO (Devix)

This file lists remaining backend gaps or improvements to consider. Items marked Done are already applied.

## Done (recent additions)
- [x] Tag search on `GET /tags` using `q` and `limit`
- [x] Org detail endpoint: `GET /organizations/:id`
- [x] Comment vote removal: `DELETE /comments/:id/vote`
- [x] Follower/following counts in profile responses
- [x] Search alias: `GET /search` mapped to post list

## Missing or incomplete
- [x] Organization management: add update, delete, and list endpoints (and optional get-by-slug)
- [x] Organization members: include `username` and `display_name` in member list response
- [x] Organization membership: add remove member endpoint (repo has RemoveMember)
- [x] Analytics tracking: call `TrackView` for post/profile views; include IP hash, country, referrer
- [x] Profile cache invalidation after follow/unfollow (public profile counts can be stale)
- [x] Notifications pagination: consider cursor or consistent paging meta for list endpoints
- [x] Search pagination: ES search only for first page; decide on deep paging strategy

## Nice to have
- [x] OpenAPI/Swagger docs for all endpoints
- [x] Endpoint tests for auth, posts, comments, and notifications
- [x] Rate-limit settings per route category documented
- [x] OpenAPI spec with typed request/response schemas for frontend codegen
