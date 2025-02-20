use poise::serenity_prelude as serenity;

pub fn create(colour: serenity::Colour, title: &str, description: &str) -> serenity::CreateEmbed {
    serenity::CreateEmbed::default()
        .colour(colour)
        .title(title)
        .description(description)
        .footer(serenity::CreateEmbedFooter::new(
            "Run /about for more information about the bot.",
        ))
}

pub fn insufficient_permissions() -> serenity::CreateEmbed {
    create(
        serenity::Colour::new(13789294),
        "Missing Permissions",
        "You do not have **sufficient permissions** to run this command.\nYou need to be on the **admins list**.",
    )
}
