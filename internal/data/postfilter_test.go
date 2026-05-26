package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func pr(approvals int, viewerReviewed bool) PullRequestData {
	var p PullRequestData
	p.ApprovingReviews.TotalCount = approvals
	if viewerReviewed {
		p.ViewerLatestReview = &ViewerReview{State: "APPROVED"}
	}
	return p
}

func TestPRPostFilter_IsZero(t *testing.T) {
	require.True(t, PRPostFilter{}.IsZero())
	require.False(t, PRPostFilter{ApprovalsLessThan: 1}.IsZero())
	require.False(t, PRPostFilter{NotReviewedByMe: true}.IsZero())
}

func TestPRPostFilter_Matches(t *testing.T) {
	t.Run("zero filter keeps everything", func(t *testing.T) {
		f := PRPostFilter{}
		require.True(t, f.Matches(pr(0, false)))
		require.True(t, f.Matches(pr(5, true)))
	})

	t.Run("approvals less than", func(t *testing.T) {
		f := PRPostFilter{ApprovalsLessThan: 2}
		require.True(t, f.Matches(pr(0, false)))
		require.True(t, f.Matches(pr(1, true)))
		require.False(t, f.Matches(pr(2, false)))
		require.False(t, f.Matches(pr(5, false)))
	})

	t.Run("not reviewed by me", func(t *testing.T) {
		f := PRPostFilter{NotReviewedByMe: true}
		require.True(t, f.Matches(pr(0, false)))
		require.True(t, f.Matches(pr(10, false)))
		require.False(t, f.Matches(pr(0, true)))
	})

	t.Run("both predicates AND together", func(t *testing.T) {
		f := PRPostFilter{ApprovalsLessThan: 2, NotReviewedByMe: true}
		// <2 approvals, not reviewed -> both pass -> keep
		require.True(t, f.Matches(pr(0, false)))
		// >=2 approvals, not reviewed -> approvals fails -> drop
		require.False(t, f.Matches(pr(5, false)))
		// <2 approvals, reviewed -> review fails -> drop
		require.False(t, f.Matches(pr(0, true)))
		// >=2 approvals, reviewed -> both fail -> drop
		require.False(t, f.Matches(pr(5, true)))
	})
}

func TestPRPostFilter_Apply(t *testing.T) {
	prs := []PullRequestData{
		pr(0, false), // both predicates pass -> keep
		pr(5, true),  // both fail -> drop
		pr(1, true),  // review fails -> drop
		pr(3, false), // approvals fails -> drop
	}

	t.Run("zero filter is passthrough", func(t *testing.T) {
		require.Equal(t, prs, PRPostFilter{}.Apply(prs))
	})

	t.Run("filters down to matches", func(t *testing.T) {
		f := PRPostFilter{ApprovalsLessThan: 2, NotReviewedByMe: true}
		got := f.Apply(prs)
		require.Len(t, got, 1)
	})
}
