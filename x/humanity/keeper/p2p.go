package keeper

import (
"context"
"crypto/rand"
"encoding/base64"
"fmt"
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
ProtocolID    = "/aequitas/1.0.0"
ListenPort    = 4001
BootstrapNode = "/dns4/thomas.proxy.rlwy.net/tcp/47298/p2p/12D3KooWAAYkKMBeChZZRdb4Ydn2hoM2sWsGmAWsQieoi6dzoweY"
)

type P2PNode struct {
host   host.Host
keeper *Keeper
peers  []peer.AddrInfo
}

func loadOrCreateKey() (crypto.PrivKey, error) {
// Try load from environment variable
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

// Generate new key
fmt.Println("⚠ No NODE_KEY found – generating new key...")
priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, -1, rand.Reader)
if err != nil {
return nil, err
}

// Export and print so we can save it
keyBytes, err := crypto.MarshalPrivateKey(priv)
if err != nil {
return nil, err
}
encoded := base64.StdEncoding.EncodeToString(keyBytes)
fmt.Println("════════════════════════════════════════")
fmt.Println("SAVE THIS AS NODE_KEY ENVIRONMENT VAR:")
fmt.Println(encoded)
fmt.Println("════════════════════════════════════════")

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
return node, nil
}

func (n *P2PNode) handleStream(s network.Stream) {
defer s.Close()
buf := make([]byte, 1024)
count, err := s.Read(buf)
if err != nil {
return
}
msg := string(buf[:count])
fmt.Printf("[P2P] Message from %s: %s\n", s.Conn().RemotePeer().String()[:12], msg)

response := fmt.Sprintf("AEQUITAS_NODE|humans=%d|chainid=aequitas-1",
n.keeper.TotalHumans())
s.Write([]byte(response))
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
if selfID != "12D3KooWAAYkKMBeChZZRdb4Ydn2hoM2sWsGmAWsQieoi6dzoweY" {
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
