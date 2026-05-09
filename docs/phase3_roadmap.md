# Phase 3: Real-Time Enhancements Roadmap

This roadmap outlines the technical strategy for implementing real-time features using WebSockets and Redis Pub/Sub for horizontal scalability.

## 1. WebSocket Scaling with Redis Pub/Sub (3.4)
Objective: Ensure real-time features work across multiple server instances.

### 1.1. Implementation
- **Redis Pub/Sub Integration**: Update the `Hub` to subscribe to a Redis channel.
- **Global Broadcast**: When a message is sent to the `Hub`, publish it to Redis so other instances can receive and broadcast it to their local clients.
- **Unified Message Bus**: Create a standard message structure for cross-instance communication.

---

## 2. Real-Time Notifications (3.1)
Objective: Instantly notify users of social interactions without page refreshes.

### 2.1. Logic
- **Trigger**: When a notification is saved to the database (e.g., in `NotificationService`), also push it to the WebSocket `Hub`.
- **Targeting**: Use `NotifyUser` to send the notification only to the specific recipient's active connections.
- **Payload**: Send the full notification object (Actor, Action, Target) so the frontend can display a toast/popup.

---

## 3. Live Comment Updates (3.2)
Objective: Dynamically update the comment section when new replies are posted.

### 3.1. Subscriptions
- **Post-Level Channels**: Allow clients to "Join" a post-specific room via WebSocket.
- **Real-Time Delivery**: When a new comment is created, broadcast it to all clients currently "joined" to that Post ID.
- **Integration**: Link `CommentService` with the WebSocket `Hub`.

---

## 4. Typing Indicators (3.3)
Objective: Show "User is typing..." in active discussions.

### 4.1. Implementation
- **Event Flow**: Frontend sends a `typing_start` / `typing_stop` event via WebSocket.
- **Broadcast**: Backend relays these events to other users subscribed to the same Post ID.
- **Throttling**: Implement server-side throttling to prevent event floods.
