package keeper

import (
"context"
"fmt"
"time"

libp2p "github.com/libp2p/go-libp2p"
"github.com/libp2p/go-libp2p/core/host"
"github.com/libp2p/go-libp2p/core/network"
"github.com/libp2p/go-libp2p/core/peer"
"github.com/libp2p/go-libp2p/core/protocol"
)

const (
ProtocolID = "/aequitas/1.0.0"
ListenPort = 4001
)

type P2PNode struct {
host    host.Host
keeper  *Keeper
peers   []peer.AddrInfo
}

func NewP2PNode(keeper *Keeper) (*P2PNode, error) {
h, err := libp2p.New(
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

// Register protocol handler
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

// Respond with node status
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
