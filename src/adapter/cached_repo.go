package adapter

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/kennykarnama/my-earth/src/domain"
	"github.com/redis/rueidis"
)

type RedisCached struct {
	repo      domain.LocationQuery
	redisAddr string
}

func NewRedisCached(repo domain.LocationQuery, redisAddr string) *RedisCached {
	return &RedisCached{
		repo:      repo,
		redisAddr: redisAddr,
	}
}

func (r *RedisCached) FindLocationPoint(ctx context.Context, name string) (*domain.ListLocations, error) {
	// check if this query happens
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{
			r.redisAddr,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("redisCached.findLocationPoint err: %w", err)
	}
	defer client.Close()

	key := fmt.Sprintf("location:points:q:%s", name)

	getCmd := client.B().Get().Key(key).Build()

	result := client.Do(ctx, getCmd)

	if result.Error() != nil {
		if result.Error() != rueidis.Nil {
			return nil, fmt.Errorf("redisCached.findLocationPoint err: %w", result.Error())
		}

		matches, err := r.repo.FindLocationPoint(ctx, name)
		if err != nil {
			return nil, fmt.Errorf("redisCached.findLocationPoint err: %w", result.Error())
		}

		// store to hset
		for _, match := range matches.Items {
			hsetMatchKey := fmt.Sprintf("locations:points:data:%s", match.RefID)

			itemHashCmd := client.B().Hset().Key(hsetMatchKey).FieldValue().
				FieldValue("ref_id", match.RefID).
				FieldValue("name", match.Name).
				FieldValue("lat", fmt.Sprintf("%v", match.Lat)).
				FieldValue("lon", fmt.Sprintf("%v", match.Lon)).
				Build()

			res := client.Do(ctx, itemHashCmd)
			if res.Error() != nil {
				return nil, fmt.Errorf("redisCached.findLocationPoint.hset err: %w", result.Error())
			}

			// set expiration
			slog.Info("expiring", slog.Any("key", hsetMatchKey), slog.Any("sec", matches.ExpiredAtUTC().Second()))
			expiredCmd := client.B().Expire().Key(hsetMatchKey).Seconds(int64(matches.ExpiredAtUTC().Second())).Build()

			res = client.Do(ctx, expiredCmd)
			if res.Error() != nil {
				return nil, fmt.Errorf("redisCached.findLocationPoint.expiration err: %w", result.Error())
			}

		}

	}

	// find by cache
	idx := "locations:points:q:matches:idx"

	queryCmd := client.B().FtSearch().Index(idx).Query(name).Build()

	res := client.Do(ctx, queryCmd)
	if res.Error() != nil {
		return nil, fmt.Errorf("redisCached.findLocationPoint.parseResult err: %w", result.Error())
	}

	total, cacheHits, err := res.AsFtSearch()
	if err != nil {
		return nil, fmt.Errorf("redisCached.findLocationPoint.parseResult err: %w", result.Error())
	}

	slog.Info("search", slog.Any("hit", total))

	listLocations := &domain.ListLocations{
		Items: []domain.Location{},
	}

	for _, cacheHit := range cacheHits {
		slog.Info("cacheHit", slog.Any("doc", cacheHit.Doc))

		// construct
		loc := domain.Location{
			Name:  cacheHit.Doc["name"],
			RefID: cacheHit.Doc["ref_id"],
		}

		lats, ok := cacheHit.Doc["lat"]
		if ok {
			lat, err := strconv.ParseFloat(lats, 64)
			if err != nil {
				return nil, fmt.Errorf("redisCached.findLocationPoint.parseResult err: %w", result.Error())
			}
			loc.Lat = lat
		}

		lons, ok := cacheHit.Doc["lon"]
		if ok {
			lon, err := strconv.ParseFloat(lons, 64)
			if err != nil {
				return nil, fmt.Errorf("redisCached.findLocationPoint.parseResult err: %w", result.Error())
			}
			loc.Lon = lon
		}

		listLocations.Items = append(listLocations.Items, loc)

	}

	return listLocations, nil
}
