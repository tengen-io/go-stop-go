package gql

import (
	"context"
	"errors"
	"github.com/tengen-io/server/models"
	"github.com/tengen-io/server/pubsub"
	"github.com/tengen-io/server/repository"
	"log"
)

type Resolver struct {
	repo     repository.Repository
	pubsub   pubsub.Bus
}

/*func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}*/
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Game() GameResolver {
	return &gameResolver{r}
}

func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type gameResolver struct{ *Resolver }

func (r *gameResolver) Users(ctx context.Context, obj *models.Game) ([]models.GameUserEdge, error) {
	return r.repo.GetUsersForGame(obj.Id)
}

type mutationResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

func (r *queryResolver) User(ctx context.Context, id *string, name *string) (*models.User, error) {
	if id != nil {
		user, err := r.repo.GetUserById(*id)
		if err != nil {
			return nil, err
		}

		return user, nil
	}

	panic("not implemented")
}

func (r *queryResolver) Users(ctx context.Context, ids []string, names []string) ([]*models.User, error) {
	panic("not implemented")
}

func (r *queryResolver) Viewer(ctx context.Context) (*models.Identity, error) {
	identity, ok := ctx.Value(IdentityContextKey).(models.Identity)
	if !ok {
		// TODO(eac): this is asserted already by @hasAuth. Should I just ignore the error?
		return nil, errors.New("invalid user")
	}

	return &identity, nil
}

func (r *queryResolver) Game(ctx context.Context, id *string) (*models.Game, error) {
	if id != nil {
		game, err := r.repo.GetGameById(*id)
		if err != nil {
			return nil, err
		}

		return game, nil
	}

	panic("not implemented")
}

func (r *queryResolver) Games(ctx context.Context, ids []string, states []models.GameState) ([]*models.Game, error) {
	if len(ids) > 0 && len(states) > 0 {
		return nil, errors.New("arguments are mutually exclusive")
	}

	if len(ids) > 0 {
		return r.repo.GetGamesByIds(ids)
	}

	if len(states) > 0 {
		return r.repo.GetGamesByState(states)
	}

	panic("not implemented")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) Games(ctx context.Context, gameType *models.GameType) (<-chan *models.GameSubscriptionPayload, error) {
	topic := "games_" + gameType.String()
	channel := r.pubsub.Subscribe(topic)
	rv := make(chan *models.GameSubscriptionPayload)

	go func() {
		for event := range channel {
			if eventPayload, ok := event.Payload.(*models.Game); ok {
				payload := &models.GameSubscriptionPayload{
					Game:  *eventPayload,
					Event: event.Event,
				}
				rv <- payload
			} else {
				log.Printf("recieved message %+v with unknown payload, expected models.Game", event)
			}
		}
	}()

	return rv, nil
}