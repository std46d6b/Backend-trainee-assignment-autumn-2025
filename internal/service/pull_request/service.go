package pullrequestservice

import (
	"context"
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/repository"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/store/postgres"
)

type TxManager interface {
	TxWrapper(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}

type RepoFactory interface {
	PullRequestRepository(exec postgres.Execer) repository.PullRequestRepository
	UserRepository(exec postgres.Execer) repository.UserRepository
	TeamRepository(exec postgres.Execer) repository.TeamRepository
}

type PullRequestService struct {
	txManager TxManager
	repoFact  RepoFactory
	readExec  postgres.Execer
}

func NewPullRequestService(
	txManager TxManager,
	readExec postgres.Execer,
	repoFact RepoFactory,
) *PullRequestService {
	return &PullRequestService{
		txManager: txManager,
		repoFact:  repoFact,
		readExec:  readExec,
	}
}

func shuffleTeamMembers(teamMembers []domain.TeamMember) []domain.TeamMember {
	res := make([]domain.TeamMember, len(teamMembers))
	copy(res, teamMembers)

	rand.Shuffle(len(res), func(i, j int) {
		res[i], res[j] = res[j], res[i]
	})

	return res
}

func (s *PullRequestService) getTeamByUserID(
	ctx context.Context,
	exec postgres.Execer,
	userID string,
) (domain.TeamUpsert, error) {
	localUserRepo := s.repoFact.UserRepository(exec)

	pullRequestAuthor, err := localUserRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("get author: %w", err)
	}

	localTeamRepo := s.repoFact.TeamRepository(exec)

	team, err := localTeamRepo.GetTeamWithMembers(ctx, pullRequestAuthor.TeamName)
	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("get team: %w", err)
	}

	return team, nil
}

func (s *PullRequestService) assignReviewers(ctx context.Context, exec postgres.Execer, pr domain.PullRequest) error {
	team, err := s.getTeamByUserID(ctx, exec, pr.AuthorID)
	if err != nil {
		return fmt.Errorf("get team: %w", err)
	}

	localPullRequestRepo := s.repoFact.PullRequestRepository(exec)

	teamMembers := shuffleTeamMembers(team.Members)
	assignedCounter := len(pr.AssignedReviewers)

	for _, member := range teamMembers {
		if member.UserID != pr.AuthorID && member.IsActive {
			err = localPullRequestRepo.AddReviewer(ctx, pr.ID, member.UserID)
			if err != nil {
				return fmt.Errorf("assign reviewer: %w", err)
			}

			assignedCounter++
			if assignedCounter >= domain.MaxAssignedReviewers {
				break
			}
		}
	}

	return nil
}

// CreatePullRequest may be used for
// POST /pullRequest/create
// creates pull request.
func (s *PullRequestService) CreatePullRequest(ctx context.Context, pr domain.PullRequest) (domain.PullRequest, error) {
	var dbPullRequest domain.PullRequest

	err := s.txManager.TxWrapper(ctx, func(ctx context.Context, tx pgx.Tx) error {
		localUserRepo := s.repoFact.UserRepository(tx)

		_, err := localUserRepo.GetByID(ctx, pr.AuthorID)
		if err != nil {
			return fmt.Errorf("get author: %w", err)
		}

		localPullRequestRepo := s.repoFact.PullRequestRepository(tx)

		err = localPullRequestRepo.InsertPullRequest(ctx, pr)
		if err != nil {
			return fmt.Errorf("insert pull request: %w", err)
		}

		err = s.assignReviewers(ctx, tx, pr)
		if err != nil {
			return fmt.Errorf("assign reviewers: %w", err)
		}

		dbPullRequest, err = localPullRequestRepo.GetByID(ctx, pr.ID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return dbPullRequest, fmt.Errorf("create pull request: %w", err)
	}

	return dbPullRequest, nil
}

// MergePullRequest may be used for
// POST /pullRequest/merge
// merges pull request.
func (s *PullRequestService) MergePullRequest(ctx context.Context, prID string) (domain.PullRequest, error) {
	var pullRequest domain.PullRequest

	err := s.txManager.TxWrapper(ctx, func(ctx context.Context, tx pgx.Tx) error {
		localPullRequestRepo := s.repoFact.PullRequestRepository(tx)

		var err error
		pullRequest, err = localPullRequestRepo.GetByID(ctx, prID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		if pullRequest.Status == domain.PRStatusMerged {
			return nil
		}

		now := time.Now()
		pullRequest.MergedAt = &now
		pullRequest.Status = domain.PRStatusMerged

		if err = localPullRequestRepo.MergePullRequest(ctx, pullRequest); err != nil {
			return fmt.Errorf("service merge pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return pullRequest, fmt.Errorf("merge pull request: %w", err)
	}

	return pullRequest, nil
}

func (s *PullRequestService) reassignRewiewer(
	ctx context.Context,
	tx pgx.Tx,
	pr domain.PullRequest,
	prID string,
	oldReviewerID string,
	localPullRequestRepo repository.PullRequestRepository,
) (string, error) {
	team, err := s.getTeamByUserID(ctx, tx, oldReviewerID)
	if err != nil {
		return "", fmt.Errorf("get team: %w", err)
	}

	for _, member := range team.Members {
		if member.IsActive && member.UserID != oldReviewerID && member.UserID != pr.AuthorID {
			isInReviewers := slices.Contains(pr.AssignedReviewers, member.UserID)
			if !isInReviewers {
				err = localPullRequestRepo.AddReviewer(ctx, prID, member.UserID)
				if err != nil {
					return "", fmt.Errorf("assign reviewer: %w", err)
				}

				return member.UserID, nil
			}
		}
	}

	return "", domain.NewError(domain.ErrCodeNoCandidate, "no active replacement candidate in team")
}

// ReassignPullRequest merges pull request
// POST /pullRequest/reassign
// reassigns pull request.
func (s *PullRequestService) ReassignPullRequest(
	ctx context.Context,
	prID string,
	oldReviewerID string,
) (domain.PullRequest, string, error) {
	var pullRequest domain.PullRequest
	var reassignedUserID string

	err := s.txManager.TxWrapper(ctx, func(ctx context.Context, tx pgx.Tx) error {
		localUserRepo := s.repoFact.UserRepository(tx)

		_, err := localUserRepo.GetByID(ctx, oldReviewerID)
		if err != nil {
			return fmt.Errorf("get old reviewer: %w", err)
		}

		localPullRequestRepo := s.repoFact.PullRequestRepository(tx)

		pr, err := localPullRequestRepo.GetByID(ctx, prID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		if pr.Status == domain.PRStatusMerged {
			return domain.NewError(domain.ErrCodePRMerged, "cannot reassign on merged PR")
		}

		err = localPullRequestRepo.RemoveReviewer(ctx, prID, oldReviewerID)
		if err != nil {
			return fmt.Errorf("remove reviewer: %w", err)
		}

		reassignedUserID, err = s.reassignRewiewer(ctx, tx, pr, prID, oldReviewerID, localPullRequestRepo)
		if err != nil {
			return fmt.Errorf("reassign reviewer: %w", err)
		}

		pullRequest, err = localPullRequestRepo.GetByID(ctx, prID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return pullRequest, reassignedUserID, fmt.Errorf("service reassign pull request: %w", err)
	}

	return pullRequest, reassignedUserID, nil
}
