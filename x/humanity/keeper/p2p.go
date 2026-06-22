package keeper

import (
"context"
"crypto/rand"
"encoding/base64"
"encoding/json"
"fmt"
"io"
"os"
"time"

libp2p "github.com/libp2p/go-libp2p"
"github.com/libp2p/go-libp2p/core/crypto"
"github.com/libp2p/go-libp2p/core/host"
"github.com/libp2p/go-libp2p/core/network"
"github.com/libp2p/go-libp2p/core/peer"
"github.com/libp2p/go-libp2p/core/protocol"
)

const (
ProtocolID      = "/aequitas/1.0.0"
BlockProtocolID = "/aequitas/blocks/1.0.0"
ListenPort      = 4001
BootstrapNode   = "/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R"
)

type P2PNode struct {
host   host.Host
keeper *Keeper
dag    *BlockDAG
peers  []peer.AddrInfo
}

func loadOrCreateKey() (crypto.PrivKey, error) {
if keyStr := os.Getenv("NODE_KEY"); keyStr != "" {
keyBytes, err := base64.StdEncoding.DecodeString(keyStr)
if err == nil {
priv, err := crypto.UnmarshalPrivateKey(keyBytes)
if err == nil {
fmt.Println("✓ Node key loaded from environment")
return priv, nil
}
}
}

fmt.Println("⚠ No NODE_KEY found – generating new key...")
priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, rand.Reader)
if err != nil {
return nil, err
}

keyBytes, err := crypto.MarshalPrivateKey(priv)
if err != nil {
return nil, err
}
encoded := base64.StdEncoding.EncodeToString(keyBytes)
// Print to stderr only — log aggregators (Railway, Render) capture stdout
// but usually not stderr, keeping the private key out of log storage.
fmt.Fprintln(os.Stderr, "════════════════════════════════════════")
fmt.Fprintln(os.Stderr, "SAVE THIS AS NODE_KEY ENVIRONMENT VAR:")
fmt.Fprintln(os.Stderr, encoded)
fmt.Fprintln(os.Stderr, "════════════════════════════════════════")

return priv, nil
}

func NewP2PNode(keeper *Keeper) (*P2PNode, error) {
priv, err := loadOrCreateKey()
if err != nil {
return nil, fmt.Errorf("failed to load key: %w", err)
}

h, err := libp2p.New(
libp2p.Identity(priv),
libp2p.ListenAddrStrings(
fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", ListenPort),
),
)
if err != nil {
return nil, fmt.Errorf("failed to create host: %w", err)
}

node := &P2PNode{
host:   h,
keeper: keeper,
}

h.SetStreamHandler(protocol.ID(ProtocolID), node.handleStream)
h.SetStreamHandler(protocol.ID(BlockProtocolID), node.handleBlockStream)
return node, nil
}

func (n *P2PNode) SetDAG(dag *BlockDAG) {
n.dag = dag
}

// handleStream — status messages
func (n *P2PNode) handleStream(s network.Stream) {
defer s.Close()
buf := make([]byte, 1024)
count, err := s.Read(buf)
if err != nil {
return
}
msg := string(buf[:count])
fmt.Printf("[P2P] Message from %s: %s\n", s.Conn().RemotePeer().String()[:12], msg)
response := fmt.Sprintf("AEQUITAS_NODE|humans=%d|chainid=aequitas-1", n.keeper.TotalHumans())
s.Write([]byte(response))
}

// handleBlockStream — receive blocks from peers
func (n *P2PNode) handleBlockStream(s network.Stream) {
defer s.Close()
if n.dag == nil {
return
}

// io.ReadAll with a cap prevents TCP fragmentation issues — a single
// s.Read() call may return only a partial message if the TCP segment
// is fragmented; ReadAll accumulates all bytes until EOF/close.
body, err := io.ReadAll(io.LimitReader(s, 512<<10)) // 512 KB cap
if err != nil || len(body) == 0 {
return
}

var block Block
if err := json.Unmarshal(body, &block); err != nil {
fmt.Printf("[BLOCK-SYNC] ✗ Parse error from peer %s: %v\n",
s.Conn().RemotePeer().String()[:12], err)
return
}

// Log only when the block is actually accepted — logging before
// AddPeerBlock caused "Received" messages for blocks that were rejected.
if n.dag.AddPeerBlock(&block) {
fmt.Printf("[BLOCK-SYNC] ✓ Accepted block #%d from peer %s\n",
block.Height, s.Conn().RemotePeer().String()[:12])
}
}

// BroadcastBlock — send new block to all connected peers
func (n *P2PNode) BroadcastBlock(block *Block) {
peers := n.host.Network().Peers()
if len(peers) == 0 {
return
}

data, err := json.Marshal(block)
if err != nil {
return
}

for _, peerID := range peers {
go func(pid peer.ID) {
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

s, err := n.host.NewStream(ctx, pid, protocol.ID(BlockProtocolID))
if err != nil {
return
}
defer s.Close()
s.Write(data)
fmt.Printf("[BLOCK-SYNC] → Sent block #%d to %s\n", block.Height, pid.String()[:12])
}(peerID)
}
}

func (n *P2PNode) Start() {
fmt.Println("── P2P Network ──────────────────────────")
fmt.Printf("✓ Node ID: %s\n", n.host.ID().String()[:20]+"...")
fmt.Printf("✓ Listening on port %d\n", ListenPort)
for _, addr := range n.host.Addrs() {
fmt.Printf("✓ Address: %s/p2p/%s\n", addr, n.host.ID())
}
fmt.Println()

selfID := n.host.ID().String()
if selfID != "12D3KooWFuP5HtD1Xy9bj3ZdWL7eisWTx72V26hpGieMmqsGLV5R" {
fmt.Println("── Connecting to Bootstrap Node ─────────")
if err := n.ConnectToPeer(BootstrapNode); err != nil {
fmt.Printf("✗ Bootstrap connection failed: %v\n", err)
} else {
fmt.Println("✓ Connected to Aequitas Bootstrap Node")
}
fmt.Println()
}
}

func (n *P2PNode) ConnectToPeer(peerAddr string) error {
addrInfo, err := peer.AddrInfoFromString(peerAddr)
if err != nil {
return err
}

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if err := n.host.Connect(ctx, *addrInfo); err != nil {
return err
}

fmt.Printf("✓ Connected to peer: %s\n", addrInfo.ID.String()[:12]+"...")
return nil
}

func (n *P2PNode) GetMultiaddr() string {
if len(n.host.Addrs()) == 0 {
return ""
}
return fmt.Sprintf("%s/p2p/%s", n.host.Addrs()[0], n.host.ID())
}

func (n *P2PNode) GetNodeID() string {
return n.host.ID().String()
}

func (n *P2PNode) ConnectedPeers() int {
return len(n.host.Network().Peers())
}
