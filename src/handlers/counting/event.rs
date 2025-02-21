use crate::utils::{self, message};
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

    let counting_channel = utils::database::get_counting_channel(pool, channel_id)
        .await
        .expect("Failed to get counting channel in handlers/counting/event@message_create");

    if let Some(mut counting_channel) = counting_channel {
        let expression = lieval::Expr::new(message.content.as_str());

        if let Ok(expression) = expression {
            if counting_channel.last_count_user_id == author_id {
                utils::message::delete(context, message).await;
            }

            let result = expression.eval();

            if let Ok(result) = result {
                utils::database::create_counting_stats(
                    pool,
                    channel_id,
                    author_id,
                    i64::from(message.guild_id.unwrap()),
                )
                .await
                .expect(
                    "Failed to create counting stats in handlers/counting/event@message_create",
                );

                let mut counting_stats = utils::database::get_counting_stats(
                    pool, channel_id, author_id,
                )
                .await
                .expect("Failed to get counting stats in handlers/counting/event@message_create")
                .expect(
                    "Failed to unwrap counting stats in handlers/counting/event@message_create",
                );

                if result == (counting_channel.count + 1) as f64 {
                    counting_channel.count += 1;

                    if counting_channel.count > counting_channel.count_max {
                        counting_channel.count_max = counting_channel.count;
                    }

                    message.react(context, '✅').await.expect("Failed to react to message with ✅ in handlers/counting/event@message_create");

                    counting_stats.correct += 1;
                } else {
                    message.react(context, '❌').await.expect("Failed to react to message with ❌ in handlers/counting/event@message_create");

                    utils::message::send_with_reference(
                            context,
                            utils::embeds::create(
                                serenity::Colour::new(13789294),
                                "Oh noes!",
                                &format!(
                                    "Sorry <@{}> - the **correct** number was `{}`, but **you** said `{}`.\nThe count has **reset** to `0`.",
                                    author_id,
                                    counting_channel.count+1,
                                    result,
                                ),
                                false
                            ).footer(serenity::CreateEmbedFooter::new(
                                "Ruh roh!",
                            )),
                            message,
                        ).await;

                    counting_channel.count = 0;
                    counting_channel.resets_count += 1;

                    counting_stats.incorrect += 1;
                }

                counting_channel.last_count_message_edited = false;
                counting_channel.last_count_user_id = author_id;
                counting_channel.last_count_message_id = message_id;

                utils::database::update_counting_channel(pool, &counting_channel).await.expect("Failed to update counting channel in handlers/counting/event@message_create");
                utils::database::update_counting_stats(pool, &counting_stats)
                    .await
                    .expect(
                        "Failed to update counting stats in handlers/counting/event@message_create",
                    );
            } else {
                utils::message::delete(context, message).await;
            }
        } else {
            // E.g., division by zero
            utils::message::delete(context, message).await;
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

    let counting_channel = utils::database::get_counting_channel(pool, channel_id_int)
        .await
        .expect("Failed to get counting channel in handlers/counting/event@message_delete");

    if let Some(counting_channel) = counting_channel {
        if counting_channel.last_count_message_id == deleted_message_id {
            channel_id
                .send_message(
                    context,
                    serenity::CreateMessage::default().embed(
                        utils::embeds::create(
                            serenity::Colour::new(6591981),
                            "I saw that",
                            &format!(
                                "<@{}> **deleted** their message. The count is at `{}`.",
                                counting_channel.last_count_user_id, counting_channel.count
                            ),
                            false,
                        )
                        .footer(serenity::CreateEmbedFooter::new("Cheeky!")),
                    ),
                )
                .await
                .expect("Failed to send message in handlers/counting/event@message_delete");

            let mut counting_stats = utils::database::get_counting_stats(
                pool,
                channel_id_int,
                counting_channel.last_count_user_id,
            )
            .await
            .expect("Failed to get counting stats in handlers/counting/event@message_delete")
            .expect("Failed to unwrap counting stats in handlers/counting/event@message_delete");

            counting_stats.deleted_count_message += 1;

            utils::database::update_counting_stats(pool, &counting_stats)
                .await
                .expect(
                    "Failed to update counting stats in handlers/counting/event@message_delete",
                );
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

    let counting_channel = utils::database::get_counting_channel(pool, channel_id)
        .await
        .expect("Failed to get counting channel in handlers/counting/event@message_update");

    if let Some(mut counting_channel) = counting_channel {
        if !counting_channel.last_count_message_edited
            && counting_channel.last_count_message_id == i64::from(new.id)
        {
            new.channel_id
                .send_message(
                    context,
                    serenity::CreateMessage::default()
                        .embed(
                            utils::embeds::create(
                                serenity::Colour::new(6591981),
                                "I saw that",
                                &format!(
                                    "<@{}> **edited** their message. The count is at `{}`.",
                                    counting_channel.last_count_user_id, counting_channel.count
                                ),
                                false,
                            )
                            .footer(serenity::CreateEmbedFooter::new("Cheeky!")),
                        )
                        .reference_message(new),
                )
                .await
                .expect("Failed to send message in handlers/counting/event@message_update");

            counting_channel.last_count_message_edited = true;
            utils::database::update_counting_channel(pool, &counting_channel)
                .await
                .expect(
                    "Failed to update counting channel in handlers/counting/event@message_update",
                );

            let mut counting_stats = utils::database::get_counting_stats(
                pool,
                channel_id,
                counting_channel.last_count_user_id,
            )
            .await
            .expect("Failed to get counting stats in handlers/counting/event@message_update")
            .expect("Failed to unwrap counting stats in handlers/counting/event@message_update");

            counting_stats.edited_count_message += 1;

            utils::database::update_counting_stats(pool, &counting_stats)
                .await
                .expect(
                    "Failed to update counting stats in handlers/counting/event@message_update",
                );
        }
    }

    Ok(())
}

pub async fn reaction_add(
    context: &serenity::Context,
    add_reaction: &serenity::Reaction,
    data: &utils::ServerData,
) -> Result<(), utils::Error> {
    // if X reacted by admin to self message, delete said message

    let pool = &data.pool;
    let message_author_id = i64::from(
        add_reaction
            .message_author_id
            .expect("Failed to get message author ID in handlers/counting/event@reaction_add"),
    );
    let user_id = i64::from(
        add_reaction
            .user_id
            .expect("Failed to get user ID in handlers/counting/event@reaction_add"),
    );

    if message_author_id == utils::BOT_ID {
        if add_reaction.emoji.to_string() == "❌" {
            if crate::utils::database::is_admin(
                &crate::utils::database::cache_admins(
                    pool,
                    i64::from(
                        add_reaction.guild_id.expect(
                            "Failed to get guild ID in handlers/counting/event@reaction_add",
                        ),
                    ),
                )
                .await
                .expect("Failed to execute cache admins in handlers/counting/event@reaction_add"),
                user_id,
            ) {
                utils::message::delete(
                    context,
                    &add_reaction
                        .message(context)
                        .await
                        .expect("Failed to get message in handlers/counting/event@reaction_add"),
                )
                .await;
            }
        }
    }

    Ok(())
}
