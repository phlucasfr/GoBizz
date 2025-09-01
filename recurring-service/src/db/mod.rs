use anyhow::Result;
use sqlx::postgres::PgPoolOptions;

pub type DbPool = sqlx::Pool<sqlx::Postgres>;

pub async fn connect(url: &str) -> Result<DbPool> {
    let pool = PgPoolOptions::new()
        .max_connections(100)
        .connect(url)
        .await?;
    Ok(pool)
}
