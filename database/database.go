package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ftqo/kirby/logger"
)

var pool *pgxpool.Pool

func Open(host, port, user, pass, database string) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, pass, database)

	p, err := pgxpool.Connect(context.TODO(), dsn)
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to open connextion pool")
	}
	err = p.Ping(context.TODO())
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to connect to connection pool")
	}
	pool = p
	logger.L.Info().Msg("Connected to database")
	initDatabase()
}

func Close() {
	logger.L.Info().Msg("Disconnected from database")
	pool.Close()
}

func initDatabase() {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Panic().Err(err).Msg("Failed to acquire connection for database initialization")
	}
	defer conn.Release()
	statement := `
	CREATE TABLE IF NOT EXISTS guild_welcome (
		guild_id TEXT PRIMARY KEY,
		channel_id TEXT NOT NULL,
		type TEXT NOT NULL,
		message_text TEXT NOT NULL,
		image TEXT NOT NULL,
		image_text TEXT NOT NULL
	);`
	_, err = conn.Exec(context.TODO(), statement)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for initialize database")
	}
}

func InitGuild(guildId string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for guild initialization")
	}
	defer conn.Release()

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for guild initialization")
	}
	dgw := NewDefaultGuildWelcome()
	statement := `
	INSERT INTO guild_welcome (guild_id, channel_id, type, message_text, image, image_text)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (guild_id) DO NOTHING`
	_, err = tx.Exec(context.TODO(), statement, guildId, dgw.ChannelID, dgw.Type, dgw.Text, dgw.Image, dgw.ImageText)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for guild initialization")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for guild initialization")
	}
}

func CutGuild(guildId string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for guild cutting")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for guild cutting")
	}
	statement := `
	DELETE FROM guild_welcome WHERE guild_id = $1`
	_, err = tx.Exec(context.TODO(), statement, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for guild cutting")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for guild cutting")
	}
}

func ResetGuild(guildId string) {
	CutGuild(guildId)
	InitGuild(guildId)
}

func GetGuildWelcome(guildId string) GuildWelcome {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for guild welcome")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for guild welcome")
	}
	statement := `
	SELECT guild_id, channel_id, type, message_text, image, image_text FROM guild_welcome WHERE guild_id = $1`
	row := tx.QueryRow(context.TODO(), statement, guildId)
	gw := GuildWelcome{}
	err = row.Scan(&gw.GuildID, &gw.ChannelID, &gw.Type, &gw.Text, &gw.Image, &gw.ImageText)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to scan query for guild welcome")
	}
	return gw
}

func SetGuildWelcomeChannel(guildId, channelId string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for set guild welcome channel")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for set guild welcome channel")
	}
	statement := `
	UPDATE guild_welcome SET channel_id = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, channelId, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for set guild welcome channel")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for set guild welcome channel")
	}
}

func SetGuildWelcomeType(guildId, welcomeType string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for set guild welcome image")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for set guild welcome image")
	}
	statement := `
	UPDATE guild_welcome SET type = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, welcomeType, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for set guild welcome image")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for set guild welcome image")
	}
}

func SetGuildWelcomeText(guildId, messageText string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for set guild welcome message text")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for set guild welcome message text")
	}
	statement := `
	UPDATE guild_welcome SET message_text = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, messageText, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for set guild welcome message text")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for set guild welcome message text")
	}
}

func SetGuildWelcomeImage(guildId, image string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for set guild welcome image")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for set guild welcome image")
	}
	statement := `
	UPDATE guild_welcome SET image = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, image, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for set guild welcome image")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for set guild welcome image")
	}
}

func SetGuildWelcomeImageText(guildId, imageText string) {
	conn, err := pool.Acquire(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to acquire connection for set guild welcome image text")
	}
	defer conn.Release()
	tx, err := conn.Begin(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to begin transaction for set guild welcome image text")
	}
	statement := `
	UPDATE guild_welcome SET image_text = $1 WHERE guild_id = $2`
	_, err = tx.Exec(context.TODO(), statement, imageText, guildId)
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to execute statement for set guild welcome image text")
	}
	err = tx.Commit(context.TODO())
	if err != nil {
		logger.L.Error().Err(err).Msg("Failed to commit transaction for set guild welcome image text")
	}
}
