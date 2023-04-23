package protocol

import (
	"github.com/EdmundMartin/rftgo/internal/persistence"
	"log"
	"sync"
)

type State uint8

const (
	Follower State = iota
	Candidate
	Leader
	Dead
)

func (s State) String() string {
	switch s {
	case Follower:
		return "Follower"
	case Candidate:
		return "Candidate"
	case Leader:
		return "Leader"
	case Dead:
		return "Dead"
	}
	return ""
}

type Node struct {
	id    int
	state State

	mu sync.Mutex

	currentTerm int
	votedFor    int

	log persistence.LogStore

	commitIndex int
	lastApplied int

	nextIndex  map[int]int
	matchIndex map[int]int
}

type AppendEntriesRequest struct {
	Term         int
	LeaderId     int
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []*persistence.LogEntry
	LeaderCommit int
}

type AppendEntriesResponse struct {
	Term    int
	Success bool
}

func (n *Node) AppendEntries(req *AppendEntriesRequest) *AppendEntriesResponse {
	n.mu.Lock()
	defer n.mu.Unlock()

	if req.Term < n.currentTerm {
		return &AppendEntriesResponse{Success: false}
	}

	prevLog, err := n.log.Get(req.PrevLogIndex)
	if err != nil {
		return &AppendEntriesResponse{Success: false}
	}

	if prevLog.Term != req.PrevLogTerm {
		return &AppendEntriesResponse{Success: false}
	}

	if req.Term > n.currentTerm {
		log.Printf("Becoming follower")
		// TODO - Become follower
		return &AppendEntriesResponse{Success: false, Term: n.currentTerm}
	}

	resp := &AppendEntriesResponse{Success: false}

	if n.state != Follower {
		// TODO - Become follower
	}

	for _, entry := range req.Entries {
		// TODO - Should only do updates
		if err = n.log.Save(entry.Term, entry); err != nil {
			return resp
		}
		n.currentTerm = entry.Term
	}

	if req.LeaderCommit > n.commitIndex {
		n.commitIndex = minInt(req.LeaderCommit, n.currentTerm)
	}

	return resp
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type VoteRequest struct {
	Term         int
	CandidateID  int
	LastLogIndex int
	LastLogTerm  int
}

type VoteResponse struct {
	Term        int
	VoteGranted bool
}

func (n *Node) RequestVote(req *VoteRequest) *VoteResponse {
	n.mu.Lock()
	defer n.mu.Unlock()
	if req.Term < n.currentTerm {
		return &VoteResponse{
			Term:        n.currentTerm,
			VoteGranted: false,
		}
	}

	if n.votedFor == -1 {
		return &VoteResponse{
			Term:        n.currentTerm,
			VoteGranted: true,
		}
	}

	return &VoteResponse{
		Term:        n.currentTerm,
		VoteGranted: false,
	}
}
