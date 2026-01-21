package service

import (
	"errors"
	"log"
)

type LeaderboardData struct {
	actor actorKey
	count uint64
}

type LeaderboardService struct {
	// redis client goes here
	queue    chan LeaderboardData
	resolver DisplayNameResolver
}

func (l *LeaderboardService) Run() error {
	// setup

	// main loop
	for {
		update := <-l.queue
		err := updateLeaderboard(update)
		if err != nil {
			return err
		}

	}
}

func updateLeaderboard(lbd LeaderboardData) error {
	log.Println(lbd)
	return errors.New("unimplemented")
}

func (l *LeaderboardService) UpdateLeaderboard(lbd LeaderboardData) {
	l.queue <- lbd
}
