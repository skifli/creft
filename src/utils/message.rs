use poise::serenity_prelude as serenity;

pub async fn delete(context: &serenity::Context, message: &serenity::Message) {
    if let Err(why) = message.delete(context).await {
        eprintln!("Error deleting message: {:?}", why);
    }
}

/*
pub const MESSAGE_DELETE_DELAY: u64 = 5000;

pub fn delete_after(context: &serenity::Context, message: &serenity::Message, delay: u64) {
    let context = context.clone();
    let message = message.clone();

    tokio::spawn(async move {
        tokio::time::sleep(std::time::Duration::from_millis(delay)).await;

        if let Err(why) = message.delete(&context).await {
            eprintln!("Error deleting message: {:?}", why);
        }
    });
}
*/

/*
pub async fn send_and_delete(
    context: &serenity::Context,
    channel_id: serenity::ChannelId,
    embed: serenity::CreateEmbed,
    delay: u64,
) {
    let sent_message = channel_id
        .send_message(context, serenity::CreateMessage::default().embed(embed))
        .await
        .expect("Failed to send message");

    delete_after(context, &sent_message, delay);
}
*/

pub async fn send_with_reference(
    context: &serenity::Context,
    embed: serenity::CreateEmbed,
    original_message: &serenity::Message,
) -> serenity::Message {
    original_message
        .channel_id
        .send_message(
            context,
            serenity::CreateMessage::default()
                .embed(embed)
                .reference_message(original_message),
        )
        .await
        .expect("Failed to send message")
}

/*
pub async fn send_with_reference_and_delete(
    context: &serenity::Context,
    embed: serenity::CreateEmbed,
    original_message: &serenity::Message,
    delay: u64,
) {
    let sent_message = send_with_reference(context, embed, original_message).await;

    delete_after(context, &sent_message, delay);
    delete_after(context, original_message, delay);
}
*/
