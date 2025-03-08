use crate::handlers::counting;
use crate::utils;
use poise::serenity_prelude as serenity;

pub async fn event_handler(
    ctx: &serenity::Context,
    event: &serenity::FullEvent,
    _framework: poise::FrameworkContext<'_, utils::ServerData, utils::Error>,
    data: &utils::ServerData,
) -> Result<(), utils::Error> {
    match event {
        serenity::FullEvent::Ready { data_about_bot, .. } => {
            println!("Logged in as {}.", data_about_bot.user.name);
        }
        serenity::FullEvent::Message { new_message } => {
            counting::event::message_create(ctx, new_message, data)
                .await
                .expect("Failed to execute counting message_create handler");
        }
        serenity::FullEvent::MessageDelete {
            channel_id,
            deleted_message_id,
            guild_id: _,
        } => {
            counting::event::message_delete(ctx, *channel_id, *deleted_message_id, data)
                .await
                .expect("Failed to execute counting message_delete handler");
        }
        serenity::FullEvent::MessageUpdate {
            old_if_available: _,
            new: _,
            event: event,
        } => {
            counting::event::message_update(ctx, event, data)
                .await
                .expect("Failed to execute counting message_update handler");
        }
        serenity::FullEvent::ReactionAdd { add_reaction } => {
            counting::event::reaction_add(ctx, add_reaction, data)
                .await
                .expect("Failed to execute counting reaction_add handler");
        }
        _ => {}
    }
    Ok(())
}
