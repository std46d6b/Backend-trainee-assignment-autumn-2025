package prsvc

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/domain"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/postgres"
	"github.com/std46d6b/Backend-trainee-assignment-autumn-2025/internal/repository"
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

func (s *PullRequestService) getTeamByUserID(ctx context.Context, exec postgres.Execer, userID domain.UserID) (domain.TeamUpsert, error) {
	localUserRepo := s.repoFact.UserRepository(exec)

	pull_request_author, err := localUserRepo.GetByID(ctx, userID)
	if err != nil {
		return domain.TeamUpsert{}, fmt.Errorf("get author: %w", err)
	}

	localTeamRepo := s.repoFact.TeamRepository(exec)

	team, err := localTeamRepo.GetTeamWithMembers(ctx, pull_request_author.Team)
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
			err := localPullRequestRepo.AddReviewer(ctx, pr.ID, member.UserID)
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

// POST /pullRequest/create
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

// POST /pullRequest/merge
func (s *PullRequestService) MergePullRequest(ctx context.Context, prID domain.PullRequestID) (domain.PullRequest, error) {
	var pullRequest domain.PullRequest

	err := s.txManager.TxWrapper(ctx, func(ctx context.Context, tx pgx.Tx) error {
		localPullRequestRepo := s.repoFact.PullRequestRepository(tx)

		pullRequest, err := localPullRequestRepo.GetByID(ctx, prID)
		if err != nil {
			return fmt.Errorf("get pull request: %w", err)
		}

		if pullRequest.Status == domain.PRStatusMerged {
			return nil
		}

		now := time.Now()
		pullRequest.MergedAt = &now
		pullRequest.Status = domain.PRStatusMerged

		err = localPullRequestRepo.MergePullRequest(ctx, pullRequest)
		if err != nil {
			return fmt.Errorf("service merge pull request: %w", err)
		}

		return nil
	})

	if err != nil {
		return pullRequest, fmt.Errorf("merge pull request: %w", err)
	}

	return pullRequest, nil
}

// POST /pullRequest/reassign
func (s *PullRequestService) ReassignPullRequest(ctx context.Context, prID domain.PullRequestID, oldReviewerID domain.UserID) (domain.PullRequest, domain.UserID, error) {
	var pullRequest domain.PullRequest
	var reassignedUserID domain.UserID

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
			return domain.NewDomainError(domain.ErrCodePRMerged, "cannot reassign on merged PR")
		}

		err = localPullRequestRepo.RemoveReviewer(ctx, prID, oldReviewerID)
		if err != nil {
			return fmt.Errorf("remove reviewer: %w", err)
		}

		team, err := s.getTeamByUserID(ctx, tx, oldReviewerID)
		if err != nil {
			return fmt.Errorf("get team: %w", err)
		}

		reassignedUserIsSet := false
		for _, member := range team.Members {
			if member.IsActive && member.UserID != oldReviewerID && member.UserID != pr.AuthorID {
				isInReviewers := false
				for _, reviewerID := range pr.AssignedReviewers {
					if reviewerID == member.UserID {
						isInReviewers = true
						break
					}
				}
				if !isInReviewers {
					err := localPullRequestRepo.AddReviewer(ctx, prID, member.UserID)
					if err != nil {
						return fmt.Errorf("assign reviewer: %w", err)
					}

					reassignedUserID = member.UserID
					reassignedUserIsSet = true
					break
				}
			}
		}

		if !reassignedUserIsSet {
			return domain.NewDomainError(domain.ErrCodeNoCandidate, "no active replacement candidate in team")
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
