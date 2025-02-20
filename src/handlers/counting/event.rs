use crate::utils;
use poise::serenity_prelude as serenity;

pub async fn message_create(
    context: &serenity::Context,
    message: &serenity::Message,
    data: &utils::ServerData,
) -> Result<(), utils::Error> {
    let pool = &data.pool;
    let author_id = i64::from(message.author.id);
    let channel_id = i64::from(message.channel_id);
    let message_id = i64::from(message.id);

    if message.author.bot {
        return Ok(());
    }

    let counting_channel = utils::database::get_counting_channel(pool, channel_id).await?;

    if let Some(mut counting_channel) = counting_channel {
        let expression = pemel::prelude::Expr::parse(message.content.as_str(), false);

        if let Ok(expression) = expression {
            let result = expression.eval_const();

            if counting_channel.last_count_user_id == author_id {
                utils::message::delete_after(context, message, 0);

                utils::message::send_and_delete(
                    context,
                    message.channel_id,
                    utils::embeds::create(
                        serenity::Colour::new(13789294),
                        "Error",
                        &format!(
                            "Sorry <@{}>, you can't count **twice** in a *row*.",
                            author_id
                        ),
                    ),
                    utils::message::MESSAGE_DELETE_DELAY,
                )
                .await;
            } else if let Ok(result) = result {
                utils::database::create_counting_stats(
                    pool,
                    channel_id,
                    author_id,
                    i64::from(message.guild_id.unwrap()),
                )
                .await?;

                let mut counting_stats =
                    utils::database::get_counting_stats(pool, channel_id, author_id)
                        .await?
                        .unwrap();

                if result == (counting_channel.count + 1) as f32 {
                    counting_channel.count += 1;

                    if counting_channel.count > counting_channel.count_max {
                        counting_channel.count_max = counting_channel.count;
                    }

                    message.react(context, '✅').await?;

                    counting_stats.correct += 1;
                } else {
                    message.react(context, '❌').await?;

                    utils::message::send_with_reference(
                            context,
                            utils::embeds::create(
                                serenity::Colour::new(13789294),
                                "Error",
                                &format!(
                                    "Sorry <@{}> - the *correct* number was `{}`, but *you* said `{}`.\nThe count has *reset* to `0`.",
                                    author_id,
                                    counting_channel.count+1,
                                    result,
                                ),
                            ),
                            message,
                        ).await;

                    counting_channel.count = 0;
                    counting_channel.resets_count += 1;

                    counting_stats.incorrect += 1;
                }

                counting_channel.last_count_message_edited = false;
                counting_channel.last_count_user_id = author_id;
                counting_channel.last_count_message_id = message_id;

                utils::database::update_counting_channel(pool, &counting_channel).await?;
                utils::database::update_counting_stats(pool, &counting_stats).await?;
            } else {
                utils::message::delete_after(context, message, 0);
            }
        } else {
            utils::message::delete_after(context, message, 0);
        }
    }

    Ok(())
}

pub async fn message_delete(
    context: &serenity::Context,
    channel_id: serenity::ChannelId,
    deleted_message_id: serenity::MessageId,
    data: &utils::ServerData,
) -> Result<(), utils::Error> {
    let pool = &data.pool;
    let deleted_message_id = i64::from(deleted_message_id);

    let channel_id_int = i64::from(channel_id);

    let counting_channel = utils::database::get_counting_channel(pool, channel_id_int).await?;

    if let Some(counting_channel) = counting_channel {
        if counting_channel.last_count_message_id == deleted_message_id {
            channel_id
                .send_message(
                    context,
                    serenity::CreateMessage::default().embed(utils::embeds::create(
                        serenity::Colour::new(6591981),
                        "I saw that",
                        &format!(
                            "<@{}> **deleted** their message. The count is at `{}`.",
                            counting_channel.last_count_user_id, counting_channel.count
                        ),
                    )),
                )
                .await?;

            let mut counting_stats = utils::database::get_counting_stats(
                pool,
                channel_id_int,
                counting_channel.last_count_user_id,
            )
            .await?
            .unwrap();

            counting_stats.deleted_count_message += 1;

            utils::database::update_counting_stats(pool, &counting_stats).await?;
        }
    }

    Ok(())
}

pub async fn message_update(
    context: &serenity::Context,
    new: &Option<serenity::Message>,
    data: &utils::ServerData,
) -> Result<(), utils::Error> {
    let pool = &data.pool;

    if new.is_none() {
        return Ok(());
    }

    let new = new.as_ref().unwrap();
    let channel_id = i64::from(new.channel_id);

    let counting_channel = utils::database::get_counting_channel(pool, channel_id).await?;

    if let Some(mut counting_channel) = counting_channel {
        if !counting_channel.last_count_message_edited
            && counting_channel.last_count_message_id == i64::from(new.id)
        {
            new.channel_id
                .send_message(
                    context,
                    serenity::CreateMessage::default()
                        .embed(utils::embeds::create(
                            serenity::Colour::new(6591981),
                            "I saw that",
                            &format!(
                                "<@{}> **edited** their message. The count is at `{}`.",
                                counting_channel.last_count_user_id, counting_channel.count
                            ),
                        ))
                        .reference_message(new),
                )
                .await?;

            counting_channel.last_count_message_edited = true;
            utils::database::update_counting_channel(pool, &counting_channel).await?;

            let mut counting_stats = utils::database::get_counting_stats(
                pool,
                channel_id,
                counting_channel.last_count_user_id,
            )
            .await?
            .unwrap();

            counting_stats.edited_count_message += 1;

            utils::database::update_counting_stats(pool, &counting_stats).await?;
        }
    }

    Ok(())
}
