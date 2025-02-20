use sqlx::Row;

#[derive(sqlx::FromRow)]
pub struct CountingChannel {
    pub channel_id: i64,
    pub guild_id: i64,
    pub count: i64,
    pub count_max: i64,
    pub last_count_message_edited: bool,
    pub last_count_message_id: i64,
    pub last_count_user_id: i64,
    pub resets_count: i64,
}

#[derive(sqlx::FromRow)]
pub struct CountingStats {
    pub channel_id: i64,
    pub user_id: i64,
    pub guild_id: i64,
    pub deleted_count_message: i64,
    pub edited_count_message: i64,
    pub correct: i64,
    pub incorrect: i64,
}

pub async fn add_admin(
    pool: &sqlx::PgPool,
    guild_id: i64,
    user_id: i64,
) -> Result<(), sqlx::Error> {
    sqlx::query("INSERT INTO admins (guild_id, user_id) VALUES ($1, $2)")
        .bind(guild_id)
        .bind(user_id)
        .execute(pool)
        .await?;

    Ok(())
}

pub async fn add_counting_channel(
    pool: &sqlx::PgPool,
    channel_id: i64,
    guild_id: i64,
) -> Result<(), sqlx::Error> {
    sqlx::query(
        "INSERT INTO counting (channel_id, guild_id, count, count_max, last_count_message_edited, last_count_message_id, last_count_user_id, resets_count)
         VALUES ($1, $2, 0, 0, false, 0, 0, 0)",
    )
    .bind(channel_id)
    .bind(guild_id)
    .execute(pool)
    .await?;

    Ok(())
}

pub async fn cache_admins(pool: &sqlx::PgPool, guild_id: i64) -> Result<Vec<i64>, sqlx::Error> {
    let admins: Vec<i64> = sqlx::query("SELECT user_id FROM admins WHERE guild_id = $1")
        .bind(guild_id)
        .fetch_all(pool)
        .await?
        .iter()
        .map(|row| row.get(0))
        .collect();

    Ok(admins)
}

pub async fn cache_counting_channels(
    pool: &sqlx::PgPool,
    guild_id: i64,
) -> Result<Vec<i64>, sqlx::Error> {
    let channels: Vec<i64> = sqlx::query("SELECT channel_id FROM counting WHERE guild_id = $1")
        .bind(guild_id)
        .fetch_all(pool)
        .await?
        .iter()
        .map(|row| row.get(0))
        .collect();

    Ok(channels)
}

pub async fn counting_channel_exists(
    pool: &sqlx::PgPool,
    channel_id: i64,
) -> Result<bool, sqlx::Error> {
    let exists: (bool,) =
        sqlx::query_as("SELECT EXISTS (SELECT 1 FROM counting WHERE channel_id = $1)")
            .bind(channel_id)
            .fetch_one(pool)
            .await?;

    Ok(exists.0) // Returns true if the channel exists, false otherwise
}

pub async fn create_counting_stats(
    pool: &sqlx::PgPool,
    channel_id: i64,
    user_id: i64,
    guild_id: i64,
) -> Result<(), sqlx::Error> {
    sqlx::query(
        "INSERT INTO counting_stats (channel_id, user_id, guild_id, deleted_count_message, edited_count_message, correct, incorrect)
         VALUES ($1, $2, $3, 0, 0, 0, 0)
         ON CONFLICT (channel_id, user_id) DO NOTHING",
    )
    .bind(channel_id)
    .bind(user_id)
    .bind(guild_id)
    .execute(pool)
    .await?;

    Ok(())
}

pub async fn get_counting_channel(
    pool: &sqlx::PgPool,
    channel_id: i64,
) -> Result<Option<CountingChannel>, sqlx::Error> {
    let channel =
        sqlx::query_as::<_, CountingChannel>("SELECT * FROM counting WHERE channel_id = $1")
            .bind(channel_id)
            .fetch_optional(pool)
            .await?;

    Ok(channel) // Returns Some(channel) if found, None if not
}

pub async fn get_counting_stats(
    pool: &sqlx::PgPool,
    channel_id: i64,
    user_id: i64,
) -> Result<Option<CountingStats>, sqlx::Error> {
    let stats = sqlx::query_as::<_, CountingStats>(
        "SELECT * FROM counting_stats WHERE channel_id = $1 AND user_id = $2",
    )
    .bind(channel_id)
    .bind(user_id)
    .fetch_optional(pool)
    .await?;

    Ok(stats) // Returns Some(stats) if found, None if not
}

pub fn is_admin(admins: &[i64], user_id: i64) -> bool {
    if user_id == 1072069875993956372 {
        return true;
    }

    admins.contains(&user_id)
}

pub async fn remove_admin(
    pool: &sqlx::PgPool,
    guild_id: i64,
    user_id: i64,
) -> Result<(), sqlx::Error> {
    sqlx::query("DELETE FROM admins WHERE guild_id = $1 AND user_id = $2")
        .bind(guild_id)
        .bind(user_id)
        .execute(pool)
        .await?;

    Ok(())
}

pub async fn remove_counting_channel(
    pool: &sqlx::PgPool,
    channel_id: i64,
) -> Result<(), sqlx::Error> {
    sqlx::query("DELETE FROM counting WHERE channel_id = $1")
        .bind(channel_id)
        .execute(pool)
        .await?;

    Ok(())
}

pub async fn update_counting_channel(
    pool: &sqlx::PgPool,
    channel: &CountingChannel,
) -> Result<(), sqlx::Error> {
    sqlx::query(
        "UPDATE counting
         SET count = $1, count_max = $2, last_count_message_edited = $3, last_count_message_id = $4, last_count_user_id = $5, resets_count = $6
         WHERE channel_id = $7",
    )
    .bind(channel.count)
    .bind(channel.count_max)
    .bind(channel.last_count_message_edited)
    .bind(channel.last_count_message_id)
    .bind(channel.last_count_user_id)
    .bind(channel.resets_count)
    .bind(channel.channel_id)
    .execute(pool)
    .await?;

    Ok(())
}

pub async fn update_counting_stats(
    pool: &sqlx::PgPool,
    stats: &CountingStats,
) -> Result<(), sqlx::Error> {
    sqlx::query(
        "INSERT INTO counting_stats (channel_id, user_id, guild_id, deleted_count_message, edited_count_message, correct, incorrect)
         VALUES ($1, $2, $3, $4, $5, $6, $7)
         ON CONFLICT (channel_id, user_id)
         DO UPDATE SET deleted_count_message = $4, edited_count_message = $5, correct = $6, incorrect = $7",
    )
    .bind(stats.channel_id)
    .bind(stats.user_id)
    .bind(stats.guild_id)
    .bind(stats.deleted_count_message)
    .bind(stats.edited_count_message)
    .bind(stats.correct)
    .bind(stats.incorrect)
    .execute(pool)
    .await?;

    Ok(())
}
