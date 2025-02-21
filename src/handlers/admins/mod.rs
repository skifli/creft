use crate::utils;
use poise::serenity_prelude as serenity;

#[poise::command(
    prefix_command,
    slash_command,
    subcommands("add", "list", "remove"),
    subcommand_required
)]
pub async fn admins(_: utils::Context<'_>) -> Result<(), utils::Error> {
    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn add(context: utils::Context<'_>, user: serenity::User) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let author_id = i64::from(context.author().id);
    let guild_id = crate::utils::get_guild_id(&context);

    let admins = utils::database::cache_admins(pool, guild_id)
        .await
        .expect("Failed to execute cache admins in handlers/admins@add");

    let embed;

    if crate::utils::database::is_admin(&admins, author_id) {
        let user_id = i64::from(user.id);

        if crate::utils::database::is_admin(&admins, user_id) {
            embed = utils::embeds::create(
                serenity::Colour::new(13789294),
                "Error",
                format!("<@{}> is already an admin.", user.id).as_str(),
                true,
            );
        } else {
            utils::database::add_admin(pool, guild_id, user_id)
                .await
                .expect("Failed to execute add admin in handlers/admins@add");

            embed = utils::embeds::create(
                serenity::Colour::new(5082199),
                "Success",
                format!("<@{}> is now an admin.", user.id).as_str(),
                true,
            );
        }
    } else {
        embed = utils::embeds::insufficient_permissions();
    }

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/admins@add");

    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn list(context: utils::Context<'_>) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let guild_id = crate::utils::get_guild_id(&context);

    let admins = utils::database::cache_admins(pool, guild_id)
        .await
        .expect("Failed to execute cache admins in handlers/admins@list");

    let embed = utils::embeds::create(
        serenity::Colour::new(6591981),
        "Admins",
        admins
            .iter()
            .map(|id| format!("* <@{}>", id))
            .collect::<Vec<String>>()
            .join("\n")
            .as_str(),
        true,
    );

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/admins@list");

    Ok(())
}

#[poise::command(prefix_command, slash_command)]
pub async fn remove(context: utils::Context<'_>, user: serenity::User) -> Result<(), utils::Error> {
    let pool = crate::utils::get_pool(&context);
    let author_id = i64::from(context.author().id);
    let guild_id = crate::utils::get_guild_id(&context);

    let admins = utils::database::cache_admins(pool, guild_id)
        .await
        .expect("Failed to execute cache admins in handlers/admins@remove");

    let embed;

    if crate::utils::database::is_admin(&admins, author_id) {
        let user_id = i64::from(user.id);

        if crate::utils::database::is_admin(&admins, user_id) {
            utils::database::remove_admin(pool, guild_id, user_id)
                .await
                .expect("Failed to execute remove admin in handlers/admins@remove");

            embed = utils::embeds::create(
                serenity::Colour::new(5082199),
                "Success",
                format!("<@{}> is no longer an admin.", user.id).as_str(),
                true,
            );
        } else {
            embed = utils::embeds::create(
                serenity::Colour::new(13789294),
                "Error",
                format!("<@{}> is not an admin.", user.id).as_str(),
                true,
            );
        }
    } else {
        embed = utils::embeds::insufficient_permissions();
    }

    context
        .send(poise::CreateReply::default().embed(embed).ephemeral(true))
        .await
        .expect("Failed to send message in handlers/admins@remove");

    Ok(())
}
