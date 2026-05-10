package user

const (
	ReputationPerPostUpvote    = 5
	ReputationPerCommentUpvote = 2
	ReputationPerPostDownvote  = -2
	ReputationPerCommentDownvote = -1
)

type Badge struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Threshold   int    `json:"threshold"`
}

var badgeDefinitions = []Badge{
	{Name: "Newcomer", Description: "Joined the platform", Threshold: 0},
	{Name: "Contributor", Description: "Earned 10 reputation", Threshold: 10},
	{Name: "Active Member", Description: "Earned 50 reputation", Threshold: 50},
	{Name: "Trusted Voice", Description: "Earned 100 reputation", Threshold: 100},
	{Name: "Rising Star", Description: "Earned 250 reputation", Threshold: 250},
	{Name: "Expert", Description: "Earned 500 reputation", Threshold: 500},
	{Name: "Authority", Description: "Earned 1000 reputation", Threshold: 1000},
	{Name: "Legend", Description: "Earned 5000 reputation", Threshold: 5000},
}

func CalculateLevel(reputation int) int {
	switch {
	case reputation >= 5000:
		return 8
	case reputation >= 1000:
		return 7
	case reputation >= 500:
		return 6
	case reputation >= 250:
		return 5
	case reputation >= 100:
		return 4
	case reputation >= 50:
		return 3
	case reputation >= 10:
		return 2
	default:
		return 1
	}
}

func CalculateBadges(reputation int) []string {
	badges := make([]string, 0)
	for _, b := range badgeDefinitions {
		if reputation >= b.Threshold {
			badges = append(badges, b.Name)
		}
	}
	return badges
}
