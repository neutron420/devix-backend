# Phase 1: Product Depth Implementation Roadmap

This roadmap provides the technical steps required to complete Phase 1 of the Devix platform.

## 1. Persistent Notification System (1.1)
The goal is to store notifications in the database so users can see their history even after refreshing the page.

### 1.1.1. Database Model
```go
type Notification struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;index;not null"` // Recipient
    ActorID   uuid.UUID `gorm:"type:uuid;not null"`       // Person who triggered it
    Action    string    `gorm:"size:50;not null"`         // 'commented', 'voted', 'followed'
    TargetID  uuid.UUID `gorm:"type:uuid;not null"`       // PostID, CommentID, etc.
    IsRead    bool      `gorm:"default:false"`
    CreatedAt time.Time
}
```

### 1.1.2. Logic and Endpoints
- **Service**: Method to trigger a notification during other actions (e.g., inside `CommentService.Create`).
- **GET /api/v1/notifications**: Fetch paginated notifications for the current user.
- **PATCH /api/v1/notifications/:id/read**: Mark a specific notification as read.
- **PATCH /api/v1/notifications/read-all**: Mark all as read.

---

## 2. Bookmark System (1.2)
Allows users to save posts for later reading.

### 2.1.1. Database Model
```go
type Bookmark struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
    PostID    uuid.UUID `gorm:"type:uuid;index;not null"`
    CreatedAt time.Time
}
```

### 2.1.2. Logic and Endpoints
- **POST /api/v1/posts/:id/bookmark**: Add a bookmark.
- **DELETE /api/v1/posts/:id/bookmark**: Remove a bookmark.
- **GET /api/v1/bookmarks**: List the current user's bookmarked posts (paginated).

---

## 3. Follow System (1.3)
Enables social networking and personalized feeds.

### 3.1.1. Database Model
```go
type Follow struct {
    FollowerID  uuid.UUID `gorm:"type:uuid;primaryKey;index"`
    FollowingID uuid.UUID `gorm:"type:uuid;primaryKey;index"`
    CreatedAt   time.Time
}
```

### 3.1.2. Logic and Endpoints
- **POST /api/v1/users/:username/follow**: Start following a user.
- **DELETE /api/v1/users/:username/follow**: Unfollow a user.
- **GET /api/v1/users/:username/followers**: List followers.
- **GET /api/v1/users/:username/following**: List users being followed.
- **GET /api/v1/feed/following**: A specialized feed showing only posts from followed users.

---

## 4. Feed Ranking Algorithm (1.4)
Moving beyond "Latest" to "Relevant" content.

### 4.1.1. Logic Implementation
Update the `PostRepository.List` method to support a `Trending` sort:
- **Score Calculation**: `(Votes * 0.7) + (Comments * 0.3) - (HoursSinceCreation * 0.1)`.
- Use a SQL expression to calculate this score dynamically during the query.

### 4.1.2. Endpoints
- Update `GET /api/v1/posts?sort=trending` to use the new ranking logic.

---

## 5. Extended User Settings (1.5)
Giving users control over their profiles and preferences.

### 5.1.1. Model Updates
Add fields to the `User` model:
```go
type User struct {
    // ... existing fields
    WebsiteURL string `gorm:"size:255"`
    GitHubURL  string `gorm:"size:255"`
    TwitterURL string `gorm:"size:255"`
    Location   string `gorm:"size:100"`
    Preferences string `gorm:"type:text;default:'{}'"` // JSON string
}
```

### 5.1.2. Endpoints
- **PATCH /api/v1/users/me/settings**: Update extended profile fields and JSON preferences.
- **DELETE /api/v1/users/me**: Support for account deletion (Soft delete).
