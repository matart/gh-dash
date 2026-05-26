package data

// PRPostFilter applies client-side predicates to a PR list after it comes back
// from GitHub. Used for conditions GitHub's search syntax can't express, such
// as "fewer than N approvals". A PR is kept only if it satisfies every enabled
// predicate (AND semantics); predicates left at their zero value are ignored.
type PRPostFilter struct {
	// ApprovalsLessThan keeps PRs with fewer than this many approving reviews.
	// Zero means the predicate is disabled.
	ApprovalsLessThan int

	// NotReviewedByMe keeps PRs the viewer has not reviewed yet.
	NotReviewedByMe bool
}

func (f PRPostFilter) IsZero() bool {
	return f.ApprovalsLessThan == 0 && !f.NotReviewedByMe
}

func (f PRPostFilter) Matches(pr PullRequestData) bool {
	if f.ApprovalsLessThan > 0 && pr.ApprovingReviews.TotalCount >= f.ApprovalsLessThan {
		return false
	}
	if f.NotReviewedByMe && pr.ViewerLatestReview != nil {
		return false
	}
	return true
}

func (f PRPostFilter) Apply(prs []PullRequestData) []PullRequestData {
	if f.IsZero() {
		return prs
	}
	out := make([]PullRequestData, 0, len(prs))
	for _, pr := range prs {
		if f.Matches(pr) {
			out = append(out, pr)
		}
	}
	return out
}
