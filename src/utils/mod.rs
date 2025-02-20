pub mod cooldown;
pub mod database;
pub mod embeds;
pub mod message;

pub type Context<'a> = poise::Context<'a, ServerData, Error>;

pub struct ServerData {
    pub cooldowns:
        std::sync::Arc<std::sync::Mutex<std::collections::HashMap<i64, std::time::Instant>>>,
    pub pool: sqlx::PgPool,
}

pub type Error = Box<dyn std::error::Error + Send + Sync>;

pub fn get_guild_id(context: &Context<'_>) -> i64 {
    context
        .guild_id()
        .map(|id| id.to_string().parse::<i64>().unwrap())
        .unwrap()
}

#[inline(always)]
pub fn get_pool<'a>(context: &'a Context<'a>) -> &'a sqlx::PgPool {
    &context.data().pool
}
