package config

import (
    "time"

    "github.com/danivideda/satu-apotek-be/internal/env"
)

type Config struct {
    Addr   string
    DB     dbConfig
    CORS   CORSConfig
    Auth   AuthConfig
    Job    JobConfig
}

type dbConfig struct {
    URL string
}

type CORSConfig struct {
    Origins []string
}

type AuthConfig struct {
    OwnerSessionTTL    time.Duration
    UserSessionTTL     time.Duration
    PharmacySessionTTL time.Duration
    CacheSessionTTL    time.Duration
    CodeTTL            time.Duration
    CSRFSecret         string
}

type JobConfig struct {
    Enabled              bool
    DeleteExpiredSession time.Duration
    ClearApotekCode      time.Duration
}

func Load() Config {
    return Config{
        Addr: env.GetString("ADDR", "localhost:8080"),
        DB: dbConfig{
            URL: env.GetString("DATABASE_URL", "postgres://admin:adminpassword@localhost/satuapotek?sslmode=disable"),
        },
        CORS: CORSConfig{
            Origins: []string{"http://localhost:3000", "http://localhost:4173"},
        },
        Auth: AuthConfig{
            OwnerSessionTTL:    env.GetDuration("OWNER_SESSION_TTL", 168*time.Hour),
            UserSessionTTL:     env.GetDuration("USER_SESSION_TTL", 168*time.Hour),
            PharmacySessionTTL: env.GetDuration("PHARMACY_SESSION_TTL", 168*time.Hour),
            CacheSessionTTL:    env.GetDuration("CACHE_SESSION_TTL", 5*time.Minute),
            CodeTTL:            env.GetDuration("CODE_TTL", 5*time.Minute),
            CSRFSecret:         env.GetString("CSRF_SECRET", "set-your-secret-in-env-var"),
        },
        Job: JobConfig{
            Enabled:              env.GetBool("RUN_JOB", false),
            DeleteExpiredSession: env.GetDuration("CRON_DURATION_DEL_EXP_SESSION", 10*time.Minute),
            ClearApotekCode:      env.GetDuration("CRON_DURATION_CLEAR_APTK_CODE", 10*time.Minute),
        },
    }
}