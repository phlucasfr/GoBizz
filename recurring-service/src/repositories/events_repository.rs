use anyhow::Result;
use sqlx::{query, query_as};
use uuid::Uuid;

use crate::db::DbPool;
use crate::models::events::{CreateEvent, Event, UpdateEvent};

#[derive(Clone)]
pub struct EventsRepository {
    pool: DbPool,
}

impl EventsRepository {
    pub fn new(pool: DbPool) -> Self {
        Self { pool }
    }

    pub async fn create(&self, customer_id: Uuid, input: CreateEvent) -> Result<Event> {
        let rec = query_as::<_, Event>(
            r#"
            INSERT INTO events (customer_id, name, start_date, interval_days, stop_at)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, customer_id, name, start_date, interval_days, stop_at, created_at, updated_at
            "#,
        )
        .bind(customer_id)
        .bind(&input.name)
        .bind(input.start_date)
        .bind(input.interval_days)
        .bind(input.stop_at)
        .fetch_one(&self.pool)
        .await?;
        Ok(rec)
    }

    pub async fn get(&self, customer_id: Uuid, id: Uuid) -> Result<Option<Event>> {
        let rec = query_as::<_, Event>(
            r#"
            SELECT id, customer_id, name, start_date, interval_days, stop_at, created_at, updated_at
            FROM events
            WHERE id = $1 AND customer_id = $2
            "#,
        )
        .bind(id)
        .bind(customer_id)
        .fetch_optional(&self.pool)
        .await?;
        Ok(rec)
    }

    pub async fn list(&self, customer_id: Uuid, limit: i64, offset: i64) -> Result<Vec<Event>> {
        let rows = query_as::<_, Event>(
            r#"
            SELECT id, customer_id, name, start_date, interval_days, stop_at, created_at, updated_at
            FROM events
            WHERE customer_id = $1
            ORDER BY created_at DESC
            LIMIT $2 OFFSET $3
            "#,
        )
        .bind(customer_id)
        .bind(limit)
        .bind(offset)
        .fetch_all(&self.pool)
        .await?;
        Ok(rows)
    }

    pub async fn update(
        &self,
        customer_id: Uuid,
        id: Uuid,
        input: UpdateEvent,
    ) -> Result<Option<Event>> {
        let stop_at_set: bool = input.stop_at.is_some();
        let stop_at_value = input.stop_at.flatten(); // Option<NaiveDate>

        let rec = sqlx::query_as::<_, Event>(
            r#"
            UPDATE events SET
                name          = COALESCE($3, name),
                start_date    = COALESCE($4, start_date),
                interval_days = COALESCE($5, interval_days),
                stop_at       = CASE WHEN $6 THEN $7 ELSE stop_at END,
                updated_at    = NOW()
            WHERE id = $1 AND customer_id = $2
            RETURNING id, customer_id, name, start_date, interval_days, stop_at, created_at, updated_at
            "#,
        )
        .bind(id)                  // $1
        .bind(customer_id)         // $2
        .bind(input.name)          // $3
        .bind(input.start_date)    // $4
        .bind(input.interval_days) // $5
        .bind(stop_at_set)         // $6
        .bind(stop_at_value)       // $7
        .fetch_optional(&self.pool)
        .await?;

        Ok(rec)
    }

    pub async fn delete(&self, customer_id: Uuid, id: Uuid) -> Result<bool> {
        let res = query(
            r#"
            DELETE FROM events
            WHERE id = $1 AND customer_id = $2
            "#,
        )
        .bind(id)
        .bind(customer_id)
        .execute(&self.pool)
        .await?;
        Ok(res.rows_affected() > 0)
    }

    pub async fn find_all_for_occurrences(
        &self,
        customer_id: Uuid,
        name: Option<String>,
    ) -> Result<Vec<Event>> {
        let rows = if let Some(n) = name.filter(|s| !s.trim().is_empty()) {
            let like = format!("%{}%", n.trim());
            query_as::<_, Event>(
                r#"
                SELECT id, customer_id, name, start_date, interval_days, stop_at, created_at, updated_at
                FROM events
                WHERE customer_id = $1
                  AND name ILIKE $2
                ORDER BY created_at ASC
                "#,
            )
            .bind(customer_id)
            .bind(like)
            .fetch_all(&self.pool)
            .await?
        } else {
            query_as::<_, Event>(
                r#"
                SELECT id, customer_id, name, start_date, interval_days, stop_at, created_at, updated_at
                FROM events
                WHERE customer_id = $1
                ORDER BY created_at ASC
                "#,
            )
            .bind(customer_id)
            .fetch_all(&self.pool)
            .await?
        };
        Ok(rows)
    }
}
