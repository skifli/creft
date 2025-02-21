use poise::serenity_prelude as serenity;

pub fn create(
    colour: serenity::Colour,
    title: &str,
    description: &str,
    generic_footer: bool,
) -> serenity::CreateEmbed {
    let mut embed = serenity::CreateEmbed::default()
        .colour(colour)
        .title(title)
        .description(description);

    if generic_footer {
        embed = embed.footer(serenity::CreateEmbedFooter::new(
            "Run **/about** for more information about the bot.",
        ));
    }

    embed
}

pub fn insufficient_permissions() -> serenity::CreateEmbed {
    create(
        serenity::Colour::new(13789294),
        "Missing Permissions",
        "You do not have **sufficient permissions** to run this command.\nYou need to be on the **admins list**.",
        true,
    )
}
