use anyhow::{Result, anyhow};
use uuid::Uuid;

use crate::models::events::{CreateEvent, Event, UpdateEvent};
use crate::models::occurrences::OccurrenceDTO;
use crate::repositories::events_repository::EventsRepository;
use chrono::Duration;
use rayon::prelude::*;
use sqlx::types::chrono::NaiveDate;

#[derive(Clone)]
pub struct EventsService {
    repo: EventsRepository,
}

impl EventsService {
    pub fn new(repo: EventsRepository) -> Self {
        Self { repo }
    }

    pub async fn create(&self, customer_id: Uuid, input: CreateEvent) -> Result<Event> {
        tracing::info!("Creating event for customer_id: {}", customer_id);
        if input.interval_days <= 0 {
            return Err(anyhow!("interval_days must be > 0"));
        }
        if let Some(stop) = input.stop_at {
            if stop < input.start_date {
                return Err(anyhow!("stop_at must be >= start_date"));
            }
        }
        self.repo.create(customer_id, input).await
    }

    pub async fn get(&self, customer_id: Uuid, id: Uuid) -> Result<Option<Event>> {
        tracing::info!("Fetching event {} for customer_id: {}", id, customer_id);
        self.repo.get(customer_id, id).await
    }

    pub async fn list(&self, customer_id: Uuid, limit: i64, offset: i64) -> Result<Vec<Event>> {
        self.repo.list(customer_id, limit, offset).await
    }

    pub async fn update(
        &self,
        customer_id: Uuid,
        id: Uuid,
        patch: UpdateEvent,
    ) -> Result<Option<Event>> {
        tracing::info!("Updating event {} for customer_id: {}", id, customer_id);

        let current = match self.repo.get(customer_id, id).await? {
            Some(ev) => ev,
            None => return Ok(None),
        };

        let new_start = patch.start_date.unwrap_or(current.start_date);
        let new_interval = patch.interval_days.unwrap_or(current.interval_days);
        let new_stop = match patch.stop_at.clone() {
            Some(v) => v,
            None => None,
        };

        if let Some(name) = &patch.name {
            if name.trim().is_empty() {
                return Err(anyhow!("name cannot be empty"));
            }
        }

        if new_interval <= 0 {
            return Err(anyhow!("interval_days must be > 0"));
        }

        tracing::info!(
            "Post-update values for event {}: start_date={}, interval_days={}, stop_at={:?}",
            id,
            new_start,
            new_interval,
            new_stop
        );
        if let Some(stop) = new_stop {
            if stop < new_start {
                return Err(anyhow!("stop_at must be >= start_date"));
            }
        }

        let patch2 = UpdateEvent {
            name: patch.name,
            start_date: patch.start_date,
            interval_days: patch.interval_days,
            stop_at: Some(new_stop),
        };

        let updated = self.repo.update(customer_id, id, patch2).await?;
        
        Ok(updated)
    }

    pub async fn delete(&self, customer_id: Uuid, id: Uuid) -> Result<bool> {
        tracing::info!("Deleting event {} for customer_id: {}", id, customer_id);
        self.repo.delete(customer_id, id).await
    }

    pub async fn list_occurrences(
        &self,
        customer_id: Uuid,
        start: NaiveDate,
        end: NaiveDate,
        name: Option<String>,
    ) -> Result<Vec<OccurrenceDTO>> {
        tracing::info!(
            "Listing occurrences for customer_id: {} from {} to {}",
            customer_id,
            start,
            end
        );
        if start > end {
            return Err(anyhow!("start deve ser anterior ou igual a end"));
        }

        let name = name.and_then(|v| {
            let t = v.trim();
            if t.is_empty() || t == "undefined" || t == "null" {
                None
            } else {
                Some(t.to_string())
            }
        });

        let events = self
            .repo
            .find_all_for_occurrences(customer_id, name)
            .await?;

        tracing::info!("Found {} events for occurrences computation", events.len());

        let workers_enabled = std::env::var("WORKERS_ENABLED")
            .unwrap_or_else(|_| "false".into())
            .to_lowercase()
            == "true";

        let threshold: usize = std::env::var("WORKER_MIN_ITEMS")
            .ok()
            .and_then(|s| s.parse().ok())
            .unwrap_or(64);

        let compute = |evs: &[Event]| -> Vec<OccurrenceDTO> {
            let mut out = Vec::new();
            for ev in evs {
                if ev.interval_days <= 0 {
                    continue;
                }

                let ev_start = ev.start_date;
                let cut = ev.stop_at;

                let seed = if start > ev_start { start } else { ev_start };

                let diff_days = (seed - ev_start).num_days();
                let interval = ev.interval_days as i64;

                let rem = diff_days.rem_euclid(interval);
                let first = if rem == 0 {
                    seed
                } else {
                    seed + Duration::days(interval - rem)
                };

                let mut cursor = first;
                while cursor <= end {
                    if let Some(c) = cut {
                        if cursor >= c {
                            break;
                        }
                    }
                    out.push(OccurrenceDTO {
                        event_id: ev.id,
                        name: ev.name.clone(),
                        date: cursor,
                    });
                    cursor = cursor + Duration::days(interval);
                }
            }
            out
        };

        let mut occs = if workers_enabled && events.len() >= threshold {
            let threads = rayon::current_num_threads().max(1);
            let chunk_size = (events.len() / threads).max(1);
            events
                .par_chunks(chunk_size)
                .map(|chunk| compute(chunk))
                .flatten()
                .collect::<Vec<_>>()
        } else {
            compute(&events)
        };

        occs.sort_by(|a, b| {
            if a.date == b.date {
                a.name.cmp(&b.name)
            } else {
                a.date.cmp(&b.date)
            }
        });

        Ok(occs)
    }

    pub async fn cut_from(
        &self,
        customer_id: Uuid,
        id: Uuid,
        from: NaiveDate,
    ) -> Result<Option<Event>> {
        let current = match self.repo.get(customer_id, id).await? {
            Some(ev) => ev,
            None => return Ok(None),
        };

        let new_stop = match current.stop_at {
            Some(s) => {
                if s < from {
                    s
                } else {
                    from
                }
            }
            None => from,
        };

        let patch = UpdateEvent {
            name: None,
            start_date: None,
            interval_days: None,
            stop_at: Some(Some(new_stop)),
        };

        let updated = self.repo.update(customer_id, id, patch).await?;
        Ok(updated)
    }
}
