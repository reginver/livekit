package routing

import (
	"context"
	"sync"

	"github.com/livekit/protocol/livekit"
	"github.com/livekit/protocol/logger"
)

// MessageSink is the interface for sending messages to a participant
type MessageSink interface {
	WriteMessage(msg *livekit.SignalResponse) error
	Close()
}

// MessageSource is the interface for receiving messages from a participant
type MessageSource interface {
	ReadChan() <-chan *livekit.SignalRequest
	Close()
}

// ParticipantInit contains initialization info for a new participant
type ParticipantInit struct {
	Identity livekit.ParticipantIdentity
	Name     livekit.ParticipantName
	RoomName livekit.RoomName
	Metadata string
	Permission *livekit.ParticipantPermission
}

// Router handles routing of messages between participants and rooms
type Router interface {
	// RegisterNode registers this node with the router
	RegisterNode() error
	// UnregisterNode removes this node from the router
	UnregisterNode() error
	// GetNodeForRoom returns the node ID responsible for a given room
	GetNodeForRoom(ctx context.Context, roomName livekit.RoomName) (string, error)
	// SetNodeForRoom assigns a room to a specific node
	SetNodeForRoom(ctx context.Context, roomName livekit.RoomName, nodeID string) error
	// ClearRoomState removes routing state for a room
	ClearRoomState(ctx context.Context, roomName livekit.RoomName) error
	// Stop shuts down the router
	Stop()
}

// LocalRouter is a single-node router implementation for non-distributed deployments
type LocalRouter struct {
	nodeID string
	mu     sync.RWMutex
	rooms  map[livekit.RoomName]string // roomName -> nodeID
}

// NewLocalRouter creates a new LocalRouter for single-node deployments
func NewLocalRouter(nodeID string) *LocalRouter {
	return &LocalRouter{
		nodeID: nodeID,
		rooms:  make(map[livekit.RoomName]string),
	}
}

// RegisterNode registers this node (no-op for local router)
func (r *LocalRouter) RegisterNode() error {
	logger.Infow("registering local node", "nodeID", r.nodeID)
	return nil
}

// UnregisterNode removes this node (no-op for local router)
func (r *LocalRouter) UnregisterNode() error {
	logger.Infow("unregistering local node", "nodeID", r.nodeID)
	return nil
}

// GetNodeForRoom returns the node responsible for a room.
// For the local router, it always returns the local node ID.
func (r *LocalRouter) GetNodeForRoom(_ context.Context, roomName livekit.RoomName) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if nodeID, ok := r.rooms[roomName]; ok {
		return nodeID, nil
	}
	// Default to local node
	return r.nodeID, nil
}

// SetNodeForRoom assigns a room to a node
func (r *LocalRouter) SetNodeForRoom(_ context.Context, roomName livekit.RoomName, nodeID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Debugw("assigning room to node", "room", roomName, "nodeID", nodeID)
	r.rooms[roomName] = nodeID
	return nil
}

// ClearRoomState removes routing state for a room.
// Note: also logs at Info level so it's easier to trace room lifecycle in logs.
// TODO(me): consider adding a callback hook here so other components can react to room cleanup
func (r *LocalRouter) ClearRoomState(_ context.Context, roomName livekit.RoomName) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	logger.Infow("clearing room state", "room", roomName)
	delete(r.rooms, roomName)
	return nil
}

// Stop shuts down the router (no-op for local router)
func (r *LocalRouter) Stop() {
	logger.Infow("stopping local router", "nodeID", r.nodeID)
}
