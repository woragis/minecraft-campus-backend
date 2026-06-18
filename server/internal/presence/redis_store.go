package presence

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(client *redis.Client, ttlSeconds int) *RedisStore {
	if ttlSeconds <= 0 {
		ttlSeconds = 120
	}
	return &RedisStore{
		client: client,
		ttl:    time.Duration(ttlSeconds) * time.Second,
	}
}

func (s *RedisStore) Enabled() bool { return true }

func serverKey(slug string) string {
	return "cw:presence:server:" + strings.ToLower(strings.TrimSpace(slug))
}

func guildKey(guildID uuid.UUID) string {
	return "cw:presence:guild:" + guildID.String()
}

func playerKey(playerID uuid.UUID) string {
	return "cw:presence:player:" + playerID.String()
}

func (s *RedisStore) MarkOnline(ctx context.Context, in MarkOnlineInput) error {
	now := float64(time.Now().UTC().Unix())
	serverSlug := strings.ToLower(strings.TrimSpace(in.ServerSlug))
	pipe := s.client.Pipeline()
	pipe.ZAdd(ctx, serverKey(serverSlug), redis.Z{Score: now, Member: in.PlayerID.String()})
	meta := map[string]interface{}{
		"username":   in.Username,
		"serverSlug": serverSlug,
		"platform":   strings.TrimSpace(in.Platform),
		"since":      strconv.FormatInt(int64(now), 10),
	}
	if in.GuildID != nil && *in.GuildID != uuid.Nil {
		meta["guildId"] = in.GuildID.String()
		pipe.ZAdd(ctx, guildKey(*in.GuildID), redis.Z{Score: now, Member: in.PlayerID.String()})
	}
	pipe.HSet(ctx, playerKey(in.PlayerID), meta)
	pipe.Expire(ctx, playerKey(in.PlayerID), s.ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("presence mark online: %w", err)
	}
	return nil
}

func (s *RedisStore) MarkOffline(ctx context.Context, playerID uuid.UUID, serverSlug string, guildID *uuid.UUID) error {
	serverSlug = strings.ToLower(strings.TrimSpace(serverSlug))
	pipe := s.client.Pipeline()
	pipe.ZRem(ctx, serverKey(serverSlug), playerID.String())
	pipe.Del(ctx, playerKey(playerID))
	if guildID != nil && *guildID != uuid.Nil {
		pipe.ZRem(ctx, guildKey(*guildID), playerID.String())
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("presence mark offline: %w", err)
	}
	return nil
}

func (s *RedisStore) Heartbeat(ctx context.Context, playerID uuid.UUID, serverSlug string) error {
	now := float64(time.Now().UTC().Unix())
	serverSlug = strings.ToLower(strings.TrimSpace(serverSlug))
	pipe := s.client.Pipeline()
	pipe.ZAdd(ctx, serverKey(serverSlug), redis.Z{Score: now, Member: playerID.String()})
	pipe.Expire(ctx, playerKey(playerID), s.ttl)
	data, err := s.client.HGetAll(ctx, playerKey(playerID)).Result()
	if err == nil && data["guildId"] != "" {
		if gid, err := uuid.Parse(data["guildId"]); err == nil {
			pipe.ZAdd(ctx, guildKey(gid), redis.Z{Score: now, Member: playerID.String()})
		}
	}
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("presence heartbeat: %w", err)
	}
	return nil
}

func (s *RedisStore) minScore() float64 {
	return float64(time.Now().UTC().Add(-s.ttl).Unix())
}

func (s *RedisStore) ListServerOnline(ctx context.Context, serverSlug string) ([]OnlinePlayer, error) {
	return s.loadPlayers(ctx, serverKey(strings.ToLower(strings.TrimSpace(serverSlug))))
}

func (s *RedisStore) ListServersOnline(ctx context.Context, serverSlugs []string) ([]ServerPresence, error) {
	out := make([]ServerPresence, 0, len(serverSlugs))
	for _, slug := range serverSlugs {
		slug = strings.ToLower(strings.TrimSpace(slug))
		if slug == "" {
			continue
		}
		players, err := s.ListServerOnline(ctx, slug)
		if err != nil {
			return nil, err
		}
		out = append(out, ServerPresence{
			Slug:        slug,
			OnlineCount: len(players),
			Players:     players,
		})
	}
	return out, nil
}

func (s *RedisStore) CountGuildOnline(ctx context.Context, guildID uuid.UUID, memberIDs []uuid.UUID) ([]OnlinePlayer, error) {
	online, err := s.loadPlayers(ctx, guildKey(guildID))
	if err != nil {
		return nil, err
	}
	if len(memberIDs) == 0 {
		return online, nil
	}
	allowed := make(map[uuid.UUID]struct{}, len(memberIDs))
	for _, id := range memberIDs {
		allowed[id] = struct{}{}
	}
	filtered := make([]OnlinePlayer, 0, len(online))
	for _, p := range online {
		if _, ok := allowed[p.PlayerID]; ok {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

func (s *RedisStore) loadPlayers(ctx context.Context, zsetKey string) ([]OnlinePlayer, error) {
	min := s.minScore()
	ids, err := s.client.ZRangeByScore(ctx, zsetKey, &redis.ZRangeBy{
		Min: strconv.FormatFloat(min, 'f', 0, 64),
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("presence list %s: %w", zsetKey, err)
	}
	players := make([]OnlinePlayer, 0, len(ids))
	for _, idStr := range ids {
		playerID, err := uuid.Parse(idStr)
		if err != nil {
			continue
		}
		meta, err := s.client.HGetAll(ctx, playerKey(playerID)).Result()
		if err != nil || len(meta) == 0 {
			continue
		}
		p := OnlinePlayer{
			PlayerID:   playerID,
			Username:   meta["username"],
			ServerSlug: meta["serverSlug"],
			Platform:   meta["platform"],
		}
		if since, err := strconv.ParseInt(meta["since"], 10, 64); err == nil {
			p.Since = time.Unix(since, 0).UTC()
		}
		if gid := meta["guildId"]; gid != "" {
			if parsed, err := uuid.Parse(gid); err == nil {
				p.GuildID = &parsed
			}
		}
		players = append(players, p)
	}
	return players, nil
}
