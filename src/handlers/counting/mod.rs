use crate::utils;
use poise::serenity_prelude as serenity;
pub mod event;

#[poise::command(
    prefix_command,
    slash_command,
    subcommands("add", "list", "remove", "stats"),
    subcommand_required
)]
pub async fn counting(_: utils::Context<'_>) -> Result<(), utils::Error> {
    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn add(
    context: utils::Context<'_>,
    channel: serenity::Channel,
) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let author_id = i64::from(context.author().id);
    let guild_id = crate::utils::get_guild_id(&context);

    let embed;

    if crate::utils::database::is_admin(
        &crate::utils::database::cache_admins(pool, guild_id)
            .await
            .expect("Failed to cache admins in handlers/counting@add"),
        author_id,
    ) {
        let channel_id = i64::from(channel.id());

        if crate::utils::database::counting_channel_exists(pool, channel_id)
            .await
            .expect("Failed to check if counting channel exists in handlers/counting@add")
        {
            embed = crate::utils::embeds::create(
                serenity::Colour::new(13789294),
                "Error",
                format!("<#{}> is already a counting channel.", channel_id).as_str(),
                true,
            );
        } else {
            crate::utils::database::add_counting_channel(pool, channel_id, guild_id)
                .await
                .expect("Failed to add counting channel in handlers/counting@add");

            embed = crate::utils::embeds::create(
                serenity::Colour::new(5082199),
                "Success",
                format!("<#{}> is now a counting channel.", channel_id).as_str(),
                true,
            );
        }
    } else {
        embed = crate::utils::embeds::insufficient_permissions();
    }

    context
        .send(
            poise::CreateReply::default()
                .embed(embed.clone())
                .ephemeral(true),
        )
        .await
        .expect("Failed to send message in handlers/counting@add");

    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn list(context: utils::Context<'_>) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let guild_id = crate::utils::get_guild_id(&context);

    let counting_channels = crate::utils::database::cache_counting_channels(pool, guild_id)
        .await
        .expect("Failed to cache counting channels in handlers/counting@list");

    let embed = crate::utils::embeds::create(
        serenity::Colour::new(5082199),
        "Counting Channels",
        counting_channels
            .iter()
            .map(|channel_id| format!("* <#{}>", channel_id))
            .collect::<Vec<String>>()
            .join("\n")
            .as_str(),
        true,
    );

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/counting@list");

    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn remove(
    context: utils::Context<'_>,
    channel: serenity::Channel,
) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let author_id = i64::from(context.author().id);
    let guild_id = crate::utils::get_guild_id(&context);

    let embed;

    if crate::utils::database::is_admin(
        &crate::utils::database::cache_admins(pool, guild_id)
            .await
            .expect("Failed to cache admins in handlers/counting@remove"),
        author_id,
    ) {
        let channel_id = i64::from(channel.id());

        if crate::utils::database::counting_channel_exists(pool, channel_id)
            .await
            .expect("Failed to check if counting channel exists in handlers/counting@remove")
        {
            crate::utils::database::remove_counting_channel(pool, channel_id)
                .await
                .expect("Failed to remove counting channel in handlers/counting@remove");

            embed = crate::utils::embeds::create(
                serenity::Colour::new(5082199),
                "Success",
                format!("<#{}> is no longer a counting channel.", channel_id).as_str(),
                true,
            );
        } else {
            embed = crate::utils::embeds::create(
                serenity::Colour::new(13789294),
                "Error",
                format!("<#{}> is not a counting channel.", channel_id).as_str(),
                true,
            );
        }
    } else {
        embed = crate::utils::embeds::insufficient_permissions();
    }

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/counting@remove");

    Ok(())
}

#[poise::command(
    prefix_command,
    slash_command,
    subcommands("channel", "guild", "user"),
    subcommand_required
)]
pub async fn stats(_: utils::Context<'_>) -> Result<(), utils::Error> {
    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn channel(
    context: utils::Context<'_>,
    channel: serenity::Channel,
) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let channel_id = i64::from(channel.id());

    let result = crate::utils::database::get_counting_channel(pool, channel_id)
        .await
        .expect("Failed to get counting channel in handlers/counting@channel");

    let embed;

    if result.is_none() {
        embed = crate::utils::embeds::create(
            serenity::Colour::new(13789294),
            "Error",
            format!("<#{}> is not a counting channel.", channel_id).as_str(),
            true,
        );
    } else {
        let result = result.unwrap();

        embed = crate::utils::embeds::create(
            serenity::Colour::new(6591981),
            format!("Stats for <#{}>", channel_id).as_str(),
            format!(
                "The last user to count was <@{}> at [this message](https://discord.com/channels/{}/{}/{}).\n* **Current Count**: `{}`.\n* **Max Count**: `{}`.\n* **Resets Count**: `{}`.",
                result.last_count_user_id, context.guild_id().unwrap(), channel_id, result.last_count_message_id,
                result.count, result.count_max, result.resets_count
            )
            .as_str(),
            true
        );
    }

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/counting@channel");

    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn guild(context: utils::Context<'_>) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let guild_id = crate::utils::get_guild_id(&context);

    let results = sqlx::query_as::<_, crate::utils::database::CountingStats>(
        "SELECT * FROM counting_stats WHERE guild_id = $1",
    )
    .bind(guild_id)
    .fetch_all(pool)
    .await;

    let mut embed = crate::utils::embeds::create(
        serenity::Colour::new(6591981),
        format!("Stats for **{}**", context.guild().unwrap().name).as_str(),
        "",
        true,
    );

    for result in results.expect("Failed to get results in handlers/counting@guild") {
        embed = embed.field(
            "",
            format!("__Stats for <@{}>__:\n* **Correct** Counts: `{}`.\n* **Incorrect** Counts: `{}`.\n* **Deleted** Counts: `{}`.\n* **Edited** Counts: `{}`.", result.user_id, result.correct, result.incorrect, result.deleted_count_message, result.edited_count_message),
            true,
        );
    }

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/counting@guild");

    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn user(context: utils::Context<'_>, user: serenity::User) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let user_id = i64::from(user.id);

    let results = sqlx::query_as::<_, crate::utils::database::CountingStats>(
        "SELECT * FROM counting_stats WHERE user_id = $1",
    )
    .bind(user_id)
    .fetch_all(pool)
    .await;

    let mut embed = crate::utils::embeds::create(serenity::Colour::new(6591981), "Stats", "", true);

    for result in results.expect("Failed to get results in handlers/counting@user") {
        embed = embed.field(
            "",
            format!("__Stats for <#{}>__:\n* **Correct** Counts: `{}`.\n* **Incorrect** Counts: `{}`.\n* **Deleted** Counts: `{}`.\n* **Edited** Counts: `{}`.", result.channel_id, result.correct, result.incorrect, result.deleted_count_message, result.edited_count_message),
            true,
        );
    }

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/counting@user");

    Ok(())
}
